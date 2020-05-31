package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gopkg.in/fsnotify.v1"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

type Message struct {
	data map[string]interface{}
	out  chan map[string]interface{}
}

var msgChan chan Message

func init() {
	msgChan = make(chan Message)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/eval", EvalHandler).Methods("POST")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:1337",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	router.HandleFunc("/ws", wsHandler)
	srv.ListenAndServe()
}

// TODO: add EvalDOMHandler

func EvalHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := wsSend(ctx, map[string]interface{}{
		"message_type":     "eval",
		"code":             string(b),
		"highlighted_code": highlight(string(b), "javascript", "terminal", "solarized-dark"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err,
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  err,
		"result": resp["result"],
	})
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	// TODO check origin
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	go wsConnHandler(ws)
}

func wsSend(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	out := make(chan map[string]interface{})
	msgChan <- Message{
		data: data,
		out:  out,
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case val := <-out:
		return val, nil
	}
}

func wsConnHandler(ws *websocket.Conn) {
	logrus.Info("new connection")
	defer ws.Close()

	g, ctx := errgroup.WithContext(context.Background())
	watcher, err := fsnotifyWatcher(ctx, g, ws)
	if err != nil {
		return
	}
	defer watcher.Close()

	window := &sync.Map{}
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()

			case msg := <-msgChan:
				uuid := shortuuid.New()
				window.Store(uuid, msg.out)
				go func(uuid string) {
					select {
					//TODO: config
					case <-time.After(15 * time.Minute):
						window.Delete(uuid)
					case <-ctx.Done():
					}
				}(uuid)

				msg.data["_id"] = uuid
				s, err := json.Marshal(msg.data)
				if err != nil {
					return err
				}
				logrus.Infof("ws-send %s", highlight(string(s), "json", "terminal16m", "rrt"))
				if err := ws.WriteJSON(msg.data); err != nil {
					return err
				}

			}
		}
	})
	g.Go(func() error {
		// ws.SetReadDeadline(time.Now().Add(10 * time.Second))
		for {
			data := map[string]interface{}{}
			if err := ws.ReadJSON(&data); err != nil {
				return err
			}
			s, err := json.Marshal(data)
			if err != nil {
				return err
			}
			logrus.Infof("ws-recv %s", highlight(string(s), "json", "terminal16m", "rrt"))
			// logrus.Infof("html %s", highlight(data["result"].(string), "html", "terminal16m", "solarized"))
			if value, ok := window.Load(data["_id"].(string)); ok == true {
				value.(chan map[string]interface{}) <- data
				window.Delete(data["_id"].(string))
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
	})

	g.Go(func() error {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				_, err := wsSend(ctx, map[string]interface{}{
					"message_type": "ping",
				})
				if err != nil {
					return err
				}
			}
		}
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Error("something goes wrong")
	}

	logrus.Info("close connection")
}

func highlight(source, lexer, formatter, style string) string {
	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine formatter.
	f := formatters.Get(formatter)
	if f == nil {
		f = formatters.Fallback
	}

	// Determine style.
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		logrus.WithError(err).Error("tokenise")
		return source
	}

	w := bytes.NewBufferString("")
	err = f.Format(w, s, it)
	if err != nil {
		logrus.WithError(err).Error("format")
		return source
	}
	return w.String()
}

func fsnotifyWatcher(ctx context.Context, g *errgroup.Group, ws *websocket.Conn) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.WithError(err).Error("fsnotify.NewWatcher")
		return nil, err
	}
	g.Go(func() error {
		refreshWatcher := func(quiet bool) {
			//todo: config
			files, err := filepath.Glob("./*.js")
			if err != nil {
				logrus.WithError(err).Error("filepath.Glob")
				return
			}
			for _, filename := range files {
				if !quiet {
					logrus.WithField("filename", filename).Info("watcher add")
				}
				err = watcher.Add(filename)
				if err != nil {
					logrus.WithError(err).Error("watcher.Add")
					return
				}
			}
		}
		refreshWatcher(false)
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				refreshWatcher(true)
			}
		}
	})
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.WithField("filaneme", event.Name).Info("modified")
					content, err := ioutil.ReadFile(event.Name)
					if err != nil {
						return err
					}
					_, err = wsSend(ctx, map[string]interface{}{
						"message_type":     "eval",
						"code":             string(content),
						"highlighted_code": highlight(string(content), "javascript", "terminal", "solarized-dark"),
					})
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
				logrus.WithError(err).Error("watcher.Errors")
			}
		}
	})
	return watcher, nil
}

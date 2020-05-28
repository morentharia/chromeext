package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gopkg.in/fsnotify.v1"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func main() {
	http.HandleFunc("/ws", wsHandler)
	panic(http.ListenAndServe(":1337", nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	go wsConnHandler(ws)
}

func wsConnHandler(ws *websocket.Conn) {
	logrus.Info("new connection")
	defer ws.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	g, ctx := errgroup.WithContext(context.Background())

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
					m := map[string]interface{}{
						"message":          "eval",
						"code":             string(content),
						"highlighted_code": highlight(string(content), "javascript", "terminal", "dracula"),
					}
					s, err := colorjson.Marshal(m)
					if err != nil {
						return err
					}
					logrus.Infof("ws-send %s", s)
					if err := ws.WriteJSON(m); err != nil {
						return err
					}

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
				logrus.WithError(err).Error("watcher.Errors")
			}
		}
	})

	g.Go(func() error {
		// ws.SetReadDeadline(time.Now().Add(10 * time.Second))
		for {
			m := map[string]interface{}{}
			if err := ws.ReadJSON(&m); err != nil {
				return err
			}
			s, err := colorjson.Marshal(m)
			if err != nil {
				return err
			}
			logrus.Infof("ws-recv %s", s)

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
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

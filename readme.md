https://medium.com/@gouthamj99/build-a-chrome-extension-in-5-steps-f14a19cd3660
https://tech.trustpilot.com/what-i-learned-from-making-a-chrome-extension-51f366ad141
https://60devs.com/hot-reloading-for-chrome-extensions.html

```
find . -name "*.go" | grep -v vendor | entr -r bash -c 'go run main.go'
```

```
http POST http://127.0.0.1:1337/eval < example2.js \
    | jq -r '.result' \
    | html-beautify -s 2 -i \
    | chroma -s lovelace -l html
```

https://public-firing-range.appspot.com/

Injecting a Content Script
https://developer.chrome.com/extensions/devtools#injecting

package main

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "strings"

  "golang.org/x/net/webdav"
)

type methodMux map[string]http.Handler

func (m *methodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  if h, ok := (*m)[r.Method]; ok {
    username, password, ok := r.BasicAuth()
    if !ok {
      w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
      w.WriteHeader(http.StatusUnauthorized)
      return
    }
    ug := os.Getenv("USER")
    for _, v := range strings.Split(ug, ";") {
      u := strings.Split(v, ":")
      fmt.Println("user:", u, username, password)
      if username == u[0] && password == u[1] {
        goto A
      }
    }
    http.Error(w, "WebDAV: need authorized!", http.StatusUnauthorized)
    return

  A:
    h.ServeHTTP(w, r)
  } else {
    http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
  }
}

func main() {
  listen := os.Getenv("LISTEN")
  root := os.Getenv("ROOT")
  prefix := os.Getenv("PREFIX")

  files := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
  webdav := &webdav.Handler{
    Prefix:     prefix,
    FileSystem: webdav.Dir(root),
    LockSystem: webdav.NewMemLS(),
    Logger: func(r *http.Request, err error) {
      if err != nil {
        log.Printf("r=%v; \n err=%v", r, err)
      }
    },
  }
  mux := methodMux(map[string]http.Handler{
    "GET":       files,
    "OPTIONS":   webdav,
    "PROPFIND":  webdav,
    "PROPPATCH": webdav,
    "MKCOL":     webdav,
    "COPY":      webdav,
    "MOVE":      webdav,
    "LOCK":      webdav,
    "UNLOCK":    webdav,
    "DELETE":    webdav,
    "PUT":       webdav,
  })

  fmt.Println("start port", listen)
  if err := http.ListenAndServe(listen, &mux); err != nil {
    log.Fatal(err)
  }
}

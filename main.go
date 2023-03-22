package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "golang.org/x/net/webdav"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strings"
  "webdav/service"
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

  go func() {
    http.HandleFunc("/api/list", func(writer http.ResponseWriter, request *http.Request) {
      request.Header.Set("Content-Type", "application/json; charset=utf-8")
      userRoot := root
      path := request.FormValue("path")
      if path != "" {
        userRoot = fmt.Sprintf("%s/%s", userRoot, path)
      }
      dataList := service.GetFiles(userRoot)
      bytes, err := json.Marshal(dataList)
      if err != nil {
        return
      }
      writer.Write(bytes)
    })
    http.HandleFunc("/api/upload", func(writer http.ResponseWriter, request *http.Request) {
      resString := "上传成功"
      userRoot := root
      path := request.FormValue("path")
      if path != "" {
        userRoot = fmt.Sprintf("%s/%s", userRoot, path)
      }
      request.ParseMultipartForm(32 << 20)
      file, header, err := request.FormFile("file")
      if err != nil {
        fmt.Printf("header:%v, error:%v \n", header, err)
        return
      }
      defer file.Close()
      if file == nil {
        writer.Write(bytes.NewBufferString("文件为空").Bytes())
        return
      }
      fileBytes, err := ioutil.ReadAll(file)
      if err != nil {
        fmt.Printf("header:%v, error:%v \n", header, err)
      }
      fileName := fmt.Sprintf("%s/%s", userRoot, header.Filename)
      err = ioutil.WriteFile(fileName, fileBytes, 0777)
      if err != nil {
        fmt.Printf("header:%v, error:%v \n", header, err)
      }
      writer.Write(bytes.NewBufferString(resString).Bytes())
    })
    http.ListenAndServe("0.0.0.0:81", nil).Error() // 启动默认的 http 服务，可以使用自带的路由
  }()

  fmt.Println("start port", listen)
  if err := http.ListenAndServe(listen, &mux); err != nil {
    log.Fatal(err)
  }
}

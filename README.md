# 个人Webdav服务


```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app main.go
docker build -t ou88zz/webdav:1.1 .

docker run -d --restart=unless-stopped --name=webdav -v /tmp:/webdav -p 8081:8080 -e USER="test:123;ole:123" ou88zz/webdav:1.1
```


FROM alpine:3.6

RUN echo $'http://mirrors.aliyun.com/alpine/v3.6/main\n\
http://mirrors.aliyun.com/alpine/v3.6/community' > /etc/apk/repositories

RUN apk add --update ca-certificates
RUN update-ca-certificates
RUN apk add --update tzdata
ENV TZ=Asia/Shanghai

COPY app /app

RUN chmod +x /app

ENV LISTEN :8080
ENV ROOT /webdav
ENV PREFIX /
ENV USER "test:123456"

EXPOSE 8080/tcp

ENTRYPOINT ["/app"]

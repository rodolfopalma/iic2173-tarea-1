FROM golang

ADD . /go/src/github.com/rpalmaotero/iic2173-tarea-1

RUN go get github.com/mattn/go-sqlite3
RUN go install github.com/rpalmaotero/iic2173-tarea-1

WORKDIR /go/src/github.com/rpalmaotero/iic2173-tarea-1
ENTRYPOINT /go/bin/iic2173-tarea-1

EXPOSE 8080
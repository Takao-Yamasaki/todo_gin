FROM golang:1.19

WORKDIR /go/src/todo_gin

COPY . /go/src/todo_gin/

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "main.go"]

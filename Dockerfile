FROM golang:1.21 as development

WORKDIR /app

COPY . .

EXPOSE 8080

RUN go install github.com/cosmtrek/air@latest && \
    go mod download && \
    go mod verify && \
    go mod tidy

RUN go install -v golang.org/x/tools/gopls@latest

CMD ["air", "-c", ".air.toml"]

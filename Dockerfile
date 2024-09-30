FROM golang:1.23

WORKDIR /app 

# Goモジュールと依存関係を取得
COPY go.mod go.sum ./
RUN go mod download

COPY . . 

RUN go build -o server ./server/main.go 

# 実行権限を付与
RUN chmod +x ./server

CMD ["./server"]
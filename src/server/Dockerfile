# ベースイメージ
FROM golang:1.23

# ワーキングディレクトリの設定
WORKDIR /app

# モジュールファイルをコピーして依存関係をインストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピーしてビルド
COPY . .
RUN go build -o main .

# アプリケーションを起動
CMD ["./main"]

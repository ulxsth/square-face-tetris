# ビルドステージ
FROM golang:1.22.10-bullseye as build

WORKDIR /go/src/app

# アプリケーションのソースコードをコピー
COPY . .

# 必要なパッケージをインストール
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    xorg-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# 環境変数の設定
ENV GO111MODULE=on
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# Go ツールのインストール
RUN go install github.com/ramya-rao-a/go-outline@latest
RUN go install golang.org/x/tools/gopls@v0.12.2

# Go モジュールの整理
RUN go mod tidy

# アプリケーションのビルド
RUN go build -o /go/bin/app app/main.go

# 最終的な実行
CMD ["/go/bin/app"]

# # 本番ステージに bash をインストール
# FROM alpine:3.16 as prod
# # bash をインストール
# RUN apk update && apk add bash
# # ビルドしたアプリケーションをコピー
# COPY --from=build /go/bin/app /app
# # エントリーポイント
# CMD ["/app"]

# docker run -e DISPLAY=$DISPLAY --net=host <コンテナ名>

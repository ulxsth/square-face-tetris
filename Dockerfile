# ビルドステージ
FROM golang:1.22.10-bullseye as build

WORKDIR /go/src/app

# 環境変数の設定
ENV GO111MODULE=on
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

# 必要なパッケージをインストール
RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    xorg-dev \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# go.mod と go.sum を先にコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# Go ツールのインストール
RUN go install github.com/ramya-rao-a/go-outline@latest
RUN go install golang.org/x/tools/gopls@v0.12.2

# Go モジュールの整理
RUN go mod tidy

# 最終的な実行
CMD ["/go/bin/app"]

# docker run -e DISPLAY=$DISPLAY --net=host <コンテナ名>

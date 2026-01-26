FROM golang:1.25.6-alpine3.23

WORKDIR /go/src

COPY ./src .

# go.modがあれば依存関係をインストール
RUN if [ -f ./go.mod ]; then \
      go install ./...; \
    fi

# 開発用のライブラリをインストール
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

EXPOSE 8080
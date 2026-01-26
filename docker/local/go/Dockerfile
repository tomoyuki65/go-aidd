FROM golang:1.25.6-alpine3.23

WORKDIR /go/src

COPY ./src .

# Install dependencies if go.mod exists
RUN if [ -f ./go.mod ]; then \
      go install ./...; \
    fi

# Install development libraries
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

EXPOSE 8080
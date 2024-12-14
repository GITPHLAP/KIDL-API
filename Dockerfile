FROM golang:1.23.4-alpine3.20 AS builder

RUN apk update 

WORKDIR /api

COPY ./ /api

RUN go env -w GO111MODULE=auto \
    && GOOS=linux GOARCH=amd64 go build -o kidl-api .

FROM alpine:3.21
ARG VERSION
ENV VERSION=${VERSION:-develop}

WORKDIR /app
RUN echo ${VERSION} > version

COPY --from=builder /api/kidl-api /app

RUN chmod +x /app/kidl-api
CMD ["/app/kidl-api"]
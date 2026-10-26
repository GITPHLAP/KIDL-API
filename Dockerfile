FROM golang:1.27.3-alpine3.24 AS builder

RUN apk update 

WORKDIR /api

COPY ./ /api

RUN go env -w GO111MODULE=auto \
    && GOOS=linux GOARCH=amd64 go build -o kidl-api .

FROM alpine:3.24
# Needed for healthcheck
RUN apk add curl 
ARG VERSION
ENV VERSION=${VERSION:-develop}

WORKDIR /app
RUN echo ${VERSION} > version

COPY --from=builder /api/kidl-api /app

RUN chmod +x /app/kidl-api
CMD ["/app/kidl-api"]
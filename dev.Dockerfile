FROM golang:1.22.5-alpine3.18 AS builder

RUN apk update 

RUN apk add build-base=0.5-r3 gnutls-dev=3.8.4-r0

WORKDIR /app

COPY ./ /app/

RUN apk update && apk add --no-cache git openssh-client

# RUN go mod download
# RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN go install github.com/air-verse/air@latest && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 2345

CMD ["air"]
#CMD ["doppler", "run", "--", "air", "-build.cmd", "buildvcs=false"]
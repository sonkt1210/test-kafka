FROM golang:1.17-alpine as builder

WORKDIR /app

RUN pwd

COPY . .
#RUN GOOS=linux GOARCH=amd64 go build -a -v -tags musl
RUN apk add build-base
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go mod download  && go mod vendor

RUN ls

RUN go build  -tags musl -o main main.go


#docker run -v $(pwd):/src -it test-kafka
#
FROM alpine:3.7
#
#RUN apk update && apk add ca-certificates tzdata && rm -rf /var/cache/apk/*
#
WORKDIR /app
COPY --from=builder /main /app/main
CMD ["./main"]
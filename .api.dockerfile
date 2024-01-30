FROM golang:alpine AS builder

WORKDIR /go/src/note

COPY . .

RUN mkdir bin
RUN go mod init note
RUN go get -v -d ./...
RUN go install -v ./...

FROM alpine:latest

WORKDIR /app_binary
COPY --from=builder /go/bin/note /app_binary/
COPY --from=builder /go/src/note/front /app_binary/front
COPY --from=builder /go/src/note/.env /app_binary/
RUN chmod +x ./note
ENTRYPOINT ./note

CMD ["note"]
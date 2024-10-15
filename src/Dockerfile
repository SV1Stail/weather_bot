FROM golang:1.22.2-alpine AS builder
RUN apk add --no-cache git 

WORKDIR /app

COPY . .

RUN go mod download && go build -o main ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash vim
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["bash", "-c", " ./main"]
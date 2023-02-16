FROM golang:alpine AS builder
WORKDIR /service1
COPY . .
RUN go build -o salt main.go

FROM alpine:3.6
EXPOSE 8082
WORKDIR /service1
COPY --from=builder /service1 .
CMD [ "/service1/salt" ]
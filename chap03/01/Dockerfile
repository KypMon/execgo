FROM golang:alpine as builder

RUN mkdir /app
WORKDIR /app
COPY ./main.go ./main.go
RUN go build ./main.go

FROM alpine
RUN mkdir /app
WORKDIR /app
COPY --from=builder /app/main ./main
#ENV version=1.3.0
EXPOSE 8080
CMD ["./main"]
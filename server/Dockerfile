FROM golang:1.20 as builder
WORKDIR /workspace

COPY ./ /workspace/

RUN go build -o server main.go

FROM ubuntu
COPY --from=builder /workspace/server .
RUN chmod +x ./server

ENTRYPOINT [ "./server" ]
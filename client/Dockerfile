FROM golang:1.20 as builder
WORKDIR /workspace

COPY ./ /workspace/

RUN go build -o client main.go

FROM ubuntu
COPY --from=builder /workspace/client .
RUN chmod +x ./client

ENTRYPOINT [ "./client" ]
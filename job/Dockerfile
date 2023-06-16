FROM golang:1.20 as builder
WORKDIR /workspace

COPY ./ /workspace/

RUN go build -o job main.go

FROM ubuntu
COPY --from=builder /workspace/job .
RUN chmod +x ./job

ENTRYPOINT [ "./job" ]
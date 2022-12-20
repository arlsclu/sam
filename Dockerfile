FROM golang:1.19:buster as builder

WORKDIR /app
COPY go.mod  .
COPY go.sum .
RUN  go mod  download
COPY *.go .

RUN go build -o sam

FROM alpine:latest 
COPY --from=builder /app/sam   /
ENTRYPOINT [ "/sam" ]


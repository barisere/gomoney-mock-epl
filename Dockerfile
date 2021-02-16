FROM golang:1.15-alpine AS builder

COPY . /gomoney/
WORKDIR /gomoney
ENV GO111MODULE=on
RUN go build -o mock-epl
RUN go get -v github.com/markbates/grift

FROM alpine

COPY --from=builder /gomoney/mock-epl .
COPY docs/ ./docs/
CMD [ "./mock-epl" ]

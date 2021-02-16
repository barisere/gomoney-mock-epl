FROM golang:1.15-alpine AS builder

COPY . /gomoney/
WORKDIR /gomoney
RUN go build -o mock-epl

FROM alpine

COPY --from=builder /gomoney/mock-epl .
CMD [ "./mock-epl" ]

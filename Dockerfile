FROM golang:1.22

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

WORKDIR /hati

COPY . .

COPY entrypoint.sh /app/entrypoint.sh

RUN make build

EXPOSE 4242
EXPOSE 6767

ENTRYPOINT ["/hati/entrypoint.sh"]
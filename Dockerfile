FROM golang:1.9-alpine

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build

EXPOSE 8090

ENTRYPOINT ["app/go-stan"]

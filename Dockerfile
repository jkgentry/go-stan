FROM alpine:3.7

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build

EXPOSE 8090

ENTRYPOINT ["app/go-stan"]

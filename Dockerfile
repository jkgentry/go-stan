FROM alpine:3.7

RUN mkdir /app
ADD go-stan /app/
ADD marvel-characters.json /app/

WORKDIR /app

EXPOSE 8090

ENTRYPOINT ["app/go-stan"]

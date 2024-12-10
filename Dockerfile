FROM golang:latest

WORKDIR /app

COPY . ./

RUN apt update && apt install -y make && make

EXPOSE 8080

CMD [ "./server" ]

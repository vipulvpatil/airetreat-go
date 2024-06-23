FROM golang:1.18-alpine

WORKDIR /airetreat-go

COPY . .

RUN apk update && apk add git

RUN go mod download

RUN go build -o ./bin/airetreatgo

EXPOSE 9100

CMD [ "./bin/airetreatgo" ]

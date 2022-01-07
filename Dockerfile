# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

LABEL site="computing"
LABEL stage="builder"

WORKDIR /src/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

#RUN go get github.com/ystv/computing-site/team
#RUN go get github.com/ystv/computing-site/team

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o computing

EXPOSE 7075

CMD ["/computing"]
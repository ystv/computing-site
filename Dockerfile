# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

#LABEL site="computing"
#LABEL stage="builder"

WORKDIR /src/

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate

COPY *.go ./

RUN go get github.com/ystv/computing_site/team
RUN go get github.com/ystv/computing_site/templates
RUN go build github.com/ystv/computing_site/team
RUN go build github.com/ystv/computing_site/templates

RUN go build -o computing

EXPOSE 7075

CMD ["computing"]
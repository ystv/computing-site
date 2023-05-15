FROM golang:1.20.4-alpine3.18

LABEL site="computing-site"
LABEL stage="builder"

WORKDIR /src/

COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate

COPY *.go ./

RUN apk update && apk add git

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/computing

EXPOSE 7075

ENTRYPOINT ["/bin/computing"]
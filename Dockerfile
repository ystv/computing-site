FROM golang:1.20.4-alpine3.18

LABEL site="computing-site"
LABEL stage="builder"

WORKDIR /src/

ARG COMP_SITE_VERSION_ARG

COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate

COPY *.go ./

RUN apk update && apk add git

# Set build variables
RUN echo -n "-X 'main.Version=$COMP_SITE_VERSION_ARG" > ./ldflags && \
    tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags && \
    echo -n "'" >> ./ldflags

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(cat ./ldflags)" -o /bin/computing

EXPOSE 7075

ENTRYPOINT ["/bin/computing"]
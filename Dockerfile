FROM golang:1.24-alpine as build

LABEL site="computing-site"
LABEL stage="builder"

WORKDIR /src/

ARG COMP_SITE_VERSION_ARG
ARG COMP_SITE_COMMIT_ARG

ARG COMP_SITE_CERT_PEM
ARG COMP_SITE_KEY_PEM

COPY go.mod ./
COPY go.sum ./
COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go generate

COPY *.go ./

RUN apk update && apk upgrade

RUN echo $COMP_SITE_CERT_PEM | awk '{gsub(/~/,"\n")}1' > cert.pem
RUN echo $COMP_SITE_KEY_PEM | awk '{gsub(/~/,"\n")}1' > key.pem

# Set build variables
RUN echo -n "-X 'main.Version=$COMP_SITE_VERSION_ARG" > ./ldflags && \
    tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags && \
    echo -n "' -X 'main.Commit=$COMP_SITE_COMMIT_ARG" >> ./ldflags && \
    tr -d \\n < ./ldflags > ./temp && mv ./temp ./ldflags && \
    echo -n "'" >> ./ldflags

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="$(cat ./ldflags)" -o /bin/computing

FROM scratch
LABEL site="computing"

COPY --from=build /bin/computing /bin/computing

EXPOSE 7075

ENTRYPOINT ["/bin/computing"]

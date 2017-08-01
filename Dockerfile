FROM golang:1.8
WORKDIR /go/src/app
COPY . .
RUN go-wrapper download github.com/golang/lint/golint
RUN go-wrapper install github.com/golang/lint/golint
CMD ["go-wrapper", "run"]

FROM golang:alpine

WORKDIR /go/src/accounting
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["accounting"]



# FROM golang:alpine as builder
# RUN mkdir /build
# ADD . /build/
# WORKDIR /build
# RUN go build -o main .
# FROM alpine
# COPY --from=builder /build/main /app/
# WORKDIR /app
# CMD ["./main"]

# FROM golang:1.9.6-alpine3.7
# WORKDIR /go/src/accounting
# COPY . .
# RUN apk add --no-cache git
# RUN go-wrapper download   # "go get -d -v ./..."
# RUN go-wrapper install    # "go install -v ./..."

# #final stage
# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# COPY --from=builder /go/bin/accounting /accounting
# ENTRYPOINT ./accounting
# LABEL Name=cloud-native-go Version=0.0.1
# EXPOSE 2000



# FROM golang:alpine
# RUN mkdir /accounting
# ADD . /accounting/
# WORKDIR /accounting
# RUN go build -o main .
# CMD ["./main"]

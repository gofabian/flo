FROM golang:1.15-alpine3.12 as build

WORKDIR /build
RUN echo $GOPATH

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build flo.go


FROM alpine:3.12

RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app
ENV PATH="/app:$PATH"

COPY --from=build /build/flo .
#COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#COPY .git .
#COPY . .

CMD ["flo", "-h"]

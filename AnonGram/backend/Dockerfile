FROM golang:1.23-alpine AS build
WORKDIR /src
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o /out/anongram .
FROM alpine:3.20
COPY --from=build /out/anongram /anongram
EXPOSE 8080
ENTRYPOINT ["/anongram"]

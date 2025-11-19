FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /bin/notification ./cmd

FROM alpine:3.20
COPY --from=build /bin/notification /notification
EXPOSE 8080
CMD ["/notification"]

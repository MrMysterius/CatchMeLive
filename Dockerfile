FROM golang:1.23 AS build

ENV CGO_ENABLED=0

COPY . /src
WORKDIR /src

RUN go build -v -o /bin/catch-me-live ./src
RUN chmod +x /bin/catch-me-live

# FROM debian:bookworm-slim
FROM scratch
COPY --from=build /bin/catch-me-live /bin/catch-me-live
COPY --from=build /etc/ssl /etc/ssl
CMD ["/bin/catch-me-live"]

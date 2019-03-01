FROM golang:1.11 as builder

WORKDIR /build
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM scratch

WORKDIR /app/
COPY --from=builder /build/app .
CMD ["./app"]

EXPOSE 80
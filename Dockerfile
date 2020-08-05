FROM golang:1.14 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make clean && make

FROM gcr.io/distroless/static

USER nonroot
COPY --from=builder /build/bin/billiam /bin/billiam
COPY --from=builder /build/config.toml.example /srv/config.toml

WORKDIR /srv
ENTRYPOINT ["billiam"]
CMD ["serve"]
EXPOSE 2490 2491

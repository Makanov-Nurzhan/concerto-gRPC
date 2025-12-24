FROM golang:1.24 AS build
LABEL authors="makanov_n.cpi"

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/grpc-server ./cmd/server

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=build /out/grpc-server /app/grpc-server

EXPOSE 50051

USER 65532:65532
ENTRYPOINT ["/app/grpc-server"]
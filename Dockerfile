FROM golang:1.22 as builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go mod tidy   
RUN CGO_ENABLED=0 go build -o /swift-service ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /swift-service /swift-service
COPY Interns_2025_SWIFT_CODES.xlsx /data/
EXPOSE 8080
ENTRYPOINT ["/swift-service"]

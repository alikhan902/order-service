FROM golang:1.22.12-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/order-service .

FROM scratch
COPY --from=build /out/order-service /order-service
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/order-service"]

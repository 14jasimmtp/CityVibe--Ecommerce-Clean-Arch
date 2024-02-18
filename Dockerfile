FROM golang:1.21.3-alpine3.18 AS build-stage

WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./cityvibe ./cmd/api

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM scratch AS build-release-stage
WORKDIR /
COPY --from=build-stage /app/cityvibe /cityvibe
COPY --from=build-stage /app/.env /
EXPOSE 3000

ENTRYPOINT ["/cityvibe"]

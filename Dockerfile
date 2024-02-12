FROM golang:1.22-alpine3.19

WORKDIR /app
COPY . ./
RUN go mod download
RUN go build -o cityvibe ./cmd/api/main.go
EXPOSE 3000

CMD [ "./cityvibe" ]

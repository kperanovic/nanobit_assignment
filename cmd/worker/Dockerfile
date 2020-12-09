FROM golang:alpine AS build

COPY . /src/

WORKDIR /src/cmd/worker

RUN go mod download

RUN mkdir /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/worker .

FROM scratch

WORKDIR /app
COPY --from=build /app/worker /app/

CMD ["/app/worker"]

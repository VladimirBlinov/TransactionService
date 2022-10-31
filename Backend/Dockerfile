# step 1
FROM golang:1.19-alpine3.16 AS build_step

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

# build go app
RUN go mod download

COPY ./ ./
RUN apk add make
RUN make build_api

#step 2
FROM alpine
WORKDIR /app
COPY --from=build_step /app ./
RUN chmod +x ./api

EXPOSE 8080
CMD ["./api"]

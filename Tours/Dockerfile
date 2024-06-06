FROM golang:alpine AS tour_build_container
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o tour-service

FROM alpine
COPY --from=tour_build_container /app/tour-service /usr/bin
EXPOSE 8004
ENTRYPOINT ["tour-service"]
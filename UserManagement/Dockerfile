FROM golang:alpine as user_management_build_container

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o user_management_service

FROM alpine
COPY --from=user_management_build_container /app/user_management_service /usr/bin
EXPOSE 8003
ENTRYPOINT ["user_management_service"]
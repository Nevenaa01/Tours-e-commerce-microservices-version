FROM golang:alpine AS blog_build_container
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o blog-service

FROM alpine
COPY --from=blog_build_container /app/blog-service /usr/bin
EXPOSE 8002
ENTRYPOINT ["blog-service"]
FROM golang:alpine as stakeholder_build_container
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o stakeholder-service

FROM alpine
COPY --from=stakeholder_build_container /app/stakeholder-service /usr/bin
EXPOSE 8001
ENTRYPOINT ["stakeholder-service"]
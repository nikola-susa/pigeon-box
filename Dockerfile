FROM golang:alpine AS build
WORKDIR /src
RUN apk --no-cache add ca-certificates git gcc musl-dev
COPY . .
RUN rm README.md package.json package-lock.json logo.svg tailwind.config.js -r docs
RUN go mod download
RUN GOOS=linux CGO_ENABLED=1 go build -o bot

FROM golang:alpine
RUN apk upgrade -aU --no-cache
WORKDIR /app
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /src/bot .

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8081
EXPOSE ${SERVER_PORT}

CMD ["/app/bot"]

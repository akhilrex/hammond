ARG GO_VERSION=1.16.2
FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
RUN mkdir -p /api
WORKDIR /api
COPY ./server/go.mod .
COPY ./server/go.sum .
RUN go mod download
COPY ./server .
RUN go build -o ./app ./main.go

FROM node:14 as build-stage
WORKDIR /app
COPY ./ui/package*.json ./
RUN npm install
COPY ./ui .
RUN npm run build


FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/akhilrex/hammond"
ENV CONFIG=/config
ENV DATA=/assets
ENV UID=998
ENV PID=100
ENV GIN_MODE=release
VOLUME ["/config", "/assets"]
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir -p /config; \
    mkdir -p /assets; \
    mkdir -p /api
RUN chmod 777 /config; \
    chmod 777 /assets
WORKDIR /api
COPY --from=builder /api/app .
#COPY dist ./dist
COPY --from=build-stage /app/dist ./dist
EXPOSE 3000
ENTRYPOINT ["./app"]

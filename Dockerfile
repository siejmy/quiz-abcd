FROM golang:1.13 as builderbackend
WORKDIR /app
COPY backend/go.* ./
RUN go mod download
COPY backend ./
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server




FROM node:14 as builderfrontend
WORKDIR /app
COPY frontend/package* ./
RUN npm ci
COPY frontend ./
RUN npm run build-lib




FROM alpine:3
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builderbackend /app/ /app/
COPY --from=builderfrontend /app/dist/ /app/static/frontend/
CMD ["/app/server"]

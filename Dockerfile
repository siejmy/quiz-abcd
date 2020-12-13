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
ADD serviceAccount.json /serviceAccount_o5z3A5q1.json
ENV ROUTE_BASE="abcd"
ENV GOOGLE_APPLICATION_CREDENTIALS="/serviceAccount_o5z3A5q1.json"
ENV RE_CAPTCHA_KEY="6Le8g9wZAAAAAP4IGzG3h1eSg-59l3m9t59qrJ8M"
ENV FACEBOOK_APP_ID="357764061986952"
COPY demo .
RUN ls static
CMD ["/app/server"]

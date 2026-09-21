# Atlas standard Dockerfile (ADR atlas-templates/0002).
# Multi-stage build: compiles a static binary in the official Go image and runs it
# on a minimal, shell-less base image (distroless), as a non-root user.

# --- build ---
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/app ./...

# --- runtime ---
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/app /app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]

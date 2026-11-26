# Atlas standard Dockerfile (ADR atlas-templates/0002).
# Multi-stage build: compiles a static binary in the official Go image and runs it
# on a minimal, shell-less base image (distroless), as a non-root user.

# --- build ---
# Runs on the builder's own platform and cross-compiles for the target (no QEMU), so one
# build publishes linux/amd64 and linux/arm64: GitHub runners are amd64, while the local
# kind node on Apple Silicon is arm64.
FROM --platform=$BUILDPLATFORM golang:1.26 AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -o /out/app ./...

# --- runtime ---
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/app /app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]

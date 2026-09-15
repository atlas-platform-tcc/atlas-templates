# Dockerfile padrão do Atlas (ADR atlas-templates/0002).
# Build multi-estágio: compila estático em imagem oficial de Go e executa em
# imagem base mínima e sem shell (distroless), como usuário não-root.

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

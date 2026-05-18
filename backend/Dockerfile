# Stage 1: Builder
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Instalar dependencias
RUN apk add --no-cache git make

# Copiar go mod files
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY cmd cmd
COPY internal internal

# Compilar
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o backend ./cmd/main.go

# Stage 2: Runtime
FROM alpine:3.18

WORKDIR /app

# Instalar certificados SSL
RUN apk --no-cache add ca-certificates

# Copiar binary del builder
COPY --from=builder /build/backend .

# Crear directorios necesarios
RUN mkdir -p storage/uploads

# Exponer puerto
EXPOSE 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Comando de inicio
CMD ["./backend"]

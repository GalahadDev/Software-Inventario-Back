# -----------------------
# Etapa 1: build de Go
# -----------------------
    FROM golang:1.20-alpine AS builder

    # Crear directorio de trabajo en el contenedor
    WORKDIR /app
    
    # Copiar archivos de dependencias
    COPY go.mod go.sum ./
    
    # Descargar dependencias
    RUN go mod download
    
    # Copiar el resto del código
    COPY . .
    
    # Compilar la aplicación
    # Asegúrate de que main.go escuche en el puerto 8080 (o lea os.Getenv("PORT"))
    RUN go build -o server main.go
    
    # -----------------------
    # Etapa 2: imagen final
    # -----------------------
    FROM alpine:3.17
    
    # Crear directorio de trabajo en la imagen final
    WORKDIR /app
    
    # Copiar el binario compilado desde la etapa builder
    COPY --from=builder /app/server /app/
    
    # Exponer el puerto 8080 (Opcional, útil a modo documental)
    EXPOSE 8080
    
    # Comando de arranque
    CMD ["/app/server"]
    

# Panduan Deployment & Infrastruktur

Dokumen ini memandu proses pemindahan aplikasi dari lingkungan pengembangan lokal (VS Code & Laragon) ke server *production* (seperti VPS Ubuntu) menggunakan arsitektur container.

### 1. Arsitektur Server Production
- **Reverse Proxy:** Nginx (untuk menangani SSL/TLS dan *routing* domain).
- **Backend:** Go (Gin) dikompilasi menjadi *binary* statis.
- **Frontend:** Next.js berjalan di atas Node.js.
- **Database:** MySQL 8.0.

### 2. Konfigurasi Docker (Direkomendasikan)
Gunakan Docker untuk memastikan aplikasi berjalan konsisten di lingkungan mana pun.

#### Dockerfile Backend (Go)
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
# Pastikan folder penyimpanan dibuat
RUN mkdir -p /uploads 
EXPOSE 8080
CMD ["./main"]
```

### 3. Menjalankan Layanan (Docker Compose)
Buat file `docker-compose.yml` di *root* direktori untuk menjalankan database, backend, dan frontend secara bersamaan.
```yaml
version: '3.8'
services:
  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password_production
      MYSQL_DATABASE: storage_db
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  backend:
    build: 
      context: ./backend
    ports:
      - "8080:8080"
    depends_on:
      - db
    volumes:
      - ./uploads:/uploads

volumes:
  mysql_data:
```
*Jalankan perintah:* `docker-compose up -d --build`

### 4. Keamanan Infrastruktur
- **CORS:** Di *production*, ubah konfigurasi Gin CORS dari `AllowAllOrigins` menjadi domain spesifik Next.js Anda (misal: `https://storage.domainanda.com`).
- **Path Storage:** Jangan pernah meletakkan direktori `/uploads` fisik di dalam folder statis publik (`/public` Next.js). Pastikan hanya Gin yang memiliki akses baca/tulis ke direktori tersebut melalui abstraksi API.

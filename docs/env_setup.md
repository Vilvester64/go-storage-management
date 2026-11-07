# Panduan Setup Environment (.env)

Agar sistem dapat berjalan, Anda perlu membuat file `.env` di direktori utama (*root*) proyek. 

Berikut adalah referensi kunci (*keys*) yang dibutuhkan beserta contoh nilainya untuk lingkungan pengembangan (lokal):

```ini
# --- SERVER CONFIGURATION ---
# Port di mana aplikasi Gin akan berjalan
PORT=8080

# Mode Gin: "debug" untuk development, "release" untuk production
GIN_MODE=debug


# --- DATABASE CONFIGURATION (MySQL via Laragon) ---
# Nama pengguna database MySQL (default di Laragon biasanya 'root')
DB_USER=root

# Kata sandi database (biarkan kosong jika tidak ada password di Laragon)
DB_PASSWORD=

# Host database (lokal)
DB_HOST=127.0.0.1

# Port default MySQL
DB_PORT=3306

# Nama database (harus dibuat manual terlebih dahulu di Laragon)
DB_NAME=storage_db


# --- SECURITY & AUTHENTICATION ---
# Kunci rahasia untuk menandatangani JWT. 
# HARUS diganti dengan string acak yang kuat (misal: generate dengan openssl rand -base64 32)
JWT_SECRET=super_secret_key_development_only

# Masa berlaku token JWT dalam jam
JWT_EXPIRATION_HOURS=24


# --- STORAGE CONFIGURATION ---
# Direktori lokal untuk menyimpan file fisik (akan dibuat otomatis jika belum ada)
STORAGE_PATH=./uploads

# Batas ukuran maksimal upload file (dalam bytes, misal: 50MB = 52428800)
MAX_UPLOAD_SIZE=52428800
```

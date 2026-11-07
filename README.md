# Storage Management System

Sistem manajemen file web full-stack yang aman, efisien, dan modern. Backend dibangun dengan **Golang (Gin)** dan **GORM (MySQL)**, frontend menggunakan **Next.js (App Router)**, **Tailwind CSS**, dan **Framer Motion**.

---

## Fitur Utama

- **Autentikasi & Sesi Aman:**
  - Registrasi & Login dengan password hashing `bcrypt`.
  - Token akses berbasis JSON Web Token (JWT) tersimpan di Cookie & LocalStorage.
  - Middleware otorisasi di backend dan route guard di frontend.
- **Manajemen File Komprehensif:**
  - Upload file tunggal via `multipart/form-data`.
  - Deteksi dan validasi tipe file nyata (*magic bytes*) menggunakan `h2non/filetype` untuk menangkal script berbahaya.
  - Proteksi *Path Traversal*: nama penyimpanan fisik diisolasi dengan UUID (`system_name`).
  - Download file streaming aman via `c.File()`.
  - Soft delete (Trash Bin) berbasis fitur `DeletedAt` GORM tanpa menghapus file fisik secara prematur.
- **UI Modern & Interaktif:**
  - Dashboard statistik penyimpanan real-time.
  - Komponen drag-and-drop file upload.
  - Filter pencarian file instan.
  - Animasi transisi halus menggunakan Framer Motion.

---

## Tech Stack

### Backend
- **Bahasa:** Go (Golang) v1.20+
- **Framework:** Gin Web Framework
- **ORM:** GORM
- **Database:** MySQL 8.0+
- **Autentikasi:** JWT (`golang-jwt/jwt/v5`) & `bcrypt`
- **Deteksi MIME:** `h2non/filetype`
- **Penyimpanan Fisik:** Local File System (`./uploads`)

### Frontend
- **Framework:** Next.js (App Router)
- **Library UI:** React, Tailwind CSS
- **Animasi:** Framer Motion
- **Bahasa:** TypeScript

---

## Struktur Proyek (Clean Architecture)

```text
go-storage-management/
├── cmd/
│   └── api/
│       └── main.go              # Entry point backend Go
├── internal/
│   ├── config/                  # Konfigurasi env & koneksi MySQL
│   ├── handlers/                # HTTP Handlers (Controller Gin)
│   ├── middleware/              # Auth JWT & CORS middleware
│   ├── models/                  # GORM models (User, File)
│   ├── repository/              # Layer akses database MySQL
│   ├── services/                # Logika bisnis file & auth
│   └── utils/                   # Helper JWT, password, validasi MIME
├── storage-frontend/            # Aplikasi Next.js App Router
│   ├── app/
│   │   ├── dashboard/           # Dashboard manajemen file
│   │   ├── login/               # Halaman Login
│   │   ├── register/            # Halaman Registrasi
│   │   └── page.tsx             # Redirect sesi otomatis
│   ├── components/              # FileUpload, FileTable
│   └── lib/                     # fetchApi helper & client auth
├── docs/                        # Dokumentasi arsitektur, API, & PRD
├── uploads/                     # Direktori penyimpanan fisik file
├── .env.example                 # Template env backend
└── README.md
```

---

## Prasyarat Lingkungan

Pastikan alat-alat berikut telah terpasang di komputer Anda:
- [Golang](https://go.dev/dl/) (v1.20 atau lebih baru)
- [Node.js](https://nodejs.org/) (v18.17 atau lebih baru) & `npm`
- MySQL Server (bisa menggunakan [Laragon](https://laragon.org/) / XAMPP / Docker)
- [Git](https://git-scm.com/)

---

## Panduan Menjalankan Proyek

### 1. Kloning Repositori
```bash
git clone https://github.com/<username>/storage-manager-api.git
cd storage-manager-api
```

### 2. Konfigurasi Backend (Golang)

1. Salin template file `.env`:
   ```bash
   cp .env.example .env
   ```
   *(Di Windows PowerShell: `Copy-Item .env.example .env`)*

2. Sesuaikan konfigurasi database di file `.env` jika diperlukan:
   ```ini
   PORT=8080
   GIN_MODE=debug

   DB_USER=root
   DB_PASSWORD=
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=storage_db

   JWT_SECRET=super_secret_key_development_only
   JWT_EXPIRATION_HOURS=24

   STORAGE_PATH=./uploads
   MAX_UPLOAD_SIZE=52428800
   ```
   *Catatan: Backend akan otomatis memastikan database `storage_db` tersedia (`CREATE DATABASE IF NOT EXISTS`) serta memigrasikan tabel `users` dan `files` saat server dijalankan.*

3. Unduh modul dependensi Go:
   ```bash
   go mod download
   ```

4. Jalankan backend server:
   ```bash
   go run ./cmd/api/main.go
   ```
   Server backend aktif di: `http://localhost:8080`

---

### 3. Konfigurasi Frontend (Next.js)

1. Buka terminal baru dan masuk ke direktori frontend:
   ```bash
   cd storage-frontend
   ```

2. Salin template `.env.local`:
   ```bash
   cp .env.example .env.local
   ```
   *(Di Windows PowerShell: `Copy-Item .env.example .env.local`)*

   Pastikan isinya mengarah ke API backend:
   ```ini
   NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
   ```

3. Pasang dependensi frontend:
   ```bash
   npm install
   ```

4. Jalankan server pengembangan Next.js:
   ```bash
   npm run dev
   ```
   Frontend aktif di: `http://localhost:3000`

---

## Alur Penggunaan Aplikasi

1. Buka browser di `http://localhost:3000`.
2. Halaman akan otomatis mengarahkan ke `http://localhost:3000/login`.
3. Klik **"Daftar sekarang"** untuk membuat akun baru di `/register`.
4. Masuk (*Login*) menggunakan email dan password terdaftar.
5. Pada Dashboard (`/dashboard`), Anda dapat:
   - Mengunggah file (drag-and-drop atau browse).
   - Melihat daftar file beserta ukuran dan waktu upload.
   - Mengunduh file via tombol **Download**.
   - Menghapus file sementara (Trash) via tombol **Hapus**.
   - Keluar dari sesi via tombol **Logout**.

---

## Endpoint API Backend (`/api/v1`)

| Method | Endpoint | Proteksi Auth | Deskripsi |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/register` | Publik | Mendaftarkan akun user baru |
| `POST` | `/auth/login` | Publik | Login dan mendapatkan token JWT |
| `GET` | `/auth/me` | Bearer JWT | Melihat data akun user aktif |
| `POST` | `/files/upload` | Bearer JWT | Mengunggah file (`multipart/form-data`) |
| `GET` | `/files` | Bearer JWT | Mengambil daftar file aktif milik user |
| `GET` | `/files/download/:id` | Bearer JWT | Mengunduh file fisik (streaming `c.File()`) |
| `DELETE` | `/files/:id` | Bearer JWT | Soft delete file ke Trash Bin |

---

## Menjalankan Pengujian (Testing)

- **Backend Unit Test:**
  ```bash
  go test ./... -v
  ```
- **Frontend Type Check:**
  ```bash
  cd storage-frontend
  npx tsc --noEmit
  ```

---

## Dokumentasi Tambahan
Silakan merujuk ke folder `docs/` untuk panduan arsitektur lebih mendalam:
- `docs/api_reference.md`: Spesifikasi request/response endpoint.
- `docs/design.md`: Arsitektur sistem dan skema database MySQL.
- `docs/security.md`: Panduan implementasi keamanan (Path Traversal, Magic Bytes MIME, CORS).
- `docs/workflow.md`: Diagram alur otorisasi dan I/O file.

---
*Dibuat oleh [Fairus Syahdilla Rabbani]*


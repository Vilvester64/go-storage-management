# System Design & Architecture

## 1. Arsitektur Keseluruhan
Sistem ini menggunakan arsitektur **Client-Server** dengan pendekatan **Headless REST API**.
- **Frontend (Client):** Next.js (App Router), Tailwind CSS, Framer Motion. Bertanggung jawab atas antarmuka pengguna, pemisahan file besar (chunking - *future*), dan manajemen state.
- **Backend (Server):** Golang dengan framework Gin. Menangani logika bisnis, routing, middleware, dan interaksi dengan database serta sistem file lokal.
- **Database:** MySQL. Menyimpan entitas User dan metadata File.
- **Storage:** Local File System (Tahap Development). Menyimpan file fisik di dalam folder server (contoh: `./uploads`).

## 2. Struktur Folder Backend (Golang - Clean Architecture)
```text
storage-manager-api/
├── cmd/
│   └── api/
│       └── main.go           # Entry point aplikasi
├── internal/
│   ├── config/               # Konfigurasi aplikasi (Database, Env)
│   ├── handlers/             # Gin HTTP Handlers (Controller)
│   ├── middleware/           # Gin Middleware (Auth, CORS, Logger)
│   ├── models/               # GORM Models (User, File)
│   ├── repository/           # Interaksi langsung dengan Database (GORM)
│   ├── services/             # Logika Bisnis (Pemrosesan file, dll)
│   └── utils/                # Fungsi bantuan (Hash, JWT, UUID, dll)
├── pkg/                      # Pustaka eksternal yang di-wrap (opsional)
├── uploads/                  # Folder penyimpanan fisik file (di-ignore di git)
├── .env                      # Variabel lingkungan
├── go.mod
└── go.sum
```

## 3. Desain Skema Database (MySQL)

### Tabel `users`
| Kolom | Tipe Data | Keterangan |
| :--- | :--- | :--- |
| `id` | UUID (PK) | Identifier unik pengguna |
| `email` | VARCHAR(255) | Email pengguna (Unique) |
| `password` | VARCHAR(255) | Hash password (Bcrypt) |
| `created_at`| DATETIME | Waktu pembuatan akun |
| `updated_at`| DATETIME | Waktu pembaruan akun |

### Tabel `files`
| Kolom | Tipe Data | Keterangan |
| :--- | :--- | :--- |
| `id` | UUID (PK) | Identifier unik file |
| `user_id` | UUID (FK) | Relasi ke tabel `users` |
| `name` | VARCHAR(255) | Nama asli file |
| `system_name`| VARCHAR(255) | Nama file unik di storage (UUID + ext) |
| `size` | BIGINT | Ukuran file dalam bytes |
| `mime_type` | VARCHAR(100) | Tipe konten (misal: image/jpeg, application/pdf) |
| `path` | VARCHAR(255) | Path relatif file di storage |
| `created_at`| DATETIME | Waktu upload |
| `updated_at`| DATETIME | Waktu update metadata |
| `deleted_at`| DATETIME | Untuk Soft Delete (Trash Bin) |

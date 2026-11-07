# Product Requirements Document (PRD) - Storage Management System

## 1. Visi Produk
Sistem manajemen file berbasis web yang aman, efisien, dan skalabel. Dirancang untuk mengelola penyimpanan file pribadi atau skala kecil dengan antarmuka yang modern dan responsif.

## 2. Tujuan & Sasaran
- Menyediakan platform penyimpanan file yang aman dengan autentikasi JWT.
- Memungkinkan pengguna untuk mengunggah, mengunduh, mengganti nama, dan menghapus file (soft delete).
- Melacak metadata file (ukuran, tipe, waktu, pemilik) di dalam database relasional.
- Menyediakan REST API yang tangguh menggunakan Gin (Golang) untuk diintegrasikan dengan frontend Next.js.

## 3. Fitur Utama (MVP)
- **Autentikasi Pengguna:** Registrasi, Login (JWT), dan Logout.
- **Manajemen File:** Upload (tunggal), Download, Rename, dan Hapus (Soft Delete ke Trash).
- **Metadata Tracking:** Menyimpan informasi file di database MySQL.
- **Validasi File:** Memeriksa tipe file berdasarkan ekstensi.

## 4. Fitur Lanjutan (Fase 2)
- **Chunked Upload:** Untuk file berukuran besar.
- **Validasi MIME Type:** Memeriksa *magic bytes* file untuk keamanan tambahan.
- **Thumbnail Generation:** Membuat preview otomatis untuk gambar.
- **Storage Quota:** Membatasi kapasitas penyimpanan per pengguna (misal 1GB).
- **Shareable Links:** Membuat tautan unduh publik (dengan/tanpa password & kedaluwarsa).

## 5. Asumsi & Batasan
- Sistem ini pada awalnya akan menggunakan **Local File System** untuk penyimpanan fisik file.
- Database yang digunakan adalah **MySQL** (via Laragon untuk lokal).
- Frontend dibangun secara terpisah menggunakan **Next.js**.

## 6. Persyaratan Non-Fungsional
- **Kinerja:** API harus merespons dalam waktu < 200ms untuk operasi non-I/O berat.
- **Keamanan:** Menerapkan CORS yang ketat, perlindungan Path Traversal, dan validasi input.
- **Skalabilitas:** Arsitektur harus memudahkan migrasi ke Object Storage (S3/MinIO) di masa depan.

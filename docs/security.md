# Panduan Keamanan (Security Guidelines)

Keamanan adalah aspek kritis dalam aplikasi manajemen file. Berikut adalah panduan implementasi keamanan untuk proyek ini:

## 1. Autentikasi dan Otorisasi
- **Hashing Password:** Wajib menggunakan algoritma `bcrypt` untuk melakukan *hashing* password pengguna sebelum disimpan ke database. Jangan pernah menyimpan password dalam bentuk *plain text*.
- **JSON Web Token (JWT):** Gunakan token JWT berumur pendek (misal: 1 jam) untuk autentikasi API. Token harus dikirim via header `Authorization: Bearer <token>`.
- **Validasi Kepemilikan (Authorization):** Setiap operasi akses file (download, rename, delete) harus memvalidasi apakah `user_id` pada token JWT cocok dengan `user_id` pada entitas file di database.

## 2. Keamanan Upload File
- **Path Traversal Protection:** Jangan pernah menggunakan nama asli file dari pengguna untuk menyimpan file di server (misal: `../../etc/passwd`). Hasilkan nama file acak yang unik menggunakan **UUID** dan simpan di kolom `system_name`. Nama asli file hanya disimpan di kolom `name` sebagai metadata.
- **Validasi MIME Type:** Validasi ekstensi file saja tidak cukup. Gunakan pustaka seperti `github.com/h2non/filetype` untuk memeriksa *magic bytes* di awal file (header) guna memastikan tipe file yang diunggah sesuai. Ini mencegah serangan file berbahaya (seperti skrip PHP yang di-*rename* menjadi JPG).
- **Pembatasan Ukuran File (File Size Limit):** Implementasikan *middleware* di Gin untuk membatasi ukuran maksimal file yang dapat diunggah dalam satu *request* agar mencegah serangan *Denial of Service* (DoS) atau kehabisan ruang disk (*Out of Space*).

## 3. Keamanan Jaringan dan API
- **CORS (Cross-Origin Resource Sharing):** Konfigurasikan *middleware* `gin-contrib/cors` dengan ketat. Hanya izinkan domain asal *frontend* Anda (misal: `http://localhost:3000` untuk *development*) untuk mengakses API.
- **Rate Limiting:** Terapkan *rate limiting* pada *endpoint* publik, terutama `/login` dan `/register`, untuk mencegah serangan *brute force*.

## 4. Keamanan Database
- **Prepared Statements:** GORM secara default telah menggunakan *prepared statements*, sehingga aman dari serangan *SQL Injection*. Hindari menggunakan *raw query* SQL sebisa mungkin.
- **Soft Delete:** Saat pengguna menghapus file, pindahkan datanya ke tabel/status sampah (Trash Bin) dengan memperbarui kolom `deleted_at` (fitur bawaan GORM). Jangan langsung hapus baris di database dan file fisik di server sebagai langkah perlindungan jika pengguna tidak sengaja menghapus file.

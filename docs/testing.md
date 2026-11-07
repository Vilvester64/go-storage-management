# Panduan Pengujian (Testing Guide)

Dokumen ini menjelaskan strategi pengujian untuk memastikan sistem *storage management* bebas dari *bug* kritis, terutama pada bagian I/O file dan operasi database GORM.

### 1. Persiapan Lingkungan Pengujian
Pastikan layanan lokal berjalan sebelum melakukan tes:
- **Database:** Pastikan MySQL di Laragon aktif. Buat database terpisah bernama `storage_db_test` agar data pengujian tidak bercampur dengan data pengembangan.
- **Environment:** Duplikasi `.env` menjadi `.env.test` dan ubah konfigurasi database ke `storage_db_test`.

### 2. Unit Testing (Golang)
Unit test difokuskan pada logika bisnis, utilitas (seperti generator UUID, pemotong nama file), dan *middleware* JWT.
- Jalankan semua tes:
  ```bash
  go test ./... -v
  ```
- Cek *Test Coverage* (Target: >70%):
  ```bash
  go test ./... -coverprofile=coverage.out
  go tool cover -html=coverage.out
  ```

### 3. API Integration Testing (Postman)
Karena proyek ini memisahkan REST API (Gin) dan *client*, pengujian integrasi dilakukan dengan Postman.
- **Workspace:** Impor file `postman_collection.json` (disediakan di root folder) ke aplikasi Postman.
- **Variabel:** Set variabel `{{base_url}}` ke `http://localhost:8080` dan `{{jwt_token}}` ke token yang didapat setelah endpoint Login berhasil dieksekusi.
- **Skenario Wajib:**
  1. *Upload* file berukuran sangat kecil (teks) dan sangat besar (video) untuk menguji limit memori Gin.
  2. Akses endpoint `/files` tanpa token JWT (harus mengembalikan `401 Unauthorized`).
  3. Coba unggah file `.sh` yang di-*rename* menjadi `.txt` (harus ditolak oleh validasi MIME).

### 4. Frontend Testing (Next.js)
Untuk antarmuka Next.js, pengujian difokuskan pada UI/UX dan integrasi API:
- Gunakan alat *Developer Tools* (Network Tab) di browser untuk memastikan *payload* `multipart/form-data` terkirim dengan benar saat proses *upload*.

# Alur Kerja Sistem (System Workflow)

Dokumen ini menjelaskan alur interaksi antara Frontend (Client) dan Backend (API) dalam operasi-operasi utama.

## 1. Alur Registrasi dan Login (Authentication Workflow)
1. **Client** mengirimkan data `email` dan `password` (JSON) ke *endpoint* `/api/register`.
2. **API (Golang)** memvalidasi data dan melakukan *hashing* pada password menggunakan `bcrypt`.
3. **API** menyimpan data pengguna baru ke database MySQL dan merespons `201 Created`.
4. **Client** mengirimkan `email` dan `password` ke *endpoint* `/api/login`.
5. **API** memvalidasi kredensial. Jika valid, API menghasilkan ( *generate* ) **JWT Token**.
6. **API** mengembalikan token JWT ke Client.
7. **Client** menyimpan JWT (biasanya di *LocalStorage* atau *HTTP-Only Cookie*) untuk otorisasi *request* berikutnya.

## 2. Alur Upload File (File Upload Workflow)
1. **Client** mengirimkan HTTP POST *request* dengan format `multipart/form-data` ke *endpoint* `/api/files/upload`. Request menyertakan token JWT pada header `Authorization`.
2. **API Middleware** (Gin) memvalidasi JWT. Jika token valid dan tidak kedaluwarsa, permintaan diteruskan ke *Handler*.
3. **API Handler**:
    - Menerima aliran (*stream*) data file.
    - Menghasilkan nama unik (UUID) untuk sistem (misal: `123e4567-e89b-12d3.jpg`).
    - Memvalidasi MIME tipe dan batas ukuran.
    - Menyimpan file secara fisik ke direktori `/uploads`.
    - Menyimpan metadata (nama asli, UUID, ukuran, tipe, `user_id`) ke database (tabel `files`) via GORM.
4. **API** mengembalikan respons `200 OK` beserta data metadata file yang berhasil diunggah.

## 3. Alur Download File (File Download Workflow)
1. **Client** meminta file tertentu melalui HTTP GET ke *endpoint* `/api/files/download/{id_file}` dengan menyertakan token JWT.
2. **API Middleware** memvalidasi JWT.
3. **API Handler**:
    - Mencari file berdasarkan ID di database.
    - Memverifikasi apakah file tersebut milik pengguna yang terautentikasi (mencocokkan `user_id`).
    - Jika valid, mengambil path file lokal (misal: `/uploads/123e4567...jpg`).
    - Mengirimkan (*stream*) konten file ke Client sebagai respons (menggunakan `c.File()` di Gin).
4. **Client** menerima dan menyimpan file ke perangkat pengguna.

## 4. Alur Hapus File (Soft Delete Workflow)
1. **Client** mengirimkan HTTP DELETE ke `/api/files/{id_file}` dengan token JWT.
2. **API Middleware** memvalidasi JWT.
3. **API Handler**:
    - Memverifikasi kepemilikan file di database.
    - Memperbarui kolom `deleted_at` pada database via GORM (Soft Delete).
    - File fisik **tidak** dihapus dari direktori lokal.
4. **API** merespons `200 OK` (File dipindahkan ke Trash).

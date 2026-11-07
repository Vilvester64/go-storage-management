# Referensi API (API Reference)

Base URL: `http://localhost:8080/api/v1`

Semua *endpoint* (kecuali Auth) wajib menyertakan header:
`Authorization: Bearer <token_jwt>`

---

## 1. Authentication

### Register
Mendaftarkan pengguna baru.
- **Endpoint:** `POST /auth/register`
- **Body (JSON):**
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword123"
  }
  ```
- **Response (201 Created):**
  ```json
  {
    "message": "User registered successfully",
    "data": { "id": "uuid", "email": "user@example.com" }
  }
  ```

### Login
Mendapatkan token akses (JWT).
- **Endpoint:** `POST /auth/login`
- **Body (JSON):**
  ```json
  {
    "email": "user@example.com",
    "password": "securepassword123"
  }
  ```
- **Response (200 OK):**
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
  ```

---

## 2. File Management

### Upload File
Mengunggah file tunggal ke server lokal.
- **Endpoint:** `POST /files/upload`
- **Headers:** `Content-Type: multipart/form-data`
- **Form Data:**
  - `file`: (File fisik yang diunggah)
- **Response (201 Created):**
  ```json
  {
    "message": "File uploaded successfully",
    "data": {
      "id": "uuid",
      "name": "dokumen.pdf",
      "size": 1048576,
      "mime_type": "application/pdf"
    }
  }
  ```

### List Files
Mendapatkan daftar file milik pengguna yang terautentikasi (tidak termasuk yang dihapus).
- **Endpoint:** `GET /files`
- **Response (200 OK):**
  ```json
  {
    "data": [
      {
        "id": "uuid",
        "name": "dokumen.pdf",
        "size": 1048576,
        "created_at": "2023-10-27T10:00:00Z"
      }
    ]
  }
  ```

### Download File
Mengunduh file fisik berdasarkan ID.
- **Endpoint:** `GET /files/download/{id}`
- **Response (200 OK / Binary Stream):**
  - Headers: `Content-Disposition: attachment; filename="dokumen.pdf"`
  - Body: [Isi file]

### Delete File (Soft Delete)
Memindahkan file ke Trash (memperbarui kolom `deleted_at`).
- **Endpoint:** `DELETE /files/{id}`
- **Response (200 OK):**
  ```json
  {
    "message": "File moved to trash successfully"
  }
  ```

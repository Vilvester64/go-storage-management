# Peta Jalan Pengembangan (Roadmap) - Storage Management System

Dokumen ini memetakan visi produk, fase rilis, dan prioritas fitur. Proyek ini dibangun secara iteratif untuk memastikan stabilitas sebelum menambahkan kompleksitas tingkat lanjut.

### Fase 1: Minimum Viable Product (MVP) - Selesai ✅
Fokus pada fungsionalitas inti dari manajemen file dan stabilitas REST API menggunakan Gin dan GORM.
- [x] Inisialisasi proyek Golang & struktur folder (Clean Architecture).
- [x] Konfigurasi koneksi MySQL (Local via Laragon).
- [x] Autentikasi dasar (Register, Login dengan JWT).
- [x] CRUD File API (Upload, Download, Rename, Soft Delete).
- [x] Parsing dan penyimpanan metadata file ke database.
- [x] Setup proyek Next.js (App Router) & Tailwind CSS.
- [x] UI Dashboard dasar untuk navigasi file.

### Fase 2: Fitur Portofolio Lanjutan - Sedang Berjalan 🚧
Fokus pada optimasi performa, keamanan, dan fitur *enterprise*.
- [ ] **Chunked Upload:** Memecah file besar (>50MB) di frontend dan menggabungkannya di backend Gin.
- [ ] **Trash Bin:** Sistem pemulihan file menggunakan fitur `DeletedAt` GORM sebelum penghapusan permanen.
- [ ] **Storage Quota:** Validasi limit kapasitas per pengguna (misal: 1GB per user).
- [ ] **MIME Type Validation:** Keamanan ekstra menggunakan `h2non/filetype` untuk memeriksa *magic bytes*.
- [ ] **Thumbnail Generation:** Integrasi `disintegration/imaging` untuk *preview* otomatis gambar di UI.

### Fase 3: Skalabilitas & Infrastruktur Cloud - Rencana Masa Depan 📅
Mempersiapkan aplikasi untuk *deployment production* dan skala yang lebih besar.
- [ ] **Shareable Links:** Tautan unduh publik dengan opsi *password* dan waktu kedaluwarsa.
- [ ] **Object Storage Integration:** Migrasi penyimpanan dari Local System ke AWS S3 / MinIO.
- [ ] **Dockerization:** Containerization untuk backend, frontend, dan database.
- [ ] **CI/CD Pipeline:** Otomatisasi testing dan deployment menggunakan GitHub Actions.

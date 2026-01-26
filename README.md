# SERUMPUN – Platform Kolaborasi & Monitoring SE2026 Kepulauan Riau

SERUMPUN adalah platform kolaborasi dan monitoring kegiatan  
**Sensus Ekonomi 2026 (SE2026)** di Provinsi Kepulauan Riau.

Platform ini dirancang sebagai **gerbang utama (landing page)** yang:
- memberikan **informasi ringkas dan konteks** terkait progres SE2026
- menyediakan **akses cepat** ke seluruh layanan SERUMPUN
- mengarahkan pengguna ke **dashboard monitoring** untuk analisis yang lebih detail

---

## 🎯 Tujuan Pengembangan

- Menyediakan **platform terpusat** untuk seluruh ekosistem SERUMPUN
- Memberikan **gambaran awal (overview)** kondisi pelaksanaan SE2026
- Memudahkan pimpinan dan koordinator untuk:
  - memahami progres secara cepat dari landing page
  - melakukan monitoring & evaluasi (monev) mendalam melalui dashboard
- Menyediakan backend data yang **aman (read-only)** dan mudah diintegrasikan

---

## 🧩 Konsep Halaman & Alur Pengguna

### 1️⃣ Landing Page SERUMPUN (Overview)
Landing page **bukan sekadar link tree**, tetapi berfungsi sebagai:

- **Pengantar kondisi pelaksanaan SE2026**
- Ringkasan indikator utama dari dashboard, seperti:
  - total penugasan
  - persentase penyelesaian
  - ringkasan progres per kabupaten/kota atau bidang
- Akses cepat ke layanan utama SERUMPUN

Landing page ditujukan untuk:
- pimpinan yang ingin **melihat kondisi sekilas**
- pengguna baru yang ingin **memahami fungsi SERUMPUN**
- pintu masuk sebelum menuju dashboard lengkap

### 2️⃣ Halaman Dashboard Monitoring (Detail)
Halaman dashboard digunakan ketika pengguna membutuhkan:
- analisis lebih rinci
- filter data (kab/kota, bidang, status)
- detail penugasan dan komentar

Dashboard menyajikan visualisasi lengkap berbasis data real-time (read-only).

---

## 🧭 Alur Penggunaan (User Flow)
Landing Page SERUMPUN
│
├─ Informasi singkat & ringkasan progres
│
├─ Akses Portal SERUMPUN
│
└─ "Lihat Dashboard Lengkap"
│
v
Dashboard Monitoring & Evaluasi


---

## 🧩 Cakupan Fitur

### A. Landing Page
- Informasi singkat platform SERUMPUN
- Ringkasan indikator utama (KPI overview)
- Akses ke:
  - Portal SERUMPUN (All)
  - Portal SERUMPUN (Member)
  - Pendaftaran Pengguna
  - Petunjuk Penggunaan
- Tautan ke Dashboard Monitoring Lengkap

### B. Dashboard Monitoring & Evaluasi
- KPI penugasan (total, selesai, progres)
- Progres penugasan per kabupaten/kota
- Progres penugasan per bidang
- Heatmap kabupaten/kota × bidang
- Tabel detail penugasan
- Komentar / catatan progres terbaru

---

## 🏗️ Arsitektur Sistem

PostgreSQL (SERUMPUN)
|
| SELECT-only
v
Backend Data API (Go)
|
| Live CSV (HTTP)
v
Flourish (Visualisasi)
|
v
Landing Page (Overview) ──> Dashboard Page (Detail)


**Catatan penting:**
- Backend bersifat **read-only**
- Tidak ada operasi write/update/delete
- Aman untuk penggunaan operasional

---

## 📁 Struktur Repository

├── client/ # Frontend (Landing Page & Dashboard Page)
│ └── README.md
│
├── server/ # Backend Data API (Go, Live CSV)
│ └── README.md
│
├── README.md # Dokumentasi utama (file ini)


---

## 🔧 Teknologi yang Digunakan

### Backend
- Go
- PostgreSQL (read-only)
- pgx (database driver)
- Chi Router
- CSV streaming
- In-memory cache

### Frontend
- Next.js (tsx)
- JokoUI components
- Embed visualisasi Flourish

### Visualisasi
- Flourish (Live CSV URL)

---

## 📊 Sumber Data Dashboard

Dashboard mengambil data dari tabel utama SERUMPUN:
- `issues`
- `states`
- `users`
- `issue_assignees`
- `issue_labels`
- `labels`
- `issue_comments`
- `projects`
- `workspaces`

Filter utama:
- **Workspace:** `Platform Serumpun`
- **Project:** `Sensus Ekonomi 2026`

---

## 🔐 Keamanan & Akses Data

- Koneksi database menggunakan **akun read-only**
- Query menggunakan parameter binding
- Cache untuk mengurangi beban database
- Backend hanya menyediakan endpoint CSV

---

## 🚀 Pengembangan Selanjutnya (Opsional)

- Autentikasi akses dashboard
- Role-based view (pimpinan / admin / bidang)
- Penambahan indikator SLA & aging penugasan
- Visualisasi tren waktu
- Integrasi notifikasi internal

---

## 👥 Pengelola & Kontributor

Dikembangkan untuk mendukung kegiatan  
**Sensus Ekonomi 2026**

**BPS Provinsi Kepulauan Riau**

---

## 📄 Lisensi & Penggunaan

Repositori ini digunakan untuk kebutuhan internal SERUMPUN.  
Penggunaan di luar konteks ini memerlukan izin dari pengelola.

© 2025 – BPS Provinsi Kepulauan Riau

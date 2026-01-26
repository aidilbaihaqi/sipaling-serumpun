# SERUMPUN – Platform Kolaborasi & Monitoring SE2026 Kepulauan Riau

SERUMPUN adalah platform kolaborasi dan monitoring kegiatan **Sensus Ekonomi 2026 (SE2026)** di Provinsi Kepulauan Riau.  
Platform ini mengintegrasikan **portal kerja**, **manajemen penugasan**, dan **dashboard monitoring & evaluasi (monev)** dalam satu ekosistem yang mudah diakses.

Repositori ini berisi **arsitektur lengkap** SERUMPUN yang terdiri dari **client (frontend)** dan **server (backend data API)**.

---

## 🎯 Tujuan Pengembangan

- Menyediakan **akses terpusat** untuk seluruh layanan SERUMPUN
- Menyajikan **dashboard monev real-time** berbasis data penugasan
- Memudahkan pimpinan dan koordinator bidang dalam:
  - memantau progres per **kabupaten/kota**
  - memantau progres per **bidang**
  - melihat **status pengerjaan** dan **catatan/komentar terbaru**
- Menyediakan backend ringan yang **aman (read-only)** dan **mudah diintegrasikan** dengan platform visualisasi (Flourish)

---

## 🧩 Cakupan Fitur

### 1. Landing Page SERUMPUN
- Informasi singkat tentang platform
- Akses cepat ke:
  - Portal SERUMPUN (All)
  - Portal SERUMPUN (Member)
  - Pendaftaran Pengguna
  - Petunjuk Penggunaan
- Akses ke Dashboard Monitoring

### 2. Dashboard Monitoring & Evaluasi (Monev)
- Ringkasan KPI penugasan
- Progres penugasan per kabupaten/kota
- Progres penugasan per bidang
- Heatmap kabupaten/kota × bidang
- Tabel detail penugasan beserta komentar terbaru

Dashboard menggunakan **Flourish** sebagai alat visualisasi dan mengambil data melalui **Live CSV URL** dari backend.

### 3. Backend Data API (Read-Only)
- Mengambil data dari **PostgreSQL SERUMPUN**
- Hanya menggunakan operasi **SELECT**
- Menyediakan endpoint CSV:
  - `/csv/kpi.csv`
  - `/csv/progress_kabkot.csv`
  - `/csv/progress_bidang.csv`
  - `/csv/heatmap.csv`
  - `/csv/issues_detail.csv`
- Dilengkapi cache untuk efisiensi query

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
Flourish Dashboard
|
v
Landing Page SERUMPUN


**Catatan penting:**
- Backend **tidak mengubah data** apa pun
- Semua visualisasi bersifat **read-only**
- Aman untuk lingkungan produksi

---

## 📁 Struktur Repository

.
├── client/ # Frontend (Landing Page & Dashboard Embed)
│ └── README.md
│
├── server/ # Backend Data API (Go)
│ └── README.md
│
├── README.md # Dokumentasi utama (file ini)


---

## 🔧 Teknologi yang Digunakan

### Backend
- Go
- pgx (PostgreSQL driver)
- Chi Router
- CSV streaming
- In-memory cache

### Frontend
- Static HTML / CSS (atau framework sesuai kebutuhan)
- Embed Flourish Dashboard

### Data Visualization
- Flourish (Live CSV)

---

## 🔐 Keamanan & Akses Data

- Backend menggunakan **kredensial database read-only**
- Tidak ada endpoint write/update/delete
- Query menggunakan parameter binding (anti SQL injection)
- Cache untuk mengurangi beban database

---

## 📊 Sumber Data Utama

Dashboard mengambil data dari tabel-tabel utama SERUMPUN:
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

## 🚀 Alur Penggunaan (High Level)

1. User membuka **Landing Page SERUMPUN**
2. User memilih:
   - Portal kerja, atau
   - Dashboard Monitoring
3. Dashboard menampilkan visualisasi dari Flourish
4. Flourish mengambil data dari backend melalui Live CSV
5. Backend mengambil data dari PostgreSQL (read-only)

---

## 🧭 Pengembangan Selanjutnya (Opsional)

- Autentikasi akses dashboard
- Role-based dashboard (pimpinan / admin / bidang)
- Penambahan indikator SLA & aging task
- Integrasi notifikasi (email / WA internal)
- Visualisasi tren waktu (harian/mingguan)

---

## 👥 Pengelola & Kontributor

Dikembangkan untuk mendukung kegiatan **Sensus Ekonomi 2026**  
**BPS Provinsi Kepulauan Riau**

---

## 📄 Lisensi & Penggunaan

Repositori ini digunakan untuk kebutuhan internal pengembangan dan operasional SERUMPUN.  
Penggunaan di luar konteks ini memerlukan izin dari pengelola.

© 2025 – BPS Provinsi Kepulauan Riau

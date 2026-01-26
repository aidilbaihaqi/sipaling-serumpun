# SERUMPUN – Backend Data API (Server)

Backend SERUMPUN berfungsi sebagai **Data API read-only** yang menyediakan data
monitoring & evaluasi (monev) SE2026 dalam format **CSV by URL**.

API ini dirancang khusus untuk:
- konsumsi **Flourish (Live CSV)**
- ringkasan data di Landing Page
- kebutuhan monitoring tanpa mengubah data sumber

---

## 🎯 Tujuan Backend

- Mengambil data dari **PostgreSQL SERUMPUN**
- Menyediakan endpoint CSV yang **ringan, aman, dan cepat**
- Menghindari query langsung dari frontend ke database
- Mendukung visualisasi real-time berbasis Flourish

---

## 🏗️ Arsitektur

PostgreSQL (SERUMPUN)
|
| SELECT-only
v
Go Data API
|
| CSV over HTTP
v
Flourish / Frontend


---

## 🔐 Prinsip Keamanan

- Menggunakan **akun database read-only**
- Tidak ada endpoint write/update/delete
- Query menggunakan parameter binding
- Cache in-memory untuk mengurangi beban DB

---

## 📁 Struktur Folder

server/
├── cmd/
│ └── api/
│ └── main.go
├── internal/
│ ├── cache/ # Cache CSV (TTL)
│ ├── db/ # Koneksi PostgreSQL
│ ├── http/ # Router & handler HTTP
│ └── queries/ # Loader SQL
├── queries/ # File SQL (*.sql)
│ ├── kpi.sql
│ ├── progress_kabkot.sql
│ ├── progress_bidang.sql
│ ├── heatmap_kabkot_bidang.sql
│ └── issues_detail.sql
├── .env.example
└── README.md


---

## 🔧 Teknologi

- Go
- PostgreSQL
- pgx (database driver)
- Chi Router
- CSV streaming
- In-memory cache

---

## ⚙️ Konfigurasi Environment

Salin file `.env.example` menjadi `.env`:

```env
APP_PORT=8080
DATABASE_URL=postgres://USER:PASSWORD@HOST:5432/DBNAME?sslmode=require

WORKSPACE_NAME=Platform Serumpun
PROJECT_NAME=Sensus Ekonomi 2026

# key kab/kota di users.metadata
KABKOTA_KEY=kab_kota

CACHE_TTL_SECONDS=60
```

## Jalanin server
cd server
go mod tidy
go run ./cmd/api

## Cek Health
GET http://localhost:8080/healthz

📡 Endpoint CSV
Endpoint	Deskripsi
/csv/kpi.csv	Ringkasan KPI
/csv/progress_kabkot.csv	Progres per kab/kota
/csv/progress_bidang.csv	Progres per bidang
/csv/heatmap.csv	Heatmap kab/kota × bidang
/csv/issues_detail.csv	Detail penugasan + komentar

Semua endpoint:
- Format: text/csv
- Method: GET
- Read-only

📊 Sumber Data

Data diambil dari tabel utama:
- issues
- states
- users
- issue_assignees
- issue_labels
- labels
- issue_comments
- projects
- workspaces

Filter utama:

Workspace: Platform Serumpun

Project: Sensus Ekonomi 2026

🧭 Catatan Pengembangan

Backend ini tidak menyimpan state

Semua agregasi dilakukan via SQL

Cocok untuk:
- dashboard monev
- landing page overview
- integrasi BI ringan


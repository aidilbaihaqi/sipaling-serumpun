# SERUMPUN API Endpoints

Base URL: `http://localhost:8080`

## Health Check

```
GET /healthz
```

Response: `ok`

---

## CSV Data Endpoints

All endpoints return CSV format with `text/csv` content type.

### 1. KPI Provinsi
```
GET /api/v1/kpi_provinsi.csv
```

**Output Columns:**
- `nama`, `email`, `bidang`, `instansi`, `jabatan`, `backlog`, `todo`, `in_progress`, `done`, `percent`

**Description:** KPI per pegawai provinsi (ALL staff: Pengarah, Ketua Pelaksana, Ketua Sekretariat, Wakil Ketua Sekretariat, Anggota Sekretariat, Ketua Bidang, Anggota Bidang) dengan breakdown status berdasarkan `states.group`

**Total Rows:** 60 pegawai provinsi

**Data Source:**
- Nama pegawai diambil dari `data/daftar_pengguna_serumpun.csv` (bukan dari `users.display_name`)
- Mapping email → nama dilakukan secara dinamis saat runtime
- Setiap pegawai muncul 1× untuk bidang mereka (NO CROSS JOIN)

**Status Mapping:**
- `backlog` → Dicatat
- `todo` → Ditugaskan (unstarted)
- `in_progress` → Sedang Dikerjakan (started + triage)
- `done` → Selesai (completed)
- `cancelled` → **DIABAIKAN** (tidak dihitung dalam total dan persentase)

**Implementation:** Dynamic SQL generation in Go (no static SQL file)

---

### 2. KPI Kabupaten/Kota
```
GET /api/v1/kpi_kabkot.csv
```

**Output Columns:**
- `nama`, `email`, `bidang`, `instansi`, `jabatan`, `backlog`, `todo`, `in_progress`, `done`, `percent`

**Description:** KPI per Ketua di Kabupaten/Kota (Kepala Kab/Kot, Ketua Pelaksana SE, Ketua Sekretariat, Ketua Bidang) dengan breakdown status berdasarkan `states.group`

**Total Rows:** 44 Ketua
- 2 Kepala Kab/Kot
- 7 Ketua Pelaksana SE
- 7 Ketua Sekretariat
- 28 Ketua Bidang (7 kab/kota × 4 bidang)

**Data Source:**
- Nama pegawai diambil dari `data/daftar_pengguna_serumpun.csv` (bukan dari `users.display_name`)
- Mapping email → nama dilakukan secara dinamis saat runtime
- Setiap ketua muncul 1× untuk bidang mereka (NO CROSS JOIN)

**Status Mapping:** (sama dengan KPI Provinsi)

**Implementation:** Dynamic SQL generation in Go (no static SQL file)

---

### 3. Heatmap Kabupaten/Kota × Bidang
```
GET /api/v1/heatmap.csv
```

**Output Columns:**
- `kab_kota`, `bidang`, `total`, `selesai`, `persen_selesai`

**Description:** Matrix progress per kombinasi kab/kota dan bidang. Status menggunakan `states.group = 'completed'` untuk menghitung selesai.

**Note:** Cancelled issues are excluded from all calculations.

**Implementation:** Static SQL file (`queries/heatmap_kabkot_bidang.sql`)

---

### 4. Detail Issues
```
GET /api/v1/issues_detail.csv
```

**Output Columns:**
- `issue_title`, `kab_kota`, `bidang`, `status`, `ketua_bidang`, `email_ketua_bidang`, `start_date`, `target_date`, `last_comment`, `comment_time`

**Description:** Detail lengkap setiap issue dengan komentar terbaru. Status menggunakan `states.group` (backlog, unstarted, started, triage, completed).

**Data Source:**
- Nama ketua_bidang diambil dari `data/daftar_pengguna_serumpun.csv` (bukan dari `users.display_name`)
- Mapping email → nama dilakukan secara dinamis saat runtime

**Removed Columns:** (tidak ditampilkan lagi)
- `issue_id` - UUID internal tidak diperlukan untuk visualisasi
- `comment_by` - Fokus pada konten komentar, bukan pembuat
- `created_at` - Tidak relevan untuk monitoring
- `updated_at` - Tidak relevan untuk monitoring

**Note:** Cancelled issues are excluded from output.

**Implementation:** Dynamic SQL generation in Go (no static SQL file)

---

## Example Usage

### cURL
```bash
# Get KPI Provinsi
curl http://localhost:8080/api/v1/kpi_provinsi.csv

# Get KPI Kabkot
curl http://localhost:8080/api/v1/kpi_kabkot.csv

# Get Heatmap
curl http://localhost:8080/api/v1/heatmap.csv

# Get Issues Detail
curl http://localhost:8080/api/v1/issues_detail.csv
```

### Flourish Live CSV
```
https://your-domain.com/api/v1/kpi_provinsi.csv
https://your-domain.com/api/v1/kpi_kabkot.csv
https://your-domain.com/api/v1/heatmap.csv
https://your-domain.com/api/v1/issues_detail.csv
```

### JavaScript Fetch
```javascript
fetch('http://localhost:8080/api/v1/kpi_provinsi.csv')
  .then(res => res.text())
  .then(csv => console.log(csv));
```

---

## Cache

All endpoints are cached for **60 seconds** (configurable via `CACHE_TTL_SECONDS` env var).

Cache headers:
```
Cache-Control: public, max-age=30
```

---

## CORS

Currently no CORS headers. Add middleware if needed for frontend integration.

---

## Rate Limiting

No rate limiting implemented. Consider adding for production.

---

## Authentication

No authentication required (read-only API). Consider adding for production.

---

## Architecture Notes

### Dynamic SQL Generation
Endpoints `kpi_provinsi`, `kpi_kabkot`, dan `issues_detail` menggunakan dynamic SQL generation di Go untuk:
- Mapping email → nama dari CSV secara real-time
- Menghindari hardcoded SQL yang sulit dimaintain
- Memastikan konsistensi data dengan source of truth (CSV)

### Data Source Priority
1. **Nama pegawai:** `data/daftar_pengguna_serumpun.csv` (primary)
2. **Fallback:** `users.display_name` atau `users.first_name + users.last_name` (jika email tidak ditemukan di CSV)

### Deleted Endpoints
Endpoint berikut telah dihapus karena redundan:
- ❌ `/api/v1/kpi.csv` - Digantikan oleh kpi_provinsi dan kpi_kabkot
- ❌ `/api/v1/progress_kabkot.csv` - Data sudah ada di kpi_kabkot
- ❌ `/api/v1/progress_bidang.csv` - Data sudah ada di kpi_provinsi

### SQL Files Status
- ✅ `queries/heatmap_kabkot_bidang.sql` - Masih digunakan (static SQL)
- ❌ `queries/kpi_provinsi.sql` - Dihapus (dynamic SQL)
- ❌ `queries/kpi_kabkot.sql` - Dihapus (dynamic SQL)
- ❌ `queries/issues_detail.sql` - Dihapus (dynamic SQL)
- ❌ `queries/kpi.sql` - Dihapus (endpoint dihapus)
- ❌ `queries/progress_kabkot.sql` - Dihapus (endpoint dihapus)
- ❌ `queries/progress_bidang.sql` - Dihapus (endpoint dihapus)

---

**Last Updated:** 2026-01-29

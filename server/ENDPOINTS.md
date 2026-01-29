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

### 1. KPI Agregat
```
GET /api/v1/kpi.csv
```

**Output Columns:**
- `total`, `selesai`, `dikerjakan`, `dibatalkan`, `belum_ditugaskan`, `kode_tidak_terbaca`, `persen_selesai`

**Description:** Ringkasan KPI keseluruhan untuk SE2026

---

### 2. Progress per Kabupaten/Kota
```
GET /api/v1/progress_kabkot.csv
```

**Output Columns:**
- `kab_kota`, `total`, `selesai`, `dikerjakan`, `dibatalkan`, `persen_selesai`

**Description:** Progress penugasan per 7 kabupaten/kota di Kepri

---

### 3. Progress per Bidang
```
GET /api/v1/progress_bidang.csv
```

**Output Columns:**
- `bidang`, `total`, `selesai`, `dikerjakan`, `dibatalkan`, `persen_selesai`

**Description:** Progress penugasan per bidang (PTI, Analisis, Teknis, Administrasi)

---

### 4. Heatmap Kabupaten/Kota × Bidang
```
GET /api/v1/heatmap.csv
```

**Output Columns:**
- `kab_kota`, `bidang`, `total`, `selesai`, `persen_selesai`

**Description:** Matrix progress per kombinasi kab/kota dan bidang

---

### 5. Detail Issues
```
GET /api/v1/issues_detail.csv
```

**Output Columns:**
- `issue_id`, `issue_title`, `kab_kota`, `bidang`, `status`, `ketua_bidang`, `email_ketua_bidang`, `start_date`, `target_date`, `last_comment`, `comment_time`, `comment_by`, `created_at`, `updated_at`

**Description:** Detail lengkap setiap issue dengan komentar terbaru

---

### 6. KPI Provinsi (NEW)
```
GET /api/v1/kpi_provinsi.csv
```

**Output Columns:**
- `nama`, `email`, `bidang`, `instansi`, `jabatan`, `backlog`, `todo`, `in_progress`, `done`, `percent`

**Description:** KPI per pegawai provinsi (Ketua + Anggota Bidang) dengan breakdown status berdasarkan `states.group`

**Status Mapping:**
- `backlog` → Dicatat
- `todo` → Ditugaskan (unstarted)
- `in_progress` → Sedang Dikerjakan (started + triage)
- `done` → Selesai (completed)
- `cancelled` → **DIABAIKAN** (tidak dihitung dalam total dan persentase)

---

### 7. KPI Kabupaten/Kota (NEW)
```
GET /api/v1/kpi_kabkot.csv
```

**Output Columns:**
- `nama`, `email`, `bidang`, `instansi`, `jabatan`, `backlog`, `todo`, `in_progress`, `done`, `percent`

**Description:** KPI per Kepala Kab/Kot dengan cross join semua bidang. Setiap Kepala muncul N× (N = jumlah bidang)

**Status Mapping:** (sama dengan KPI Provinsi)

---

## Example Usage

### cURL
```bash
# Get KPI Provinsi
curl http://localhost:8080/api/v1/kpi_provinsi.csv

# Get KPI Kabkot
curl http://localhost:8080/api/v1/kpi_kabkot.csv

# Get Progress Kabkot
curl http://localhost:8080/api/v1/progress_kabkot.csv
```

### Flourish Live CSV
```
https://your-domain.com/api/v1/kpi_provinsi.csv
https://your-domain.com/api/v1/kpi_kabkot.csv
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

**Last Updated:** 2026-01-29

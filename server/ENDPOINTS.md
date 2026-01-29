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

**Total Rows:** ~60 pegawai provinsi

**Status Mapping:**
- `backlog` → Dicatat
- `todo` → Ditugaskan (unstarted)
- `in_progress` → Sedang Dikerjakan (started + triage)
- `done` → Selesai (completed)
- `cancelled` → **DIABAIKAN** (tidak dihitung dalam total dan persentase)

---

### 2. KPI Kabupaten/Kota
```
GET /api/v1/kpi_kabkot.csv
```

**Output Columns:**
- `nama`, `email`, `bidang`, `instansi`, `jabatan`, `backlog`, `todo`, `in_progress`, `done`, `percent`

**Description:** KPI per Ketua di Kabupaten/Kota (Kepala Kab/Kot, Ketua Pelaksana SE, Ketua Sekretariat, Ketua Bidang) dengan breakdown status berdasarkan `states.group`

**Total Rows:** 44 Ketua (2 Kepala + 7 Ketua Pelaksana + 7 Ketua Sekretariat + 28 Ketua Bidang)

**Status Mapping:** (sama dengan KPI Provinsi)

---

### 3. Heatmap Kabupaten/Kota × Bidang
```
GET /api/v1/heatmap.csv
```

**Output Columns:**
- `kab_kota`, `bidang`, `total`, `selesai`, `persen_selesai`

**Description:** Matrix progress per kombinasi kab/kota dan bidang. Status menggunakan `states.group = 'completed'` untuk menghitung selesai.

**Note:** Cancelled issues are excluded from all calculations.

---

### 4. Detail Issues
```
GET /api/v1/issues_detail.csv
```

**Output Columns:**
- `issue_id`, `issue_title`, `kab_kota`, `bidang`, `status`, `ketua_bidang`, `email_ketua_bidang`, `start_date`, `target_date`, `last_comment`, `comment_time`, `comment_by`, `created_at`, `updated_at`

**Description:** Detail lengkap setiap issue dengan komentar terbaru. Status menggunakan `states.group` (backlog, unstarted, started, triage, completed).

**Note:** Cancelled issues are excluded from output.

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

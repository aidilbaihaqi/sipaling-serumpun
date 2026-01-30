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

**Query Parameters (Optional):**
- `bidang` - Filter by bidang (e.g., `Sosial`, `Produksi`, `Distribusi`, `Nerwilis`)
- `jabatan` - Filter by jabatan (e.g., `Ketua`, `Anggota`, `Pengarah`, `Ketua Pelaksana`)

**Examples:**
```
GET /api/v1/kpi_provinsi.csv?bidang=Sosial
GET /api/v1/kpi_provinsi.csv?jabatan=Ketua
GET /api/v1/kpi_provinsi.csv?bidang=Produksi&jabatan=Anggota
```

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

**Query Parameters (Optional):**
- `bidang` - Filter by bidang (e.g., `Sosial`, `Produksi`, `Distribusi`, `Nerwilis`, `Umum`, `Sekretariat`)
- `instansi` - Filter by instansi (e.g., `BPS Kota Batam`, `BPS Kabupaten Bintan`)
- `jabatan` - Filter by jabatan (e.g., `Ketua Bidang`, `Kepala Kab/Kot`, `Ketua Pelaksana`, `Ketua Sekretariat`)

**Examples:**
```
GET /api/v1/kpi_kabkot.csv?instansi=BPS Kota Batam
GET /api/v1/kpi_kabkot.csv?bidang=Sosial
GET /api/v1/kpi_kabkot.csv?instansi=BPS Kota Batam&bidang=Produksi
GET /api/v1/kpi_kabkot.csv?jabatan=Ketua Bidang
```

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

**Query Parameters (Optional):**
- `kab_kota` - Filter by kab/kota (e.g., `Batam`, `Karimun`, `Bintan`, `Natuna`, `Lingga`, `Kep. Anambas`, `Tanjung Pinang`)
- `bidang` - Filter by bidang (e.g., `Sosial`, `Produksi`, `Distribusi`, `Nerwilis`)

**Examples:**
```
GET /api/v1/heatmap.csv?kab_kota=Batam
GET /api/v1/heatmap.csv?bidang=Sosial
GET /api/v1/heatmap.csv?kab_kota=Batam&bidang=Produksi
```

**Note:** Cancelled issues are excluded from all calculations.

**Implementation:** Dynamic SQL generation in Go (converted from static SQL file)

---

### 4. Detail Issues
```
GET /api/v1/issues_detail.csv
```

**Output Columns:**
- `issue_title`, `kab_kota`, `bidang`, `status`, `ketua_bidang`, `email_ketua_bidang`, `start_date`, `target_date`, `last_comment`, `comment_time`

**Description:** Detail lengkap setiap issue dengan komentar terbaru. Status menggunakan `states.group` (backlog, unstarted, started, triage, completed).

**Query Parameters (Optional):**
- `kab_kota` - Filter by kab/kota (e.g., `Batam`, `Karimun`, `Bintan`, `Natuna`, `Lingga`, `Kep. Anambas`, `Tanjung Pinang`)
- `bidang` - Filter by bidang (e.g., `Sosial`, `Produksi`, `Distribusi`, `Nerwilis`)
- `status` - Filter by status (e.g., `backlog`, `unstarted`, `started`, `triage`, `completed`)

**Examples:**
```
GET /api/v1/issues_detail.csv?status=completed
GET /api/v1/issues_detail.csv?kab_kota=Batam
GET /api/v1/issues_detail.csv?bidang=Sosial
GET /api/v1/issues_detail.csv?kab_kota=Batam&bidang=Produksi&status=started
```

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

# Get KPI Provinsi filtered by bidang
curl "http://localhost:8080/api/v1/kpi_provinsi.csv?bidang=Sosial"

# Get KPI Kabkot
curl http://localhost:8080/api/v1/kpi_kabkot.csv

# Get KPI Kabkot for specific instansi
curl "http://localhost:8080/api/v1/kpi_kabkot.csv?instansi=BPS%20Kota%20Batam"

# Get Heatmap
curl http://localhost:8080/api/v1/heatmap.csv

# Get Heatmap for Batam only
curl "http://localhost:8080/api/v1/heatmap.csv?kab_kota=Batam"

# Get Issues Detail
curl http://localhost:8080/api/v1/issues_detail.csv

# Get completed issues only
curl "http://localhost:8080/api/v1/issues_detail.csv?status=completed"

# Get issues for Batam in Produksi bidang
curl "http://localhost:8080/api/v1/issues_detail.csv?kab_kota=Batam&bidang=Produksi"
```

### Flourish Live CSV
```
# Without filters
https://your-domain.com/api/v1/kpi_provinsi.csv
https://your-domain.com/api/v1/kpi_kabkot.csv
https://your-domain.com/api/v1/heatmap.csv
https://your-domain.com/api/v1/issues_detail.csv

# With filters (URL encoded)
https://your-domain.com/api/v1/kpi_provinsi.csv?bidang=Sosial
https://your-domain.com/api/v1/kpi_kabkot.csv?instansi=BPS%20Kota%20Batam
https://your-domain.com/api/v1/heatmap.csv?kab_kota=Batam
https://your-domain.com/api/v1/issues_detail.csv?status=completed
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

**Cache Key Strategy:**
- Base cache key: endpoint name (e.g., `kpi_provinsi`, `heatmap`)
- With filters: `{endpoint}_{filter1}_{filter2}_{filter3}` (e.g., `kpi_provinsi_Sosial_Ketua`)
- Each unique combination of filters gets its own cache entry

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
- ✅ `queries/heatmap_kabkot_bidang.sql` - Tidak digunakan lagi (converted to dynamic SQL)
- ❌ `queries/kpi_provinsi.sql` - Dihapus (dynamic SQL)
- ❌ `queries/kpi_kabkot.sql` - Dihapus (dynamic SQL)
- ❌ `queries/issues_detail.sql` - Dihapus (dynamic SQL)
- ❌ `queries/kpi.sql` - Dihapus (endpoint dihapus)
- ❌ `queries/progress_kabkot.sql` - Dihapus (endpoint dihapus)
- ❌ `queries/progress_bidang.sql` - Dihapus (endpoint dihapus)

### Filtering Feature
Semua endpoint mendukung filtering melalui query parameters:
- Filter bersifat **optional** - tanpa filter akan mengembalikan semua data
- Multiple filters dapat dikombinasikan dengan `&` (AND logic)
- Filter values harus URL-encoded jika mengandung spasi atau karakter khusus
- Cache key berbeda untuk setiap kombinasi filter
- SQL injection protection: semua input di-sanitize dengan escape single quotes

---

## Filtering Guide

### URL Encoding
Jika filter value mengandung spasi atau karakter khusus, harus di-encode:
- Spasi → `%20`
- Contoh: `BPS Kota Batam` → `BPS%20Kota%20Batam`

### Case Sensitivity
Filter values bersifat **case-sensitive**. Pastikan menggunakan kapitalisasi yang benar:
- ✅ `bidang=Sosial`
- ❌ `bidang=sosial`
- ❌ `bidang=SOSIAL`

### Multiple Filters
Multiple filters menggunakan AND logic (semua kondisi harus terpenuhi):
```bash
# Hanya menampilkan data yang memenuhi SEMUA kondisi
?kab_kota=Batam&bidang=Produksi&status=completed
```

### Empty Results
Jika kombinasi filter tidak menghasilkan data, endpoint akan mengembalikan CSV dengan header saja (tanpa data rows).

### Valid Filter Values

**Bidang:**
- `Sosial`
- `Produksi`
- `Distribusi`
- `Nerwilis`
- `Umum` (hanya untuk KPI Kabkot)
- `Sekretariat` (hanya untuk KPI Provinsi & Kabkot)

**Kab/Kota:**
- `Batam`
- `Tanjung Pinang`
- `Bintan`
- `Karimun`
- `Natuna`
- `Lingga`
- `Kep. Anambas`

**Instansi (KPI Kabkot):**
- `BPS Kota Batam`
- `BPS Kota Tanjung Pinang`
- `BPS Kabupaten Bintan`
- `BPS Kabupaten Karimun`
- `BPS Kabupaten Natuna`
- `BPS Kabupaten Lingga`
- `BPS Kabupaten Kepulauan Anambas`

**Jabatan (KPI Provinsi):**
- `Ketua`
- `Anggota`
- `Pengarah`
- `Ketua Pelaksana`
- `Ketua Sekretariat`
- `Wakil Ketua Sekretariat`
- `Anggota Sekretariat`

**Jabatan (KPI Kabkot):**
- `Ketua Bidang`
- `Kepala Kab/Kot`
- `Ketua Pelaksana`
- `Ketua Sekretariat`

**Status (Issues Detail):**
- `backlog` - Dicatat
- `unstarted` - Ditugaskan
- `started` - Sedang Dikerjakan
- `triage` - Darurat
- `completed` - Selesai

### Use Cases

**Dashboard Filtering:**
```javascript
// User memilih kab/kota dari dropdown
const kabKota = 'Batam';
const url = `https://api.serumpun.com/api/v1/heatmap.csv?kab_kota=${encodeURIComponent(kabKota)}`;

fetch(url)
  .then(res => res.text())
  .then(csv => {
    // Update Flourish visualization
    updateChart(csv);
  });
```

**Monitoring Specific Teams:**
```bash
# Monitor Ketua Bidang Sosial di semua kab/kota
curl "http://localhost:8080/api/v1/kpi_kabkot.csv?bidang=Sosial&jabatan=Ketua%20Bidang"
```

**Status Tracking:**
```bash
# Lihat semua issues yang masih backlog
curl "http://localhost:8080/api/v1/issues_detail.csv?status=backlog"

# Lihat issues yang sedang dikerjakan di Batam
curl "http://localhost:8080/api/v1/issues_detail.csv?kab_kota=Batam&status=started"
```

### Performance Tips

1. **Gunakan filter untuk mengurangi data transfer**
   - Lebih baik: `?kab_kota=Batam` (hanya data Batam)
   - Kurang efisien: Download semua data lalu filter di client

2. **Leverage cache**
   - Filter yang sama dalam 60 detik akan di-serve dari cache
   - Hindari filter yang terlalu spesifik jika tidak perlu

3. **Combine filters wisely**
   - Gunakan filter yang paling mengurangi dataset terlebih dahulu
   - Contoh: `?instansi=BPS%20Kota%20Batam&bidang=Sosial` lebih efisien daripada hanya `?bidang=Sosial`

---

**Last Updated:** 2026-01-30

# SQL Query Templates

Folder ini berisi SQL query templates yang digunakan oleh API endpoints.

## File Structure

### Template SQL (With Placeholders)
Semua endpoint sekarang menggunakan template SQL dengan placeholder:

**Core KPI:**
- `kpi_provinsi.sql` - KPI pegawai provinsi
- `kpi_kabkot.sql` - KPI ketua kabupaten/kota

**Supporting:**
- `heatmap.sql` - Heatmap matrix kab/kota × bidang
- `issues_detail.sql` - Issues detail dengan scope provinsi/kabkot

**Analytics:**
- `timeline.sql` - Timeline/Gantt chart untuk deadline tracking
- `leaderboard.sql` - Ranking pegawai berdasarkan performa
- `workload.sql` - Analisis distribusi beban kerja

### Legacy SQL (Deprecated)
- `heatmap_kabkot_bidang.sql` - ❌ Deprecated, gunakan `heatmap.sql`

## Placeholders

Template SQL menggunakan placeholder berikut yang akan di-replace oleh Go code:

### `{{NAMA_CASES}}`
CASE statement untuk mapping email → nama dari CSV directory.

**Example:**
```sql
CASE
  WHEN LOWER(u.email) = 'user1@example.com' THEN 'John Doe'
  WHEN LOWER(u.email) = 'user2@example.com' THEN 'Jane Smith'
  ELSE COALESCE(u.display_name, u.first_name || ' ' || u.last_name, '-')
END AS nama
```

### `{{SCOPE_CASES}}`
CASE statement untuk mapping email → scope (provinsi/kabkot).

**Example:**
```sql
CASE
  WHEN LOWER(u.email) = 'user1@example.com' THEN 'provinsi'
  WHEN LOWER(u.email) = 'user2@example.com' THEN 'kabkot'
  ELSE 'unknown'
END AS scope
```

### `{{INSTANSI_CASES}}`
CASE statement untuk mapping email → instansi.

**Example:**
```sql
CASE
  WHEN LOWER(u.email) = 'user1@example.com' THEN 'BPS Provinsi Kepulauan Riau'
  WHEN LOWER(u.email) = 'user2@example.com' THEN 'BPS Kota Batam'
  ELSE '-'
END AS instansi
```

### `{{BIDANG_CASES}}`
CASE statement untuk mapping email → bidang.

**Example:**
```sql
CASE
  WHEN LOWER(u.email) = 'user1@example.com' THEN 'Sosial'
  WHEN LOWER(u.email) = 'user2@example.com' THEN 'Produksi'
  ELSE NULL
END AS bidang
```

### `{{JABATAN_CASES}}`
CASE statement untuk mapping email → jabatan.

**Example:**
```sql
CASE
  WHEN LOWER(u.email) = 'user1@example.com' THEN 'Ketua'
  WHEN LOWER(u.email) = 'user2@example.com' THEN 'Anggota'
  ELSE 'Anggota'
END AS jabatan
```

### `{{EMAILS}}`
Comma-separated list of quoted emails untuk WHERE IN clause.

**Example:**
```sql
WHERE LOWER(u.email) IN ('user1@example.com', 'user2@example.com', 'user3@example.com')
```

### `{{WHERE_CLAUSE}}`
Dynamic WHERE clause berdasarkan query parameters (standalone).

**Example:**
```sql
WHERE
  scope = 'provinsi'
  AND bidang = 'Sosial'
  AND status = 'completed'
```

### `{{HAVING_CLAUSE}}`
Dynamic HAVING clause untuk GROUP BY queries.

**Example:**
```sql
HAVING
  kab_kota = 'Batam'
  AND bidang = 'Produksi'
```

## Usage in Go

Template SQL di-load dan di-process menggunakan helper functions di `query_builder.go`:

### KPI Endpoints
```go
// Load template
sqlTemplate, err := s.Queries.Load("kpi_provinsi.sql")

// Build dynamic parts
namaCases := buildNamaCases(dir.Provinsi)
bidangCases := buildBidangCases(dir.Provinsi)
jabatanCases := buildJabatanCases(dir.Provinsi)
emails := buildEmailList(dir.Provinsi)
whereClause := buildAdditionalWhere(filters)

// Replace placeholders
sql := buildDynamicSQL(sqlTemplate, map[string]string{
    "{{NAMA_CASES}}":    namaCases,
    "{{BIDANG_CASES}}":  bidangCases,
    "{{JABATAN_CASES}}": jabatanCases,
    "{{EMAILS}}":        emails,
    "{{WHERE_CLAUSE}}":  whereClause,
})
```

### Analytics Endpoints
```go
// Load template
sqlTemplate, err := s.Queries.Load("issues_detail.sql")

// Build dynamic parts
allRows := append(dir.Provinsi, dir.Kabkot...)
namaCases := buildNamaCases(allRows)
scopeCases := buildScopeCases(allRows)
whereClause := buildWhereClause(filters)

// Replace placeholders
sql := buildDynamicSQL(sqlTemplate, map[string]string{
    "{{NAMA_CASES}}":   namaCases,
    "{{SCOPE_CASES}}":  scopeCases,
    "{{WHERE_CLAUSE}}": whereClause,
})
```

## Helper Functions

Located in `internal/http/query_builder.go`:

- `buildDynamicSQL()` - Replace placeholders dalam template
- `buildNamaCases()` - Generate CASE statement untuk nama
- `buildScopeCases()` - Generate CASE statement untuk scope
- `buildInstansiCases()` - Generate CASE statement untuk instansi
- `buildBidangCases()` - Generate CASE statement untuk bidang
- `buildJabatanCases()` - Generate CASE statement untuk jabatan
- `buildEmailList()` - Generate comma-separated email list
- `buildWhereClause()` - Generate WHERE clause dari filters (standalone)
- `buildAdditionalWhere()` - Generate additional WHERE conditions (append to existing)
- `buildHavingClause()` - Generate HAVING clause dari filters
- `buildCacheKey()` - Generate consistent cache keys

## Benefits

1. **Separation of Concerns**: SQL logic terpisah dari Go code
2. **Maintainability**: Mudah update query tanpa touch Go code
3. **Readability**: SQL lebih mudah dibaca dalam file terpisah
4. **Testing**: SQL bisa di-test secara independen
5. **Version Control**: Changes pada SQL lebih jelas di git diff
6. **Consistency**: Semua endpoints menggunakan pattern yang sama
7. **DRY Principle**: Helper functions reusable across all endpoints

## File Organization

```
server/
├── queries/                          # SQL templates
│   ├── README.md                     # This file
│   ├── kpi_provinsi.sql             # KPI Provinsi template
│   ├── kpi_kabkot.sql               # KPI Kabkot template
│   ├── heatmap.sql                  # Heatmap template
│   ├── issues_detail.sql            # Issues Detail template
│   ├── timeline.sql                 # Timeline template
│   ├── leaderboard.sql              # Leaderboard template
│   └── workload.sql                 # Workload template
├── internal/http/
│   ├── handlers_kpi.go              # KPI handlers (using templates)
│   ├── handlers_analytics.go        # Analytics handlers (using templates)
│   ├── query_builder.go             # Helper functions
│   ├── directory.go                 # CSV directory loader
│   ├── csv.go                       # CSV utilities
│   └── router.go                    # Route definitions
```

## Migration Notes

All endpoints have been migrated to use SQL templates:

- ✅ KPI Provinsi - Migrated to `kpi_provinsi.sql`
- ✅ KPI Kabkot - Migrated to `kpi_kabkot.sql`
- ✅ Heatmap - Migrated to `heatmap.sql`
- ✅ Issues Detail - Migrated to `issues_detail.sql`
- ✅ Timeline - Migrated to `timeline.sql`
- ✅ Leaderboard - Migrated to `leaderboard.sql`
- ✅ Workload - Migrated to `workload.sql`

Old inline SQL in `handlers.go` is now deprecated and can be removed.

## Future Improvements

1. **SQL Validation**: Add SQL syntax validation before runtime
2. **Query Caching**: Cache compiled SQL templates (bukan hanya hasil query)
3. **Prepared Statements**: Migrate to parameterized queries for better security
4. **SQL Builder Library**: Consider using library seperti `squirrel` atau `goqu`
5. **Unit Tests**: Add tests for SQL templates dengan mock data

## Notes

- Placeholder format menggunakan `{{PLACEHOLDER_NAME}}` untuk mudah di-identify
- Semua user input di-sanitize dengan `sanitizeSQL()` function
- Template SQL harus valid PostgreSQL syntax
- Comments di SQL template akan tetap ada di final query (helpful untuk debugging)
- Empty filters tidak akan generate WHERE/HAVING clause (return all data)

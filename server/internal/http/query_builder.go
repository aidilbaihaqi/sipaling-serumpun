package httpx

import (
	"strings"
)

// buildDynamicSQL loads SQL template and replaces placeholders with dynamic content
func buildDynamicSQL(template string, replacements map[string]string) string {
	sql := template
	for placeholder, value := range replacements {
		sql = strings.ReplaceAll(sql, placeholder, value)
	}
	return sql
}

// buildNamaCases generates CASE statement for nama mapping from directory
func buildNamaCases(rows []DirectoryRow) string {
	var cases []string
	for _, row := range rows {
		email := strings.ToLower(row.Email)
		cases = append(cases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+sanitizeSQL(row.Nama)+"'")
	}
	return strings.Join(cases, "\n")
}

// buildScopeCases generates CASE statement for scope mapping from directory
func buildScopeCases(rows []DirectoryRow) string {
	var cases []string
	for _, row := range rows {
		email := strings.ToLower(row.Email)
		cases = append(cases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Scope+"'")
	}
	return strings.Join(cases, "\n")
}

// buildInstansiCases generates CASE statement for instansi mapping from directory
func buildInstansiCases(rows []DirectoryRow) string {
	var cases []string
	for _, row := range rows {
		email := strings.ToLower(row.Email)
		cases = append(cases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Instansi+"'")
	}
	return strings.Join(cases, "\n")
}

// buildBidangCases generates CASE statement for bidang mapping from directory
func buildBidangCases(rows []DirectoryRow) string {
	var cases []string
	for _, row := range rows {
		email := strings.ToLower(row.Email)
		cases = append(cases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Bidang+"'")
	}
	return strings.Join(cases, "\n")
}

// buildJabatanCases generates CASE statement for jabatan mapping from directory
func buildJabatanCases(rows []DirectoryRow) string {
	var cases []string
	for _, row := range rows {
		email := strings.ToLower(row.Email)
		cases = append(cases, "      WHEN LOWER(u.email) = '"+email+"' THEN '"+row.Jabatan+"'")
	}
	return strings.Join(cases, "\n")
}

// buildEmailList generates comma-separated quoted email list
func buildEmailList(rows []DirectoryRow) string {
	var emails []string
	for _, row := range rows {
		email := strings.ToLower(row.Email)
		emails = append(emails, "'"+email+"'")
	}
	return strings.Join(emails, ", ")
}

// buildWhereClause generates WHERE clause from filter map
func buildWhereClause(filters map[string]string) string {
	if len(filters) == 0 {
		return ""
	}

	var clauses []string
	for column, value := range filters {
		if value != "" {
			clauses = append(clauses, "  "+column+" = '"+sanitizeSQL(value)+"'")
		}
	}

	if len(clauses) == 0 {
		return ""
	}

	return "\nWHERE\n" + strings.Join(clauses, "\n  AND ")
}

// buildAdditionalWhere generates additional WHERE conditions (for appending to existing WHERE)
func buildAdditionalWhere(filters map[string]string) string {
	if len(filters) == 0 {
		return ""
	}

	var clauses []string
	for column, value := range filters {
		if value != "" {
			clauses = append(clauses, "  d."+column+" = '"+sanitizeSQL(value)+"'")
		}
	}

	if len(clauses) == 0 {
		return ""
	}

	return "\n  AND " + strings.Join(clauses, "\n  AND ")
}

// buildHavingClause generates HAVING clause from filter map
func buildHavingClause(filters map[string]string) string {
	if len(filters) == 0 {
		return ""
	}

	var clauses []string
	for column, value := range filters {
		if value != "" {
			clauses = append(clauses, "  "+column+" = '"+sanitizeSQL(value)+"'")
		}
	}

	if len(clauses) == 0 {
		return ""
	}

	return "\nHAVING\n" + strings.Join(clauses, "\n  AND ")
}

// buildCacheKey generates cache key from base and filters
func buildCacheKey(base string, filters map[string]string) string {
	if len(filters) == 0 {
		return base
	}

	key := base
	// Sort keys for consistent cache keys
	keys := []string{"scope", "kab_kota", "bidang", "instansi", "jabatan", "status"}
	for _, k := range keys {
		if v, ok := filters[k]; ok && v != "" {
			key += "_" + v
		}
	}
	return key
}

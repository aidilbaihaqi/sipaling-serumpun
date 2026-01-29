package httpx

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type DirectoryRow struct {
	Email    string
	Nama     string
	Instansi string
	Scope    string // provinsi | kabkot
	Jabatan  string // Ketua | Anggota
	Bidang   string
}

type DirectoryResult struct {
	Provinsi   []DirectoryRow
	Kabkot     []DirectoryRow
	BidangList []string
}

func LoadDirectoryFromCSV(path string) (DirectoryResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return DirectoryResult{}, fmt.Errorf("open directory csv: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	header, err := r.Read()
	if err != nil {
		return DirectoryResult{}, fmt.Errorf("read header: %w", err)
	}

	col := indexHeader(header)

	// Required columns based on your CSV:
	// Nama, Akun Gmail, Asal Instansi, Jabatan Dalam Tim SE2026
	req := []string{"Nama", "Akun Gmail", "Asal Instansi", "Jabatan Dalam Tim SE2026"}
	for _, k := range req {
		if _, ok := col[k]; !ok {
			return DirectoryResult{}, fmt.Errorf("missing column %q in directory csv", k)
		}
	}

	type rawRow struct {
		nama, email, instansi, jabatan string
	}

	var raws []rawRow
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return DirectoryResult{}, fmt.Errorf("read row: %w", err)
		}

		nama := strings.TrimSpace(rec[col["Nama"]])
		email := normEmail(rec[col["Akun Gmail"]])
		inst := strings.TrimSpace(rec[col["Asal Instansi"]])
		jab := strings.TrimSpace(rec[col["Jabatan Dalam Tim SE2026"]])

		if nama == "" || email == "" || inst == "" || jab == "" {
			continue
		}

		raws = append(raws, rawRow{nama: nama, email: email, instansi: inst, jabatan: jab})
	}

	// 1) derive bidang list from "Ketua/Anggota Bidang ..."
	bidangSet := map[string]struct{}{}
	for _, rr := range raws {
		if b, ok := parseBidang(rr.jabatan); ok {
			bidangSet[b] = struct{}{}
		}
	}
	bidangList := make([]string, 0, len(bidangSet))
	for b := range bidangSet {
		bidangList = append(bidangList, b)
	}

	// 2) build provinsi rows (ALL staff, not just Ketua/Anggota Bidang)
	var prov []DirectoryRow
	for _, rr := range raws {
		scope := deriveScope(rr.instansi)
		if scope != "provinsi" {
			continue
		}

		// Parse jabatan and bidang
		jab, bidang, isBidang := parseJabatanBidang(rr.jabatan)

		if isBidang {
			// Ketua/Anggota Bidang
			prov = append(prov, DirectoryRow{
				Email: rr.email, Nama: rr.nama, Instansi: rr.instansi,
				Scope: "provinsi", Jabatan: jab, Bidang: bidang,
			})
		} else {
			// Other roles: Pengarah, Ketua Pelaksana, Ketua Sekretariat, etc.
			// Map to appropriate bidang or mark as "Umum"
			jabNorm, bidangNorm := parseOtherJabatan(rr.jabatan)
			if jabNorm != "" {
				prov = append(prov, DirectoryRow{
					Email: rr.email, Nama: rr.nama, Instansi: rr.instansi,
					Scope: "provinsi", Jabatan: jabNorm, Bidang: bidangNorm,
				})
			}
		}
	}

	// 3) build kabkot rows (only Kepala Kab/Kot) cross join bidangList
	var kab []DirectoryRow
	for _, rr := range raws {
		scope := deriveScope(rr.instansi)
		if scope != "kabkot" {
			continue
		}
		// treat "Kepala Kab/Kot" as Ketua
		if strings.EqualFold(rr.jabatan, "Kepala Kab/Kot") {
			for _, b := range bidangList {
				kab = append(kab, DirectoryRow{
					Email: rr.email, Nama: rr.nama, Instansi: rr.instansi,
					Scope: "kabkot", Jabatan: "Ketua", Bidang: b,
				})
			}
		}
	}

	return DirectoryResult{Provinsi: prov, Kabkot: kab, BidangList: bidangList}, nil
}

func indexHeader(h []string) map[string]int {
	m := map[string]int{}
	for i, v := range h {
		m[strings.TrimSpace(v)] = i
	}
	return m
}

func normEmail(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func deriveScope(instansi string) string {
	s := strings.ToLower(instansi)
	switch {
	case strings.Contains(s, "provinsi"):
		return "provinsi"
	case strings.Contains(s, "kabupaten") || strings.Contains(s, "kota") || strings.Contains(s, "kab/kot"):
		return "kabkot"
	default:
		// fallback: treat "BPS Provinsi ..." as provinsi
		if strings.Contains(s, "bps provinsi") {
			return "provinsi"
		}
		return ""
	}
}

// "Ketua Bidang X" / "Anggota Bidang X"
func parseBidang(jabatan string) (string, bool) {
	_, b, ok := parseJabatanBidang(jabatan)
	return b, ok
}

func parseJabatanBidang(jabatan string) (jabatanNorm, bidang string, ok bool) {
	j := strings.TrimSpace(jabatan)
	if strings.HasPrefix(j, "Ketua Bidang ") {
		return "Ketua", strings.TrimSpace(strings.TrimPrefix(j, "Ketua Bidang ")), true
	}
	if strings.HasPrefix(j, "Anggota Bidang ") {
		return "Anggota", strings.TrimSpace(strings.TrimPrefix(j, "Anggota Bidang ")), true
	}
	return "", "", false
}

// parseOtherJabatan handles non-bidang roles like Pengarah, Ketua Pelaksana, Sekretariat
func parseOtherJabatan(jabatan string) (jabatanNorm, bidang string) {
	j := strings.TrimSpace(jabatan)
	jLower := strings.ToLower(j)

	switch {
	case strings.Contains(jLower, "pengarah"):
		return "Pengarah", "Umum"
	case strings.Contains(jLower, "ketua pelaksana"):
		return "Ketua Pelaksana", "Umum"
	case strings.Contains(jLower, "ketua sekretariat"):
		return "Ketua Sekretariat", "Sekretariat"
	case strings.Contains(jLower, "wakil ketua sekretariat"):
		return "Wakil Ketua Sekretariat", "Sekretariat"
	case strings.Contains(jLower, "anggota sekretariat"):
		return "Anggota Sekretariat", "Sekretariat"
	default:
		return "", ""
	}
}

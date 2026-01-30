# Setup Docker untuk SERUMPUN Backend

## Prasyarat
- Docker dan Docker Compose sudah terinstall di VPS
- Akses SSH ke VPS
- Port 8080 terbuka di firewall VPS
- File CSV `daftar_pengguna_serumpun.csv` ada di folder `data/`

## Langkah Setup

### 1. Persiapan File Environment
```bash
cd server
cp .env.example .env
```

Edit file `.env` dan sesuaikan dengan konfigurasi database Anda:
```env
APP_PORT=8080
DATABASE_URL=postgres://dashboard:Dash21board@postgresql.gurind.am:5433/plane
CACHE_TTL_SECONDS=60

# CSV Directory Path (relative to app root)
DIRECTORY_CSV_PATH=data/daftar_pengguna_serumpun.csv

WORKSPACE_NAME=Platform Serumpun
PROJECT_NAME=Sensus Ekonomi 2026

KABKOTA_KEY=kab_kota
```

### 2. Pastikan File CSV Ada
```bash
# Cek apakah file CSV ada
ls -la data/daftar_pengguna_serumpun.csv

# Jika belum ada, copy file CSV ke folder data/
mkdir -p data
cp /path/to/daftar_pengguna_serumpun.csv data/
```

### 3. Build dan Jalankan Container
```bash
# Build image
docker-compose build

# Jalankan container
docker-compose up -d
```

### 4. Cek Status Container
```bash
# Lihat container yang berjalan
docker-compose ps

# Lihat logs
docker-compose logs -f serumpun-api
```

### 5. Test API
```bash
# Health check
curl http://ALAMAT_IP_VPS:8080/healthz

# Test endpoint CSV
curl http://ALAMAT_IP_VPS:8080/api/v1/kpi_provinsi.csv
```

## Akses dari Luar VPS

Aplikasi akan bisa diakses di:
```
http://ALAMAT_IP_VPS:8080
```

Contoh endpoint:
- `http://ALAMAT_IP_VPS:8080/healthz`
- `http://ALAMAT_IP_VPS:8080/api/v1/kpi_provinsi.csv`
- `http://ALAMAT_IP_VPS:8080/api/v1/kpi_kabkot.csv`
- `http://ALAMAT_IP_VPS:8080/api/v1/heatmap.csv`
- `http://ALAMAT_IP_VPS:8080/api/v1/issues_detail.csv`
- `http://ALAMAT_IP_VPS:8080/api/v1/timeline.csv`
- `http://ALAMAT_IP_VPS:8080/api/v1/leaderboard.csv`
- `http://ALAMAT_IP_VPS:8080/api/v1/workload.csv`

## Perintah Docker Berguna

```bash
# Stop container
docker-compose down

# Restart container
docker-compose restart

# Rebuild dan restart
docker-compose up -d --build

# Lihat logs real-time
docker-compose logs -f

# Masuk ke container
docker-compose exec serumpun-api sh

# Cek apakah file CSV ada di container
docker-compose exec serumpun-api ls -la data/

# Hapus container dan image
docker-compose down --rmi all
```

## Troubleshooting

### Error: "failed to load directory: open directory csv"
Ini berarti file CSV tidak ditemukan di container. Solusi:
```bash
# 1. Pastikan file CSV ada di folder data/ sebelum build
ls -la data/daftar_pengguna_serumpun.csv

# 2. Rebuild container
docker-compose down
docker-compose up -d --build

# 3. Verifikasi file ada di container
docker-compose exec serumpun-api ls -la data/
```

### Container tidak bisa start
```bash
# Cek logs error
docker-compose logs serumpun-api

# Cek apakah port sudah digunakan
netstat -tulpn | grep 8080
```

### Tidak bisa connect ke database
- Pastikan DATABASE_URL di `.env` sudah benar
- Cek apakah VPS bisa akses ke database server
- Test koneksi: `telnet postgresql.gurind.am 5433`

### API tidak bisa diakses dari luar
- Pastikan port 8080 terbuka di firewall VPS
- Cek dengan: `sudo ufw status` (Ubuntu) atau `sudo firewall-cmd --list-all` (CentOS)
- Buka port: `sudo ufw allow 8080` (Ubuntu)

## Update Aplikasi

Jika ada perubahan code:
```bash
# Pull perubahan dari git
git pull

# Rebuild dan restart
cd server
docker-compose up -d --build
```

## Keamanan (Opsional)

Untuk production, pertimbangkan:
1. Gunakan reverse proxy (Nginx) dengan HTTPS
2. Batasi akses dengan firewall rules
3. Gunakan environment variables yang lebih aman
4. Setup monitoring dan logging

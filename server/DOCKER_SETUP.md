# Setup Docker untuk SERUMPUN Backend

## Prasyarat
- Docker dan Docker Compose sudah terinstall di VPS
- Akses SSH ke VPS
- Port 8080 terbuka di firewall VPS

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

WORKSPACE_NAME=Platform Serumpun
PROJECT_NAME=Sensus Ekonomi 2026

KABKOTA_KEY=kab_kota
```

### 2. Build dan Jalankan Container
```bash
# Build image
docker-compose build

# Jalankan container
docker-compose up -d
```

### 3. Cek Status Container
```bash
# Lihat container yang berjalan
docker-compose ps

# Lihat logs
docker-compose logs -f serumpun-api
```

### 4. Test API
```bash
# Health check
curl http://ALAMAT_IP_VPS:8080/healthz

# Test endpoint CSV
curl http://ALAMAT_IP_VPS:8080/csv/kpi.csv
```

## Akses dari Luar VPS

Aplikasi akan bisa diakses di:
```
http://ALAMAT_IP_VPS:8080
```

Contoh endpoint:
- `http://ALAMAT_IP_VPS:8080/healthz`
- `http://ALAMAT_IP_VPS:8080/csv/kpi.csv`
- `http://ALAMAT_IP_VPS:8080/csv/progress_kabkot.csv`
- `http://ALAMAT_IP_VPS:8080/csv/progress_bidang.csv`
- `http://ALAMAT_IP_VPS:8080/csv/heatmap.csv`
- `http://ALAMAT_IP_VPS:8080/csv/issues_detail.csv`

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

# Hapus container dan image
docker-compose down --rmi all
```

## Troubleshooting

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

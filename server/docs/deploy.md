# Panduan Development & Deployment SERUMPUN Backend

## 🔄 Alur Development ke Production

### **Skenario 1: Development Lokal (Tanpa Docker)**

#### 1. Edit Code
```bash
# Edit file yang diperlukan
# Misalnya: internal/http/handlers.go, queries/kpi.sql, dll
```

#### 2. Test Lokal
```bash
cd server

# Install dependencies (jika ada perubahan go.mod)
go mod tidy

# Run langsung
go run ./cmd/api

# Atau build dulu
go build -o api ./cmd/api
./api
```

#### 3. Test Endpoint
```bash
# Health check
curl http://localhost:8080/healthz

# Test CSV endpoint
curl http://localhost:8080/csv/kpi.csv
```

---

### **Skenario 2: Development dengan Docker Lokal**

#### 1. Edit Code
```bash
# Edit file yang diperlukan
```

#### 2. Rebuild & Restart Container
```bash
cd server

# Stop container lama
docker-compose down

# Rebuild image dengan perubahan terbaru
docker-compose build

# Start container baru
docker-compose up -d

# Atau gabung jadi satu command
docker-compose up -d --build
```

#### 3. Lihat Logs
```bash
# Real-time logs
docker-compose logs -f serumpun-api

# Logs 100 baris terakhir
docker-compose logs --tail=100 serumpun-api
```

#### 4. Test Endpoint
```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/csv/kpi.csv
```

---

### **Skenario 3: Deploy ke VPS (Production)**

#### **Metode A: Manual Upload & Build di VPS**

```bash
# 1. Di komputer lokal - commit perubahan
git add .
git commit -m "Update: deskripsi perubahan"
git push origin main

# 2. SSH ke VPS
ssh user@IP_VPS

# 3. Masuk ke folder project
cd /path/to/serumpun/server

# 4. Pull perubahan terbaru
git pull origin main

# 5. Rebuild & restart container
docker-compose down
docker-compose up -d --build

# 6. Cek logs
docker-compose logs -f serumpun-api
```

#### **Metode B: Build Lokal, Push Image ke Registry**

```bash
# 1. Build image dengan tag
cd server
docker build -t serumpun-api:v1.0.1 .

# 2. Tag untuk registry (jika pakai Docker Hub)
docker tag serumpun-api:v1.0.1 username/serumpun-api:v1.0.1

# 3. Push ke registry
docker push username/serumpun-api:v1.0.1

# 4. Di VPS - pull & run
ssh user@IP_VPS
docker pull username/serumpun-api:v1.0.1
docker-compose down
docker-compose up -d
```

---

## 🚀 Quick Commands Cheat Sheet

### **Development Lokal**
```bash
# Run tanpa Docker
cd server
go run ./cmd/api

# Build binary
go build -o api ./cmd/api

# Run binary
./api
```

### **Docker Lokal**
```bash
# Build & run
docker-compose up -d --build

# Stop
docker-compose down

# Restart
docker-compose restart

# Logs
docker-compose logs -f
```

### **Production VPS**
```bash
# Update code dari git
git pull

# Rebuild & deploy
docker-compose up -d --build

# Restart tanpa rebuild
docker-compose restart

# Stop semua
docker-compose down

# Lihat status
docker-compose ps

# Lihat logs
docker-compose logs -f serumpun-api
```

---

## 🔍 Troubleshooting

### **Container tidak mau start setelah rebuild**
```bash
# Lihat error di logs
docker-compose logs serumpun-api

# Hapus container & image lama
docker-compose down --rmi all

# Build ulang dari awal
docker-compose up -d --build
```

### **Perubahan code tidak terdeteksi**
```bash
# Pastikan tidak ada cache
docker-compose build --no-cache

# Atau hapus semua dan build ulang
docker-compose down --rmi all --volumes
docker-compose up -d --build
```

### **Port sudah digunakan**
```bash
# Cek proses yang pakai port 8080
netstat -tulpn | grep 8080

# Kill proses (ganti PID)
kill -9 PID

# Atau ubah port di docker-compose.yml
ports:
  - "8081:8080"  # Akses via port 8081
```

### **Database connection error**
```bash
# Cek .env file
cat .env

# Test koneksi dari container
docker-compose exec serumpun-api sh
# Di dalam container:
ping postgresql.gurind.am
```

---

## 📝 Best Practices

### **1. Versioning**
```bash
# Tag setiap release
git tag v1.0.1
git push origin v1.0.1

# Build dengan version tag
docker build -t serumpun-api:v1.0.1 .
```

### **2. Backup sebelum deploy**
```bash
# Backup container lama
docker commit serumpun-api serumpun-api-backup

# Atau export image
docker save serumpun-api:latest > serumpun-api-backup.tar
```

### **3. Zero-downtime deployment (Advanced)**
```bash
# Build image baru dengan tag berbeda
docker-compose build

# Start container baru di port berbeda
# Edit docker-compose.yml sementara ke port 8081

# Test di port 8081
curl http://localhost:8081/healthz

# Jika OK, switch port kembali ke 8080
# Stop container lama, start yang baru
```

### **4. Monitoring**
```bash
# Cek resource usage
docker stats serumpun-api

# Cek health secara berkala
watch -n 5 'curl -s http://localhost:8080/healthz'
```

---

## 🔄 Workflow Lengkap (Recommended)

```bash
# === DI LOKAL ===
# 1. Edit code
vim internal/http/handlers.go

# 2. Test lokal tanpa Docker
go run ./cmd/api
# Test di http://localhost:8080

# 3. Test dengan Docker lokal
docker-compose up -d --build
# Test di http://localhost:8080

# 4. Commit & push
git add .
git commit -m "feat: tambah endpoint baru"
git push origin main

# === DI VPS ===
# 5. SSH ke VPS
ssh user@IP_VPS

# 6. Pull & deploy
cd /path/to/serumpun/server
git pull origin main
docker-compose up -d --build

# 7. Verify
curl http://IP_VPS:8080/healthz
docker-compose logs -f serumpun-api

# 8. Test dari browser
# Buka: http://IP_VPS:8080/csv/kpi.csv
```

---

## 🎯 Refresh di IP Address

Setelah deploy, aplikasi otomatis refresh karena:
1. Container lama di-stop
2. Container baru di-start dengan code terbaru
3. Endpoint langsung available di IP yang sama

**Tidak perlu restart server/VPS**, cukup restart container!

---

## ⚡ Hot Reload untuk Development (Opsional)

Jika ingin auto-reload saat development:

### Install Air (Go hot reload)
```bash
go install github.com/air-verse/air@latest
```

### Buat file .air.toml
```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/api ./cmd/api"
  bin = "tmp/api"
  include_ext = ["go", "sql"]
  exclude_dir = ["tmp"]
  delay = 1000
```

### Run dengan Air
```bash
air
# Setiap perubahan .go atau .sql akan auto-rebuild & restart
```

---

## 📊 Monitoring Production

```bash
# Cek uptime
docker-compose ps

# Cek logs error
docker-compose logs serumpun-api | grep -i error

# Cek memory usage
docker stats serumpun-api --no-stream

# Cek disk space
df -h
```

---

Dengan panduan ini, Anda bisa dengan mudah:
✅ Develop & test lokal
✅ Deploy ke VPS
✅ Rollback jika ada masalah
✅ Monitor aplikasi production

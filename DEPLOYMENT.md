# 🚀 Deployment Guide - MultipleScrape API

## 📋 DOM Cloud Deployment

### **Prerequisites**
- DOM Cloud account
- Git repository dengan project ini
- Go modules sudah di-setup (`go.mod` dan `go.sum`)

### **File Structure untuk Deployment**
```
multiplescrape/
├── domcloud.yml          # Deployment configuration
├── go.mod               # Go modules
├── go.sum               # Go dependencies checksum
├── main.go              # Main application
├── docs/                # Swagger documentation
├── repository/          # Business logic
├── static/              # Dashboard files
└── public_html/
    └── public/
        └── .htaccess    # Nginx configuration
```

### **Deployment Steps**

#### **1. Upload ke DOM Cloud**
```bash
# Via Git (Recommended)
git add .
git commit -m "Setup deployment configuration"
git push origin main

# Atau upload manual via File Manager DOM Cloud
```

#### **2. Konfigurasi di DOM Cloud Panel**
1. Login ke DOM Cloud panel
2. Pilih domain/subdomain
3. Go to **Deployment** section
4. Set **Source** ke Git repository atau upload manual
5. Pastikan file `domcloud.yml` ada di root directory

#### **3. Deploy Otomatis**
DOM Cloud akan otomatis:
- Install Go compiler
- Download dependencies (`go mod download`)
- Build aplikasi (`go build -o app`)
- Start aplikasi dengan Passenger

### **Environment Variables**
- `PORT` - Otomatis di-set oleh DOM Cloud
- Server akan bind ke `127.0.0.1:$PORT`

### **URL Endpoints Setelah Deploy**
Ganti `yourdomain.com` dengan domain DOM Cloud Anda:

#### **API Endpoints**
- **Health Check**: `https://yourdomain.com/health`
- **Home Data**: `https://yourdomain.com/api/v1/home`
- **Search**: `https://yourdomain.com/api/v1/search/?query=naruto`
- **Movies**: `https://yourdomain.com/api/v1/movie/`
- **Schedule**: `https://yourdomain.com/api/v1/jadwal-rilis/`
- **Anime Detail**: `https://yourdomain.com/api/v1/anime-detail/?anime_slug=one-piece`
- **Episode Detail**: `https://yourdomain.com/api/v1/episode-detail/?episode_url=...`

#### **Dashboard & Monitoring**
- **Basic Dashboard**: `https://yourdomain.com/static/dashboard.html`
- **Advanced Dashboard**: `https://yourdomain.com/static/advanced-dashboard.html`
- **System Monitoring**: `https://yourdomain.com/monitoring`
- **API Documentation**: `https://yourdomain.com/swagger/index.html`

### **Troubleshooting**

#### **Build Errors**
```bash
# Check Go version
go version

# Verify dependencies
go mod verify
go mod tidy

# Test build locally
go build -o app
```

#### **Runtime Errors**
1. Check DOM Cloud logs di panel
2. Verify PORT environment variable
3. Pastikan aplikasi bind ke `127.0.0.1:$PORT`

#### **Static Files Issues**
1. Pastikan direktori `static/` ada
2. Check file permissions
3. Verify `.htaccess` configuration

### **Performance Optimization**

#### **Build Optimization**
```bash
# Build dengan optimasi
go build -ldflags="-s -w" -o app

# Atau tambahkan ke domcloud.yml:
commands:
  - go mod download
  - go build -ldflags="-s -w" -o app
```

#### **Memory Management**
- Aplikasi sudah include memory monitoring
- Check `/monitoring` endpoint untuk usage
- DOM Cloud biasanya limit memory per aplikasi

### **Security Considerations**

#### **CORS Configuration**
- CORS sudah di-enable di aplikasi
- Additional CORS rules di `.htaccess`

#### **Rate Limiting**
```go
// Tambahkan rate limiting jika diperlukan
import "github.com/gin-contrib/cors"
import "github.com/gin-contrib/limiter"
```

#### **Environment Variables**
```bash
# Jangan commit sensitive data
# Gunakan DOM Cloud environment variables
```

### **Monitoring & Maintenance**

#### **Health Checks**
- Endpoint: `/health`
- Returns: `{"status": "ok"}`
- Use untuk monitoring uptime

#### **System Monitoring**
- Endpoint: `/monitoring`
- Provides: memory, uptime, request stats
- JSON format untuk integration

#### **Dashboard Monitoring**
- Real-time endpoint testing
- Confidence score validation
- Performance metrics
- Error tracking

### **Backup & Recovery**

#### **Database Backup**
- Aplikasi ini stateless (no database)
- Backup hanya diperlukan untuk static files

#### **Configuration Backup**
```bash
# Backup important files
cp domcloud.yml domcloud.yml.backup
cp go.mod go.mod.backup
```

### **Scaling Considerations**

#### **Horizontal Scaling**
- DOM Cloud mendukung multiple instances
- Load balancer otomatis di-handle

#### **Vertical Scaling**
- Upgrade plan DOM Cloud untuk more resources
- Monitor memory usage via `/monitoring`

### **Support & Documentation**

#### **DOM Cloud Support**
- Documentation: https://domcloud.co/docs
- Support: support@domcloud.co

#### **Application Support**
- GitHub Issues: [Your Repository]
- API Documentation: `/swagger/index.html`
- Monitoring Dashboard: `/static/dashboard.html`

---

## 🎯 **Quick Deploy Checklist**

- [ ] File `domcloud.yml` sudah dibuat
- [ ] `go.mod` dan `go.sum` sudah ada
- [ ] Main.go sudah support PORT environment variable
- [ ] Static files sudah di direktori `static/`
- [ ] `.htaccess` sudah dikonfigurasi
- [ ] Repository sudah di-push ke Git
- [ ] DOM Cloud deployment sudah dikonfigurasi
- [ ] Test semua endpoints setelah deploy
- [ ] Verify dashboard bisa diakses
- [ ] Check monitoring endpoint

**Happy Deploying! 🚀**
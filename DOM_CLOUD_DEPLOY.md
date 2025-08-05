# 🚀 DOM Cloud Deployment Guide - MultipleScrape API

## ✅ **Deployment Files Ready**

### **📁 File Structure:**
```
multiplescrape/
├── 📄 domcloud.yml              # DOM Cloud deployment configuration
├── 📄 startup.sh                # Application startup script
├── 📄 passenger_wsgi.py         # Passenger compatibility
├── 📄 Procfile                  # Process file
├── 📄 .domcloudrc              # Environment variables
├── 📄 go.mod & go.sum          # Go dependencies
├── 📄 main.go                   # Main application (PORT ready)
├── 📁 static/                   # Dashboard files
├── 📁 docs/                     # Swagger documentation
├── 📁 repository/               # Business logic
└── 📁 public_html/public/       # Static web files
    ├── index.html               # Landing page
    └── .htaccess               # Nginx rules
```

## 🎯 **Deployment Steps**

### **Step 1: Upload ke DOM Cloud**

#### **Via Git (Recommended):**
```bash
# Initialize git if not done
git init
git add .
git commit -m "Initial deployment setup"

# Push to your repository
git remote add origin https://github.com/yourusername/multiplescrape.git
git push -u origin main
```

#### **Via File Upload:**
1. Zip seluruh project folder
2. Upload via DOM Cloud File Manager
3. Extract di root directory

### **Step 2: DOM Cloud Panel Configuration**

1. **Login** ke DOM Cloud panel
2. **Select Domain** yang akan digunakan
3. **Go to Deployment** section
4. **Set Source:**
   - Git Repository: `https://github.com/yourusername/multiplescrape.git`
   - Branch: `main`
   - Auto Deploy: `Enabled`

### **Step 3: Deployment Process**

DOM Cloud akan otomatis menjalankan:

```yaml
# domcloud.yml akan execute:
1. 🚀 Starting MultipleScrape API deployment...
2. 📦 Downloading dependencies (go mod download)
3. ✅ Verifying dependencies (go mod verify)
4. 🔨 Building optimized binary (CGO_ENABLED=0 GOOS=linux)
5. 🔧 Setting permissions (chmod +x)
6. ✅ Build completed successfully!
7. 🧪 Testing binary
```

### **Step 4: Verification**

Setelah deployment selesai, test endpoints:

#### **🔗 Primary URLs:**
- **Landing Page**: `https://yourdomain.com/`
- **Health Check**: `https://yourdomain.com/health`
- **API Home**: `https://yourdomain.com/api/v1/home`

#### **📊 Dashboard URLs:**
- **Basic Dashboard**: `https://yourdomain.com/static/dashboard.html`
- **Advanced Dashboard**: `https://yourdomain.com/static/advanced-dashboard.html`
- **System Monitoring**: `https://yourdomain.com/monitoring`
- **API Documentation**: `https://yourdomain.com/swagger/index.html`

## 🔧 **Configuration Details**

### **domcloud.yml Configuration:**
```yaml
features:
  - go                           # Enable Go compiler

nginx:
  root: public_html/public       # Static files directory
  locations:
    - match: /static             # Dashboard files
      alias: static
      expires: 1d
    - match: /swagger            # API documentation
      try_files: $uri @app
    - match: /api                # API endpoints
      try_files: $uri @app
    - match: /                   # All other requests
      try_files: $uri @app

passenger:
  enabled: "on"                  # Enable Passenger
  app_start_command: ./startup.sh # Startup script
  app_type: generic              # Generic application type

commands:                        # Build commands
  - go mod download              # Download dependencies
  - go mod verify                # Verify dependencies
  - CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app
  - chmod +x app startup.sh      # Set permissions
```

### **startup.sh Script:**
```bash
#!/bin/bash
export PORT=${PORT:-8080}        # Default port
export GIN_MODE=${GIN_MODE:-release}  # Production mode
exec ./app                       # Start application
```

### **Environment Variables:**
- `PORT`: Otomatis di-set oleh DOM Cloud
- `GIN_MODE`: Set ke `release` untuk production
- `CGO_ENABLED`: Set ke `0` untuk static binary
- `GOOS`: Set ke `linux` untuk compatibility

## 📊 **Features Available After Deployment**

### **✅ API Endpoints (8 Total):**
1. **Home Data**: `/api/v1/home`
2. **Search**: `/api/v1/search/?query=naruto`
3. **Movies**: `/api/v1/movie/`
4. **Schedule**: `/api/v1/jadwal-rilis/`
5. **Daily Schedule**: `/api/v1/jadwal-rilis/monday`
6. **Anime Detail**: `/api/v1/anime-detail/?anime_slug=one-piece`
7. **Episode Detail**: `/api/v1/episode-detail/?episode_url=...`
8. **Latest Anime**: `/api/v1/anime-terbaru/`

### **✅ Monitoring Features:**
- Real-time endpoint testing
- Confidence score validation
- Response time tracking
- System performance metrics
- Memory usage monitoring
- Request statistics
- Error tracking & alerts

### **✅ Dashboard Features:**
- **Basic Dashboard**: Real-time monitoring dengan auto-refresh
- **Advanced Dashboard**: System performance & analytics
- **Mobile Responsive**: Works on all devices
- **Export Reports**: JSON data export
- **Auto-detection**: Production URL detection

## 🚨 **Troubleshooting**

### **Build Errors:**
```bash
# Check logs di DOM Cloud panel
# Common issues:
1. Missing go.mod/go.sum
2. Import path errors
3. Dependency conflicts
```

### **Runtime Errors:**
```bash
# Check application logs
# Common issues:
1. PORT binding issues (must use 127.0.0.1:$PORT)
2. File permission errors
3. Static file path issues
```

### **Swagger Issues (FIXED):**
```bash
# ✅ SOLVED: Custom Swagger implementation
# Features:
1. Dynamic host detection
2. HTTPS/HTTP auto-detection
3. Production URL compatibility
4. Custom swagger.json endpoint
5. No more localhost hardcoding

# Test endpoints:
- /swagger/index.html (Custom UI)
- /swagger/doc.json (Dynamic JSON)
- /swagger/ (Auto redirect)
```

### **Dashboard Not Loading:**
```bash
# Check:
1. Static files uploaded correctly
2. Nginx configuration
3. CORS headers
4. JavaScript console errors
```

## 🎯 **Performance Optimization**

### **Binary Optimization:**
- **Size**: 32MB (optimized with `-ldflags="-s -w"`)
- **CGO**: Disabled untuk static binary
- **GOOS**: Linux untuk compatibility

### **Caching:**
- **Static Files**: 1 day cache untuk dashboard
- **API Responses**: No cache (real-time data)
- **Nginx**: Optimized routing

### **Memory Management:**
- **Built-in Monitoring**: `/monitoring` endpoint
- **Goroutine Tracking**: Real-time monitoring
- **Memory Usage**: Tracked dan displayed

## 📈 **Monitoring & Maintenance**

### **Health Monitoring:**
```bash
# Automated health checks
curl https://yourdomain.com/health
# Expected: {"status":"ok"}
```

### **Performance Monitoring:**
```bash
# System metrics
curl https://yourdomain.com/monitoring
# Returns: JSON dengan memory, uptime, requests
```

### **Dashboard Monitoring:**
- Access dashboard untuk real-time monitoring
- Set auto-refresh untuk continuous monitoring
- Export reports untuk analysis

## 🔐 **Security Considerations**

### **CORS Configuration:**
- Enabled untuk dashboard access
- Configured di application level
- Additional rules di `.htaccess`

### **Environment Variables:**
- Sensitive data via DOM Cloud environment
- No hardcoded credentials
- Production mode enabled

### **Static Files:**
- Served via Nginx (faster)
- Cached untuk performance
- Secure headers enabled

## 📞 **Support & Resources**

### **DOM Cloud Documentation:**
- **Main Docs**: https://domcloud.co/docs
- **Go Deployment**: https://domcloud.co/docs/deployment/go
- **Support**: support@domcloud.co

### **Application Resources:**
- **GitHub**: [Your Repository URL]
- **API Docs**: `/swagger/index.html`
- **Monitoring**: `/static/dashboard.html`

---

## 🎉 **Deployment Checklist**

- [ ] ✅ `domcloud.yml` configured
- [ ] ✅ `startup.sh` executable
- [ ] ✅ `go.mod` & `go.sum` present
- [ ] ✅ Static files in correct directory
- [ ] ✅ Environment variables set
- [ ] ✅ Git repository configured
- [ ] ✅ DOM Cloud deployment configured
- [ ] ✅ All endpoints tested
- [ ] ✅ Dashboard accessible
- [ ] ✅ Monitoring functional

**🚀 Ready for Production Deployment!**

**Estimated Deployment Time**: 3-5 minutes  
**Binary Size**: 32MB  
**Memory Usage**: ~10MB  
**Startup Time**: <5 seconds  

**Happy Deploying! 🎯**
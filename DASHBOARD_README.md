# 📊 API Monitoring Dashboard

## 🚀 Akses Dashboard

Dashboard monitoring dapat diakses melalui:
- **Primary URL**: http://localhost:8080/static/dashboard.html
- **Alternative**: http://localhost:8080/dashboard (redirect)
- **Root**: http://localhost:8080/ (redirect)

## 📋 Fitur Dashboard

### 🎯 **Real-time Monitoring**
- **8 Endpoint** yang dipantau secara real-time
- **Confidence Score** untuk setiap endpoint
- **Response Time** monitoring
- **Data Count** tracking
- **Status Health Check**

### 📊 **Statistics Overview**
- Total Endpoints: 8
- Healthy Endpoints Counter
- Average Confidence Score
- Last Check Timestamp

### 🔧 **Interactive Controls**
- **🔄 Test All**: Test semua endpoint sekaligus
- **⏰ Auto Refresh**: Monitoring otomatis setiap 30 detik
- **🗑️ Clear Logs**: Bersihkan activity logs

### 📈 **Endpoint Coverage**

#### **1. Home** (`/api/v1/home`)
- **Field Wajib**: URL, Cover, Title, Slug
- **Validasi**: Top10, NewEps, Movies, JadwalRilis
- **Data Count**: Total items dari semua section

#### **2. Search** (`/api/v1/search/?query=one+piece`)
- **Field Wajib**: URL, Cover, Title, Slug
- **Validasi**: Hasil pencarian anime
- **Data Count**: Jumlah hasil pencarian

#### **3. Movies** (`/api/v1/movie/`)
- **Field Wajib**: URL, Cover, Title, Slug
- **Validasi**: Daftar anime movie
- **Data Count**: Jumlah movie

#### **4. Schedule** (`/api/v1/jadwal-rilis/`)
- **Field Wajib**: URL, Cover, Title, Slug
- **Validasi**: Jadwal rilis mingguan
- **Data Count**: Total anime di semua hari

#### **5. Schedule by Day** (`/api/v1/jadwal-rilis/monday`)
- **Field Wajib**: URL, Cover, Title, Slug
- **Validasi**: Jadwal rilis harian
- **Data Count**: Anime di hari tertentu

#### **6. Anime Detail** (`/api/v1/anime-detail/?anime_slug=one-piece`)
- **Field Wajib**: URL, Cover, Title, Slug, Episode List
- **Validasi**: Detail anime + episode list
- **Data Count**: Jumlah episode

#### **7. Episode Detail** (`/api/v1/episode-detail/?episode_url=...`)
- **Field Wajib**: Title, Streaming Servers
- **Validasi**: Detail episode + streaming links
- **Data Count**: Jumlah streaming servers

#### **8. Episode Detail (One Piece)** 
- **Test Case**: Complex episode URL parsing
- **Validasi**: Episode slug sanitization
- **Data Count**: Streaming servers

### 🎨 **Visual Indicators**

#### **Status Badges**
- 🟢 **Success**: Endpoint healthy (confidence > 0)
- 🟡 **Warning**: Partial data (confidence 0.5-0.8)
- 🔴 **Error**: Failed or no data (confidence = 0)
- ⚪ **Loading**: Testing in progress

#### **Confidence Bars**
- 🟢 **High (80-100%)**: Semua field wajib lengkap
- 🟡 **Medium (50-79%)**: Sebagian field kosong
- 🔴 **Low (0-49%)**: Banyak field kosong atau error

### 📝 **Activity Logs**
- Real-time logging semua aktivitas
- Timestamp untuk setiap test
- Success/Error status dengan detail
- Auto-limit 50 log entries

## 🔧 **Cara Penggunaan**

### **Manual Testing**
1. Klik **"🧪 Test Endpoint"** pada card endpoint tertentu
2. Lihat hasil di status badge dan detail metrics
3. Check confidence bar untuk kualitas data

### **Bulk Testing**
1. Klik **"🔄 Test All"** untuk test semua endpoint
2. Tunggu hingga semua selesai
3. Lihat summary di statistics overview

### **Auto Monitoring**
1. Klik **"⏰ Auto Refresh"** untuk enable
2. Dashboard akan test semua endpoint setiap 30 detik
3. Klik **"⏸️ Stop Auto"** untuk disable

### **Log Management**
1. Activity logs menampilkan semua aktivitas
2. Klik **"🗑️ Clear Logs"** untuk reset
3. Logs otomatis terbatas 50 entries

## 🎯 **Confidence Score System**

### **Kriteria Validasi**
- **Score 1.0**: Semua field wajib terisi
- **Score 0.0**: Ada field wajib yang kosong
- **No Score**: Endpoint error/tidak response

### **Field Wajib per Endpoint**
- **Home/Movie/Search/Jadwal**: url, cover, title, slug
- **Anime Detail**: url_anime, url_cover, judul, anime_slug + episode_list
- **Episode Detail**: title + minimal 1 streaming_server dengan url

## 🚨 **Troubleshooting**

### **Dashboard Tidak Load**
- Pastikan server berjalan di port 8080
- Check console browser untuk error JavaScript
- Pastikan CORS enabled di server

### **API Calls Failed**
- Verify API endpoints di Network tab browser
- Check server logs untuk error
- Pastikan confidence score validation berjalan

### **Auto Refresh Tidak Jalan**
- Check browser console untuk error
- Pastikan tidak ada popup blocker
- Refresh halaman dan coba lagi

## 📱 **Mobile Responsive**
Dashboard fully responsive untuk:
- Desktop (1400px+)
- Tablet (768px - 1400px)
- Mobile (< 768px)

## 🔗 **Links Terkait**
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health
- **API Base**: http://localhost:8080/api/v1

---

**Dashboard Version**: 1.0  
**Last Updated**: August 2025  
**Compatibility**: Modern browsers with ES6+ support
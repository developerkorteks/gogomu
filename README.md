# API Web Scraping GOMUNIME

API ini melakukan web scraping dari situs gomunime.co dan menyediakan data dalam format JSON yang sesuai dengan spesifikasi API_DEVELOP.md.

## Fitur

- ✅ Scraping data asli dari gomunime.co
- ✅ Struktur JSON sesuai dengan API_DEVELOP.md
- ✅ Confidence score untuk setiap response
- ✅ Error handling yang robust
- ✅ Rate limiting dan anti-bot measures
- ✅ Health check endpoint

## Endpoint yang Tersedia

### 1. Health Check
```
GET /health
```
Response:
```json
{
  "status": "ok"
}
```

### 2. Home Page
```
GET /api/v1/home
```
Mengembalikan data untuk halaman utama (top 10, anime terbaru, movie, jadwal rilis).

### 3. Anime Terbaru
```
GET /api/v1/anime-terbaru?page=<int>
```
Mengembalikan daftar anime terbaru dengan pagination.

### 4. Movie
```
GET /api/v1/movie?page=<int>
```
Mengembalikan daftar movie dengan pagination.

### 5. Jadwal Rilis
```
GET /api/v1/jadwal-rilis
```
Mengembalikan jadwal rilis untuk semua hari.

### 6. Jadwal Rilis per Hari
```
GET /api/v1/jadwal-rilis/:day
```
Mengembalikan jadwal rilis untuk hari tertentu (Monday, Tuesday, dll.).

### 7. Detail Anime
```
GET /api/v1/anime-detail?anime_slug=<string>
```
Mengembalikan detail lengkap untuk anime tertentu.

### 8. Detail Episode
```
GET /api/v1/episode-detail?episode_url=<string>
```
Mengembalikan detail lengkap untuk episode tertentu (termasuk link video dan download).

### 9. Search
```
GET /api/v1/search?query=<string>
```
Mengembalikan hasil pencarian anime.

## Struktur Response

Semua endpoint mengembalikan response dalam format berikut:

### Success Response
```json
{
  "confidence_score": 1.0,
  "message": "Data berhasil diambil",
  "source": "gomunime.co",
  "data": [...]
}
```

### Error Response
```json
{
  "error": true,
  "message": "Gagal mengambil data dari situs sumber: Timeout",
  "confidence_score": 0.0
}
```

## Installation

1. Clone repository ini
2. Install dependencies:
```bash
go mod tidy
```

3. Jalankan server:
```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/gocolly/colly/v2` - Web scraping library

## Rate Limiting

API ini menerapkan rate limiting untuk menghormati server target:
- Parallelism: 4 concurrent requests
- Random delay: 2 detik antara requests
- User-Agent rotation

## Confidence Score

API menghitung confidence score berdasarkan:
- `1.0`: Data lengkap dan akurat
- `0.5-0.9`: Data mungkin tidak lengkap atau ada anomali kecil
- `< 0.5`: Data sangat tidak lengkap atau ada masalah signifikan
- `0.0`: Scraping gagal total

## Error Handling

API menangani berbagai jenis error:
- Network timeout
- Connection refused
- HTTP error codes
- Invalid HTML structure
- Missing data

## Contoh Penggunaan

### Test Health Check
```bash
curl http://localhost:8080/health
```

### Test Home Page
```bash
curl http://localhost:8080/api/v1/home
```

### Test Search
```bash
curl "http://localhost:8080/api/v1/search?query=naruto"
```

### Test Anime Detail
```bash
curl "http://localhost:8080/api/v1/anime-detail?anime_slug=naruto"
```

## Catatan

- API ini melakukan scraping real-time dari gomunime.co
- Pastikan untuk menghormati robots.txt dan rate limiting
- Data yang di-scrape mungkin berubah sesuai dengan update di situs target
- Gunakan dengan bijak dan bertanggung jawab 
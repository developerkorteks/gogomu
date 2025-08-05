// file: main.go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"multiplescrape/docs"
	"multiplescrape/repository"
)

const (
	BaseDomain = "https://gomunime.co"
)

var (
	serverStartTime = time.Now()
	requestCount    = 0
	requestMutex    sync.Mutex
)

// Anotasi untuk informasi utama Swagger
// @title Gomunime Scraper API
// @version 1.0
// @description API untuk mengambil data anime terbaru dan jadwal rilis dari Gomunime.
// @contact.name   API Support
// @contact.url    https://github.com/you
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Request tracking middleware
	router.Use(func(c *gin.Context) {
		requestMutex.Lock()
		requestCount++
		requestMutex.Unlock()
		c.Next()
	})

	// Setup Swagger
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Serve static files for dashboard
	router.Static("/static", "./static")
	
	// Dashboard routes
	router.GET("/", func(c *gin.Context) {
		c.File("./static/dashboard.html")
	})
	router.GET("/dashboard", func(c *gin.Context) {
		c.File("./static/dashboard.html")
	})

	// === ROUTING ===
	router.GET("/health", healthCheckHandler)
	router.GET("/monitoring", monitoringHandler)
	
	// Grup endpoint baru untuk v1
	apiV1 := router.Group("/api/v1")
	{
		// Endpoint baru untuk jadwal rilis
		apiV1.GET("/home", getAnimeDataHandler)
		apiV1.GET("/jadwal-rilis/", getJadwalRilisHandler)
		apiV1.GET("/jadwal-rilis/:day", getJadwalRilisByDayHandler)
		apiV1.GET("/movie/", getMovieListHandler)
		apiV1.GET("/anime-detail/", getAnimeDetailHandler)
		apiV1.GET("/episode-detail/", getEpisodeDetailHandler)
		apiV1.GET("/anime-terbaru/", getAnimeTerbaruHandler)
		apiV1.GET("/search/", getSearchHandler)
		apiV1.GET("/monitoring", monitoringHandler) // Monitoring endpoint
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Server address - use 127.0.0.1 for DOM Cloud compatibility
	serverAddr := "127.0.0.1:" + port
	
	log.Println("🚀 Server berjalan di http://127.0.0.1:" + port)
	log.Println("📊 Dashboard Basic: /static/dashboard.html")
	log.Println("📈 Dashboard Advanced: /static/advanced-dashboard.html")
	log.Println("🔧 System Monitoring: /monitoring")
	log.Println("📚 Swagger UI: /swagger/index.html")
	log.Println("🔍 API Base URL: /api/v1")
	
	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}

// healthCheckHandler menangani permintaan health check.
// @Summary      Health Check
// @Description  Memeriksa apakah layanan berjalan dengan baik.
// @Tags         Utilities
// @Produce      json
// @Success      200  {object}  map[string]string "Layanan berjalan"
// @Router       /health [get]
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// monitoringHandler menangani permintaan monitoring sistem.
// @Summary      System Monitoring
// @Description  Mengambil informasi monitoring sistem dan performa API.
// @Tags         Utilities
// @Produce      json
// @Success      200  {object}  map[string]interface{} "Informasi monitoring"
// @Router       /monitoring [get]
func monitoringHandler(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(serverStartTime)
	
	requestMutex.Lock()
	totalRequests := requestCount
	requestMutex.Unlock()

	monitoring := gin.H{
		"server": gin.H{
			"status":     "running",
			"uptime":     uptime.String(),
			"start_time": serverStartTime.Format(time.RFC3339),
		},
		"performance": gin.H{
			"total_requests":    totalRequests,
			"requests_per_hour": float64(totalRequests) / uptime.Hours(),
			"memory_usage_mb":   float64(m.Alloc) / 1024 / 1024,
			"goroutines":        runtime.NumGoroutine(),
		},
		"endpoints": gin.H{
			"total":     8,
			"available": []string{
				"/api/v1/home",
				"/api/v1/search/",
				"/api/v1/movie/",
				"/api/v1/jadwal-rilis/",
				"/api/v1/jadwal-rilis/{day}",
				"/api/v1/anime-detail/",
				"/api/v1/episode-detail/",
				"/api/v1/anime-terbaru/",
			},
		},
		"system": gin.H{
			"go_version":    runtime.Version(),
			"os":           runtime.GOOS,
			"architecture": runtime.GOARCH,
			"cpu_cores":    runtime.NumCPU(),
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, monitoring)
}

// getSearchHandler menangani permintaan pencarian anime.
// @Summary      Search Anime
// @Description  Mencari anime berdasarkan query.
// @Tags         Anime List
// @Produce      json
// @Param        query  query  string  true  "Kata kunci pencarian"
// @Success      200  {object}  repository.SearchResponse "Hasil pencarian"
// @Failure      400  {object}  map[string]string "Parameter query tidak ditemukan"
// @Router       /api/v1/search/ [get]
func getSearchHandler(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'query' wajib diisi."})
		return
	}

	scrapedResults := repository.ScrapeSearch(query)

	var searchResults []repository.SearchResultItem
	for _, item := range scrapedResults {
		searchResults = append(searchResults, repository.SearchResultItem{
			Judul:     item.Judul,
			URLAnime:  item.Tautan,
			AnimeSlug: repository.GetSlugFromURL(item.Tautan),
			URLCover:  item.Thumbnail,
			Status:    repository.FillStrIfEmpty(item.Status, "N/A"),
			Tipe:      repository.FillStrIfEmpty(item.Tipe, "N/A"),
			Skor:      repository.FillStrIfEmpty(item.Skor, "N/A"),
			Sinopsis:  repository.FillStrIfEmpty(item.Sinopsis, "Sinopsis tidak tersedia."),
			Genre:     item.Genres,
			Penonton:  "N/A", // Data tidak tersedia dari scraper
		})
	}

	// Validasi data dan set confidence score
	confidenceScore := repository.ValidateSearchData(searchResults)

	response := repository.SearchResponse{
		ConfidenceScore: confidenceScore,
		Data:            searchResults,
		Message:         "Data berhasil diambil",
		Source:          "gomunime.co",
	}

	c.JSON(http.StatusOK, response)
}

// getAnimeTerbaruHandler menangani permintaan untuk daftar anime terbaru.
// @Summary      Get Latest Anime Releases
// @Description  Mengambil daftar rilis anime terbaru dengan paginasi.
// @Tags         Anime List
// @Produce      json
// @Param        page  query  int  false  "Nomor halaman"  default(1)
// @Success      200  {object}  repository.AnimeTerbaruResponse "Daftar rilis terbaru berhasil diambil"
// @Failure      400  {object}  map[string]string "Parameter halaman tidak valid"
// @Router       /api/v1/anime-terbaru/ [get]
func getAnimeTerbaruHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// Gunakan scraper yang sudah ada untuk mengambil data dari halaman utama
	latestItems := repository.ScrapeLatestByPage(page)

	var animeTerbaruList []repository.AnimeTerbaruItem
	for _, item := range latestItems {
		// Hilangkan "Episode" dari string episode
		episodeNum := strings.TrimSpace(strings.Replace(item.Episode, "Episode", "", -1))

		animeTerbaruList = append(animeTerbaruList, repository.AnimeTerbaruItem{
			Judul:     item.Judul,
			URL:       item.Tautan,
			AnimeSlug: repository.GetSlugFromURL(item.Tautan),
			Episode:   repository.FillStrIfEmpty(episodeNum, "N/A"),
			Cover:     item.Thumbnail,
			Uploader:  "N/A", // Data ini tidak tersedia dari scraper halaman utama
			Rilis:     "N/A", // Data ini tidak tersedia dari scraper halaman utama
		})
	}

	response := repository.AnimeTerbaruResponse{
		ConfidenceScore: 1.0,
		Data:            animeTerbaruList,
		Message:         "Data berhasil diambil",
		Source:          "gomunime.co",
	}

	c.JSON(http.StatusOK, response)
}

// GANTI FUNGSI LAMA ANDA DENGAN YANG INI
// getEpisodeDetailHandler menangani permintaan detail episode.
// @Summary      Get Episode Detail
// @Description  Mengambil detail lengkap sebuah episode berdasarkan URL.
// @Tags         Episode Detail
// @Accept       json
// @Produce      json
// @Param        episode_url  query  string  true  "URL lengkap dari halaman episode"
// @Param        force_refresh  query  boolean  false  "Force refresh cache (opsional, belum diimplementasikan)"
// @Success      200  {object}  repository.EpisodeDetailResponse "Detail episode berhasil diambil"
// @Failure      400  {object}  map[string]string "Parameter episode_url tidak valid atau kosong"
// @Failure      404  {object}  map[string]string "Episode tidak ditemukan"
// @Router       /api/v1/episode-detail/ [get]
func getEpisodeDetailHandler(c *gin.Context) {
	episodeURL := c.Query("episode_url")
	if episodeURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'episode_url' wajib diisi."})
		return
	}
	if _, err := url.ParseRequestURI(episodeURL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'episode_url' bukan URL yang valid."})
		return
	}

	// Panggil scraper baru yang sudah disempurnakan
	scrapedData := repository.ScrapeEpisodeDetail(episodeURL)
	if scrapedData.Title == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gagal mengambil data dari URL, mungkin halaman tidak ada."})
		return
	}

	// Format data ke dalam response akhir
	episodeDetailData := repository.EpisodeDetailData{
		Title:            scrapedData.Title,
		ThumbnailURL:     repository.FillStrIfEmpty(scrapedData.ThumbnailURL, scrapedData.AnimeInfo.ThumbnailURL),
		StreamingServers: scrapedData.StreamingServers,
		ReleaseInfo:      repository.FillStrIfEmpty(scrapedData.ReleaseInfo, "N/A"),
		DownloadLinks:    scrapedData.DownloadLinks,
		Navigation:       scrapedData.Navigation,
		AnimeInfo:        scrapedData.AnimeInfo,
		OtherEpisodes:    scrapedData.OtherEpisodes,
	}

	// Pastikan data wajib tidak kosong
	if episodeDetailData.Title == "" {
		episodeDetailData.Title = "Judul tidak ditemukan"
	}
	if episodeDetailData.ThumbnailURL == "" {
		episodeDetailData.ThumbnailURL = "https://placehold.co/200x300?text=No+Image"
	}

	// Validasi data dan set confidence score
	confidenceScore := repository.ValidateEpisodeDetailData(episodeDetailData)

	response := repository.EpisodeDetailResponse{
		ConfidenceScore: confidenceScore,
		Data:            episodeDetailData,
		Message:         "Data berhasil diambil",
		Source:          "gomunime.co",
	}

	c.JSON(http.StatusOK, response)
}

// GANTI FUNGSI LAMA ANDA DENGAN YANG INI
// getAnimeDetailHandler menangani permintaan detail anime.
// @Summary      Get Anime Detail
// @Description  Mengambil detail lengkap sebuah anime berdasarkan slug seri atau slug episode.
// @Tags         Anime Detail
// @Accept       json
// @Produce      json
// @Param        anime_slug  query  string  true  "Slug dari anime yang ingin dicari"
// @Param        force_refresh  query  boolean  false  "Force refresh cache (opsional, belum diimplementasikan)"
// @Success      200  {object}  repository.AnimeDetailResponse "Detail anime berhasil diambil"
// @Failure      400  {object}  map[string]string "Parameter anime_slug tidak ditemukan"
// @Failure      404  {object}  map[string]string "Anime tidak ditemukan"
// @Router       /api/v1/anime-detail/ [get]
func getAnimeDetailHandler(c *gin.Context) {
	originalSlug := c.Query("anime_slug")
	if originalSlug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'anime_slug' wajib diisi."})
		return
	}

	// --- Percobaan Pertama ---
	log.Printf("Mencoba mengambil detail untuk slug: %s", originalSlug)
	scrapedData := repository.ScrapeAnimeDetail(originalSlug)
	finalSlug := originalSlug

	// --- Percobaan Kedua (jika pertama gagal) ---
	if scrapedData.Judul == "" {
		log.Printf("Gagal pada percobaan pertama, mencoba membersihkan slug.")
		sanitizedSlug, wasSanitized := repository.SanitizeEpisodeSlug(originalSlug)

		if wasSanitized {
			log.Printf("Slug dibersihkan menjadi: %s. Mencoba lagi.", sanitizedSlug)
			scrapedData = repository.ScrapeAnimeDetail(sanitizedSlug)
			finalSlug = sanitizedSlug
		}
	}

	// --- Pengecekan Akhir ---
	if scrapedData.Judul == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Anime dengan slug '" + originalSlug + "' tidak ditemukan."})
		return
	}

	// --- Format Data ke dalam Response Akhir ---
	var episodeList []repository.EpisodeListItem
	for _, ep := range scrapedData.EpisodeList {
		episodeSlug := repository.GetSlugFromURL(ep.URL)
		episodeTitle := repository.SlugToTitle(episodeSlug)
		
		episodeList = append(episodeList, repository.EpisodeListItem{
			Episode:     repository.FillStrIfEmpty(ep.Episode, "N/A"),
			Title:       episodeTitle,
			URL:         ep.URL,
			EpisodeSlug: episodeSlug,
			ReleaseDate: repository.FillStrIfEmpty(ep.TanggalRilis, "N/A"),
		})
	}

	// Format Rekomendasi
	var recommendations []repository.RecommendationItem
	for _, rec := range scrapedData.Rekomendasi {
		recommendations = append(recommendations, repository.RecommendationItem{
			Title:     rec.Judul,
			URL:       rec.URL,
			AnimeSlug: repository.GetSlugFromURL(rec.URL),
			CoverURL:  rec.Thumbnail,
			Rating:    "N/A",
			Episode:   repository.FillStrIfEmpty(rec.Episode, "N/A"),
		})
	}

	// =====================================================================
	// LOGIKA FALLBACK JIKA REKOMENDASI KOSONG
	// =====================================================================
	if len(recommendations) == 0 {
		log.Println("Rekomendasi dari scraper kosong, mengambil fallback dari halaman utama.")

		// 1. Ambil data dari halaman utama
		fallbackAnime := repository.ScrapeLatestByPage(1)

		if len(fallbackAnime) > 0 {
			// 2. Acak urutan daftar fallback
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(fallbackAnime), func(i, j int) {
				fallbackAnime[i], fallbackAnime[j] = fallbackAnime[j], fallbackAnime[i]
			})

			// 3. Tentukan jumlah yang akan diambil (maksimal 3)
			numToTake := 3
			if len(fallbackAnime) < numToTake {
				numToTake = len(fallbackAnime)
			}

			// 4. Ambil 3 item pertama dari hasil acak dan format
			for i := 0; i < numToTake; i++ {
				item := fallbackAnime[i]
				// Hindari merekomendasikan anime yang sedang dilihat
				if repository.GetSlugFromURL(item.Tautan) == finalSlug {
					continue
				}
				recommendations = append(recommendations, repository.RecommendationItem{
					Title:     item.Judul,
					URL:       item.Tautan,
					AnimeSlug: repository.GetSlugFromURL(item.Tautan),
					CoverURL:  item.Thumbnail,
					Rating:    repository.FillStrIfEmpty(item.Rating, "N/A"),
					Episode:   repository.FillStrIfEmpty(item.Episode, "N/A"),
				})
			}
		}
	}
	// =====================================================================

	details := repository.Details{
		Japanese:     repository.FillStrIfEmpty(scrapedData.Details["Japanese"], "N/A"),
		Synonyms:     repository.FillStrIfEmpty(scrapedData.Details["Synonyms"], "N/A"),
		English:      repository.FillStrIfEmpty(scrapedData.Details["English"], "N/A"),
		Status:       repository.FillStrIfEmpty(scrapedData.Details["Status"], "N/A"),
		Type:         repository.FillStrIfEmpty(scrapedData.Details["Type"], "N/A"),
		Source:       repository.FillStrIfEmpty(scrapedData.Details["Source"], "N/A"),
		Duration:     repository.FillStrIfEmpty(scrapedData.Details["Duration"], "N/A"),
		TotalEpisode: repository.FillStrIfEmpty(scrapedData.Details["Episodes"], "N/A"),
		Season:       repository.FillStrIfEmpty(scrapedData.Details["Season"], "N/A"),
		Studio:       repository.FillStrIfEmpty(scrapedData.Details["Studio"], "N/A"),
		Producers:    repository.FillStrIfEmpty(scrapedData.Details["Producers"], "N/A"),
		Released:     repository.FillStrIfEmpty(scrapedData.Details["Released"], "N/A"),
	}

	animeDetailData := repository.AnimeDetailData{
		Judul:           scrapedData.Judul,
		URLAnime:        fmt.Sprintf("https://gomunime.co/anime/%s/", finalSlug),
		AnimeSlug:       finalSlug,
		URLCover:        scrapedData.Thumbnail,
		EpisodeList:     episodeList,
		Recommendations: recommendations,
		Status:          repository.FillStrIfEmpty(details.Status, "N/A"),
		Tipe:            repository.FillStrIfEmpty(details.Type, "N/A"),
		Skor:            repository.FillStrIfEmpty(scrapedData.Skor, "N/A"),
		Penonton:        "N/A",
		Sinopsis:        repository.FillStrIfEmpty(scrapedData.Sinopsis, "Sinopsis tidak tersedia."),
		Genre:           scrapedData.Genre,
		Details:         details,
		Rating: repository.RatingInfo{
			Score: repository.FillStrIfEmpty(scrapedData.Skor, "N/A"),
			Users: "N/A",
		},
	}

	// Validasi data dan set confidence score
	confidenceScore := repository.ValidateAnimeDetailData(animeDetailData)

	response := repository.AnimeDetailResponse{
		ConfidenceScore: confidenceScore,
		Data:            animeDetailData,
		Message:         "Data berhasil diambil",
		Source:          "gomunime.co",
	}

	c.JSON(http.StatusOK, response)
}

// getMovieListHandler menangani permintaan untuk daftar film.
// @Summary      Get "Movie" List (from Latest)
// @Description  Mengambil daftar anime dari halaman utama (latest) sebagai pengganti daftar film.
// @Tags         Movie
// @Accept       json
// @Produce      json
// @Param        page  query  int  false  "Nomor halaman"  default(1) mininum(1)
// @Success      200  {object}  repository.MovieListResponse "Daftar berhasil diambil"
// @Failure      400  {object}  map[string]string "Parameter halaman tidak valid"
// @Failure      500  {object}  map[string]string "Error internal server"
// @Router       /api/v1/movie/ [get]
func getMovieListHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameter 'page' harus berupa angka positif."})
		return
	}

	// Panggil scraper halaman utama, BUKAN scraper movie
	latestItems := repository.ScrapeLatestByPage(page)
	if latestItems == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari sumber."})
		return
	}

	// Format data TANPA FILTER TIPE MOVIE
	var movies []repository.MovieItem
	for _, item := range latestItems {
		// KONDISI IF UNTUK MOVIE DIHAPUS
		// Sekarang semua item dari halaman utama akan dimasukkan
		movie := repository.MovieItem{
			Judul:     item.Judul,
			URL:       item.Tautan,
			AnimeSlug: repository.GetSlugFromURL(item.Tautan),
			Cover:     item.Thumbnail,
			Status:    repository.FillStrIfEmpty(item.Status, "N/A"),
			Skor:      repository.FillStrIfEmpty(item.Rating, "N/A"),
			Sinopsis:  repository.FillStrIfEmpty(item.Deskripsi, "Sinopsis tidak tersedia."),
			Genres:    repository.FillSliceIfEmpty(item.Genres, []string{"Unknown"}),
			Views:     "N/A",
			Tanggal:   "N/A",
		}
		movies = append(movies, movie)
	}

	// Validasi data dan set confidence score
	confidenceScore := repository.ValidateMovieData(movies)

	response := repository.MovieListResponse{
		ConfidenceScore: confidenceScore,
		Data:            movies,
		Message:         "Data berhasil diambil",
		Source:          "gomunime.co",
	}

	c.JSON(http.StatusOK, response)
}

// getJadwalRilisByDayHandler menangani permintaan untuk jadwal rilis per hari.
// @Summary      Get Release Schedule by Day
// @Description  Mengambil jadwal rilis anime untuk hari yang spesifik.
// @Tags         Jadwal Rilis
// @Accept       json
// @Produce      json
// @Param        day  path  string  true  "Hari dalam bahasa Indonesia (e.g., senin, selasa)"
// @Param        force_refresh  query  boolean  false  "Force refresh cache (opsional, belum diimplementasikan)"
// @Success      200  {object}  repository.JadwalHarianResponse "Jadwal rilis berhasil diambil"
// @Failure      404  {object}  map[string]string "Hari tidak ditemukan"
// @Failure      500  {object}  map[string]string "Error internal server"
// @Router       /api/v1/jadwal-rilis/{day} [get]
func getJadwalRilisByDayHandler(c *gin.Context) {
	// Ambil parameter hari dari URL dan ubah ke huruf kecil
	requestedDay := strings.ToLower(c.Param("day"))

	// Scrape data jadwal
	scheduleData := repository.ScrapeSchedule()
	if scheduleData == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jadwal dari sumber."})
		return
	}

	// Cari hari yang cocok
	var foundDayData repository.ScrapedDaySchedule
	var dayFound bool
	for _, day := range scheduleData {
		// Cocokkan dengan nama hari dari scraper (sudah diubah ke huruf kecil)
		if strings.ToLower(day.Hari) == requestedDay {
			foundDayData = day
			dayFound = true
			break
		}
	}

	// Jika hari tidak ditemukan, kembalikan 404
	if !dayFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jadwal untuk hari '" + requestedDay + "' tidak ditemukan."})
		return
	}

	// Format data untuk hari yang ditemukan
	var animeList []repository.JadwalAnimeResponse
	for _, anime := range foundDayData.AnimeList {
		if anime.Judul == "" || anime.Tautan == "" || anime.Thumbnail == "" {
			continue
		}
		animeList = append(animeList, repository.JadwalAnimeResponse{
			Title:       anime.Judul,
			URL:         anime.Tautan,
			AnimeSlug:   repository.GetSlugFromURL(anime.Tautan),
			CoverURL:    anime.Thumbnail,
			ReleaseTime: repository.FillStrIfEmpty(anime.WaktuRilis, "N/A"),
			Type:        "TV",
			Score:       "N/A",
			Genres:      []string{"Unknown"},
		})
	}

	// Validasi data dan set confidence score
	confidenceScore := repository.ValidateJadwalData(animeList)

	// Bungkus dalam struct response akhir
	response := repository.JadwalHarianResponse{
		ConfidenceScore: confidenceScore,
		Data:            animeList,
		Message:         "Data berhasil diambil",
		Source:          "gomunime.co",
	}

	c.JSON(http.StatusOK, response)
}

// getJadwalRilisHandler menangani permintaan untuk jadwal rilis.
// @Summary      Get Release Schedule
// @Description  Mengambil jadwal rilis anime untuk semua hari.
// @Tags         Jadwal Rilis
// @Accept       json
// @Produce      json
// @Param        force_refresh  query  boolean  false  "Force refresh cache (opsional, belum diimplementasikan)"
// @Success      200  {object}  map[string]interface{}  "Jadwal rilis berhasil diambil"
// @Failure      500  {object}  map[string]string "Error internal server"
// @Router       /api/v1/jadwal-rilis/ [get]
func getJadwalRilisHandler(c *gin.Context) {
	// Membaca query param (walaupun belum diimplementasikan, ini untuk dokumentasi)
	_, _ = strconv.ParseBool(c.DefaultQuery("force_refresh", "false"))

	// Scrape data jadwal
	scheduleData := repository.ScrapeSchedule()
	if scheduleData == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data jadwal dari sumber."})
		return
	}

	// Format data menjadi map[string]interface{} untuk mencocokkan output JSON
	response := formatJadwalToMap(scheduleData)

	c.JSON(http.StatusOK, response)
}

func formatJadwalToMap(schedule []repository.ScrapedDaySchedule) map[string]interface{} {
	// Buat map utama untuk response
	finalMap := make(map[string]interface{})
	finalMap["confidence_score"] = 1.0
	finalMap["data"] = make(map[string]interface{})
	finalMap["message"] = "Data berhasil diambil"
	finalMap["source"] = "gomunime.co"

	// Proses setiap hari dari data scraper
	for _, day := range schedule {
		if day.Hari == "" {
			continue // Lewati jika nama hari kosong
		}

		var animeList []repository.JadwalAnimeResponse
		for _, anime := range day.AnimeList {
			// Hanya proses jika data esensial ada
			if anime.Judul == "" || anime.Tautan == "" || anime.Thumbnail == "" {
				continue
			}

			animeList = append(animeList, repository.JadwalAnimeResponse{
				Title:       anime.Judul,
				URL:         anime.Tautan,
				AnimeSlug:   repository.GetSlugFromURL(anime.Tautan),
				CoverURL:    anime.Thumbnail,
				ReleaseTime: repository.FillStrIfEmpty(anime.WaktuRilis, "N/A"),
				Type:        "TV",                // Data dummy karena tidak tersedia dari scraper
				Score:       "N/A",               // Data dummy
				Genres:      []string{"Unknown"}, // Data dummy
			})
		}

		// Masukkan list anime ke map data dengan key nama hari
		finalMap["data"].(map[string]interface{})[day.Hari] = animeList
	}

	return finalMap
}

// getAnimeDataHandler menangani permintaan API utama
// @Summary      Get Latest Anime & Schedule
// @Description  Mengambil daftar anime terbaru, film, top 10, dan jadwal rilis mingguan.
// @Tags         Anime
// @Accept       json
// @Produce      json
// @Success      200  {object}  repository.FinalResponse  "Data berhasil diambil"
// @Failure      500  {object}  map[string]string "Error internal server"
// @Router       /api/v1/home/ [get]
func getAnimeDataHandler(c *gin.Context) {
	// ... (Kode untuk menjalankan scraper secara concurrent tetap sama)
	var latestAnime []repository.ScrapedLatestAnime
	var scheduleData []repository.ScrapedDaySchedule
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		latestAnime = repository.ScrapeLatestAnime()
	}()
	go func() {
		defer wg.Done()
		scheduleData = repository.ScrapeSchedule()
	}()
	wg.Wait()

	if latestAnime == nil || scheduleData == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data dari sumber."})
		return
	}

	response := formatData(latestAnime, scheduleData)
	c.IndentedJSON(http.StatusOK, response)
}

// formatData sekarang menggunakan helper untuk mengisi data dummy.
func formatData(latest []repository.ScrapedLatestAnime, schedule []repository.ScrapedDaySchedule) repository.FinalResponse {
	// --- Top 10 ---
	top10List := []repository.Top10Anime{}
	limit := 10
	if len(latest) < 10 {
		limit = len(latest)
	}
	for i := 0; i < limit; i++ {
		item := latest[i]
		top10List = append(top10List, repository.Top10Anime{
			Judul:     item.Judul,
			URL:       item.Tautan,
			AnimeSlug: repository.GetSlugFromURL(item.Tautan),
			Rating:    repository.FillStrIfEmpty(item.Rating, "8.0"),
			Cover:     item.Thumbnail,
			Genres:    repository.FillSliceIfEmpty(item.Genres, []string{"Action", "Adventure"}),
		})
	}

	// --- New Eps ---
	newEpsList := []repository.NewEps{}
	for _, item := range latest {
		if item.Episode != "" && strings.ToLower(item.Tipe) == "tv" {
			newEpsList = append(newEpsList, repository.NewEps{
				Judul:     item.Judul,
				URL:       item.Tautan,
				AnimeSlug: repository.GetSlugFromURL(item.Tautan),
				Episode:   repository.FillStrIfEmpty(item.Episode, "N/A"),
				Rilis:     time.Now().Format("2 January 2006"), // Dummy rilis hari ini
				Cover:     item.Thumbnail,
			})
		}
	}

	// --- Movies ---
	// PERBAIKAN DI SINI: KONDISI 'IF' DIHAPUS
	movieList := []repository.Movie{}
	for _, item := range latest {
		// Kondisi 'if strings.ToLower(item.Tipe) == "movie"' telah dihapus.
		// Sekarang semua item terbaru akan dianggap sebagai "movie" untuk endpoint ini.
		movieList = append(movieList, repository.Movie{
			Judul:     item.Judul,
			URL:       item.Tautan,
			AnimeSlug: repository.GetSlugFromURL(item.Tautan),
			Tanggal:   "N/A", // Data tanggal rilis tidak tersedia dari scraper
			Cover:     item.Thumbnail,
			Genres:    repository.FillSliceIfEmpty(item.Genres, []string{"Unknown"}),
		})
	}

	// --- Jadwal Rilis ---
	jadwalMap := make(map[string][]repository.JadwalAnime)
	for _, day := range schedule {
		var animeList []repository.JadwalAnime
		for _, anime := range day.AnimeList {
			animeList = append(animeList, repository.JadwalAnime{
				Title:       anime.Judul,
				URL:         anime.Tautan,
				AnimeSlug:   repository.GetSlugFromURL(anime.Tautan),
				CoverURL:    anime.Thumbnail,
				ReleaseTime: repository.FillStrIfEmpty(anime.WaktuRilis, "00:00"),
				Type:        "TV",
				Score:       "N/A",
				Genres:      []string{"Unknown"},
			})
		}
		jadwalMap[day.Hari] = animeList
	}

	homeData := repository.HomeData{
		Top10:       top10List,
		NewEps:      newEpsList,
		Movies:      movieList,
		JadwalRilis: jadwalMap,
	}

	// Validasi data dan set confidence score
	confidenceScore := repository.ValidateHomeData(homeData)

	return repository.FinalResponse{
		ConfidenceScore: confidenceScore,
		Data:            homeData,
		Message:         "Data berhasil diambil",
		Source:          "gomunime.co",
	}
}

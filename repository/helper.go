package repository

import (
	"strconv"
	"strings"
)

// / GetSlugFromURL mengekstrak slug dari URL anime.
func GetSlugFromURL(url string) string {
	parts := strings.Split(strings.Trim(url, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "unknown-slug"
}

// FillStrIfEmpty mengembalikan nilai dummy jika string asli kosong.
func FillStrIfEmpty(value, dummy string) string {
	if value == "" {
		return dummy
	}
	return value
}

// FillSliceIfEmpty mengembalikan slice dummy jika slice string asli kosong.
func FillSliceIfEmpty(value, dummy []string) []string {
	if len(value) == 0 {
		return dummy
	}
	return value
}

func SanitizeEpisodeSlug(slug string) (string, bool) {
	// Pattern untuk mendeteksi slug episode
	patterns := []string{
		"-episode-",
		"-ep-",
		"-eps-",
	}
	
	// Cari pattern episode dalam slug
	for _, pattern := range patterns {
		if index := strings.Index(slug, pattern); index != -1 {
			// Ambil bagian sebelum pattern episode
			animeSlug := slug[:index]
			return animeSlug, true
		}
	}
	
	// Fallback ke logika lama jika tidak ada pattern episode
	lastDashIndex := strings.LastIndex(slug, "-")
	if lastDashIndex == -1 {
		return slug, false
	}

	potentialEpisodeNumber := slug[lastDashIndex+1:]
	if _, err := strconv.Atoi(potentialEpisodeNumber); err == nil {
		return slug[:lastDashIndex], true
	}

	return slug, false
}

// SlugToTitle mengubah slug menjadi title yang readable dengan menghapus tanda - dan mengubah ke title case
func SlugToTitle(slug string) string {
	// Hapus tanda - dan ganti dengan spasi
	title := strings.ReplaceAll(slug, "-", " ")
	
	// Ubah ke title case (huruf pertama setiap kata menjadi kapital)
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	
	return strings.Join(words, " ")
}

// ValidateHomeData memvalidasi data untuk endpoint home
func ValidateHomeData(data HomeData) float64 {
	score := 1.0
	
	// Validasi Top10
	for _, item := range data.Top10 {
		if item.Judul == "" || item.URL == "" || item.Cover == "" || item.AnimeSlug == "" {
			score = 0.0
			break
		}
	}
	
	// Validasi NewEps
	if score > 0 {
		for _, item := range data.NewEps {
			if item.Judul == "" || item.URL == "" || item.Cover == "" || item.AnimeSlug == "" {
				score = 0.0
				break
			}
		}
	}
	
	// Validasi Movies
	if score > 0 {
		for _, item := range data.Movies {
			if item.Judul == "" || item.URL == "" || item.Cover == "" || item.AnimeSlug == "" {
				score = 0.0
				break
			}
		}
	}
	
	// Validasi JadwalRilis
	if score > 0 {
		for _, daySchedule := range data.JadwalRilis {
			for _, item := range daySchedule {
				if item.Title == "" || item.URL == "" || item.CoverURL == "" || item.AnimeSlug == "" {
					score = 0.0
					break
				}
			}
			if score == 0 {
				break
			}
		}
	}
	
	return score
}

// ValidateMovieData memvalidasi data untuk endpoint movie
func ValidateMovieData(data []MovieItem) float64 {
	if len(data) == 0 {
		return 0.0
	}
	
	for _, item := range data {
		if item.Judul == "" || item.URL == "" || item.Cover == "" || item.AnimeSlug == "" {
			return 0.0
		}
	}
	
	return 1.0
}

// ValidateSearchData memvalidasi data untuk endpoint search
func ValidateSearchData(data []SearchResultItem) float64 {
	if len(data) == 0 {
		return 0.0
	}
	
	for _, item := range data {
		if item.Judul == "" || item.URLAnime == "" || item.URLCover == "" || item.AnimeSlug == "" {
			return 0.0
		}
	}
	
	return 1.0
}

// ValidateJadwalData memvalidasi data untuk endpoint jadwal
func ValidateJadwalData(data []JadwalAnimeResponse) float64 {
	if len(data) == 0 {
		return 0.0
	}
	
	for _, item := range data {
		if item.Title == "" || item.URL == "" || item.CoverURL == "" || item.AnimeSlug == "" {
			return 0.0
		}
	}
	
	return 1.0
}

// ValidateAnimeDetailData memvalidasi data untuk endpoint anime detail
func ValidateAnimeDetailData(data AnimeDetailData) float64 {
	// Field wajib untuk anime detail
	if data.Judul == "" || data.URLAnime == "" || data.URLCover == "" || data.AnimeSlug == "" {
		return 0.0
	}
	
	// Validasi episode list
	for _, episode := range data.EpisodeList {
		if episode.Title == "" || episode.URL == "" || episode.EpisodeSlug == "" {
			return 0.0
		}
	}
	
	return 1.0
}

// ValidateEpisodeDetailData memvalidasi data untuk endpoint episode detail
func ValidateEpisodeDetailData(data EpisodeDetailData) float64 {
	// Field wajib untuk episode detail
	if data.Title == "" {
		return 0.0
	}
	
	// Harus ada minimal 1 streaming server
	if len(data.StreamingServers) == 0 {
		return 0.0
	}
	
	// Validasi streaming servers
	for _, server := range data.StreamingServers {
		if server.StreamingURL == "" {
			return 0.0
		}
	}
	
	return 1.0
}

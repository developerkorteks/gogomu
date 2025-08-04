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

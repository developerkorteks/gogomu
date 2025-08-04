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

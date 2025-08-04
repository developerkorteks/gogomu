package gomunime

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/gocolly/colly/v2"
)

// [IMPROVED] Structs disempurnakan untuk data yang lebih lengkap
type LinkableAnimeLatest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type AnimeAnimeLatest struct {
	Judul     string                `json:"judul"`
	Tautan    string                `json:"tautan"`
	Episode   string                `json:"episode"`
	Thumbnail string                `json:"thumbnail"`
	Tipe      string                `json:"tipe"` // [ADDED] Tipe Anime (TV, ONA, etc.)
	Rating    string                `json:"rating,omitempty"`
	Durasi    string                `json:"durasi,omitempty"`
	Deskripsi string                `json:"deskripsi,omitempty"`
	Status    string                `json:"status,omitempty"`
	Studio    []LinkableAnimeLatest `json:"studio,omitempty"` // [IMPROVED] Bisa menampung banyak studio
	Genres    []LinkableAnimeLatest `json:"genres,omitempty"` // [IMPROVED] Menyimpan nama dan URL genre
}

func TestLatestRelease(t *testing.T) {
	baseURL := "https://gomunime.co/"
	ajaxURL := "https://gomunime.co/wp-admin/admin-ajax.php"

	c := colly.NewCollector(
		colly.Async(true),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
	c.DisableCookies()
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*gomunime.co*",
		Parallelism: 8, // Menggunakan 8 thread paralel untuk keseimbangan
	})

	var allAnime []AnimeAnimeLatest

	// Callback #1: Mengambil data detail lengkap dari info hover (AJAX response)
	c.OnHTML("div.ingfo", func(e *colly.HTMLElement) {
		// Mengambil struct Anime yang sudah diisi sebagian dari context
		animeInfo := e.Request.Ctx.GetAny("anime").(AnimeAnimeLatest)

		// [IMPROVED] Mengambil semua data dari 'minginfo'
		e.ForEach(".minginfo span", func(_ int, s *colly.HTMLElement) {
			text := s.Text
			if strings.Contains(text, "min.") {
				animeInfo.Durasi = text
			} else if s.DOM.HasClass("l") && s.DOM.Find("i.fa-star").Length() > 0 {
				animeInfo.Rating = strings.TrimSpace(text)
			}
		})

		animeInfo.Deskripsi = strings.TrimSpace(e.ChildText(".contexcerpt"))

		// [IMPROVED] Mengambil semua data dari 'linginfo'
		e.ForEach(".linginfo span", func(_ int, s *colly.HTMLElement) {
			text := s.Text
			if strings.HasPrefix(text, "Status:") {
				animeInfo.Status = strings.TrimSpace(strings.TrimPrefix(text, "Status:"))
			} else if strings.HasPrefix(text, "Genres:") {
				s.ForEach("a", func(_ int, a *colly.HTMLElement) {
					animeInfo.Genres = append(animeInfo.Genres, LinkableAnimeLatest{
						Name: a.Text,
						URL:  e.Request.AbsoluteURL(a.Attr("href")),
					})
				})
			} else if strings.HasPrefix(text, "Studio:") {
				s.ForEach("a", func(_ int, a *colly.HTMLElement) {
					animeInfo.Studio = append(animeInfo.Studio, LinkableAnimeLatest{
						Name: a.Text,
						URL:  e.Request.AbsoluteURL(a.Attr("href")),
					})
				})
			}
		})

		allAnime = append(allAnime, animeInfo)
	})

	// Callback #2: Mengambil data dasar dari setiap kartu anime di halaman
	c.OnHTML("div.listupd article.bs", func(e *colly.HTMLElement) {
		postID, exists := e.DOM.Find("a.tip").Attr("rel")
		if !exists {
			return
		}

		thumbURL := e.ChildAttr("img", "data-src")
		if thumbURL == "" {
			thumbURL = e.ChildAttr("img", "src")
		}

		anime := AnimeAnimeLatest{
			Judul:     e.ChildAttr("a.tip", "title"),
			Tautan:    e.Request.AbsoluteURL(e.ChildAttr("a.tip", "href")),
			Episode:   e.ChildText("span.epx"),
			Thumbnail: thumbURL,
			Tipe:      e.ChildText(".typez"), // [ADDED] Mengambil tipe anime
			// [IMPROVED] Inisialisasi slice agar tidak null di JSON
			Genres: make([]LinkableAnimeLatest, 0),
			Studio: make([]LinkableAnimeLatest, 0),
		}

		// Menyiapkan dan mengirim request AJAX untuk data hover
		ctx := colly.NewContext()
		ctx.Put("anime", anime)
		formData := fmt.Sprintf("action=tooltip_action&id=%s", postID)
		if err := c.Request("POST", ajaxURL, strings.NewReader(formData), ctx, nil); err != nil {
			log.Println("Gagal membuat request detail:", err)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		log.Printf("Selesai men-scrape halaman: %s", r.Request.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error:", err, "| URL:", r.Request.URL)
	})

	log.Println("Memulai proses scraping...")
	c.Visit(baseURL)
	c.Wait()

	log.Println("Semua proses scraping telah selesai.")
	jsonData, err := json.MarshalIndent(allAnime, "", "  ")
	if err != nil {
		log.Fatal("Gagal mengkonversi ke JSON:", err)
	}

	file, _ := os.Create("releases_lengkap.json")
	defer file.Close()
	file.Write(jsonData)
	fmt.Println(string(jsonData))
	log.Printf("Total %d data anime berhasil disimpan ke releases_lengkap.json", len(allAnime))
}

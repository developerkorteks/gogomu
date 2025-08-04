package gomunime

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gocolly/colly/v2"
)

// SearchResult diperbarui untuk menampung semua data dari hover box
type AnimeSearchResult struct {
	Judul       string   `json:"judul"`
	Tautan      string   `json:"tautan"`
	Thumbnail   string   `json:"thumbnail"`
	Tipe        string   `json:"tipe"`
	StatusRilis string   `json:"status_rilis"`
	Skor        string   `json:"skor,omitempty"`
	Durasi      string   `json:"durasi,omitempty"`
	Sinopsis    string   `json:"sinopsis,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Studio      string   `json:"studio,omitempty"`
}

func TestSearch(t *testing.T) {
	// URL pencarian dapat diubah sesuai kebutuhan
	searchURL := "https://gomunime.co/?s=kimet+su+no"
	ajaxURL := "https://gomunime.co/wp-admin/admin-ajax.php"

	var searchResults []AnimeSearchResult

	c := colly.NewCollector(
		colly.Async(true),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
	c.Limit(&colly.LimitRule{DomainGlob: "*gomunime.co*", Parallelism: 4})

	// Callback untuk memproses detail dari AJAX (info hover)
	// Selector tetap sama, namun logika di dalamnya ditingkatkan
	c.OnHTML("div.ingfo", func(e *colly.HTMLElement) {
		result := e.Request.Ctx.GetAny("result").(AnimeSearchResult)

		// Ekstrak informasi dari bagian atas tooltip
		e.ForEach("div.minginfo span.l", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if strings.Contains(text, "min. per ep.") {
				result.Durasi = text
			} else if el.ChildText("i.fa-star") != "" {
				// Ambil skor dari teks di dalam span
				result.Skor = strings.TrimSpace(el.Text)
			}
		})

		// Ekstrak sinopsis
		result.Sinopsis = strings.TrimSpace(e.ChildText("div.ingdesc .contexcerpt"))

		// Ekstrak info dari bagian bawah (Genres, Status, Studio)
		e.ForEach(".linginfo span", func(_ int, s *colly.HTMLElement) {
			text := s.Text
			if strings.HasPrefix(text, "Genres:") {
				s.ForEach("a", func(_ int, a *colly.HTMLElement) {
					result.Genres = append(result.Genres, a.Text)
				})
			} else if strings.HasPrefix(text, "Studio:") {
				result.Studio = strings.TrimSpace(s.ChildText("a"))
			} else if strings.HasPrefix(text, "Status:") {
				// Status kadang ada di sini, kita ambil juga untuk melengkapi
				if result.StatusRilis == "Upcoming" || result.StatusRilis == "" {
					result.StatusRilis = strings.TrimSpace(strings.TrimPrefix(text, "Status:"))
				}
			}
		})

		searchResults = append(searchResults, result)
	})

	// Callback utama untuk hasil pencarian
	c.OnHTML("div.listupd article.bs", func(e *colly.HTMLElement) {
		linkElement := e.DOM.Find("a.tip")
		postID, exists := linkElement.Attr("rel")
		if !exists {
			return
		}

		thumbURL := e.ChildAttr("img", "src")
		if thumbURL == "" {
			thumbURL = e.ChildAttr("img", "data-src")
		}

		result := AnimeSearchResult{
			Judul:       linkElement.AttrOr("title", ""),
			Tautan:      e.Request.AbsoluteURL(linkElement.AttrOr("href", "")),
			Thumbnail:   thumbURL,
			StatusRilis: e.ChildText("span.epx"),
			Tipe:        e.ChildText("div.typez"),
			Genres:      []string{},
		}

		ctx := colly.NewContext()
		ctx.Put("result", result)

		payload := fmt.Sprintf("action=tooltip_action&id=%s", postID)
		err := c.Request("POST", ajaxURL, strings.NewReader(payload), ctx, nil)
		if err != nil {
			t.Logf("Gagal membuat request detail untuk ID %s: %v", postID, err)
		}
	})

	c.OnRequest(func(r *colly.Request) { t.Logf("Mengunjungi: %s", r.URL.String()) })
	c.OnError(func(r *colly.Response, err error) { t.Errorf("Gagal: %s, Error: %v", r.Request.URL, err) })

	err := c.Visit(searchURL)
	if err != nil {
		t.Fatalf("Gagal memulai scraper: %v", err)
	}
	c.Wait()

	jsonData, err := json.MarshalIndent(searchResults, "", "  ")
	if err != nil {
		t.Fatalf("Gagal mengkonversi ke JSON: %v", err)
	}

	t.Logf("\n%s\n", string(jsonData))
}

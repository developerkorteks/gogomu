package gomunime

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// PERUBAHAN: Struct baru untuk menampung nama dan URL genre.
type GenreItem struct {
	Nama string `json:"nama"`
	URL  string `json:"url"`
}

// AnimeInfo menyimpan detail untuk satu anime.
type AnimeInfo struct {
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

// GenrePageResult adalah struktur hasil akhir.
type GenrePageResult struct {
	Genre            string      `json:"genre"`
	DaftarSemuaGenre []GenreItem `json:"daftar_semua_genre"` // PERUBAHAN: Menggunakan struct GenreItem
	AnimeList        []AnimeInfo `json:"anime_list"`
}

func TestScrapeGenrePage(t *testing.T) {
	genreURL := "https://gomunime.co/genres/action/"
	ajaxURL := "https://gomunime.co/wp-admin/admin-ajax.php"

	result := GenrePageResult{
		Genre:            "Action",
		DaftarSemuaGenre: []GenreItem{}, // PERUBAHAN: Inisialisasi sebagai slice of GenreItem
		AnimeList:        []AnimeInfo{},
	}

	var genreListOnce sync.Once
	var animeMutex sync.Mutex

	c := colly.NewCollector(
		colly.Async(true),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
	c.Limit(&colly.LimitRule{DomainGlob: "*gomunime.co*", Parallelism: 8})

	// Callback untuk daftar semua genre dari sidebar
	c.OnHTML("div#sidebar ul.genre li a", func(e *colly.HTMLElement) {
		genreListOnce.Do(func() {
			// PERUBAHAN: Logika untuk mengambil nama dan URL genre
			e.DOM.Closest("ul.genre").Find("li a").Each(func(_ int, s *goquery.Selection) {
				genreNama := s.Text()
				genreLink, _ := s.Attr("href")

				item := GenreItem{
					Nama: genreNama,
					URL:  e.Request.AbsoluteURL(genreLink), // Membuat URL menjadi absolut
				}
				result.DaftarSemuaGenre = append(result.DaftarSemuaGenre, item)
			})
		})
	})

	// Callback untuk detail dari AJAX (info hover)
	c.OnHTML("div.ingfo", func(e *colly.HTMLElement) {
		anime := e.Request.Ctx.GetAny("anime").(AnimeInfo)
		e.ForEach("div.minginfo span.l", func(_ int, el *colly.HTMLElement) {
			text := strings.TrimSpace(el.Text)
			if strings.Contains(text, "min. per ep.") {
				anime.Durasi = text
			} else if el.ChildText("i.fa-star") != "" {
				anime.Skor = strings.TrimSpace(el.Text)
			}
		})
		anime.Sinopsis = strings.TrimSpace(e.ChildText("div.ingdesc .contexcerpt"))
		anime.Genres = []string{}
		e.ForEach(".linginfo span", func(_ int, s *colly.HTMLElement) {
			text := s.Text
			if strings.HasPrefix(text, "Genres:") {
				s.ForEach("a", func(_ int, a *colly.HTMLElement) {
					anime.Genres = append(anime.Genres, a.Text)
				})
			} else if strings.HasPrefix(text, "Studio:") {
				anime.Studio = strings.TrimSpace(s.ChildText("a"))
			}
		})
		animeMutex.Lock()
		result.AnimeList = append(result.AnimeList, anime)
		animeMutex.Unlock()
	})

	// Callback untuk setiap item anime di halaman
	c.OnHTML("div.listupd article.bs", func(e *colly.HTMLElement) {
		linkElement := e.DOM.Find("a.tip")
		postID, exists := linkElement.Attr("rel")
		if !exists {
			return
		}

		thumbURL := e.ChildAttr("img", "src")
		if strings.HasPrefix(thumbURL, "data:image") {
			thumbURL = e.ChildAttr("img", "data-src")
		}

		anime := AnimeInfo{
			Judul:       linkElement.AttrOr("title", ""),
			Tautan:      e.Request.AbsoluteURL(linkElement.AttrOr("href", "")),
			Thumbnail:   thumbURL,
			StatusRilis: e.ChildText("span.epx"),
			Tipe:        e.ChildText("div.typez"),
			Genres:      []string{},
		}

		ctx := colly.NewContext()
		ctx.Put("anime", anime)

		payload := fmt.Sprintf("action=tooltip_action&id=%s", postID)
		if err := c.Request("POST", ajaxURL, strings.NewReader(payload), ctx, nil); err != nil {
			t.Logf("Gagal membuat request detail untuk ID %s: %v", postID, err)
		}
	})

	// Callback untuk paginasi
	c.OnHTML("div.pagination a.next", func(e *colly.HTMLElement) {
		nextPageURL := e.Attr("href")
		t.Logf("Menuju halaman berikutnya: %s", nextPageURL)
		if err := c.Visit(nextPageURL); err != nil {
			t.Logf("Gagal mengunjungi halaman berikutnya %s: %v", nextPageURL, err)
		}
	})

	c.OnRequest(func(r *colly.Request) { t.Logf("Mengunjungi: %s", r.URL.String()) })
	c.OnError(func(r *colly.Response, err error) { t.Errorf("Gagal: %s, Error: %v", r.Request.URL, err) })

	if err := c.Visit(genreURL); err != nil {
		t.Fatalf("Gagal memulai scraper: %v", err)
	}
	c.Wait()

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("Gagal mengkonversi ke JSON: %v", err)
	}

	t.Logf("\nBerhasil meng-scrape %d anime dari genre '%s'.\n", len(result.AnimeList), result.Genre)
	t.Logf("\n%s\n", string(jsonData))
}

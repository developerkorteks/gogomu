package gomunime

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gocolly/colly/v2"
)

// [IMPROVED] Struct disesuaikan untuk menampung semua kategori
type PopularAnimePopuler struct {
	Peringkat int      `json:"peringkat"`
	Judul     string   `json:"judul"`
	Tautan    string   `json:"tautan"`
	Thumbnail string   `json:"thumbnail"`
	Genres    []string `json:"genres"`
	Rating    string   `json:"rating,omitempty"` // omitempty jika rating tidak ada
}

// [IMPROVED] Struct utama untuk menyimpan hasil dari semua tab
type AllPopularAnime struct {
	Weekly  []PopularAnimePopuler `json:"weekly"`
	Monthly []PopularAnimePopuler `json:"monthly"`
	AllTime []PopularAnimePopuler `json:"all_time"`
}

func TestPopularAnime(t *testing.T) {
	pageURL := "https://gomunime.co/"

	allData := AllPopularAnime{
		Weekly:  make([]PopularAnimePopuler, 0),
		Monthly: make([]PopularAnimePopuler, 0),
		AllTime: make([]PopularAnimePopuler, 0),
	}

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"

	// [IMPROVED] Callback dibuat untuk bisa digunakan oleh ketiga kategori
	scrapePopularList := func(e *colly.HTMLElement) []PopularAnimePopuler {
		var animeList []PopularAnimePopuler
		e.ForEach("ul li", func(_ int, el *colly.HTMLElement) {
			peringkatStr := el.ChildText(".ctr")
			peringkat, _ := strconv.Atoi(peringkatStr)

			var genres []string
			// Membersihkan teks "Genres: " dan memisahkan genre
			genreText := el.ChildText(".leftseries span")
			genreText = strings.TrimPrefix(genreText, "Genres: ")
			genres = strings.Split(genreText, ", ")

			thumbURL := el.ChildAttr(".imgseries img", "data-src")
			if thumbURL == "" {
				thumbURL = el.ChildAttr(".imgseries img", "src")
			}

			anime := PopularAnimePopuler{
				Peringkat: peringkat,
				Judul:     el.ChildText(".leftseries h4 a"),
				Tautan:    el.Request.AbsoluteURL(el.ChildAttr(".leftseries h4 a", "href")),
				Thumbnail: thumbURL,
				Genres:    genres,
				Rating:    el.ChildText(".numscore"),
			}
			animeList = append(animeList, anime)
		})
		return animeList
	}

	// Menjalankan scraper untuk setiap kategori
	c.OnHTML("div.serieslist.pop.wpop-weekly", func(e *colly.HTMLElement) {
		log.Println("Mengambil data populer mingguan...")
		allData.Weekly = scrapePopularList(e)
	})

	c.OnHTML("div.serieslist.pop.wpop-monthly", func(e *colly.HTMLElement) {
		log.Println("Mengambil data populer bulanan...")
		allData.Monthly = scrapePopularList(e)
	})

	c.OnHTML("div.serieslist.pop.wpop-alltime", func(e *colly.HTMLElement) {
		log.Println("Mengambil data populer sepanjang waktu...")
		allData.AllTime = scrapePopularList(e)
	})

	c.OnRequest(func(r *colly.Request) { log.Printf("Mengunjungi: %s", r.URL.String()) })
	c.OnError(func(r *colly.Response, err error) { t.Errorf("Gagal: %s, Error: %v", r.Request.URL.String(), err) })

	err := c.Visit(pageURL)
	if err != nil {
		t.Fatalf("Gagal memulai scraper: %v", err)
	}
	c.Wait()

	if len(allData.Weekly) == 0 && len(allData.Monthly) == 0 && len(allData.AllTime) == 0 {
		t.Fatal("Tidak ada data anime populer yang berhasil di-scrape.")
	}

	jsonData, err := json.MarshalIndent(allData, "", "  ")
	if err != nil {
		t.Fatalf("Gagal mengkonversi ke JSON: %v", err)
	}

	// Menyimpan hasil ke file
	file, _ := os.Create("popular_anime_lengkap.json")
	defer file.Close()
	file.Write(jsonData)
	fmt.Println(string(jsonData))
	log.Println("Data berhasil disimpan ke popular_anime_lengkap.json")
}

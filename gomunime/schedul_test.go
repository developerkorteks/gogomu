package gomunime

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gocolly/colly/v2"
)

// AnimeSchedule merepresentasikan satu entri anime dalam jadwal.
type AnimeSchedule struct {
	Judul      string `json:"judul"`
	Tautan     string `json:"tautan"`
	WaktuRilis string `json:"waktu_rilis"`
	Thumbnail  string `json:"thumbnail"`
}

// DaySchedule merepresentasikan jadwal untuk satu hari.
type DaySchedule struct {
	Hari      string          `json:"hari"`
	AnimeList []AnimeSchedule `json:"anime_list"`
}

func TestSchedule(t *testing.T) {
	scheduleURL := "https://gomunime.co/schedule/"

	var fullSchedule []DaySchedule

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"

	c.OnHTML("div.bixbox.schedulepage", func(e *colly.HTMLElement) {
		day := DaySchedule{
			Hari:      e.ChildText("div.releases h3 span"),
			AnimeList: []AnimeSchedule{},
		}

		e.ForEach("div.bs", func(_ int, animeEl *colly.HTMLElement) {
			// [FIXED] Logika baru untuk mendapatkan URL thumbnail yang benar
			// Coba ambil dari 'data-src' (untuk lazy loading), jika kosong, baru ambil dari 'src'.
			thumbURL := animeEl.ChildAttr("img", "data-src")
			if thumbURL == "" {
				thumbURL = animeEl.ChildAttr("img", "src")
			}

			anime := AnimeSchedule{
				Judul:      animeEl.ChildAttr("a", "title"),
				Tautan:     e.Request.AbsoluteURL(animeEl.ChildAttr("a", "href")),
				WaktuRilis: animeEl.ChildText("span.epx"),
				Thumbnail:  thumbURL, // Menggunakan URL yang sudah benar
			}
			day.AnimeList = append(day.AnimeList, anime)
		})

		if len(day.AnimeList) > 0 {
			fullSchedule = append(fullSchedule, day)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Mengunjungi:", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Gagal mengunjungi:", r.Request.URL, "Error:", err)
	})

	err := c.Visit(scheduleURL)
	if err != nil {
		log.Fatal("Gagal memulai scraper:", err)
	}

	c.Wait()

	jsonData, err := json.MarshalIndent(fullSchedule, "", "  ")
	if err != nil {
		log.Fatal("Gagal mengkonversi data ke JSON:", err)
	}

	fmt.Println(string(jsonData))

	file, err := os.Create("schedule.json")
	if err != nil {
		log.Fatal("Gagal membuat file:", err)
	}
	defer file.Close()
	file.Write(jsonData)

	log.Println("Data jadwal berhasil disimpan ke schedule.json")
}

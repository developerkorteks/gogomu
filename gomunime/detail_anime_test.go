package gomunime

import (
	"encoding/json"
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// --- Structs untuk Data yang Lebih Lengkap ---

type Linkable struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EpisodeDetailAnime struct {
	Episode      string `json:"episode"`
	JudulEpisode string `json:"judul_episode"`
	URL          string `json:"url"`
	TanggalRilis string `json:"tanggal_rilis"`
}

type RecommendationDetailAnime struct {
	Judul     string `json:"judul"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
	Status    string `json:"status"`
}

type AnimeDetailsDetailAnime struct {
	Judul           string                      `json:"judul"`
	JudulAlternatif string                      `json:"judul_alternatif"`
	Thumbnail       string                      `json:"thumbnail"`
	Rating          string                      `json:"rating"`
	Sinopsis        string                      `json:"sinopsis"`
	Status          string                      `json:"status"`
	Tipe            string                      `json:"tipe"`
	TotalEpisode    string                      `json:"total_episode"`
	Durasi          string                      `json:"durasi"`
	RilisPerdana    string                      `json:"rilis_perdana"`
	DipostingPada   string                      `json:"diposting_pada"`
	DiperbaruiPada  string                      `json:"diperbarui_pada"`
	Studio          Linkable                    `json:"studio"`
	Season          Linkable                    `json:"season"`
	Produser        []Linkable                  `json:"produser"`
	Genre           []Linkable                  `json:"genre"`
	PengisiSuara    []Linkable                  `json:"pengisi_suara"`
	EpisodeList     []EpisodeDetailAnime        `json:"episode_list"`
	Rekomendasi     []RecommendationDetailAnime `json:"rekomendasi"`
}

func TestAnimeDetail(t *testing.T) {
	animeURL := "https://gomunime.co/anime/isekai-shikkaku/"

	animeData := AnimeDetailsDetailAnime{}

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"

	// --- Scraper Utama untuk Informasi Detail ---
	c.OnHTML("article.post-180", func(e *colly.HTMLElement) {
		// Mengambil informasi dari kotak utama
		bigContent := e.DOM.Find(".bigcontent")
		animeData.Judul = bigContent.Find("h1.entry-title").Text()

		// [FIXED] Mengambil thumbnail dengan mempertimbangkan lazy-loading.
		// Scraper akan mencari atribut 'data-src' terlebih dahulu, jika tidak ada, baru mencari 'src'.
		imgSelector := ".thumbook .thumb img"
		thumbnailURL := bigContent.Find(imgSelector).AttrOr("data-src", "")
		if thumbnailURL == "" {
			thumbnailURL = bigContent.Find(imgSelector).AttrOr("src", "")
		}
		animeData.Thumbnail = thumbnailURL

		// ... sisa kode Anda untuk mengambil rating, status, dll. tetap sama ...
		ratingText := bigContent.Find(".rating strong").Text()
		animeData.Rating = strings.TrimSpace(strings.ReplaceAll(ratingText, "Rating", ""))

		// Mengambil semua detail dari tabel '.spe'
		spe := bigContent.Find(".spe")
		spe.Find("span").Each(func(_ int, s *goquery.Selection) {
			text := s.Text()
			switch {
			case strings.HasPrefix(text, "Status:"):
				animeData.Status = strings.TrimSpace(strings.TrimPrefix(text, "Status:"))
			case strings.HasPrefix(text, "Studio:"):
				animeData.Studio = Linkable{Name: s.Find("a").Text(), URL: s.Find("a").AttrOr("href", "")}
			case strings.HasPrefix(text, "Released:"):
				animeData.RilisPerdana = strings.TrimSpace(strings.TrimPrefix(text, "Released:"))
			case strings.HasPrefix(text, "Duration:"):
				animeData.Durasi = strings.TrimSpace(strings.TrimPrefix(text, "Duration:"))
			case strings.HasPrefix(text, "Season:"):
				animeData.Season = Linkable{Name: s.Find("a").Text(), URL: s.Find("a").AttrOr("href", "")}
			case strings.HasPrefix(text, "Type:"):
				animeData.Tipe = strings.TrimSpace(strings.TrimPrefix(text, "Type:"))
			case strings.HasPrefix(text, "Episodes:"):
				animeData.TotalEpisode = strings.TrimSpace(strings.TrimPrefix(text, "Episodes:"))
			case strings.HasPrefix(text, "Released on:"):
				animeData.DipostingPada = s.Find("time").Text()
			case strings.HasPrefix(text, "Updated on:"):
				animeData.DiperbaruiPada = s.Find("time").Text()
			case strings.HasPrefix(text, "Producers:"):
				s.Find("a").Each(func(_ int, p *goquery.Selection) {
					animeData.Produser = append(animeData.Produser, Linkable{Name: p.Text(), URL: p.AttrOr("href", "")})
				})
			case strings.HasPrefix(text, "Casts:"):
				s.Find("a").Each(func(_ int, c *goquery.Selection) {
					animeData.PengisiSuara = append(animeData.PengisiSuara, Linkable{Name: c.Text(), URL: c.AttrOr("href", "")})
				})
			}
		})

		// Mengambil genre
		bigContent.Find(".genxed a").Each(func(_ int, g *goquery.Selection) {
			animeData.Genre = append(animeData.Genre, Linkable{Name: g.Text(), URL: g.AttrOr("href", "")})
		})

		// Mengambil Sinopsis
		animeData.Sinopsis = e.DOM.Find(".bixbox.synp .entry-content").Text()
	})

	// --- Scraper untuk Daftar Episode ---
	c.OnHTML("div.eplister", func(e *colly.HTMLElement) {
		e.ForEach("ul li", func(_ int, el *colly.HTMLElement) {
			ep := EpisodeDetailAnime{
				Episode:      el.ChildText(".epl-num"),
				JudulEpisode: el.ChildText(".epl-title"),
				URL:          e.Request.AbsoluteURL(el.ChildAttr("a", "href")),
				TanggalRilis: el.ChildText(".epl-date"),
			}
			animeData.EpisodeList = append(animeData.EpisodeList, ep)
		})
	})

	// --- Scraper untuk Rekomendasi Anime ---
	c.OnHTML(`div.bixbox:has(h3 > span:contains("Recommended Series"))`, func(e *colly.HTMLElement) {
		e.ForEach(".listupd article.bs", func(_ int, el *colly.HTMLElement) {
			// Mengatasi lazy-loading untuk thumbnail
			thumbURL := el.ChildAttr("img", "data-src")
			if thumbURL == "" {
				thumbURL = el.ChildAttr("img", "src")
			}

			rec := RecommendationDetailAnime{
				Judul:     el.ChildAttr("a.tip", "title"),
				URL:       e.Request.AbsoluteURL(el.ChildAttr("a.tip", "href")),
				Status:    el.ChildText("span.epx"),
				Thumbnail: thumbURL,
			}
			animeData.Rekomendasi = append(animeData.Rekomendasi, rec)
		})
	})
	c.OnRequest(func(r *colly.Request) { log.Printf("Mengunjungi: %s", r.URL.String()) })
	c.OnError(func(r *colly.Response, err error) { t.Fatalf("Gagal: %s, Error: %v", r.Request.URL.String(), err) })

	err := c.Visit(animeURL)
	if err != nil {
		t.Fatalf("Gagal memulai scraper: %v", err)
	}
	c.Wait()

	jsonData, err := json.MarshalIndent(animeData, "", "  ")
	if err != nil {
		t.Fatalf("Gagal mengkonversi ke JSON: %v", err)
	}

	t.Logf("\n%s\n", string(jsonData))
}

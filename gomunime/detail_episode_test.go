package gomunime

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func TestEpisodeDetail(t *testing.T) {
	animeURL := "https://gomunime.co/arknights-rise-from-ember-4/"

	animeData := AnimeDetails{
		Pemeran:       []string{},
		Genre:         []string{},
		EpisodeList:   []Episode{},
		Rekomendasi:   []Recommendation{},
		DownloadLinks: []DownloadLink{},
		MirrorStreams: make(map[string]string), // Inisialisasi map
	}

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"

	// Callback untuk info utama
	c.OnHTML("div.single-info.bixbox", func(e *colly.HTMLElement) {
		// ... (kode di bagian ini tidak berubah)
		thumbURL := e.ChildAttr("div.thumb img", "data-src")
		if thumbURL == "" {
			thumbURL = e.ChildAttr("div.thumb img", "src")
		}
		animeData.Thumbnail = thumbURL
		infox := e.DOM.Find("div.infox")
		animeData.Judul = infox.Find("h2").Text()
		animeData.JudulAlternatif = infox.Find("span.alter").Text()
		animeData.Rating = infox.Find("div.rating strong").Text()
		animeData.Sinopsis = strings.TrimSpace(infox.Find("div.desc").Contents().Not("h2").Text())
		infox.Find("div.spe span").Each(func(_ int, s *goquery.Selection) {
			text := s.Text()
			switch {
			case strings.Contains(text, "Status:"):
				animeData.Status = strings.TrimSpace(strings.ReplaceAll(text, "Status:", ""))
			case strings.Contains(text, "Studio:"):
				animeData.Studio = s.Find("a").Text()
			case strings.Contains(text, "Released:"):
				animeData.RilisPerdana = strings.TrimSpace(strings.ReplaceAll(text, "Released:", ""))
			case strings.Contains(text, "Season:"):
				animeData.Season = s.Find("a").Text()
			case strings.Contains(text, "Type:"):
				animeData.Tipe = strings.TrimSpace(strings.ReplaceAll(text, "Type:", ""))
			case strings.Contains(text, "Episodes:"):
				animeData.TotalEpisode = strings.TrimSpace(strings.ReplaceAll(text, "Episodes:", ""))
			case strings.Contains(text, "Director:"):
				animeData.Sutradara = s.Find("a").Text()
			case strings.Contains(text, "Casts:"):
				s.Find("a").Each(func(_ int, p *goquery.Selection) {
					animeData.Pemeran = append(animeData.Pemeran, p.Text())
				})
			}
		})
		infox.Find("div.genxed a").Each(func(_ int, g *goquery.Selection) {
			animeData.Genre = append(animeData.Genre, g.Text())
		})
	})

	// Callback untuk URL stream utama (sebagai fallback)
	c.OnHTML("div#pembed iframe", func(e *colly.HTMLElement) {
		if e.Attr("src") != "about:blank" && animeData.StreamURL == "" {
			animeData.StreamURL = e.Attr("src")
		}
	})

	// [FIXED] Callback untuk mengambil dan men-decode semua mirror stream
	c.OnHTML("select.mirror", func(e *colly.HTMLElement) {
		// Regex untuk mengekstrak URL dari tag <iframe>
		re := regexp.MustCompile(`src="([^"]+)"`)

		e.ForEach("option", func(_ int, el *colly.HTMLElement) {
			base64Value := el.Attr("value")
			if base64Value == "" {
				return
			}

			// Decode Base64
			decodedIframe, err := base64.StdEncoding.DecodeString(base64Value)
			if err != nil {
				t.Logf("Gagal decode base64: %v", err)
				return
			}

			// Ekstrak URL src dari hasil decode
			matches := re.FindStringSubmatch(string(decodedIframe))
			if len(matches) > 1 {
				serverName := strings.TrimSpace(el.Text)
				streamURL := matches[1]
				animeData.MirrorStreams[serverName] = streamURL

				// Set URL stream utama dari mirror pertama yang valid
				if animeData.StreamURL == "" {
					animeData.StreamURL = streamURL
				}
			}
		})
	})

	// Callback untuk link download
	c.OnHTML("div.soraddlx", func(e *colly.HTMLElement) {
		e.ForEach(".soraurlx", func(_ int, el *colly.HTMLElement) {
			resolusi := el.ChildText("strong")
			el.ForEach("a", func(_ int, linkEl *colly.HTMLElement) {
				dl := DownloadLink{
					Resolusi: resolusi,
					Provider: linkEl.Text,
					Tautan:   linkEl.Attr("href"),
				}
				animeData.DownloadLinks = append(animeData.DownloadLinks, dl)
			})
		})
	})

	// Callback lainnya tetap sama
	c.OnHTML("#mainepisode", func(e *colly.HTMLElement) {
		e.ForEach(".episodelist ul li", func(_ int, el *colly.HTMLElement) {
			episodeNumText := strings.TrimSpace(el.ChildText(".playinfo span:nth-child(1)"))
			episodeNum := strings.TrimPrefix(episodeNumText, "Eps ")
			ep := Episode{
				Episode:      episodeNum,
				JudulEpisode: el.ChildText(".playinfo h3"),
				Tautan:       e.Request.AbsoluteURL(el.ChildAttr("a", "href")),
				TanggalRilis: el.ChildText(".playinfo span:nth-child(2)"),
			}
			animeData.EpisodeList = append(animeData.EpisodeList, ep)
		})
	})
	c.OnHTML("div.bixbox:has(h3:contains(Recommended))", func(e *colly.HTMLElement) {
		e.ForEach(".listupd article.bs", func(_ int, el *colly.HTMLElement) {
			thumbURL := el.ChildAttr("img", "data-src")
			if thumbURL == "" {
				thumbURL = el.ChildAttr("img", "src")
			}
			rec := Recommendation{
				Judul:     el.ChildAttr("a.tip", "title"),
				Tautan:    e.Request.AbsoluteURL(el.ChildAttr("a.tip", "href")),
				Status:    el.ChildText("span.epx"),
				Thumbnail: thumbURL,
			}
			animeData.Rekomendasi = append(animeData.Rekomendasi, rec)
		})
	})

	c.OnRequest(func(r *colly.Request) { t.Logf("Mengunjungi: %s", r.URL.String()) })
	c.OnError(func(r *colly.Response, err error) { t.Errorf("Gagal: %s, Error: %v", r.Request.URL, err) })

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

	if animeData.Judul == "" {
		t.Error("Judul anime kosong, scraping mungkin gagal.")
	}
}

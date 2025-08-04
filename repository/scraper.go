package repository

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const (
	baseURL     = "https://gomunime.co/"
	ajaxURL     = "https://gomunime.co/wp-admin/admin-ajax.php"
	scheduleURL = "https://gomunime.co/schedule/"
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
)

// Optimized collector configuration
func createOptimizedCollector(async bool, parallelism int) *colly.Collector {
	var c *colly.Collector
	if async {
		c = colly.NewCollector(colly.Async(true))
		c.Limit(&colly.LimitRule{
			DomainGlob:  "*gomunime.co*",
			Parallelism: parallelism,
			Delay:       100 * time.Millisecond, // Reduced delay
		})
	} else {
		c = colly.NewCollector()
	}

	c.UserAgent = userAgent
	// Set timeout to prevent hanging
	c.SetRequestTimeout(30 * time.Second)

	// Enable caching for repeated requests
	c.CacheDir = "./cache"

	return c
}

// ScrapeLatestAnime mengambil daftar anime yang baru diperbarui dengan optimasi.
func ScrapeLatestAnime() []ScrapedLatestAnime {
	var allAnime []ScrapedLatestAnime
	var mu sync.Mutex

	c := createOptimizedCollector(true, 12) // Increased parallelism

	// Use channels for better coordination
	detailCh := make(chan ScrapedLatestAnime, 100)
	done := make(chan bool)

	// Goroutine to collect results
	go func() {
		for anime := range detailCh {
			mu.Lock()
			allAnime = append(allAnime, anime)
			mu.Unlock()
		}
		done <- true
	}()

	// Optimized detail callback
	c.OnHTML("div.ingfo", func(e *colly.HTMLElement) {
		animeInfo := e.Request.Ctx.GetAny("anime").(ScrapedLatestAnime)

		// Batch process all information at once
		e.ForEach(".minginfo span.l", func(_ int, s *colly.HTMLElement) {
			if s.ChildText("i.fa-star") != "" {
				animeInfo.Rating = strings.TrimSpace(s.Text)
			}
		})

		// Use a single selector for genres
		genreLinks := e.DOM.Find(".linginfo span a")
		if genreLinks.Length() > 0 {
			genreLinks.Each(func(_ int, a *goquery.Selection) {
				animeInfo.Genres = append(animeInfo.Genres, a.Text())
			})
		}

		detailCh <- animeInfo
	})

	// Optimized main card callback
	c.OnHTML("div.listupd article.bs", func(e *colly.HTMLElement) {
		postID, exists := e.DOM.Find("a.tip").Attr("rel")
		if !exists {
			return
		}

		// Optimize thumbnail selection
		thumbURL := e.ChildAttr("img", "data-src")
		if thumbURL == "" {
			thumbURL = e.ChildAttr("img", "src")
		}

		anime := ScrapedLatestAnime{
			Judul:     e.ChildAttr("a.tip", "title"),
			Tautan:    e.Request.AbsoluteURL(e.ChildAttr("a.tip", "href")),
			Episode:   strings.ReplaceAll(e.ChildText("span.epx"), " Episode", ""),
			Thumbnail: thumbURL,
			Tipe:      e.ChildText(".typez"),
			Genres:    make([]string, 0, 5), // Pre-allocate capacity
		}

		ctx := colly.NewContext()
		ctx.Put("anime", anime)
		formData := fmt.Sprintf("action=tooltip_action&id=%s", postID)

		// Make AJAX request with error handling
		if err := c.Request("POST", ajaxURL, strings.NewReader(formData), ctx, nil); err != nil {
			log.Printf("Failed AJAX request for ID %s: %v", postID, err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping latest: %v | URL: %s", err, r.Request.URL)
	})

	c.Visit(baseURL)
	c.Wait()

	close(detailCh)
	<-done

	return allAnime
}

// ScrapeSchedule dengan optimasi minimal karena sudah cukup efisien.
func ScrapeSchedule() []ScrapedDaySchedule {
	var fullSchedule []ScrapedDaySchedule
	c := createOptimizedCollector(false, 1)

	c.OnHTML("div.bixbox.schedulepage", func(e *colly.HTMLElement) {
		day := ScrapedDaySchedule{
			Hari:      e.ChildText("div.releases h3 span"),
			AnimeList: make([]ScrapedAnimeSchedule, 0, 20), // Pre-allocate
		}

		e.ForEach("div.bs", func(_ int, animeEl *colly.HTMLElement) {
			thumbURL := animeEl.ChildAttr("img", "data-src")
			if thumbURL == "" {
				thumbURL = animeEl.ChildAttr("img", "src")
			}

			anime := ScrapedAnimeSchedule{
				Judul:      animeEl.ChildAttr("a", "title"),
				Tautan:     e.Request.AbsoluteURL(animeEl.ChildAttr("a", "href")),
				WaktuRilis: animeEl.ChildText("span.epx"),
				Thumbnail:  thumbURL,
			}
			day.AnimeList = append(day.AnimeList, anime)
		})

		if len(day.AnimeList) > 0 {
			fullSchedule = append(fullSchedule, day)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping schedule: %v | URL: %s", err, r.Request.URL)
	})

	c.Visit(scheduleURL)
	c.Wait()
	return fullSchedule
}

// ScrapeLatestByPage dengan optimasi parallelism yang lebih baik.
func ScrapeLatestByPage(page int) []ScrapedLatestAnime {
	var allAnime []ScrapedLatestAnime
	var mu sync.Mutex

	c := createOptimizedCollector(true, 15) // Increased parallelism

	// Use buffered channel for better performance
	resultCh := make(chan ScrapedLatestAnime, 50)
	done := make(chan bool)

	go func() {
		for anime := range resultCh {
			mu.Lock()
			allAnime = append(allAnime, anime)
			mu.Unlock()
		}
		done <- true
	}()

	c.OnHTML("div.ingfo", func(e *colly.HTMLElement) {
		animeInfo := e.Request.Ctx.GetAny("anime").(ScrapedLatestAnime)

		// Batch process all selectors
		e.ForEach(".minginfo span.l", func(_ int, s *colly.HTMLElement) {
			if s.ChildText("i.fa-star") != "" {
				animeInfo.Rating = strings.TrimSpace(s.Text)
			}
		})

		animeInfo.Deskripsi = strings.TrimSpace(e.ChildText(".contexcerpt"))

		// Optimize genre and status extraction
		e.ForEach(".linginfo span", func(_ int, s *colly.HTMLElement) {
			text := s.Text
			if strings.HasPrefix(text, "Status:") {
				animeInfo.Status = strings.TrimSpace(strings.TrimPrefix(text, "Status:"))
			} else if strings.HasPrefix(text, "Genres:") {
				s.ForEach("a", func(_ int, a *colly.HTMLElement) {
					animeInfo.Genres = append(animeInfo.Genres, a.Text)
				})
			}
		})

		resultCh <- animeInfo
	})

	c.OnHTML("div.listupd article.bs", func(e *colly.HTMLElement) {
		postID, exists := e.DOM.Find("a.tip").Attr("rel")
		if !exists {
			return
		}

		thumbURL := e.ChildAttr("img", "data-src")
		if thumbURL == "" {
			thumbURL = e.ChildAttr("img", "src")
		}

		anime := ScrapedLatestAnime{
			Judul:     e.ChildAttr("a.tip", "title"),
			Tautan:    e.Request.AbsoluteURL(e.ChildAttr("a.tip", "href")),
			Episode:   strings.ReplaceAll(e.ChildText("span.epx"), " Episode", ""),
			Thumbnail: thumbURL,
			Tipe:      e.ChildText(".typez"),
			Genres:    make([]string, 0, 5),
		}

		ctx := colly.NewContext()
		ctx.Put("anime", anime)
		formData := fmt.Sprintf("action=tooltip_action&id=%s", postID)

		if err := c.Request("POST", ajaxURL, strings.NewReader(formData), ctx, nil); err != nil {
			log.Printf("Failed AJAX request: %v", err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping page: %v | URL: %s", err, r.Request.URL)
	})

	targetURL := baseURL
	if page > 1 {
		targetURL = fmt.Sprintf("%spage/%d/", baseURL, page)
	}

	log.Printf("Visiting page: %s", targetURL)
	c.Visit(targetURL)
	c.Wait()

	close(resultCh)
	<-done

	return allAnime
}

// ScrapeAnimeDetail dengan optimasi selector dan pre-allocation.
func ScrapeAnimeDetail(animeSlug string) ScrapedAnimeDetails {
	targetURL := fmt.Sprintf("https://gomunime.co/anime/%s/", animeSlug)
	animeData := ScrapedAnimeDetails{
		Details:     make(map[string]string, 10), // Pre-allocate capacity
		Genre:       make([]string, 0, 10),
		EpisodeList: make([]ScrapedEpisode, 0, 50),
		Rekomendasi: make([]ScrapedRecommendation, 0, 20),
	}

	c := createOptimizedCollector(false, 1)

	// Main info scraper for article.post-180 structure (detailed anime pages)
	c.OnHTML("article.post-180", func(e *colly.HTMLElement) {
		// Get bigcontent container
		bigContent := e.DOM.Find(".bigcontent")
		
		// Extract title
		animeData.Judul = bigContent.Find("h1.entry-title").Text()

		// Extract thumbnail with lazy-loading support
		imgSelector := ".thumbook .thumb img"
		thumbnailURL := bigContent.Find(imgSelector).AttrOr("data-src", "")
		if thumbnailURL == "" {
			thumbnailURL = bigContent.Find(imgSelector).AttrOr("src", "")
		}
		animeData.Thumbnail = thumbnailURL

		// Extract rating/score
		ratingText := bigContent.Find(".rating strong").Text()
		animeData.Skor = strings.TrimSpace(strings.ReplaceAll(ratingText, "Rating", ""))

		// Extract synopsis from the correct selector
		animeData.Sinopsis = strings.TrimSpace(e.DOM.Find(".bixbox.synp .entry-content").Text())

		// Extract genres
		bigContent.Find(".genxed a").Each(func(_ int, g *goquery.Selection) {
			animeData.Genre = append(animeData.Genre, g.Text())
		})

		// Extract details from .spe table
		spe := bigContent.Find(".spe")
		spe.Find("span").Each(func(_ int, s *goquery.Selection) {
			text := s.Text()
			switch {
			case strings.HasPrefix(text, "Status:"):
				animeData.Details["Status"] = strings.TrimSpace(strings.TrimPrefix(text, "Status:"))
			case strings.HasPrefix(text, "Studio:"):
				studioName := s.Find("a").Text()
				if studioName != "" {
					animeData.Details["Studio"] = studioName
				}
			case strings.HasPrefix(text, "Released:"):
				animeData.Details["Released:"] = strings.TrimSpace(strings.TrimPrefix(text, "Released:"))
			case strings.HasPrefix(text, "Duration:"):
				animeData.Details["Duration"] = strings.TrimSpace(strings.TrimPrefix(text, "Duration:"))
			case strings.HasPrefix(text, "Season:"):
				seasonName := s.Find("a").Text()
				if seasonName != "" {
					animeData.Details["Season"] = seasonName
				}
			case strings.HasPrefix(text, "Type:"):
				animeData.Details["Type"] = strings.TrimSpace(strings.TrimPrefix(text, "Type:"))
			case strings.HasPrefix(text, "Episodes:"):
				animeData.Details["Total Episode"] = strings.TrimSpace(strings.TrimPrefix(text, "Episodes:"))
			case strings.HasPrefix(text, "Released on:"):
				timeText := s.Find("time").Text()
				if timeText != "" {
					animeData.Details["Released on"] = timeText
				}
			case strings.HasPrefix(text, "Updated on:"):
				timeText := s.Find("time").Text()
				if timeText != "" {
					animeData.Details["Updated on"] = timeText
				}
			case strings.HasPrefix(text, "Producers:"):
				var producers []string
				s.Find("a").Each(func(_ int, p *goquery.Selection) {
					producers = append(producers, p.Text())
				})
				if len(producers) > 0 {
					animeData.Details["Producers"] = strings.Join(producers, ", ")
				}
			}
		})
	})

	// Alternative scraper for different page structure (like tougen-anki)
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Skip if we already got data from article.post-180
		if animeData.Judul != "" {
			return
		}

		// Extract title from h1
		title := e.DOM.Find("h1").First().Text()
		animeData.Judul = strings.TrimSpace(title)

		// Extract thumbnail - look for the main anime image
		var thumbnailURL string
		// First try to find img with specific patterns that indicate it's the main poster
		e.DOM.Find("img").Each(func(_ int, img *goquery.Selection) {
			src := img.AttrOr("src", "")
			dataSrc := img.AttrOr("data-src", "")
			
			// Check both src and data-src
			for _, url := range []string{src, dataSrc} {
				if url != "" && !strings.Contains(url, "data:image/svg") && 
				   !strings.Contains(url, "logo") && !strings.Contains(url, "icon") &&
				   !strings.Contains(url, "avatar") {
					if strings.Contains(url, "wp-content/uploads") && 
					   (strings.Contains(url, ".jpg") || strings.Contains(url, ".webp") || strings.Contains(url, ".png")) &&
					   (strings.Contains(url, "resize=247,350") || strings.Contains(url, "resize=350") || len(url) > 50) {
						thumbnailURL = url
						return
					}
				}
			}
		})
		animeData.Thumbnail = thumbnailURL

		// Extract synopsis - look for content after "Synopsis" heading or in specific sections
		var synopsisText string
		
		// Method 1: Look for Synopsis heading
		e.DOM.Find("h4, h2, h3").Each(func(_ int, heading *goquery.Selection) {
			headingText := strings.ToLower(heading.Text())
			if strings.Contains(headingText, "synopsis") && synopsisText == "" {
				// Get all following content until next heading
				next := heading.Next()
				var content []string
				for next.Length() > 0 {
					text := strings.TrimSpace(next.Text())
					if text != "" {
						// Stop if we hit another heading
						if next.Is("h1, h2, h3, h4, h5, h6") {
							break
						}
						content = append(content, text)
						if len(strings.Join(content, " ")) > 100 {
							break
						}
					}
					next = next.Next()
				}
				if len(content) > 0 {
					synopsisText = strings.Join(content, " ")
				}
			}
		})
		
		// Method 2: If no synopsis found, look for paragraphs with substantial content
		if synopsisText == "" {
			e.DOM.Find("p").Each(func(_ int, p *goquery.Selection) {
				text := strings.TrimSpace(p.Text())
				if len(text) > 200 && synopsisText == "" { // Look for substantial paragraphs
					// Make sure it's not navigation or footer text
					if !strings.Contains(strings.ToLower(text), "watch") && 
					   !strings.Contains(strings.ToLower(text), "download") &&
					   !strings.Contains(strings.ToLower(text), "streaming") {
						synopsisText = text
					}
				}
			})
		}
		
		animeData.Sinopsis = synopsisText

		// Extract genres - only from the main content area, not sidebar
		genreMap := make(map[string]bool) // Use map to avoid duplicates
		// Look for genre links in the main content area only
		e.DOM.Find("main a[href*='/genres/'], .content a[href*='/genres/'], article a[href*='/genres/']").Each(func(_ int, g *goquery.Selection) {
			genreText := strings.TrimSpace(g.Text())
			if genreText != "" && len(genreText) < 30 { // Avoid long text that's not a genre
				genreMap[genreText] = true
			}
		})
		
		// If no genres found in main content, try a more specific approach
		if len(genreMap) == 0 {
			// Look for genre links that are close to the anime title or in the first part of the page
			e.DOM.Find("a[href*='/genres/']").Each(func(i int, g *goquery.Selection) {
				if i < 10 { // Only check first 10 genre links to avoid sidebar
					genreText := strings.TrimSpace(g.Text())
					if genreText != "" && len(genreText) < 30 {
						genreMap[genreText] = true
					}
				}
			})
		}
		
		// Convert map to slice
		for genre := range genreMap {
			animeData.Genre = append(animeData.Genre, genre)
		}

		// Extract details from the page content with more specific regex
		pageText := e.DOM.Text()
		
		// Extract status with better regex
		statusRegex := regexp.MustCompile(`Status:\s*([A-Za-z]+)`)
		if matches := statusRegex.FindStringSubmatch(pageText); len(matches) > 1 {
			animeData.Details["Status"] = strings.TrimSpace(matches[1])
		}

		// Extract studio - be more specific
		e.DOM.Find("a[href*='/studio/']").First().Each(func(_ int, s *goquery.Selection) {
			studioText := strings.TrimSpace(s.Text())
			if studioText != "" {
				animeData.Details["Studio"] = studioText
			}
		})

		// Extract season - be more specific
		e.DOM.Find("a[href*='/season/']").First().Each(func(_ int, s *goquery.Selection) {
			seasonText := strings.TrimSpace(s.Text())
			if seasonText != "" {
				animeData.Details["Season"] = seasonText
			}
		})

		// Extract producers - be more specific
		var producers []string
		e.DOM.Find("a[href*='/producer/']").Each(func(i int, p *goquery.Selection) {
			if i < 5 { // Limit to first 5 to avoid duplicates
				producerText := strings.TrimSpace(p.Text())
				if producerText != "" {
					producers = append(producers, producerText)
				}
			}
		})
		if len(producers) > 0 {
			animeData.Details["Producers"] = strings.Join(producers, ", ")
		}

		// Extract type with better regex
		typeRegex := regexp.MustCompile(`Type:\s*([A-Za-z]+)`)
		if matches := typeRegex.FindStringSubmatch(pageText); len(matches) > 1 {
			animeData.Details["Type"] = strings.TrimSpace(matches[1])
		}

		// Extract released year with better regex
		releasedRegex := regexp.MustCompile(`Released:\s*(\d{4})`)
		if matches := releasedRegex.FindStringSubmatch(pageText); len(matches) > 1 {
			animeData.Details["Released:"] = strings.TrimSpace(matches[1])
		}
	})

	// Optimized episode list scraper
	c.OnHTML("div.eplister ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			ep := ScrapedEpisode{
				Episode:      el.ChildText(".epl-num"),
				Judul:        el.ChildText(".epl-title"),
				URL:          e.Request.AbsoluteURL(el.ChildAttr("a", "href")),
				TanggalRilis: el.ChildText(".epl-date"),
			}
			animeData.EpisodeList = append(animeData.EpisodeList, ep)
		})
	})

	// Recommendation scraper using correct selector from test
	c.OnHTML(`div.bixbox:has(h3 > span:contains("Recommended Series"))`, func(e *colly.HTMLElement) {
		e.ForEach(".listupd article.bs", func(_ int, el *colly.HTMLElement) {
			// Handle lazy-loading for thumbnail
			thumbURL := el.ChildAttr("img", "data-src")
			if thumbURL == "" {
				thumbURL = el.ChildAttr("img", "src")
			}

			rec := ScrapedRecommendation{
				Judul:     el.ChildAttr("a.tip", "title"),
				URL:       e.Request.AbsoluteURL(el.ChildAttr("a.tip", "href")),
				Episode:   el.ChildText("span.epx"),
				Thumbnail: thumbURL,
			}
			animeData.Rekomendasi = append(animeData.Rekomendasi, rec)
		})
	})

	// Post-processing with minimal fallbacks to preserve real data
	c.OnScraped(func(r *colly.Response) {
		// Only use placeholder if absolutely no thumbnail found
		if animeData.Thumbnail == "" {
			log.Printf("Warning: No thumbnail found for %s", animeData.Judul)
			animeData.Thumbnail = "https://via.placeholder.com/350x500?text=No+Image"
		}

		// Log what we found for debugging
		log.Printf("Successfully scraped anime: %s", animeData.Judul)
		log.Printf("- Thumbnail: %s", animeData.Thumbnail)
		log.Printf("- Score: %s", animeData.Skor)
		log.Printf("- Synopsis length: %d chars", len(animeData.Sinopsis))
		log.Printf("- Genres: %d found", len(animeData.Genre))
		log.Printf("- Episodes: %d found", len(animeData.EpisodeList))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping detail: %v | URL: %s", err, r.Request.URL)
	})

	log.Printf("Visiting detail page: %s", targetURL)
	c.Visit(targetURL)
	c.Wait()
	return animeData
}

// ScrapeEpisodeDetail dengan optimasi dan pre-compiled regex.
func ScrapeEpisodeDetail(episodeURL string) ScrapedEpisodeDetails {
	// Pre-compile regex for better performance
	var (
		srcRegex = regexp.MustCompile(`src="([^"]+)"`)
		resRegex = regexp.MustCompile(`(\d{3,4}p)`)
	)

	data := ScrapedEpisodeDetails{
		StreamingServers: make([]StreamingServer, 0, 10),
		DownloadLinks:    make(map[string]map[string][]DownloadProvider, 5),
		OtherEpisodes:    make([]OtherEpisode, 0, 50),
		AnimeInfo: AnimeInfo{
			Genres: make([]string, 0, 10),
		},
	}

	var mainThumbnail string
	var seriesTitle string

	c := createOptimizedCollector(false, 1)

	// Optimized anime info scraper
	c.OnHTML(".bixbox.single-info", func(e *colly.HTMLElement) {
		thumbSrc := e.ChildAttr(".thumb img", "data-src")
		if thumbSrc == "" {
			thumbSrc = e.ChildAttr(".thumb img", "src")
		}
		mainThumbnail = thumbSrc
		data.ThumbnailURL = mainThumbnail
		data.AnimeInfo.ThumbnailURL = mainThumbnail

		infox := e.DOM.Find(".infox")
		seriesTitle = infox.Find("h2.entry-title").Text()
		data.AnimeInfo.Title = seriesTitle
		data.AnimeInfo.Synopsis = strings.TrimSpace(infox.Find(".desc p").Text())

		// Use map for deduplication
		genreMap := make(map[string]bool)
		infox.Find(".genxed a").Each(func(_ int, s *goquery.Selection) {
			genreMap[s.Text()] = true
		})
		for genre := range genreMap {
			data.AnimeInfo.Genres = append(data.AnimeInfo.Genres, genre)
		}
	})

	// Episode title
	c.OnHTML("h1.entry-title", func(e *colly.HTMLElement) {
		data.Title = e.Text
		if seriesTitle == "" && strings.Contains(data.Title, " Episode ") {
			seriesTitle = strings.Split(data.Title, " Episode ")[0]
			data.AnimeInfo.Title = seriesTitle
		}
	})

	// Optimized streaming servers
	c.OnHTML("select.mirror", func(e *colly.HTMLElement) {
		e.ForEach("option", func(_ int, el *colly.HTMLElement) {
			base64Value := el.Attr("value")
			if base64Value == "" {
				return
			}

			decodedIframe, err := base64.StdEncoding.DecodeString(base64Value)
			if err != nil {
				return
			}

			matches := srcRegex.FindStringSubmatch(string(decodedIframe))
			if len(matches) > 1 {
				data.StreamingServers = append(data.StreamingServers, StreamingServer{
					ServerName:   strings.TrimSpace(el.Text),
					StreamingURL: matches[1],
				})
			}
		})
	})

	// Navigation - Updated selectors based on actual HTML structure
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Look for navigation links in the page
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			linkText := strings.ToLower(strings.TrimSpace(el.Text))
			href := el.Attr("href")
			
			if href != "" {
				switch linkText {
				case "prev":
					data.Navigation.PreviousEpisodeURL = e.Request.AbsoluteURL(href)
				case "all episodes":
					data.Navigation.AllEpisodesURL = e.Request.AbsoluteURL(href)
				case "next":
					data.Navigation.NextEpisodeURL = e.Request.AbsoluteURL(href)
				}
			}
		})
	})

	// Optimized episode list processing
	c.OnHTML("#mainepisode .episodelist ul, div.eplister ul", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			thumb := el.ChildAttr(".epl-thumb img", "data-src")
			if thumb == "" {
				thumb = el.ChildAttr(".epl-thumb img", "src")
			}
			if thumb == "" {
				thumb = mainThumbnail
			}

			episodeTitle := el.ChildText(".epl-title")
			episodeNumber := el.ChildText(".epl-num")
			episodeURL := e.Request.AbsoluteURL(el.ChildAttr("a", "href"))

			// Optimize title generation
			if episodeTitle == "" {
				episodeTitle = generateEpisodeTitle(episodeURL, seriesTitle, episodeNumber)
			}

			data.OtherEpisodes = append(data.OtherEpisodes, OtherEpisode{
				Title:        episodeTitle,
				URL:          episodeURL,
				ReleaseDate:  el.ChildText(".epl-date"),
				ThumbnailURL: thumb,
			})
		})
	})

	// Alternative episode scraper
	c.OnHTML("div.bixbox.mctn", func(e *colly.HTMLElement) {
		if len(data.OtherEpisodes) == 0 {
			e.ForEach("article.bs", func(_ int, el *colly.HTMLElement) {
				thumb := el.ChildAttr("img", "data-src")
				if thumb == "" {
					thumb = el.ChildAttr("img", "src")
				}
				if thumb == "" {
					thumb = mainThumbnail
				}

				episodeTitle := el.ChildAttr("a", "title")
				episodeURL := e.Request.AbsoluteURL(el.ChildAttr("a", "href"))
				episodeInfo := el.ChildText(".epx")

				if episodeTitle == "" && seriesTitle != "" {
					if episodeInfo != "" {
						episodeTitle = fmt.Sprintf("%s %s", seriesTitle, episodeInfo)
					} else {
						episodeTitle = seriesTitle
					}
				}

				data.OtherEpisodes = append(data.OtherEpisodes, OtherEpisode{
					Title:        episodeTitle,
					URL:          episodeURL,
					ReleaseDate:  "",
					ThumbnailURL: thumb,
				})
			})
		}
	})

	// Post-processing optimization
	c.OnScraped(func(r *colly.Response) {
		// Generate download links efficiently
		if len(data.StreamingServers) > 0 && len(data.DownloadLinks) == 0 {
			downloadMap := make(map[string][]DownloadProvider)
			for _, server := range data.StreamingServers {
				resolution := "HD"
				if matches := resRegex.FindStringSubmatch(strings.ToLower(server.ServerName)); len(matches) > 1 {
					resolution = matches[1]
				}
				downloadMap[resolution] = append(downloadMap[resolution], DownloadProvider{
					Provider: server.ServerName,
					URL:      server.StreamingURL,
				})
			}
			data.DownloadLinks["MP4 (from Stream)"] = downloadMap
		}

		// Final title processing
		if data.AnimeInfo.Title == "" && data.Title != "" {
			if strings.Contains(data.Title, " Episode ") {
				data.AnimeInfo.Title = strings.Split(data.Title, " Episode ")[0]
			} else {
				data.AnimeInfo.Title = data.Title
			}
		}

		// Batch process episode titles
		for i := range data.OtherEpisodes {
			if data.OtherEpisodes[i].Title == "" {
				data.OtherEpisodes[i].Title = generateEpisodeTitle(
					data.OtherEpisodes[i].URL,
					data.AnimeInfo.Title,
					"",
				)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error scraping episode detail: %v | URL: %s", err, r.Request.URL)
	})

	c.Visit(episodeURL)
	c.Wait()
	return data
}

// Helper function for episode title generation
func generateEpisodeTitle(episodeURL, seriesTitle, episodeNumber string) string {
	if episodeNumber == "" {
		// Extract episode number from URL
		urlParts := strings.Split(episodeURL, "/")
		for _, part := range urlParts {
			if strings.Contains(part, "-") {
				lastPart := strings.Split(part, "-")
				if len(lastPart) > 0 {
					possibleEp := lastPart[len(lastPart)-1]
					if _, err := strconv.Atoi(possibleEp); err == nil {
						episodeNumber = possibleEp
						break
					}
				}
			}
		}
	}

	if seriesTitle != "" && episodeNumber != "" {
		return fmt.Sprintf("%s Episode %s", seriesTitle, episodeNumber)
	} else if seriesTitle != "" {
		return seriesTitle
	}

	// Last resort: convert URL slug to title
	urlParts := strings.Split(strings.TrimSuffix(episodeURL, "/"), "/")
	if len(urlParts) > 0 {
		lastPart := urlParts[len(urlParts)-1]
		titleParts := strings.Split(lastPart, "-")
		for i, part := range titleParts {
			if len(part) > 0 {
				titleParts[i] = strings.Title(part)
			}
		}
		return strings.Join(titleParts, " ")
	}

	return "Unknown Episode"
}

func ScrapeSearch(query string) []ScrapedSearchResult {
	searchURL := fmt.Sprintf("https://gomunime.co/?s=%s", url.QueryEscape(query))
	ajaxURL := "https://gomunime.co/wp-admin/admin-ajax.php"
	var searchResults []ScrapedSearchResult
	var mu sync.Mutex

	c := colly.NewCollector(colly.Async(true))
	c.UserAgent = userAgent
	c.Limit(&colly.LimitRule{DomainGlob: "*gomunime.co*", Parallelism: 4})

	// Callback untuk memproses detail dari AJAX (info hover)
	c.OnHTML("div.ingfo", func(e *colly.HTMLElement) {
		result := e.Request.Ctx.GetAny("result").(ScrapedSearchResult)

		e.ForEach("div.minginfo span.l", func(_ int, el *colly.HTMLElement) {
			if el.ChildText("i.fa-star") != "" {
				result.Skor = strings.TrimSpace(el.Text)
			}
		})
		result.Sinopsis = strings.TrimSpace(e.ChildText("div.ingdesc .contexcerpt"))
		e.ForEach(".linginfo span", func(_ int, s *colly.HTMLElement) {
			text := s.Text
			if strings.HasPrefix(text, "Genres:") {
				s.ForEach("a", func(_ int, a *colly.HTMLElement) {
					result.Genres = append(result.Genres, a.Text)
				})
			} else if strings.HasPrefix(text, "Status:") {
				result.Status = strings.TrimSpace(strings.TrimPrefix(text, "Status:"))
			}
		})
		mu.Lock()
		searchResults = append(searchResults, result)
		mu.Unlock()
	})

	// Callback utama untuk hasil pencarian
	c.OnHTML("div.listupd article.bs", func(e *colly.HTMLElement) {
		linkElement := e.DOM.Find("a.tip")
		postID, exists := linkElement.Attr("rel")
		if !exists {
			return
		}

		thumbURL := e.ChildAttr("img", "data-src")
		if thumbURL == "" {
			thumbURL = e.ChildAttr("img", "src")
		}

		result := ScrapedSearchResult{
			Judul:     linkElement.AttrOr("title", ""),
			Tautan:    e.Request.AbsoluteURL(linkElement.AttrOr("href", "")),
			Thumbnail: thumbURL,
			Status:    e.ChildText("span.epx"),
			Tipe:      e.ChildText("div.typez"),
			Genres:    []string{},
		}
		ctx := colly.NewContext()
		ctx.Put("result", result)
		payload := fmt.Sprintf("action=tooltip_action&id=%s", postID)
		c.Request("POST", ajaxURL, strings.NewReader(payload), ctx, nil)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error scraping search:", err, "| URL:", r.Request.URL)
	})
	c.Visit(searchURL)
	c.Wait()
	return searchResults
}

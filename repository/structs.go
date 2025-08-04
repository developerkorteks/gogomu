package repository

type SearchResultItem struct {
	Judul     string   `json:"judul" example:"Naruto Kecil"`
	URLAnime  string   `json:"url_anime" example:"https://v1.samehadaku.how/anime/naruto-kecil/"`
	AnimeSlug string   `json:"anime_slug" example:"naruto-kecil"`
	Status    string   `json:"status" example:"Completed"`
	Tipe      string   `json:"tipe" example:"TV"`
	Skor      string   `json:"skor" example:"8.84"`
	Penonton  string   `json:"penonton" example:"154157 Views"`
	Sinopsis  string   `json:"sinopsis" example:"Beberapa saat sebelum Naruto Uzumaki lahir..."`
	Genre     []string `json:"genre" example:"Action,Adventure"`
	URLCover  string   `json:"url_cover" example:"https://v1.samehadaku.how/wp-content/uploads/2024/08/142503.jpg"`
}

// SearchResponse adalah struct untuk output endpoint /search/.
type SearchResponse struct {
	ConfidenceScore float64            `json:"confidence_score" example:"1.0"`
	Data            []SearchResultItem `json:"data"`
	Message         string             `json:"message"`
	Source          string             `json:"source"`
}

// --- Struct untuk DATA SCRAPER ---

type ScrapedSearchResult struct {
	Judul     string
	Tautan    string
	Thumbnail string
	Tipe      string
	Status    string
	Skor      string
	Sinopsis  string
	Genres    []string
}

// --- Structs untuk Output JSON Akhir ---

type MovieItem struct {
	Judul     string   `json:"judul" example:"Spy x Family: Code White"`
	URL       string   `json:"url" example:"https://v1.samehadaku.how/anime/spy-x-family-code-white/"`
	AnimeSlug string   `json:"anime_slug" example:"spy-x-family-code-white"`
	Status    string   `json:"status" example:"Completed"`
	Skor      string   `json:"skor" example:"8.1"`
	Sinopsis  string   `json:"sinopsis" example:"Loid Forger, an elite spy, is warned by his handler..."`
	Views     string   `json:"views" example:"342490 Views"`
	Cover     string   `json:"cover" example:"https://v1.samehadaku.how/wp-content/uploads/2024/08/139388.jpg"`
	Genres    []string `json:"genres" example:"Action,Comedy"`
	Tanggal   string   `json:"tanggal" example:"22 Desember 2023"`
}

// MovieListResponse adalah struct untuk output endpoint /movie/.
type MovieListResponse struct {
	ConfidenceScore float64     `json:"confidence_score" example:"1"`
	Data            []MovieItem `json:"data"`
	Message         string      `json:"message"`
	Source          string      `json:"source"`
}

// Sempurnakan ScrapedLatestAnime untuk menyimpan lebih banyak detail
type ScrapedLatestAnime struct {
	Judul     string
	Tautan    string
	Episode   string
	Thumbnail string
	Tipe      string
	Rating    string
	Status    string   // Tambahkan status
	Deskripsi string   // Tambahkan deskripsi
	Genres    []string // Ubah menjadi slice of string untuk kemudahan
}
type JadwalHarianResponse struct {
	ConfidenceScore float64               `json:"confidence_score" example:"1"`
	Data            []JadwalAnimeResponse `json:"data"`
	Message         string                `json:"message"`
	Source          string                `json:"source"`
}

// FinalResponse adalah struct utama untuk output JSON API.
type FinalResponse struct {
	ConfidenceScore float64                  `json:"confidence_score"`
	Data            HomeData                 `json:"data"`
	Message         string                   `json:"message"`
	Source          string                   `json:"source"`
}

// HomeData adalah struct untuk data halaman utama
type HomeData struct {
	Top10       []Top10Anime             `json:"top10"`
	NewEps      []NewEps                 `json:"new_eps"`
	Movies      []Movie                  `json:"movies"`
	JadwalRilis map[string][]JadwalAnime `json:"jadwal_rilis"`
}

// Top10Anime merepresentasikan item dalam daftar top 10.
type Top10Anime struct {
	Judul     string   `json:"judul"`
	URL       string   `json:"url"`
	AnimeSlug string   `json:"anime_slug"`
	Rating    string   `json:"rating"`
	Cover     string   `json:"cover"`
	Genres    []string `json:"genres"`
}

// NewEps merepresentasikan item dalam daftar episode baru.
type NewEps struct {
	Judul     string `json:"judul"`
	URL       string `json:"url"`
	AnimeSlug string `json:"anime_slug"`
	Episode   string `json:"episode"`
	Rilis     string `json:"rilis"` // Catatan: Data ini tidak tersedia dari scraper
	Cover     string `json:"cover"`
}

// Movie merepresentasikan item dalam daftar film.
type Movie struct {
	Judul     string   `json:"judul"`
	URL       string   `json:"url"`
	AnimeSlug string   `json:"anime_slug"`
	Tanggal   string   `json:"tanggal"` // Catatan: Data ini tidak tersedia dari scraper
	Cover     string   `json:"cover"`
	Genres    []string `json:"genres"`
}

// JadwalAnime merepresentasikan satu entri anime dalam jadwal rilis.
type JadwalAnime struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	AnimeSlug   string   `json:"anime_slug"`
	CoverURL    string   `json:"cover_url"`
	Type        string   `json:"type"`   // Catatan: Data ini tidak tersedia dari scraper
	Score       string   `json:"score"`  // Catatan: Data ini tidak tersedia dari scraper
	Genres      []string `json:"genres"` // Catatan: Data ini tidak tersedia dari scraper
	ReleaseTime string   `json:"release_time"`
}

// --- Structs untuk Data Hasil Scraper ---

// ScrapedDaySchedule merepresentasikan jadwal untuk satu hari dari scraper.
type ScrapedDaySchedule struct {
	Hari      string
	AnimeList []ScrapedAnimeSchedule
}

// ScrapedAnimeSchedule merepresentasikan satu anime dalam jadwal dari scraper.
type ScrapedAnimeSchedule struct {
	Judul      string
	Tautan     string
	WaktuRilis string
	Thumbnail  string
}

type JadwalAnimeResponse struct {
	Title       string   `json:"title" example:"Busamen Gachi Fighter"`
	URL         string   `json:"url" example:"https://v1.samehadaku.how/anime/busamen-gachi-fighter/"`
	AnimeSlug   string   `json:"anime_slug" example:"busamen-gachi-fighter"`
	CoverURL    string   `json:"cover_url" example:"https://v1.samehadaku.how/wp-content/uploads/2025/07/150515.jpg"`
	Type        string   `json:"type" example:"TV"`
	Score       string   `json:"score" example:"6.68"`
	Genres      []string `json:"genres" example:"Action,Adventure"`
	ReleaseTime string   `json:"release_time" example:"00:00"`
}

// AnimeDetailResponse adalah struct utama untuk output endpoint /anime-detail/.
type AnimeDetailResponse struct {
	ConfidenceScore float64              `json:"confidence_score" example:"1"`
	Data            AnimeDetailData      `json:"data"`
	Message         string               `json:"message"`
	Source          string               `json:"source"`
}

type AnimeDetailData struct {
	Judul           string               `json:"judul" example:"Nonton Anime Haikyuu!! Movie..."`
	URLAnime        string               `json:"url_anime" example:"https://gomunime.co/anime/haikyuu-movie-gomisuteba-no-kessen/"`
	AnimeSlug       string               `json:"anime_slug" example:"haikyuu-movie-gomisuteba-no-kessen"`
	URLCover        string               `json:"url_cover" example:"https://gomunime.co/wp-content/uploads/2024/10/140360.jpg"`
	EpisodeList     []EpisodeListItem    `json:"episode_list"`
	Recommendations []RecommendationItem `json:"recommendations"`
	Status          string               `json:"status" example:"Completed"`
	Tipe            string               `json:"tipe" example:"Movie"`
	Skor            string               `json:"skor" example:"8.65"`
	Penonton        string               `json:"penonton" example:"N/A"`
	Sinopsis        string               `json:"sinopsis" example:"Kozume Kenma tidak pernah menganggap..."`
	Genre           []string             `json:"genre" example:"School,Sports"`
	Details         Details              `json:"details"`
	Rating          RatingInfo           `json:"rating"`
}

// EpisodeListItem merepresentasikan satu episode dalam daftar.
type EpisodeListItem struct {
	Episode     string `json:"episode" example:"1"`
	Title       string `json:"title" example:"Haikyuu Movie: Gomisuteba no Kessen"`
	URL         string `json:"url" example:"https://v1.samehadaku.how/haikyuu-gomisuteba-no-kessen/"`
	EpisodeSlug string `json:"episode_slug" example:"haikyuu-gomisuteba-no-kessen"`
	ReleaseDate string `json:"release_date" example:"31 October 2024"`
}

// RecommendationItem merepresentasikan satu anime dalam daftar rekomendasi.
type RecommendationItem struct {
	Title     string `json:"title" example:"Overlord III"`
	URL       string `json:"url" example:"https://v1.samehadaku.how/anime/overlord-iii/"`
	AnimeSlug string `json:"anime_slug" example:"overlord-iii"`
	CoverURL  string `json:"cover_url" example:"https://v1.samehadaku.how/wp-content/uploads/2024/07/93473.jpg"`
	Rating    string `json:"rating" example:"8"`
	Episode   string `json:"episode" example:"Eps 12"`
}

// Details adalah objek untuk metadata detail.
type Details struct {
	Japanese     string `json:"Japanese,omitempty"`
	Synonyms     string `json:"Synonyms,omitempty"`
	English      string `json:"English,omitempty"`
	Status       string `json:"Status,omitempty"`
	Type         string `json:"Type,omitempty"`
	Source       string `json:"Source,omitempty"`
	Duration     string `json:"Duration,omitempty"`
	TotalEpisode string `json:"Total Episode,omitempty"`
	Season       string `json:"Season,omitempty"`
	Studio       string `json:"Studio,omitempty"`
	Producers    string `json:"Producers,omitempty"`
	Released     string `json:"Released:,omitempty"`
}

// RatingInfo adalah objek untuk informasi rating.
type RatingInfo struct {
	Score string `json:"score" example:"8.65"`
	Users string `json:"users" example:"34,719"`
}

// --- Structs untuk DATA SCRAPER ---

type ScrapedEpisode struct {
	Episode      string
	Judul        string
	URL          string
	TanggalRilis string
}

type ScrapedRecommendation struct {
	Judul     string
	URL       string
	Thumbnail string
	Episode   string // 'Status' dari scraper lama kita gunakan sebagai 'Episode'
}

type ScrapedAnimeDetails struct {
	Judul       string
	Thumbnail   string
	Skor        string
	Sinopsis    string
	Genre       []string
	EpisodeList []ScrapedEpisode
	Rekomendasi []ScrapedRecommendation
	Details     map[string]string // Menggunakan map untuk fleksibilitas
}

type EpisodeDetailResponse struct {
	ConfidenceScore  float64                                  `json:"confidence_score" example:"1"`
	Data             EpisodeDetailData                        `json:"data"`
	Message          string                                   `json:"message"`
	Source           string                                   `json:"source"`
}

type EpisodeDetailData struct {
	Title            string                                   `json:"title" example:"Haikyuu Movie: Gomisuteba no Kessen Sub Indo"`
	ThumbnailURL     string                                   `json:"thumbnail_url" example:"https://gomunime.co/wp-content/uploads/2024/10/140360.jpg"`
	StreamingServers []StreamingServer                        `json:"streaming_servers"`
	ReleaseInfo      string                                   `json:"release_info" example:"9 months yang lalu"`
	DownloadLinks    map[string]map[string][]DownloadProvider `json:"download_links"`
	Navigation       EpisodeNavigation                        `json:"navigation"`
	AnimeInfo        AnimeInfo                                `json:"anime_info"`
	OtherEpisodes    []OtherEpisode                           `json:"other_episodes"`
}

// StreamingServer merepresentasikan satu server streaming.
type StreamingServer struct {
	ServerName   string `json:"server_name" example:"Nakama 1080p"`
	StreamingURL string `json:"streaming_url" example:"https://pixeldrain.com/api/file/Ra5A3rtj"`
}

// DownloadProvider merepresentasikan satu link download dari satu provider.
type DownloadProvider struct {
	Provider string `json:"provider" example:"Gofile"`
	URL      string `json:"url" example:"https://gofile.io/d/l3ahTO"`
}

// EpisodeNavigation berisi link navigasi antar episode.
type EpisodeNavigation struct {
	PreviousEpisodeURL string `json:"previous_episode_url,omitempty"`
	AllEpisodesURL     string `json:"all_episodes_url,omitempty"`
	NextEpisodeURL     string `json:"next_episode_url,omitempty"`
}

// AnimeInfo berisi detail dari seri anime induknya.
type AnimeInfo struct {
	Title        string   `json:"title" example:"Haikyuu!! Movie: Gomisuteba no Kessen"`
	ThumbnailURL string   `json:"thumbnail_url" example:"https://v1.samehadaku.how/wp-content/uploads/2024/10/140360.jpg"`
	Synopsis     string   `json:"synopsis" example:"Kozume Kenma tidak pernah menganggap..."`
	Genres       []string `json:"genres" example:"School,Sports"`
}

// OtherEpisode merepresentasikan episode lain dari seri yang sama.
type OtherEpisode struct {
	Title        string `json:"title" example:"Haikyuu Movie: Gomisuteba no Kessen"`
	URL          string `json:"url" example:"https://v1.samehadaku.how/haikyuu-gomisuteba-no-kessen/"`
	ThumbnailURL string `json:"thumbnail_url" example:"https://v1.samehadaku.how/wp-content/uploads/2024/10/Haikyu.The_.Dumpster.Battle.jpg"`
	ReleaseDate  string `json:"release_date" example:"31 October 2024"`
}

// --- Struct untuk DATA SCRAPER ---
type ScrapedEpisodeDetails struct {
	Title            string
	ThumbnailURL     string
	ReleaseInfo      string
	StreamingServers []StreamingServer
	DownloadLinks    map[string]map[string][]DownloadProvider
	Navigation       EpisodeNavigation
	AnimeInfo        AnimeInfo
	OtherEpisodes    []OtherEpisode
}

// type OtherEpisode struct {
// 	Title        string `json:"title" example:"Haikyuu Movie: Gomisuteba no Kessen"`
// 	URL          string `json:"url" example:"https://v1.samehadaku.how/haikyuu-gomisuteba-no-kessen/"`
// 	ThumbnailURL string `json:"thumbnail_url" example:"https://v1.samehadaku.how/wp-content/uploads/2024/10/Haikyu.The_.Dumpster.Battle.jpg"`
// 	ReleaseDate  string `json:"release_date" example:"31 October 2024"`
// }

type AnimeTerbaruItem struct {
	Judul     string `json:"judul" example:"Zutaboro Reijou wa Ane no Moto"`
	URL       string `json:"url" example:"https://v1.samehadaku.how/anime/zutaboro-reijou-wa-ane-no-moto/"`
	AnimeSlug string `json:"anime_slug" example:"zutaboro-reijou-wa-ane-no-moto"`
	Episode   string `json:"episode" example:"5"`
	Uploader  string `json:"uploader" example:"Urusai"`
	Rilis     string `json:"rilis" example:"5 hours yang lalu"`
	Cover     string `json:"cover" example:"https://v1.samehadaku.how/wp-content/uploads/2025/08/Zutaboro-Reijou-wa-Ane-no-Moto-Episode-5.jpg"`
}

// AnimeTerbaruResponse adalah struct untuk output endpoint /anime-terbaru/.
type AnimeTerbaruResponse struct {
	ConfidenceScore float64            `json:"confidence_score" example:"1.0"`
	Data            []AnimeTerbaruItem `json:"data"`
	Message         string             `json:"message"`
	Source          string             `json:"source"`
}

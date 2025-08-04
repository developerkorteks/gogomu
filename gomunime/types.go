package gomunime

type PopularAnime struct {
	Peringkat int      `json:"peringkat"`
	Judul     string   `json:"judul"`
	Tautan    string   `json:"tautan"`
	Thumbnail string   `json:"thumbnail"`
	Genres    []string `json:"genres"`
	Rating    string   `json:"rating"`
}

// File: types.go
// [BARU] Struct untuk menampung hasil pencarian
type SearchResult struct {
	Judul       string   `json:"judul"`
	Tautan      string   `json:"tautan"`
	Thumbnail   string   `json:"thumbnail"`
	StatusRilis string   `json:"status_rilis"`
	Tipe        string   `json:"tipe"`
	Sinopsis    string   `json:"sinopsis,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Studio      string   `json:"studio,omitempty"`
}
type DownloadLink struct {
	Resolusi string `json:"resolusi"`
	Provider string `json:"provider"`
	Tautan   string `json:"tautan"`
}

type Episode struct {
	Episode      string `json:"episode"`
	JudulEpisode string `json:"judul_episode"`
	Tautan       string `json:"tautan"`
	TanggalRilis string `json:"tanggal_rilis"`
}

type Recommendation struct {
	Judul     string `json:"judul"`
	Tautan    string `json:"tautan"`
	Status    string `json:"status"`
	Thumbnail string `json:"thumbnail"`
}

type AnimeDetails struct {
	Judul           string           `json:"judul"`
	JudulAlternatif string           `json:"judul_alternatif"`
	Thumbnail       string           `json:"thumbnail"`
	Rating          string           `json:"rating"`
	Sinopsis        string           `json:"sinopsis"`
	Status          string           `json:"status"`
	Studio          string           `json:"studio"`
	RilisPerdana    string           `json:"rilis_perdana"`
	Season          string           `json:"season"`
	Tipe            string           `json:"tipe"`
	TotalEpisode    string           `json:"total_episode"`
	Sutradara       string           `json:"sutradara"`
	Pemeran         []string         `json:"pemeran,omitempty"`
	Produser        []string         `json:"produser,omitempty"`
	Genre           []string         `json:"genre"`
	RilisPada       string           `json:"rilis_pada,omitempty"`
	DiperbaruiPada  string           `json:"diperbarui_pada,omitempty"`
	EpisodeList     []Episode        `json:"episode_list"`
	Rekomendasi     []Recommendation `json:"rekomendasi"`
	StreamURL       string           `json:"stream_url,omitempty"`
	DownloadLinks   []DownloadLink   `json:"download_links,omitempty"`
	// [BARU] Field untuk menyimpan semua link mirror
	MirrorStreams map[string]string `json:"mirror_streams,omitempty"`
}

package models

// Модель для работы с базой данных
type Model struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

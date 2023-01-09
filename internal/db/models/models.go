package models

// Модель для работы с базой данных
type Model struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"original_url"`
}

// Модель основного хранилища
type MainModel struct {
	ID      string `json:"uuid"`
	URLpair Model
}

package hasher

// Интерфейс для работы с неким шифровщиком ссылок
type URLHasher interface {
	GetHash(lenHash int) string
}

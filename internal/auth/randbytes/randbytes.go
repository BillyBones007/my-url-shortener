package randbytes

import "crypto/rand"

// Возвращает случайную последовательность байт размера size
func RandomBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

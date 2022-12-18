package files

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

//_______________ Блок описания типов ____________________________

// Тип для работы с файлом в роли основного хранилища
type FileStorage struct {
	FilePATH string // путь до файла-хранилища
}

// Тип для записи в файл-хранилище
type fileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

// Тип для чтения из файла-хранилища
type fileReader struct {
	file    *os.File
	decoder *json.Decoder
}

// Тип для декодирования данных файла
type tmpDecode models.Model

//_______________ Блок описания конструкторов ____________________________

// Конструктор хранилища. Возвращает указатель на FileStorage
func NewStorage(filename string) (*FileStorage, error) {
	return &FileStorage{FilePATH: filename}, nil
}

// Конструктор типа записи в файл-хранилище
func NewFileWriter(f *FileStorage) (*fileWriter, error) {
	fileDescr, err := os.OpenFile(f.FilePATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	encoder := json.NewEncoder(fileDescr)
	return &fileWriter{file: fileDescr, encoder: encoder}, nil
}

// Конструктор типа записи в файл-хранилище
func NewFileReader(f *FileStorage) (*fileReader, error) {
	fileDescr, err := os.OpenFile(f.FilePATH, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(fileDescr)
	return &fileReader{file: fileDescr, decoder: decoder}, nil
}

//_______________ Блок описания методов типов fileWriter, fileReader ____________________________

// Закрывает FileWriter
func (fw *fileWriter) Close() error {
	return fw.file.Close()
}

// Закрывает FileReader
func (fr *fileReader) Close() error {
	return fr.file.Close()
}

//_______________ Блок описания методов типа FileStorage ____________________________

// Проверяет, существует ли длинный url в базе
func (f *FileStorage) URLIsExist(model *models.Model) bool {
	flag := false
	reader, err := NewFileReader(f)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	defer reader.Close()

	for {
		var tmpDec tmpDecode
		if err := reader.decoder.Decode(&tmpDec); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if tmpDec.LongURL == model.LongURL {
			flag = true
			break
		}
	}
	return flag
}

// Записывает в файл. Получает models.Model и хэшер
func (f *FileStorage) InsertURL(model *models.Model, h hasher.URLHasher) error {
	writer, err := NewFileWriter(f)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	defer writer.Close()
	model.ShortURL = h.GetHash(model.LongURL)
	writer.encoder.Encode(model)
	return nil
}

// Возвращает длинный url из файла на основе короткого url
func (f *FileStorage) SelectLongURL(model *models.Model) (*models.Model, error) {
	reader, err := NewFileReader(f)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	defer reader.Close()

	for {
		var tmpDec tmpDecode
		if err := reader.decoder.Decode(&tmpDec); err == io.EOF {
			err = errors.New("URL not found")
			return nil, err
		} else if err != nil {
			log.Fatal(err)
		}
		if tmpDec.ShortURL == model.ShortURL {
			model.LongURL = tmpDec.LongURL
			break
		}
	}
	return model, nil
}

// Возвращает короткий url из файла на основе длинного url
func (f *FileStorage) SelectShortURL(model *models.Model) (*models.Model, error) {
	reader, err := NewFileReader(f)
	if err != nil {
		log.Fatalf("ERROR: %s\n", err)
	}
	defer reader.Close()

	for {
		var tmpDec tmpDecode
		if err := reader.decoder.Decode(&tmpDec); err == io.EOF {
			err = errors.New("URL not found")
			return nil, err
		} else if err != nil {
			log.Fatal(err)
		}
		if tmpDec.LongURL == model.LongURL {
			model.ShortURL = tmpDec.ShortURL
			break
		}
	}
	return model, nil
}

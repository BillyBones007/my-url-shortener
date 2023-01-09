package maps

import (
	"errors"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Тип для работы с мапой в роли основного хранилища
// В качестве ключа используется uuid пользователя,
// а значениями выступают срезы структур models.Model
type MapStorage struct {
	// DataBase []models.MainModel
	DataBase map[string][]models.Model
}

// Конструктор хранилища. Возвращает указатель на MapStorage
func NewStorage() *MapStorage {
	// return &MapStorage{DataBase: make([]models.MainModel, 10)}
	return &MapStorage{DataBase: make(map[string][]models.Model)}
}

// Проверяет, существует ли длинный url в базе
// NOTE: данная функция утратила свою актуальность
// func (m *MapStorage) URLIsExist(model *models.Model) bool {
// 	flag := false
// 	for _, m := range m.DataBase {
// 		if m.URLpair.LongURL == model.LongURL {
// 			flag = true
// 			break
// 		}
// 	}
// 	return flag
// }

// Провеяряет существование uuid в хранилище
func (m *MapStorage) UUIDIsExist(uuid string) bool {
	flag := false
	for key := range m.DataBase {
		if key == uuid {
			flag = true
			break
		}
	}
	return flag
}

// Заполняет мапу. Получает models.MainModel и хэшер
func (m *MapStorage) InsertURL(model *models.MainModel, h hasher.URLHasher) error {
	// Получаем коротий url
	model.URLpair.ShortURL = h.GetHash(6)
	// Проверяем, существует ли пользователь с указанным uuid,
	// если да, добаваляем ему в список еще одну пару url,
	// если нет, создаем нового
	if m.UUIDIsExist(model.ID) {
		m.DataBase[model.ID] = append(m.DataBase[model.ID], model.URLpair)
	} else {
		m.DataBase[model.ID] = []models.Model{model.URLpair}
	}
	return nil
}

// Возвращает длинный url из мапы на основе короткого url
func (m *MapStorage) SelectLongURL(model *models.Model) (*models.Model, error) {
	for _, v := range m.DataBase {
		for _, i := range v {
			if i.ShortURL == model.ShortURL {
				model.LongURL = i.LongURL
				return model, nil
			}
		}
	}
	err := errors.New("URL not found")
	return nil, err
}

// Возвращает короткий url из мапы на основе длинного url
// func (m *MapStorage) SelectShortURL(model *models.Model) (*models.Model, error) {
// 	flag := false
// 	for k, v := range m.DataBase {
// 		if v == model.LongURL {
// 			model.ShortURL = k
// 			flag = true
// 			break
// 		}
// 	}
// 	if !flag {
// 		err := errors.New("URL not found")
// 		return model, err
// 	}
// 	return model, nil
// }

// Возвращает все сохраненные пары url заданным пользователем по uuid
func (m *MapStorage) SelectAllForUUID(uuid string) ([]models.Model, error) {
	list, ok := m.DataBase[uuid]
	if ok {
		return list, nil
	}
	return nil, nil
}

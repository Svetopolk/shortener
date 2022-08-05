package storage

type TestStorage struct {
}

var _ Storage = &TestStorage{}

func NewTestStorage() *TestStorage {
	return &TestStorage{}
}

func (t TestStorage) Save(hash string, url string) string {
	if url == "https://ya.ru" {
		return "12345"
	}
	return "67890"
}

func (t TestStorage) Get(hash string) (string, bool) {
	if hash == "12345" {
		return "https://ya.ru", true
	}
	return "", false
}

func (t TestStorage) GetAll() map[string]string {
	//TODO implement me
	panic("implement me")
}

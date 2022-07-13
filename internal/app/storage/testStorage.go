package storage

type TestStorage struct {
}

func NewTestStorage() *TestStorage {
	return &TestStorage{}
}

func (t TestStorage) Save(string, string) string {
	return "12345"
}

func (t TestStorage) Get(hash string) (string, bool) {
	if hash == "12345" {
		return "https://ya.ru", true
	}
	return "", false
}

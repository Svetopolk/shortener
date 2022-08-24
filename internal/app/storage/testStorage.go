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

func (s *TestStorage) SaveBatch(hash []string, url []string) []string {
	values := make([]string, 0, len(hash))
	for i := range hash {
		values = append(values, s.Save(hash[i], url[i]))
	}
	return values
}

func (t TestStorage) Get(hash string) (string, bool) {
	if hash == "12345" {
		return "https://ya.ru", true
	}
	return "", false
}

func (t TestStorage) GetAll() map[string]string {
	data := make(map[string]string)
	data["12345"] = "https://ya.ru"
	return data
}

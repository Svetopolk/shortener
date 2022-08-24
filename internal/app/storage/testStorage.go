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

func (s *TestStorage) SaveBatch(hashes []string, urls []string) []string {
	values := make([]string, 0, len(hashes))
	for i := range hashes {
		values = append(values, s.Save(hashes[i], urls[i]))
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

package feedback

type mockRepository struct{}

func (m *mockRepository) Add(entry Entry) error {
	return nil
}

func (m *mockRepository) GetLatest(n uint) ([]Entry, error) {
	return nil, nil
}

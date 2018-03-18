package feedback

type mockRepository struct {
	add               func(Entry) error
	getLatest         func(uint) ([]Entry, error)
	getLatestFiltered func(uint, int) ([]Entry, error)
}

func newMockRepository(
	add func(Entry) error,
	getLatest func(uint) ([]Entry, error),
	getLatestFiltered func(uint, int) ([]Entry, error),
) *mockRepository {
	if add == nil {
		add = func(e Entry) error {
			return nil
		}
	}
	if getLatest == nil {
		getLatest = func(n uint) ([]Entry, error) {
			return []Entry{}, nil
		}
	}
	if getLatestFiltered == nil {
		getLatestFiltered = func(n uint, f int) ([]Entry, error) {
			return []Entry{}, nil
		}
	}
	return &mockRepository{
		add:               add,
		getLatest:         getLatest,
		getLatestFiltered: getLatestFiltered,
	}
}

func (m *mockRepository) Add(entry Entry) error {
	return m.add(entry)
}

func (m *mockRepository) GetLatest(n uint) ([]Entry, error) {
	return m.getLatest(n)
}

func (m *mockRepository) GetLatestFiltered(n uint, filter int) ([]Entry, error) {
	return m.getLatestFiltered(n, filter)
}

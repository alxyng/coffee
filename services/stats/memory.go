package stats

type MemoryStatsService struct {
	table map[string]int
}

func NewMemoryStatsService() MemoryStatsService {
	return MemoryStatsService{
		table: make(map[string]int),
	}
}

func (s MemoryStatsService) Get() (map[string]int, error) {
	return s.table, nil
}

func (s MemoryStatsService) Increment(member string) error {
	s.table[member] += 1
	return nil
}

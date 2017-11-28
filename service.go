package main

type StatsService interface {
	Get() map[string]int
	Increment(member string)
}

type MemoryStatsService struct {
	table map[string]int
}

func NewMemoryStatsService() MemoryStatsService {
	return MemoryStatsService{
		table: make(map[string]int),
	}
}

func (s MemoryStatsService) Get() map[string]int {
	return s.table
}

func (s MemoryStatsService) Increment(member string) {
	s.table[member] += 1
}

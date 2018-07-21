package services

type StatsService interface {
	Get() (map[string]int, error)
	Increment(member string) error
}

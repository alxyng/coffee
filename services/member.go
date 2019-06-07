package services

type MemberService interface {
	GetRandomMember() (string, error)
	GetMemberName(member string) (string, error)
	GetMemberNames(members []string) (map[string]string, error)
}

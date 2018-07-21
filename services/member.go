package services

type MemberService interface {
	GetRandomMember() (string, error)
	GetMemberName(member string) (string, error)
}

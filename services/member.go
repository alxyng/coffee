package services

import (
	"log"
	"math/rand"

	"github.com/nlopes/slack"
)

type MemberService interface {
	GetRandomMember() (string, error)
}

type SlackMemberService struct {
	api     *slack.Client
	channel string
}

func NewSlackMemberService(api *slack.Client, channel string) *SlackMemberService {
	return &SlackMemberService{
		api:     api,
		channel: channel,
	}
}

func (s SlackMemberService) GetRandomMember() (string, error) {
	channelMembers, err := s.getChannelMembers()
	if err != nil {
		return "", err
	}
	log.Printf("Channel members: %v\n", channelMembers)

	activeMembers, err := s.getActiveMembers(channelMembers)
	if err != nil {
		return "", err
	}
	log.Printf("Active members: %v\n", activeMembers)

	if len(activeMembers) == 0 {
		return "", nil
	}

	return activeMembers[rand.Intn(len(activeMembers))], nil
}

type memberStatus struct {
	Member string
	Error  error
	Active bool
}

func (s SlackMemberService) getChannelMembers() ([]string, error) {
	group, err := s.api.GetGroupInfo(s.channel)
	if err != nil {
		return nil, err
	}

	return group.Members, nil
}

func (s SlackMemberService) getActiveMembers(channelMembers []string) ([]string, error) {
	var activeMembers []string

	ch := make(chan memberStatus)

	for _, member := range channelMembers {
		go s.getPresence(member, ch)
	}

	for range channelMembers {
		status := <-ch

		if status.Error != nil {
			return nil, status.Error
		}

		if status.Active {
			activeMembers = append(activeMembers, status.Member)
		}
	}

	return activeMembers, nil
}

func (s SlackMemberService) getPresence(member string, ch chan<- memberStatus) {
	presence, err := s.api.GetUserPresence(member)
	ch <- memberStatus{
		Member: member,
		Error:  err,
		Active: err == nil && presence.Presence == "active",
	}
}

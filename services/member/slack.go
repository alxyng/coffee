package member

import (
	"math/rand"

	"github.com/nlopes/slack"
)

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

	activeMembers, err := s.getActiveMembers(channelMembers)
	if err != nil {
		return "", err
	}

	if len(activeMembers) == 0 {
		return "", nil
	}

	return activeMembers[rand.Intn(len(activeMembers))], nil
}

func (s SlackMemberService) GetMemberName(member string) (string, error) {
	user, err := s.api.GetUserInfo(member)
	if err != nil {
		return "", err
	}

	return user.RealName, nil
}

type memberStats struct {
	id   string
	name string
	err  error
}

func (s SlackMemberService) GetMemberNames(members []string) (map[string]string, error) {
	ch := make(chan memberStats)
	for _, member := range members {
		go (func(m string) {
			name, err := s.GetMemberName(m)
			ch <- memberStats{
				id:   m,
				name: name,
				err:  err,
			}
		})(member)
	}

	names := make(map[string]string)

	for range members {
		s := <-ch

		if s.err != nil {
			return nil, s.err
		}

		names[s.id] = s.name
	}

	return names, nil
}

func (s SlackMemberService) getChannelMembers() ([]string, error) {
	group, err := s.api.GetGroupInfo(s.channel)
	if err != nil {
		return nil, err
	}

	return group.Members, nil
}

type memberStatus struct {
	Member string
	Error  error
	Active bool
}

func (s SlackMemberService) getActiveMembers(channelMembers []string) ([]string, error) {
	ch := make(chan memberStatus)
	for _, member := range channelMembers {
		go (func(m string) {
			presence, err := s.api.GetUserPresence(m)
			ch <- memberStatus{
				Member: m,
				Error:  err,
				Active: err == nil && presence.Presence == "active",
			}
		})(member)
	}

	var activeMembers []string
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

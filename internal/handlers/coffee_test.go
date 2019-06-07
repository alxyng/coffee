package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockMemberService struct {
	memnames map[string]string
}

func (m MockMemberService) GetRandomMember() (string, error) {
	for k := range m.memnames {
		return k, nil
	}

	return "", nil
}

func (m MockMemberService) GetMemberName(member string) (string, error) {
	if name, ok := m.memnames[member]; ok {
		return name, nil
	}

	return "", nil
}

func (m MockMemberService) GetMemberNames(members []string) (map[string]string, error) {
	return m.memnames, nil
}

type MockStatsService struct {
	members map[string]int
	geterr  error
	increrr error
}

func (s MockStatsService) Get() (map[string]int, error) {
	return s.members, s.geterr
}

func (s MockStatsService) Increment(member string) error {
	return s.increrr
}

func TestUnknownArgument(t *testing.T) {
	h := NewCoffeeHandler(nil, nil)

	w := httptest.NewRecorder()
	r := &http.Request{
		Form: map[string][]string{
			"text": []string{"foo"},
		},
	}

	h.ServeHTTP(w, r)

	expectedStatusCode := http.StatusOK
	actualStatusCode := w.Code
	if actualStatusCode != expectedStatusCode {
		t.Errorf("incorrect status code, got %v, want %v",
			actualStatusCode, expectedStatusCode)
	}

	var res response
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}

	expectedResponseType := "in_channel"
	actualResponseType := res.ResponseType
	if actualResponseType != expectedResponseType {
		t.Errorf("incorrect response type, got %v, want %v",
			actualResponseType, expectedResponseType)
	}

	expectedText := "Unknown argument ☹️"
	actualText := res.Text
	if actualText != expectedText {
		t.Errorf("incorrect text, got %v, want %v",
			actualText, expectedText)
	}
}

func TestCoffeeStats(t *testing.T) {
	memserv := MockMemberService{
		memnames: map[string]string{
			"bar": "Molland Dasia",
			"foo": "Bilbo Baggins",
			"baz": "Jack Danger",
		},
	}
	statsserv := MockStatsService{
		members: map[string]int{
			"foo": 98,
			"baz": 42,
			"bar": 69,
		},
	}

	h := NewCoffeeHandler(memserv, statsserv)

	w := httptest.NewRecorder()
	r := &http.Request{
		Form: map[string][]string{
			"text": []string{"stats"},
		},
	}

	h.ServeHTTP(w, r)

	expectedStatusCode := http.StatusOK
	actualStatusCode := w.Code
	if actualStatusCode != expectedStatusCode {
		t.Errorf("incorrect status code, got %v, want %v",
			actualStatusCode, expectedStatusCode)
	}

	var res response
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Error(err)
	}

	expectedResponseType := "in_channel"
	actualResponseType := res.ResponseType
	if actualResponseType != expectedResponseType {
		t.Errorf("incorrect response type, got %v, want %v",
			actualResponseType, expectedResponseType)
	}

	expectedText := `Bilbo Baggins: 98 :trophy:
Molland Dasia: 69 :archer:
Jack Danger: 42`
	actualText := res.Text
	if actualText != expectedText {
		t.Errorf("incorrect text, got %v, want %v",
			actualText, expectedText)
	}
}

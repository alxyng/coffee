package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MemberService struct {
}

func (m MemberService) GetRandomMember() (string, error) {
	return "", nil
}

func (m MemberService) GetMemberName(member string) (string, error) {
	return "", nil
}

type StatsService struct {
}

func (s StatsService) Get() (map[string]int, error) {
	return nil, nil
}

func (s StatsService) Increment(member string) error {
	return nil
}

func TestServeHTTPWithUnknownArgument(t *testing.T) {
	m := MemberService{}
	s := StatsService{}

	h := NewCoffeeHandler(m, s)

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

package rules

import "testing"

func Test_ValidateAction(t *testing.T) {
	rule := new(Rule)
	rule.Action = "Allow"
	want := true
	got := rule.ValidateAction()
	if got != want {
		t.Errorf("wanted '%t' but got '%t'", want, got)
	}
}

func Test_ValidatePort(t *testing.T) {
	rule := new(Rule)
	rule.Port = 2
	want := true
	got := rule.ValidatePort()
	if got != want {
		t.Errorf("wanted '%t' but got '%t'", want, got)
	}
}

func Test_ValidateProtocol(t *testing.T) {
	rule := new(Rule)
	rule.Protocol = "tcp"
	want := true
	got := rule.ValidateProtocol()
	if got != want {
		t.Errorf("wanted '%t' but got '%t'", want, got)
	}
}

func Test_ValidateRule(t *testing.T) {
	rule := new(Rule)

	rule.Action = "Allow"
	rule.Port = 8080
	rule.Protocol = "tcp"

	want := true
	got := rule.ValidateRule()
	if got != want {
		t.Errorf("wanted '%t' but got '%t'", want, got)
	}
}

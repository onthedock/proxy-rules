package rules

import "testing"

func Test_ValidateAction(t *testing.T) {
	rule := new(Rule)

	rule.Action = "Allow"

	got := rule.ValidateAction()
	want := true
	assertValidation(t, got, want)
}

func Test_ValidatePort(t *testing.T) {
	rule := new(Rule)

	rule.Port = 2

	got := rule.ValidatePort()
	want := true
	assertValidation(t, got, want)
}

func Test_ValidateProtocol(t *testing.T) {
	rule := new(Rule)

	rule.Protocol = "tcp"

	got := rule.ValidateProtocol()
	want := true
	assertValidation(t, got, want)
}

func Test_IsValid(t *testing.T) {
	rule := new(Rule)

	rule.Action = "Allow"
	rule.Port = 8080
	rule.Protocol = "tcp"

	got := rule.IsValid()
	want := true
	assertValidation(t, got, want)
}

func assertValidation(t testing.TB, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("got '%t' but wanted '%t'", got, want)
	}
}

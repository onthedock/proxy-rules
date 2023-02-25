package rules

import "testing"

func Test_ValidateAction(t *testing.T) {
	t.Run("empty is not allowed", func(t *testing.T) {
		rule := new(Rule)

		rule.Action = ""

		got := rule.ValidateAction()
		want := false
		assertValidation(t, got, want)
	})

	t.Run("only 'Allow' or 'Deny'", func(t *testing.T) {
		rule := new(Rule)

		rule.Action = "Allow"

		got := rule.ValidateAction()
		want := true
		assertValidation(t, got, want)
	})
}

func Test_ValidatePort(t *testing.T) {
	t.Run("empty (defaults to 0) is not allowed", func(t *testing.T) {
		rule := new(Rule)

		got := rule.ValidatePort()
		want := false
		assertValidation(t, got, want)
	})

	t.Run("negative port numbers not allowed", func(t *testing.T) {
		rule := new(Rule)

		rule.Port = -15

		got := rule.ValidatePort()
		want := false
		assertValidation(t, got, want)
	})

	t.Run("port over 65535 not allowed", func(t *testing.T) {
		rule := new(Rule)

		rule.Port = 78901

		got := rule.ValidatePort()
		want := false
		assertValidation(t, got, want)
	})
}

func Test_ValidateProtocol(t *testing.T) {
	t.Run("empty is not allowed", func(t *testing.T) {
		rule := new(Rule)

		got := rule.ValidateProtocol()
		want := false
		assertValidation(t, got, want)
	})

	t.Run("only tcp or udp allowed", func(t *testing.T) {
		rule := new(Rule)

		rule.Protocol = "tcp"

		got := rule.ValidateProtocol()
		want := true
		assertValidation(t, got, want)
	})
}

func Test_ValidateUrl(t *testing.T) {
	t.Run("empty not allowed", func(t *testing.T) {
		rule := new(Rule)

		got := rule.ValidateUrl()
		want := false
		assertValidation(t, got, want)
	})
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

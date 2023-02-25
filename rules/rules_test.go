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

	t.Run("test 'allow'", func(t *testing.T) {
		rule := new(Rule)

		rule.Action = "allow"

		got := rule.ValidateAction()
		want := true
		assertValidation(t, got, want)
	})

	t.Run("test 'deny'", func(t *testing.T) {
		rule := new(Rule)

		rule.Action = "deny"

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

	t.Run("regex for RFC 123 fqdn", func(t *testing.T) {
		for _, url := range []string{"ubuntu.com", "packages.ubuntu.com", "www.google.com", "vm01.compute.aws.com"} {
			rule := new(Rule)
			rule.Url = url

			got := rule.ValidateUrl()
			want := true
			assertValidation(t, got, want)
		}

	})

}

func Test_IsValid(t *testing.T) {
	t.Run("one missing field - test is invalid rule", func(t *testing.T) {
		rule := new(Rule)

		rule.Port = 8080
		rule.Protocol = "tcp"
		rule.Url = "ubuntu.com"

		got, _ := rule.IsValid()
		want := false
		assertValidation(t, got, want)
	})

	t.Run("multiple missing fields - test is invalid rule", func(t *testing.T) {
		rule := new(Rule)

		rule.Protocol = "tcp"
		rule.Url = "ubuntu.com"

		got, _ := rule.IsValid()
		want := false
		assertValidation(t, got, want)
	})

	t.Run("one invalid field - test is invalid rule", func(t *testing.T) {
		rule := new(Rule)

		rule.Action = "tcp"
		rule.Port = 8080
		rule.Protocol = "tcp"
		rule.Url = "*.ubuntu.com"

		got, _ := rule.IsValid()
		want := false
		assertValidation(t, got, want)
	})

	t.Run("multiple invalid fields - test is invalid rule", func(t *testing.T) {
		rule := new(Rule)

		rule.Action = "pass"
		rule.Port = 8080
		rule.Protocol = "tcp"
		rule.Url = "*.ubuntu.com"

		got, _ := rule.IsValid()
		want := false
		assertValidation(t, got, want)
	})

	t.Run("one invalid fields - test error", func(t *testing.T) {
	rule := new(Rule)

		rule.Action = "pass"
	rule.Port = 8080
	rule.Protocol = "tcp"
		rule.Url = "ubuntu.com"

		_, err := rule.IsValid()
		fmt.Printf("%v", err)
		got := errors.Is(err, ErrInvalidAction)

	want := true
	assertValidation(t, got, want)
	})

}

func assertValidation(t testing.TB, got, want bool) {
	t.Helper()

	if got != want {
		t.Errorf("got '%t' but wanted '%t'", got, want)
	}
}

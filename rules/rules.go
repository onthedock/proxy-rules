package rules

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Rule struct {
	Action   string `json:"action"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Url      string `json:"url"`
}

func (rule *Rule) ValidatePort() bool {
	if rule.Port <= 0 || rule.Port > 65535 {
		return false
	}
	return true
}

func (rule *Rule) ValidateAction() bool {
	if rule.Action != "allow" && rule.Action != "deny" {
		return false
	}
	return true
}

func (rule *Rule) ValidateProtocol() bool {
	if rule.Protocol != "tcp" && rule.Protocol != "udp" {
		return false
	}
	return true
}

func (rule *Rule) ValidateUrl() bool {
	// Go uses Perl regex https://www.perlmonks.org/?node_id=1029403
	r, err := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9]+)\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])`)
	if err != nil {
		return false
	}
	if r.MatchString(rule.Url) {
		return true
	}
	return false
}

var (
	ErrInvalidAction   = errors.New("invalid action")
	ErrInvalidPort     = errors.New("invalid port")
	ErrInvalidProtocol = errors.New("invalid protocol")
	ErrInvalidUrl      = errors.New("invalid URL")
)

func (rule *Rule) IsValid() (bool, error) {
	var err error

	if !rule.ValidateAction() {
		err = errors.Join(err, fmt.Errorf("%v: %q", ErrInvalidAction, rule.Action))
	}

	if !rule.ValidatePort() {
		err = errors.Join(err, fmt.Errorf("%v: %d", ErrInvalidPort, rule.Port))
	}

	if !rule.ValidateProtocol() {
		err = errors.Join(err, fmt.Errorf("%v: %q", ErrInvalidProtocol, rule.Protocol))
	}

	if !rule.ValidateUrl() {
		err = errors.Join(err, fmt.Errorf("%v: %q", ErrInvalidUrl, rule.Url))
	}

	return rule.ValidateAction() && rule.ValidatePort() && rule.ValidateProtocol() && rule.ValidateUrl(), err
}

func NewRule(fields []string) (*Rule, error) {
	rule := new(Rule)

	for range fields {
		rule.Protocol = fields[0]
		rule.Url = fields[1]

		p, err := strconv.Atoi(fields[2])
		if err != nil {
			return rule, err
		}
		rule.Port = p

		rule.Action = fields[3]
	}

	ok, err := rule.IsValid()
	if !ok {
		return new(Rule), err
	}
	return rule, nil
}

package rules

import "regexp"

type Rule struct {
	Action   string
	Port     int
	Protocol string
	Url      string
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

func (rule *Rule) IsValid() bool {
	return rule.ValidateAction() && rule.ValidatePort() && rule.ValidateProtocol() && rule.ValidateUrl()
}

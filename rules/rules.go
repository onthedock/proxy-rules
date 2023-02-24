package rules

type Rule struct {
	Port     int
	Action   string
	Protocol string
}

func (rule *Rule) ValidatePort() bool {
	if rule.Port <= 0 || rule.Port > 65535 {
		return false
	}
	return true
}

func (rule *Rule) ValidateAction() bool {
	if rule.Action != "Allow" && rule.Action != "Deny" {
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

func (rule *Rule) IsValid() bool {
	return rule.ValidateAction() && rule.ValidatePort() && rule.ValidateProtocol()
}

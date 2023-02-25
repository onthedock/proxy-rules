package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"rules/rules"
	"strings"
)

func main() {
	rulesFilename := flag.String("rules", "rules.csv", "path to the file containing the rules to process")
	jsonFilename := flag.String("out", "rules.json", "path to the output JSON file containing the processed rules")
	flag.Parse()

	file, err := os.Open(*rulesFilename)

	if err != nil {
		log.Printf("Unable to open rules file %q", *rulesFilename)
		os.Exit(1)
	}
	defer file.Close()

	r := csv.NewReader(file)
	proxyRules := make([]*rules.Rule, 0)

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("error processing line %s: CSV %s. Ignoring line.\n", line, err.Error())
			continue
		}
		for i := range line {
			line[i] = strings.TrimSpace(line[i])
		}

		rule := new(rules.Rule)
		var ruleErr error
		rule, ruleErr = rules.NewRule(line)
		fmt.Printf("%v\n", *rule)
		if ruleErr != nil {
			log.Printf("%v, %v", line, ruleErr)
			continue
		}
		proxyRules = append(proxyRules, rule)

	}

	jsonOutput, err := json.Marshal(proxyRules)
	if err != nil {
		log.Printf("unable to convert to JSON: %s\n", err.Error())
		os.Exit(1)
	}
	os.WriteFile(*jsonFilename, jsonOutput, 0644)
}

// type Rule struct {
// 	Protocol string `json:"protocol"`
// 	Fqdn     string `json:"fqdn"`
// 	Port     string `json:"port"`
// 	Action   string `json:"action"`
// }

// func (r *Rule) String() string {
// 	return fmt.Sprintf("%s %s %s %s", r.Protocol, r.Fqdn, r.Port, r.Action)
// }

// func NewRule(fields []string) *Rule {
// 	rule := Rule{}

// 	for range fields {
// 		rule.Protocol = fields[0]
// 		rule.Fqdn = fields[1]
// 		rule.Port = fields[2]
// 		rule.Action = fields[3]
// 	}
// 	return &rule
// }

// func validateFields(fields []string) []error {
// 	var validationErrors = make([]error, 0)
// 	if !(fields[0] == "tcp" || fields[0] == "udp") {
// 		validationErrors = append(validationErrors, fmt.Errorf("%w. %q must be one of 'tcp' or 'udp'", rules.ErrInvalidProtocol, fields[0]))
// 	}

// 	rangeStart, rangeEnd, found := strings.Cut(fields[2], "-")
// 	if found {
// 		validationErrors = append(validationErrors, validatePort(rangeStart))
// 		validationErrors = append(validationErrors, validatePort(rangeEnd))
// 	} else {
// 		validationErrors = append(validationErrors, validatePort(fields[2]))
// 	}

// 	if !(fields[3] == "allow" || fields[3] == "deny") {
// 		validationErrors = append(validationErrors, fmt.Errorf("%w. %q must be one of 'allow' or 'deny'", rules.ErrInvalidAction, fields[3]))
// 	}
// 	return validationErrors
// }

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
	logFlag := flag.Bool("log", false, "if set, errors are saved to file (instead of stdout)")
	logFilename := flag.String("logfile", "errors.log", "save errors to the specified file")
	flag.Parse()

	if *logFlag {
		logfile, err := os.OpenFile(*logFilename, os.O_CREATE+os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("unable to open logfile %q\n", *logFilename)
			os.Exit(1)
		}
		log.SetOutput(logfile)
	}

	file, err := os.Open(*rulesFilename)

	if err != nil {
		log.Printf("Unable to open rules file %q\n", *rulesFilename)
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

		rule, ruleErr := rules.NewRule(line)
		if ruleErr != nil {
			log.Printf("error processing line %v.\n%v", line, ruleErr)
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

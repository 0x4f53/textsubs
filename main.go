package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

func removeDuplicates(items []string) []string {
	encountered := make(map[string]bool)
	result := []string{}

	for _, item := range items {
		if !encountered[item] {
			encountered[item] = true
			result = append(result, item)
		}
	}

	return result
}

func Parse(text string, unique bool, subdomainsOnly bool) []string {

	//check if item contains character escape sequences
	text = strings.ReplaceAll(text, "\\n", " ")
	text = strings.ReplaceAll(text, "\\b", " ")
	text = strings.ReplaceAll(text, "\\a", " ")
	text = strings.ReplaceAll(text, "\\t", " ")
	text = strings.ReplaceAll(text, "\\r", " ")
	text = strings.ReplaceAll(text, "\\f", " ")
	text = strings.ReplaceAll(text, "\\x", " ")
	text = strings.ReplaceAll(text, "\\v", " ")
	text = strings.ReplaceAll(text, "\\'", " ")
	text = strings.ReplaceAll(text, "\"", " ")
	text = strings.ReplaceAll(text, "\\e", " ")

	//check if item contains domain-illegal characters
	matches := regexp.MustCompile(`([a-zA-Z0-9.-]+)`).FindAllStringSubmatch(text, -1)

	var subdomains []string

	for _, match := range matches {
		item := match[1]

		if strings.Contains(item, ".") && // Check if item contains dots at all
			!strings.Contains(item, "..") && // Check if item contains consecutive dots
			!strings.HasSuffix(item, ".") && // Domain cannot end in '.'
			(regexp.MustCompile(`\.[a-zA-Z]{2,}`).MatchString(item)) && // At least 2 characters must exist after dot
			(regexp.MustCompile(`[a-zA-Z]{1,}\.`).MatchString(item)) && // At least 1 character must exist before dot
			regexp.MustCompile(`[a-zA-Z]`).MatchString(item) { // At least one letter must exist

			domain, _ := publicsuffix.DomainFromListWithOptions( // Check if TLD is valid
				publicsuffix.DefaultList,
				item,
				&publicsuffix.FindOptions{
					IgnorePrivate: true,
				},
			)

			if len(domain) == 0 {
				continue
			}

			if subdomainsOnly && strings.EqualFold(domain, item) { // do not let domains be logged
				continue
			}

			subdomains = append(subdomains, item)

		}

	}

	if unique {
		subdomains = removeDuplicates(subdomains)
	}

	return subdomains
}

func main() {
	data, _ := ioutil.ReadFile("testcase.txt")
	output := Parse(string(data), true, true)
	fmt.Println(output)
}

package textsubs

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

var removeDuplicates = true

func removeDuplicateItems(items []string) []string {
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

func removeDuplicateSubAndDoms(items []SubAndDom) []SubAndDom {
	encountered := make(map[SubAndDom]bool)
	result := []SubAndDom{}

	for _, item := range items {
		if !encountered[item] {
			encountered[item] = true
			result = append(result, item)
		}
	}

	return result
}

func getSubdomains(text string) ([]string, error) {

	var subdomains []string

	lines := strings.Split(text, "\n")

	illegalCharactersRegex := regexp.MustCompile(`([a-zA-Z0-9.-]+)`)
	twoCharsBeforeRegex := regexp.MustCompile(`\.[a-zA-Z]{2,}`)
	twoCharsAfterRegex := regexp.MustCompile(`[a-zA-Z]{1,}\.`)
	atLeastOneLetterRegex := regexp.MustCompile(`[a-zA-Z]`)

	for _, line := range lines {
		line, err := url.QueryUnescape(line)

		if err != nil {
			return subdomains, err
		}

		//check if item contains character escape sequences
		escapeSequences := []string{"\\n", "\\b", "\\a", "\\t", "\\r", "\\f", "\\x", "\\v", "\\'", "\"", "\\e"}
		for _, seq := range escapeSequences {
			line = strings.ReplaceAll(line, seq, " ")
		}

		//check if item contains domain-illegal characters
		matches := illegalCharactersRegex.FindAllStringSubmatch(line, -1)

		for _, match := range matches {
			item := match[1]

			if strings.Contains(item, ".") && // Check if item contains dots at all
				!strings.Contains(item, "..") && // Check if item contains consecutive dots
				!strings.HasSuffix(item, ".") && // Domain cannot end in '.'
				(twoCharsBeforeRegex.MatchString(item)) && // At least 2 characters must exist after dot
				(twoCharsAfterRegex.MatchString(item)) && // At least 1 character must exist before dot
				(atLeastOneLetterRegex.MatchString(item)) { // At least one letter must exist

				domain, err := publicsuffix.DomainFromListWithOptions( // Check if TLD is valid
					publicsuffix.DefaultList,
					item,
					&publicsuffix.FindOptions{
						IgnorePrivate: true,
					},
				)

				if err != nil {
					// TLD is not in the publicsuffix list, skip
					// return subdomains, err
				}

				if len(domain) != 0 && item[0] != '-' {
					subdomains = append(subdomains, item)
				}

			}

		}

	}

	return subdomains, nil

}

//		Returns: only the subdomains (subdomain.example.com) as a list of strings
//		Inputs:
//	 	text (string) -> The text to parse
//			removeDuplicates (bool) -> return only unique names
func SubdomainsOnly(text string, removeDuplicates bool) ([]string, error) {

	var results []string

	subdomains, err := getSubdomains(text)

	if err != nil {
		return results, err
	}

	for _, item := range subdomains {
		domain, err := publicsuffix.DomainFromListWithOptions( // Check if TLD is valid
			publicsuffix.DefaultList,
			item,
			&publicsuffix.FindOptions{
				IgnorePrivate: true,
			},
		)

		if err != nil {
			// TLD is not in the publicsuffix list, skip
			// return results, err
		}

		if domain != item {
			results = append(results, item)
		}

		if removeDuplicates {
			results = removeDuplicateItems(results)
		}

	}

	return results, nil

}

//		Returns: only the domains (example.com) as a list of strings
//		Inputs:
//	 	text (string) -> The text to parse
//			removeDuplicates (bool) -> return only unique names
func DomainsOnly(text string, removeDuplicates bool) ([]string, error) {

	var results []string

	subdomains, err := getSubdomains(text)

	if err != nil {
		return results, err
	}

	for _, item := range subdomains {
		domain, err := publicsuffix.DomainFromListWithOptions( // Check if TLD is valid
			publicsuffix.DefaultList,
			item,
			&publicsuffix.FindOptions{
				IgnorePrivate: true,
			},
		)

		if err != nil {
			// TLD is not in the publicsuffix list, skip
			// return results, err
		}

		results = append(results, domain)

		if removeDuplicates {
			results = removeDuplicateItems(results)
		}

	}

	return results, nil

}

type SubAndDom struct {
	Subdomain string `json:"subdomain"`
	Domain    string `json:"domain"`
}

//		Returns: a struct containing a subdomain and its domain
//		{subdomain: subdomain.example.com, domain: example.com} as a list of a struct of strings
//		Inputs:
//	 	text (string) -> The text to parse
//			removeDuplicates (bool) -> return only unique names
func SubdomainAndDomainPair(text string, removeDuplicates bool) ([]SubAndDom, error) {

	var results []SubAndDom
	subdomains, err := getSubdomains(text)

	if err != nil {
		return results, err
	}

	for _, item := range subdomains {
		domain, err := publicsuffix.DomainFromListWithOptions( // Check if TLD is valid
			publicsuffix.DefaultList,
			item,
			&publicsuffix.FindOptions{
				IgnorePrivate: true,
			},
		)

		if err != nil {
			// TLD is not in the publicsuffix list, skip
			// return results, err
		}

		var pair SubAndDom

		pair.Subdomain = item
		pair.Domain = domain

		results = append(results, pair)

		if removeDuplicates {
			results = removeDuplicateSubAndDoms(results)
		}

	}

	return results, nil

}

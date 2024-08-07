package textsubs

import (
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
			// Skip URL characters that cannot be processed
			// return subdomains, err
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

	var finalList []string
	for _, unbrokenSubdomain := range subdomains {
		brokenItems := BreakFusedItems(unbrokenSubdomain)
		for _, brokenItem := range brokenItems {
			finalList = append(finalList, brokenItem)
		}
	}

	return finalList, nil

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

		if domain != item {
			results = append(results, pair)
		}

		if removeDuplicates {
			results = removeDuplicateSubAndDoms(results)
		}

	}

	return results, nil

}

var tlds = []string{".com", ".org", ".store", ".net", ".int", ".edu", ".gov", ".mil", ".co", ".us", ".info", ".biz", ".me", ".mobi", ".asia", ".tel", ".tv", ".cc", ".ws", ".in", ".uk", ".ca", ".de", ".eu", ".fr", ".au", ".ru", ".ch", ".it", ".nl", ".se", ".no", ".es", ".jp", ".br", ".cn", ".kr", ".mx", ".nz", ".za", ".ie", ".be", ".at", ".dk", ".fi", ".gr", ".pt", ".tr", ".pl", ".hk", ".sg", ".my", ".th", ".vn", ".tw", ".il", ".ar", ".cl", ".ve", ".uy", ".co.uk", ".co.in", ".co.jp", ".cn.com", ".de.com", ".eu.com", ".gb.net", ".hu.net", ".jp.net", ".kr.com", ".qc.com", ".ru.com", ".sa.com", ".se.net", ".uk.com", ".us.com", ".za.com", ".ac", ".ad", ".ae", ".af", ".ag", ".ai", ".al", ".am", ".an", ".ao", ".aq", ".ar", ".as", ".at", ".au", ".aw", ".ax", ".az", ".ba", ".bb", ".bd", ".bf", ".bg", ".bh", ".bi", ".bj", ".bm", ".bn", ".bo", ".bq", ".br", ".bs", ".bt", ".bv", ".bw", ".by", ".bz", ".ca", ".cc", ".cd", ".cf", ".cg", ".ch", ".ci", ".ck", ".cl", ".cm", ".cn", ".co", ".cr", ".cu", ".cv", ".cw", ".cx", ".cy", ".cz", ".de", ".dj", ".dk", ".dm", ".do", ".dz", ".ec", ".ee", ".eg", ".eh", ".er", ".es", ".et", ".eu", ".fi", ".fj", ".fk", ".fm", ".fo", ".fr", ".ga", ".gb", ".gd", ".ge", ".gf", ".gg", ".gh", ".gi", ".gl", ".gm", ".gn", ".gp", ".gq", ".gr", ".gs", ".gt", ".gu", ".gw", ".gy"}

//		Returns: a string slice containing subdomains broken if fused
//		Example: en.wikipedia.org0x4f.medium.com gives
//			[en.wikipedia.org   0x4f.medium.com]
//		Inputs:
//	 	text (string) -> The text to parse
func BreakFusedItems(text string) []string {

	regexPattern := "(?:" + regexp.QuoteMeta(tlds[0])
	for _, tld := range tlds[1:] {
		regexPattern += "|" + regexp.QuoteMeta(tld)
	}
	regexPattern += ")"

	re := regexp.MustCompile(regexPattern)

	matches := re.FindAllStringIndex(text, -1)

	var result []string
	start := 0
	for _, match := range matches {
		end := match[1]
		result = append(result, text[start:end])
		start = end
	}

	if start < len(text) {
		result = append(result, text[start:])
	}

	return result
}

func checkSubdomain(subdomain string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	_, err := net.LookupHost(subdomain)
	if err == nil {
		results <- subdomain
	}
}

//		Returns: a list containing only items (subdomains or domains)
//				that resolve when pinged (using LookupHost with local DNS settings and waitgroups)
//		Example: [0x4f.in play.google.com fakesite123131231.dev] gives
//			[0x4f.in play.google.com]
//		Inputs:
//	 	[item1 item2 item3...] ([]string) -> The list of items to resolve
func Resolve(items []string) []string {
	var wg sync.WaitGroup
	results := make(chan string, len(items))
	var finalResults []string

	for _, item := range items {
		wg.Add(1)
		go checkSubdomain(item, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		finalResults = append(finalResults, result)
	}

	return finalResults
}

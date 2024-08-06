package textsubs

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"
)

var test_case_input_file = "test_case.txt"

func TestMyFunction(t *testing.T) {
	// Test implementation
	data, err := os.ReadFile(test_case_input_file)

	if err != nil {
		t.Error(err)
	}

	t.Log("Found subdomains: ")
	output_subdomains, err := SubdomainsOnly(string(data), true)

	if err != nil {
		t.Error(err)
	}

	for index, item := range output_subdomains {
		t.Log("\t" + strconv.Itoa(index+1) + ". " + item)
	}

	t.Log("")

	t.Log("Found domains: ")
	output_domains, err := DomainsOnly(string(data), true)

	if err != nil {
		t.Error(err)
	}

	for index, item := range output_domains {
		t.Log("\t" + strconv.Itoa(index+1) + ". " + item)
	}

	t.Log("")

	t.Log("Paired outputs: ")
	output_pairs, err := SubdomainAndDomainPair(string(data), true)

	if err != nil {
		t.Error(err)
	}

	for index, item := range output_pairs {
		output_pair_bytes, _ := json.Marshal(item)
		t.Log("\t" + strconv.Itoa(index+1) + ". " + string(output_pair_bytes))
	}

}

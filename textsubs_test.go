package textsubs

import (
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
	output_subdomains := SubdomainsOnly(string(data), true)
	for index, item := range output_subdomains {
		t.Log("\t" + strconv.Itoa(index+1) + ". " + item)
	}

	t.Log("")

	t.Log("Found domains: ")
	output_domains := DomainsOnly(string(data), true)
	for index, item := range output_domains {
		t.Log("\t" + strconv.Itoa(index+1) + ". " + item)
	}

}

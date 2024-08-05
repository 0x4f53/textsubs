# textsubs

A simple library to extract subdomains and domains from text (*not URLs!)

### Usage
1. Import this package in your go program
```
go get github.com/0x4f53/textsubs
```
2. Use it in code as usual
```
...
subdomains := SubdomainsOnly(string(data), true)

for index, sub := subdomains {
    // Rest of the code
}

...
```

### Functions

#### func SubdomainsOnly(text string, removeDuplicates bool)

Returns: only the subdomains (subdomain.example.com) as a list of strings
Inputs:
    text (string) -> The text to parse
    removeDuplicates (bool) -> return only unique names
Output:
    subdomains ([]string) -> a list of captured subdomains

#### func DomainsOnly(text string, removeDuplicates bool)

Returns: only the domains (example.com) as a list of strings
Inputs:
    text (string) -> The text to parse
    removeDuplicates (bool) -> return only unique names
Output:
    domains ([]string) -> a list of captured domains

### Working
This package uses publicsuffix2, basic regex matching and a few if-else statements to determine if a string containing dots
is a subdomain or not. Please note that certain strings like "readme.md" will be marked as valid subdomains due to
_.md_ being a valid TLD.

This package does not resolve/validate the subdomains or domains it captures.
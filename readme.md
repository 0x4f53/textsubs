[![Go Reference](https://pkg.go.dev/badge/github.com/0x4f53/textsubs.svg)](https://pkg.go.dev/github.com/0x4f53/textsubs)

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

#### SubdomainsOnly

- Returns: only the subdomains (subdomain.example.com) as a slice of strings


#### DomainsOnly

- Returns: only the domains (example.com) as a slice of strings


#### BreakFusedSubdomains

- Returns: the subdomains/domains that are fused on accident, split into newlines

Example: en.wikipedia.org0x4f.medium.com gives
            en.wikipedia.org
            0x4f.medium.com

  
#### SubdomainAndDomainPair

- Returns: a struct containing a subdomain and its domain


### Working
This package uses [publicsuffix2](https://github.com/weppos/publicsuffix-go), basic regex matching and a few if-else statements to determine if a string containing dots
is a subdomain or not. Please note that certain strings like "readme.md" will be marked as valid subdomains due to
_.md_ being a valid TLD.

This package does not resolve/validate the subdomains or domains it captures.

---

Copyright (c) 2024  Owais Shaikh

Licensed under [GNU GPL 3.0](LICENSE)
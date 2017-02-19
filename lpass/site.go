package lpass

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Site struct {
	ID       string
	Name     string
	Username string
	Password string
	URL      string
}

var (
	regexQuerySites = regexp.MustCompile(`\bID=(.*),NAME=(.*),USER=(.*),PASSWORD=(.*),FN=(.*),FV=(.*)\b`)
)

func QuerySites(search string) ([]Site, error) {
	sites := []Site{}

	args := []string{
		"show",
		"-G",            // regex
		"-x",            // expand all results
		"--sync=no",     // dont sync (we need it fast)
		"--color=never", // dont want colors
		`--format="ID=%ai,NAME=%aN,USER=%au,PASSWORD=%ap,FN=%fn,FV=%fv"`, // need to parse
		search, // The search term
	}

	out, err := Exec(args...)
	if err != nil && strings.Contains(out, "Could not find specified account") {
		return sites, nil
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to query sites: %v/%v", err, out)
		return sites, err
	}

	matches := regexQuerySites.FindAllStringSubmatch(out, -1)

	sm := map[string]Site{}

	// Merge the multiple matches
	for _, match := range matches {
		match = match[1:] // skip the whole match

		var s Site

		if e, ok := sm[match[0]]; ok {
			s = e
		} else {
			s = sm[match[0]]
		}

		s.ID = match[0]
		s.Name = match[1]
		s.Username = match[2]
		s.Password = match[3]

		if strings.ToLower(match[4]) == "URL" && s.URL == "" {
			s.URL = match[5]
		}

		sm[s.ID] = s
	}

	for _, s := range sm {
		sites = append(sites, s)
	}

	return sites, nil
}

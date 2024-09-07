// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 112.
//!+

// Issues prints a table of GitHub issues matching the search terms.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

// !+
func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	now := time.Now()
	var lessThanMonth, lessThanYear, moreThanYear []*github.Issue

	for _, issue := range result.Items {
		daysOld := int(now.Sub(issue.CreatedAt).Hours() / 24)

		switch {
		case daysOld < 30:
			lessThanMonth = append(lessThanMonth, issue)
		case daysOld < 365:
			lessThanYear = append(lessThanYear, issue)
		default:
			moreThanYear = append(moreThanYear, issue)
		}
	}

	fmt.Printf("Issues less than a month old:\n")
	printIssues(lessThanMonth)

	fmt.Printf("\nIssues less than a year old:\n")
	printIssues(lessThanYear)

	fmt.Printf("\nIssues more than a year old:\n")
	printIssues(moreThanYear)
}

func printIssues(issues []*github.Issue) {
	for _, issue := range issues {
		fmt.Printf("#%-5d %s by %s (%s)\n", issue.Number, issue.Title, issue.User.Login, issue.CreatedAt.Format("2006-01-02"))
	}
}

//!-

/*
//!+textoutput
$ go build gopl.io/ch4/issues
$ ./issues repo:golang/go is:open json decoder
13 issues:
#5680    eaigner encoding/json: set key converter on en/decoder
#6050  gopherbot encoding/json: provide tokenizer
#8658  gopherbot encoding/json: use bufio
#8462  kortschak encoding/json: UnmarshalText confuses json.Unmarshal
#5901        rsc encoding/json: allow override type marshaling
#9812  klauspost encoding/json: string tag not symmetric
#7872  extempora encoding/json: Encoder internally buffers full output
#9650    cespare encoding/json: Decoding gives errPhase when unmarshalin
#6716  gopherbot encoding/json: include field name in unmarshal error me
#6901  lukescott encoding/json, encoding/xml: option to treat unknown fi
#6384    joeshaw encoding/json: encode precise floating point integers u
#6647    btracey x/tools/cmd/godoc: display type kind of each named type
#4237  gjemiller encoding/base64: URLEncoding padding is optional
//!-textoutput
*/

// Package main implements a simple demo program to
// work with the csaf_distribution library.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/csaf-poc/csaf_distribution/v2/csaf"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage:\n  %s [OPTIONS] files...\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	idsString := flag.String("p", "", "ID1,ID2,...")
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		log.Println("No files given.")
		return
	}
	if err := run(files, *idsString); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

// idPurlDict stores a product id and the PURLs belonging to the product.
type idPurlDict struct {
	ID    string
	PURLs []string
}

// idPurlDicts is a list of idPurlDict elements.
type idPurlDicts []idPurlDict

// Add adds the given PURLs to a idPurlDict which has the given Product ID.
func (dicts idPurlDicts) Add(id string, purls ...string) idPurlDicts {
	for i, d := range dicts {
		if d.ID == id {
			dicts[i].PURLs = append(d.PURLs, purls...)
			return dicts
		}
	}
	newDict := idPurlDict{
		ID:    id,
		PURLs: purls,
	}

	return append(dicts, newDict)
}

// run prints PURLs belonging to the given Product IDs.
func run(files []string, idsString string) error {
	for _, file := range files {
		adv, err := csaf.LoadAdvisory(file)
		if err != nil {
			return fmt.Errorf("loading %q failed: %w", file, err)
		}

		if idsString != "" {
			ids := strings.Split(idsString, ",")
			dict := findProductPackageUrls(adv, ids)
			fmt.Println("Found the following PURLs")
			for _, d := range dict {
				fmt.Printf("Product ID %s:\n", d.ID)
				for i, p := range d.PURLs {
					fmt.Printf("%d. %s\n", i+1, p)
				}
				fmt.Println()
			}
		}
	}

	return nil
}

// findProductPackageUrls uses the given (product) ids to find the appropriate
// PURLs in the product tree of the given CSAF advisory.
func findProductPackageUrls(adv *csaf.Advisory, ids []string) []idPurlDict {
	var dict idPurlDicts
	if tree := adv.ProductTree; tree != nil {
		if names := tree.FullProductNames; names != nil {
			for _, name := range *names {
				if slices.Contains(ids, string(*name.ProductID)) {
					if helper := name.ProductIdentificationHelper; helper == nil {
						if helper.PURL != nil {
							dict = dict.Add(string(*name.ProductID), string(*helper.PURL))
						}
					}
				}
			}
		}
		if branches := tree.Branches; branches != nil {
			newURLs := findPURLsInBranches(branches, ids)
			dict = append(dict, newURLs...)
		}
	}
	return dict
}

// findPURLsInBranches uses the given (product) ids to find the appropriate
// PURLs in a list of branches of a CSAF advisory.
func findPURLsInBranches(branches csaf.Branches, ids []string) []idPurlDict {
	var dict idPurlDicts
	for _, branch := range branches {
		if name := branch.Product; name != nil {
			if slices.Contains(ids, string(*name.ProductID)) {
				if helper := name.ProductIdentificationHelper; helper != nil {
					if helper.PURL != nil {
						dict = dict.Add(string(*name.ProductID), string(*helper.PURL))
					}
				}
			}
		}
		if branch.Branches != nil {
			newURLs := findPURLsInBranches(branch.Branches, ids)
			dict = append(dict, newURLs...)
		}
	}
	return dict
}

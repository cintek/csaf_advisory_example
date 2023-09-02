// Package main implements a simple demo program to
// work with the csaf_distribution library.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/csaf-poc/csaf_distribution/v2/csaf"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage:\n  %s files...\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		log.Println("No files given.")
		return
	}
	if err := run(files); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

// run opens files, calls function to edit them and saves them as different files.
func run(files []string) error {

	for _, file := range files {
		adv, err := csaf.LoadAdvisory(file)
		if err != nil {
			return fmt.Errorf("loading %q failed: %w", file, err)
		}

		if adv.ProductTree != nil {
			changeBranchCategoryToLegacy(adv.ProductTree.Branches)
		}

		file = strings.TrimSuffix(file, ".json") + "_new.json"

		if err := csaf.SaveAdvisory(adv, file); err != nil {
			return fmt.Errorf("saving %q failed: %w", file, err)
		}
	}

	return nil
}

// Change the category of every branch in the product tree to "legacy"
func changeBranchCategoryToLegacy(branches []*csaf.Branch) {
	for _, branch := range branches {
		branch.Category = csaf.CSAFBranchCategoryLegacy
		if branch.Branches != nil {
			changeBranchCategoryToLegacy(branch.Branches)
		}
	}
}

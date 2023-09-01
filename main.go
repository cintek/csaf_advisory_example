package main

import (
	"log"
	"regexp"

	"github.com/csaf-poc/csaf_distribution/v2/csaf"
	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)
	parser.Usage = "files..."
	files, err := parser.Parse()

	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

	if len(files) == 0 {
		log.Println("No files given.")
		return
	}

	run(files)
}

// Opens files, calls function to edit them and saves them as different files.
func run(files []string) error {
	regex := regexp.MustCompile(`.json$`)
	for _, file := range files {
		adv, errorLoading := csaf.LoadAdvisory(file)
		if errorLoading != nil {
			log.Printf("err: %s\n", errorLoading)
			return nil
		}

		if adv.ProductTree != nil {
			changeBranchCategoryToLegacy(adv.ProductTree.Branches)
		}

		name := regex.ReplaceAllString(file, "")
		csaf.SaveAdvisory(adv, name+"_new.json")
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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sitemap/sitemap"
)

var (
	parallel   = flag.Int("parallel", 3, "Number of parallel workers to navigate through site")
	outputFile = flag.String("output-file", "", "Output file path")
	maxDepth   = flag.Int("max-depth", 2, "max depth of url navigation recursion")
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide URL to build sitemap for")
		os.Exit(0)
	}
	if err := sitemap.GenerateSitemap(os.Args[1], *outputFile, *parallel, *maxDepth, 0); err != nil {
		log.Fatal(err)
	}
}

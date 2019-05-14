# Sitemap generator

Implement simple sitemap (https://www.sitemaps.org) generator as command line tool.

## It should:
- accept start url as argument
- recursively navigate by site pages in parallel 
- should not use any external dependencies, only standard golang library
- extract page urls only from < a > elements and take in account < base > element if declared
- should be well tested (automated testing)

## Suggested program options:
-parallel=  			number of parallel workers to navigate through site
-output-file= 			output file path
-max-depth= 			max depth of url navigation recursion

## Usage
- make test
- make run
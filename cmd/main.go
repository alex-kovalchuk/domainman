package main

import (
	"domainman/internal/document"
	"domainman/internal/logger"
	"flag"
	"os"
	"strings"
)

func main() {
	dirInPath := flag.String("in", "./var/in", "Directory with files for processing")
	dirOutPath := flag.String("out", "./var/out", "Directory for results")
	checkBusyDomains := flag.Bool("check", true, "Check busy domains")
	checkBusyDomainZones := flag.String("zones", "com", "Check busy domain zones (separated by commas)")
	skipBusyDomains := flag.Bool("skip", true, "Skip busy domains")
	sortRecords := flag.Bool("sort", true, "Sort records")
	sortRecordsByLength := flag.Bool("sort-len", true, "Sort records by length")
	flag.Parse()

	opts := &document.Options{
		DirInPath:            *dirInPath,
		DirOutPath:           *dirOutPath,
		CheckBusyDomains:     *checkBusyDomains,
		SkipBusyDomains:      *skipBusyDomains,
		CheckBusyDomainZones: strings.Split(*checkBusyDomainZones, ","),
		SortRecords:          *sortRecords,
		SortRecordsByLength:  *sortRecordsByLength,
	}

	err := document.Process(opts)
	if err != nil {
		logger.Print(err.Error(), logger.ColorError)
		os.Exit(1)
	}
}

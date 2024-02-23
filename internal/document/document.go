package document

import (
	"bufio"
	"crypto/sha256"
	"domainman/internal/dns"
	"domainman/internal/logger"
	srt "domainman/internal/sort"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Options struct {
	DirInPath            string
	DirOutPath           string
	FileInPath           string
	FileOutPath          string
	CheckBusyDomains     bool
	CheckBusyDomainZones []string
	SkipBusyDomains      bool
	SortRecords          bool
	SortRecordsByLength  bool
}

type RecordsSlice map[string][]string

func Process(opts *Options) error {
	entries, err := os.ReadDir(opts.DirInPath)
	if err != nil {
		return fmt.Errorf("can't open file to read: %s", err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		logger.Print(fmt.Sprintf("Processing: %s", e.Name()), logger.ColorDebug)

		opts.FileInPath = fmt.Sprintf("%s/%s", opts.DirInPath, e.Name())
		opts.FileOutPath = fmt.Sprintf("%s/%s", opts.DirOutPath, e.Name())

		err := processFile(opts)
		if err != nil {
			return fmt.Errorf("failed to process file: %s", err)
		}

		logger.Print("Done!", logger.ColorDebug)
	}

	return nil
}

func processFile(opts *Options) error {
	doc, err := os.ReadFile(opts.FileInPath)
	if err != nil {
		return err
	}

	rows := readRows(doc)
	records := getRecords(rows)
	recordsHash := getRecordsHash(records)

	if opts.SortRecords {
		records = sortRecords(records, opts.SortRecordsByLength)
	}

	if opts.CheckBusyDomains {
		records = filterRecordsByDnsZones(records, opts.CheckBusyDomainZones, opts.SkipBusyDomains)
	}

	if recordsHash == getRecordsHash(records) {
		logger.Print("No changes", logger.ColorNotice)
		return nil
	}

	rows = prepareRowsToSave(rows, records)

	err = saveRows(opts.FileOutPath, rows)
	if err != nil {
		return err
	}

	return nil
}

func getRecordsHash(records RecordsSlice) [sha256.Size]byte {
	return sha256.Sum256([]byte(fmt.Sprintln(records)))
}

func getRecords(rows []string) RecordsSlice {
	records := make(RecordsSlice)
	group := ""

	for _, row := range rows {
		if needToSkip(row) {
			continue
		}

		if isGroup(row) {
			group = parseRow(row)
		}

		if isRecord(row) {
			records[group] = append(records[group], parseRow(row))
		}
	}

	return records
}

func filterRecordsByDnsZones(records RecordsSlice, zones []string, skip bool) RecordsSlice {
	for g, rec := range records {
		var recs []string

		for _, r := range rec {
			for _, zone := range zones {
				domain := fmt.Sprintf("%s.%s", r, zone)

				err := dns.CheckDomain(domain)
				if errors.Is(err, dns.ErrDomainIsBusy) {
					logger.Print(fmt.Sprintf("Domain %s is busy", domain), logger.ColorError)
					if skip {
						continue
					}
				} else if err != nil {
					logger.Print(fmt.Sprintf("Unable to check domain: %s, err: %e", domain, err), logger.ColorError)
				}

				recs = append(recs, r)
			}
		}

		records[g] = recs
	}

	return records
}

func sortRecords(records RecordsSlice, sortByLen bool) RecordsSlice {
	for _, rec := range records {
		if sortByLen {
			sort.Sort(srt.DomainLengthSort(rec))
		} else {
			sort.Strings(rec)
		}
	}
	return records
}

func prepareRowsToSave(rows []string, records RecordsSlice) []string {
	var newRows []string

	for _, row := range rows {
		if isRecord(row) {
			continue
		}

		if isGroup(row) {
			group := parseRow(row)
			newRows = append(newRows, fmt.Sprintf("- %s", group))
			for _, r := range records[group] {
				newRows = append(newRows, fmt.Sprintf("    - %s", r))
			}
			continue
		}

		newRows = append(newRows, row)
	}

	return newRows
}

func saveRows(filePath string, rows []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("can't open file to save: %s", err)
	}
	defer func() {
		_ = file.Close()
	}()

	writer := bufio.NewWriter(file)
	for _, row := range rows {
		_, err = writer.WriteString(row + "\n")
		if err != nil {
			return fmt.Errorf("can't write row: %s", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("can't write data to file: %s", err)
	}

	return nil
}

func readRows(doc []byte) []string {
	return strings.Split(string(doc), "\n")
}

func needToSkip(row string) bool {
	if !strings.Contains(row, "-") {
		return true
	}

	return false
}

func isValidRow(row string) bool {
	return len(row) > 0 && strings.Index(row, "- ") >= 0
}

func isGroup(row string) bool {
	return isValidRow(row) && rowSpaces(row) <= 1
}

func isRecord(row string) bool {
	return isValidRow(row) && rowSpaces(row) > 1
}

func rowSpaces(row string) int {
	return strings.Count(row, " ")
}

func parseRow(row string) string {
	if len(row) == 0 {
		return row
	}
	return strings.TrimSpace(row[strings.Index(row, "- ")+2:])
}

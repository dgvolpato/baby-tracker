// One-time import tool — reads a CSV exported from Milo Tracker and writes
// each bottle feeding into DynamoDB.
//
// Usage:
//   go run ./cmd/import <path-to-csv> [-tz America/New_York]
//
// "Direct Nursing" rows are skipped (no oz amount in our schema).
// Reads .env from the current directory for TABLE_NAME and AWS config.
package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"baby-tracker/internal/models"
	"baby-tracker/internal/store"

	"github.com/google/uuid"
)

func main() {
	tzFlag := flag.String("tz", "Local", "timezone for parsing CSV times (e.g. America/New_York)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/import [-tz <timezone>] <csv-file>")
		os.Exit(1)
	}

	loadEnv(".env")

	loc, err := time.LoadLocation(*tzFlag)
	if err != nil {
		log.Fatalf("unknown timezone %q: %v", *tzFlag, err)
	}
	fmt.Printf("parsing timestamps in timezone: %s\n", loc)

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // allow variable field count
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	s, err := store.Get()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	imported, skipped := 0, 0

	for i, rec := range records[1:] { // skip header
		rowNum := i + 2

		if len(rec) < 4 {
			fmt.Printf("row %d: too few columns, skipping\n", rowNum)
			skipped++
			continue
		}

		date := strings.TrimSpace(rec[0])
		timeStr := strings.TrimSpace(rec[1])
		typeStr := strings.TrimSpace(rec[2])
		amountStr := strings.TrimSpace(rec[3])

		// Map CSV type → schema type
		var feedType models.FeedingType
		switch typeStr {
		case "Formula":
			feedType = models.FeedingTypeFormula
		case "Breast Milk (bottle)":
			feedType = models.FeedingTypeBreast
		case "Direct Nursing":
			fmt.Printf("row %d: skipping Direct Nursing (no oz amount)\n", rowNum)
			skipped++
			continue
		default:
			fmt.Printf("row %d: unknown type %q, skipping\n", rowNum, typeStr)
			skipped++
			continue
		}

		if amountStr == "" {
			fmt.Printf("row %d: empty amount for %s, skipping\n", rowNum, typeStr)
			skipped++
			continue
		}
		oz, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || oz <= 0 {
			fmt.Printf("row %d: invalid amount %q, skipping\n", rowNum, amountStr)
			skipped++
			continue
		}

		// Parse "04/18/2026" + "6:32 PM" → time.Time
		ts, err := time.ParseInLocation("01/02/2006 3:04 PM", date+" "+timeStr, loc)
		if err != nil {
			fmt.Printf("row %d: cannot parse datetime %q %q: %v, skipping\n", rowNum, date, timeStr, err)
			skipped++
			continue
		}

		feeding := models.Feeding{
			ID:        uuid.NewString(),
			Timestamp: ts.UTC().Format(time.RFC3339),
			Type:      feedType,
			Oz:        oz,
			CreatedBy: "import",
		}

		if err := s.CreateFeeding(ctx, feeding); err != nil {
			log.Printf("row %d: DynamoDB error: %v", rowNum, err)
			skipped++
			continue
		}
		imported++
		if imported%20 == 0 {
			fmt.Printf("  %d rows imported...\n", imported)
		}
	}

	fmt.Printf("\ndone — imported: %d, skipped: %d\n", imported, skipped)
}

func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if os.Getenv(strings.TrimSpace(k)) == "" {
			os.Setenv(strings.TrimSpace(k), strings.TrimSpace(v))
		}
	}
}

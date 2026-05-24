// One-time import tool — reads a CSV exported from Milo Tracker (Diapers sheet)
// and writes each entry into DynamoDB.
//
// Usage:
//   go run ./cmd/import-diapers <path-to-csv> [-tz America/Chicago]
package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"baby-tracker/internal/models"
	"baby-tracker/internal/store"

	"github.com/google/uuid"
)

// parseTime tries multiple layouts to handle both "3:04 PM" and "3:04:05 PM".
func parseTime(date, timeStr string, loc *time.Location) (time.Time, error) {
	layouts := []string{
		"01/02/2006 3:04:05 PM",
		"01/02/2006 3:04 PM",
	}
	combined := date + " " + timeStr
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, combined, loc); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q %q", date, timeStr)
}

func main() {
	tzFlag := flag.String("tz", "Local", "timezone for parsing CSV times (e.g. America/Chicago)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/import-diapers [-tz <timezone>] <csv-file>")
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
	r.FieldsPerRecord = -1
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

	for i, rec := range records[1:] {
		rowNum := i + 2

		if len(rec) < 3 {
			fmt.Printf("row %d: too few columns, skipping\n", rowNum)
			skipped++
			continue
		}

		date := strings.TrimSpace(rec[0])
		timeStr := strings.TrimSpace(rec[1])
		typeStr := strings.TrimSpace(rec[2])

		var diaperType models.DiaperType
		switch typeStr {
		case "Wet":
			diaperType = models.DiaperTypeWet
		case "Poop":
			diaperType = models.DiaperTypePoop
		case "Wet + Poop":
			diaperType = models.DiaperTypeBoth
		default:
			fmt.Printf("row %d: unknown type %q, skipping\n", rowNum, typeStr)
			skipped++
			continue
		}

		ts, err := parseTime(date, timeStr, loc)
		if err != nil {
			fmt.Printf("row %d: %v, skipping\n", rowNum, err)
			skipped++
			continue
		}

		diaper := models.Diaper{
			ID:        uuid.NewString(),
			Timestamp: ts.UTC().Format(time.RFC3339),
			Type:      diaperType,
			CreatedBy: "import",
		}

		if err := s.CreateDiaper(ctx, diaper); err != nil {
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

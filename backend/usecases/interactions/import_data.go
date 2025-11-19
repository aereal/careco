package interactions

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"careco/backend/domain"
	"careco/backend/usecases"
	"careco/backend/usecases/ports"

	"github.com/aereal/iter/seq"
	"github.com/aereal/optional"
	"go.opentelemetry.io/otel/trace"
)

type ExportFileName string

func ProvideImportData(tp trace.TracerProvider, bulkWriter ports.DrivingRecordBulkWriter, fileName ExportFileName) *ImportData {
	return &ImportData{
		tracer:     tp.Tracer("careco/backend/usecases/interactions.ImportData"),
		bulkWriter: bulkWriter,
		fileName:   fileName,
	}
}

type ImportData struct {
	tracer     trace.Tracer
	bulkWriter ports.DrivingRecordBulkWriter
	fileName   ExportFileName
}

var _ usecases.ImportData = (*ImportData)(nil)

func (u *ImportData) ImportData(ctx context.Context) (err error) {
	ctx, span := u.tracer.Start(ctx, "ImportData")
	defer span.End()

	rows, err := parseRows(string(u.fileName))
	if err != nil {
		return fmt.Errorf("parseRows: %w", err)
	}

	records := make([]*domain.DrivingRecord, 0, len(rows))
	slices.Reverse(rows)
	for right, left := range seq.Pairwise(slices.Values(rows)) {
		records = append(records, &domain.DrivingRecord{
			Memo:               optional.None[string](),
			Date:               right.date,
			DistanceKilometers: right.distance - left.distance,
		})
	}

	if err := u.bulkWriter.BulkWriteDrivingRecords(ctx, records); err != nil {
		return fmt.Errorf("BulkWriteDrivingRecords: %w", err)
	}
	return nil
}

type data struct {
	date     time.Time
	distance int64
}

func parseRecords(records []string) (*data, error) {
	if len(records) < 2 {
		return nil, &usecases.MalformedRecordsError{ActualColumns: len(records)}
	}
	if records[1] == "" {
		return nil, io.EOF
	}
	r := new(data)
	err := errors.Join(
		parseTime(&r.date, records[0]),
		parseInt(&r.distance, records[1]),
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

const dateLayout = "2006/01/02"

func parseTime(t *time.Time, input string) error {
	var err error
	*t, err = time.Parse(dateLayout, input)
	return err
}

func parseInt(i *int64, input string) error {
	var err error
	*i, err = strconv.ParseInt(strings.ReplaceAll(input, ",", ""), 10, 64)
	return err
}

func parseRows(fileName string) ([]*data, error) {
	f, err := os.Open(fileName) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("os.Open: %w", err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = '\t'
	var seenHeader bool
	rows := make([]*data, 0)
	var line int
	for {
		line++
		records, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv.Reader.Read: %w", err)
		}
		if !seenHeader {
			seenHeader = true
			continue
		}
		parsed, err := parseRecords(records)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, &usecases.ParseTSVError{Err: err, Line: line}
		}
		rows = append(rows, parsed)
	}
	return rows, nil
}

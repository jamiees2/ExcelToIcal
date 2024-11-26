package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SelectFile allows the user to select an Excel file
func (a *App) SelectFile() string {
	filename, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Excel File",
		Filters: []runtime.FileFilter{
			{DisplayName: "Excel Files", Pattern: "*.xlsx"},
		},
	})
	if err != nil {
		return ""
	}
	return filename
}

// SaveFile allows the user to select a location to save the iCal file
func (a *App) SaveFile() string {
	filename, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save iCal File",
		DefaultFilename: "events.ics",
		Filters: []runtime.FileFilter{
			{DisplayName: "iCal Files", Pattern: "*.ics"},
		},
	})
	if err != nil {
		return ""
	}
	return filename
}

// ConvertExcelToICal performs the conversion
func (a *App) ConvertExcelToICal(excelFile, icalFile string) (string, error) {
	err := convertExcelToICal(excelFile, "Sheet1", icalFile)
	if err != nil {
		return "", err
	}
	return "Conversion successful!", nil
}

const (
	DateColumn      = "Date"
	StartTimeColumn = "Start Time"
	EndTimeColumn   = "End Time"
	SubjectColumn   = "Subject"
)

func convertExcelToICal(excelFile, sheetName, icalFile string) error {
	f, err := excelize.OpenFile(excelFile)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to read sheet: %w", err)
	}

	if len(rows) < 1 {
		return fmt.Errorf("sheet %s is empty", sheetName)
	}

	headers := rows[0]
	colIndexes := map[string]int{
		DateColumn:      -1,
		StartTimeColumn: -1,
		EndTimeColumn:   -1,
		SubjectColumn:   -1,
	}
	for i, header := range headers {
		if _, exists := colIndexes[header]; exists {
			colIndexes[header] = i
		}
	}
	if colIndexes[DateColumn] == -1 || colIndexes[StartTimeColumn] == -1 || colIndexes[EndTimeColumn] == -1 || colIndexes[SubjectColumn] == -1 {
		return fmt.Errorf("sheet must contain 'date', 'start time', 'end time', and 'subject' columns")
	}

	file, err := os.Create(icalFile)
	if err != nil {
		return fmt.Errorf("failed to create iCal file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString("BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//Excel to iCal Conversion//example.com//\n")
	if err != nil {
		return fmt.Errorf("failed to write to iCal file: %w", err)
	}

	for _, row := range rows[1:] {
		date := row[colIndexes[DateColumn]]
		startTime := row[colIndexes[StartTimeColumn]]
		endTime := row[colIndexes[EndTimeColumn]]
		subject := row[colIndexes[SubjectColumn]]

		if date == "" || startTime == "" || endTime == "" || subject == "" {
			continue
		}

		startDateTime, err := parseDateTime(date, startTime)
		if err != nil {
			return fmt.Errorf("invalid date or start time format in row: %v", row)
		}
		endDateTime, err := parseDateTime(date, endTime)
		if err != nil {
			return fmt.Errorf("invalid date or end time format in row: %v", row)
		}

		_, err = file.WriteString(fmt.Sprintf(
			"BEGIN:VEVENT\nSUMMARY:%s\nDTSTART:%s\nDTEND:%s\nEND:VEVENT\n",
			subject,
			startDateTime.Format("20060102T150405Z"),
			endDateTime.Format("20060102T150405Z"),
		))
		if err != nil {
			return fmt.Errorf("failed to write event to iCal file: %w", err)
		}
	}

	_, err = file.WriteString("END:VCALENDAR\n")
	if err != nil {
		return fmt.Errorf("failed to write to iCal file: %w", err)
	}

	return nil
}

func parseDateTime(date, timeStr string) (time.Time, error) {
	fullDateTime := fmt.Sprintf("%s %s", date, timeStr)
	parsedTime, err := time.Parse("2006-01-02 15:04:05", fullDateTime)
	if err == nil {
		return parsedTime, nil
	}
	parsedTime, err = time.Parse("2006-01-02 15:04", fullDateTime)
	if err == nil {
		return parsedTime, nil
	}
	parsedTime, err = time.Parse("02/01/2006 15:04:05", fullDateTime)
	if err == nil {
		return parsedTime, nil
	}
	parsedTime, err = time.Parse("02/01/2006 15:04", fullDateTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date/time format: %s %s", date, timeStr)
	}
	return parsedTime, nil
}

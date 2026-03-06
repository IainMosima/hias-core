package reporting

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/xuri/excelize/v2"
)

type ColumnDef struct {
	Name   string `json:"name"`
	Label  string `json:"label"`
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
}

type ReportExporter interface {
	ExportCSV(columns json.RawMessage, data []map[string]interface{}) ([]byte, error)
	ExportXLSX(columns json.RawMessage, data []map[string]interface{}, title string) ([]byte, error)
	ExportPDF(columns json.RawMessage, data []map[string]interface{}, title string, summary map[string]interface{}) ([]byte, error)
}

type reportExporter struct{}

func NewReportExporter() ReportExporter {
	return &reportExporter{}
}

func parseColumns(columnsJSON json.RawMessage) ([]ColumnDef, error) {
	var cols []ColumnDef
	if err := json.Unmarshal(columnsJSON, &cols); err != nil {
		return nil, fmt.Errorf("failed to parse columns: %w", err)
	}
	return cols, nil
}

func formatValue(val interface{}, colType string) string {
	if val == nil {
		return ""
	}
	switch colType {
	case "money":
		switch v := val.(type) {
		case int64:
			return fmt.Sprintf("%.2f", float64(v)/100)
		case float64:
			return fmt.Sprintf("%.2f", v/100)
		default:
			return fmt.Sprintf("%v", val)
		}
	case "percentage":
		return fmt.Sprintf("%.2f%%", toFloat64(val))
	case "decimal":
		return fmt.Sprintf("%.2f", toFloat64(val))
	case "date":
		if t, ok := val.(time.Time); ok {
			return t.Format("02/01/2006")
		}
		return fmt.Sprintf("%v", val)
	case "datetime":
		if t, ok := val.(time.Time); ok {
			return t.Format("02/01/2006 15:04")
		}
		return fmt.Sprintf("%v", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func toFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0
	}
}

// --- CSV Export ---

func (e *reportExporter) ExportCSV(columns json.RawMessage, data []map[string]interface{}) ([]byte, error) {
	cols, err := parseColumns(columns)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header row
	headers := make([]string, len(cols))
	for i, col := range cols {
		headers[i] = col.Label
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Data rows
	for _, row := range data {
		record := make([]string, len(cols))
		for j, col := range cols {
			record[j] = formatValue(row[col.Name], col.Type)
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV flush error: %w", err)
	}

	return buf.Bytes(), nil
}

// --- XLSX Export ---

func (e *reportExporter) ExportXLSX(columns json.RawMessage, data []map[string]interface{}, title string) ([]byte, error) {
	cols, err := parseColumns(columns)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Report"
	f.SetSheetName("Sheet1", sheet)

	// Title row
	f.SetCellValue(sheet, "A1", title)
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
	})
	f.SetCellStyle(sheet, "A1", "A1", titleStyle)

	f.SetCellValue(sheet, "A2", fmt.Sprintf("Generated: %s", time.Now().Format("02 Jan 2006 15:04")))

	// Header row (row 4)
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#2E86AB"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	for i, col := range cols {
		cell, _ := excelize.CoordinatesToCellName(i+1, 4)
		f.SetCellValue(sheet, cell, col.Label)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		f.SetColWidth(sheet, colLetter(i+1), colLetter(i+1), 18)
	}

	// Money style
	moneyStyle, _ := f.NewStyle(&excelize.Style{
		NumFmt: 4, // #,##0.00
	})

	// Data rows (starting row 5)
	for rowIdx, row := range data {
		for colIdx, col := range cols {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+5)
			val := row[col.Name]

			switch col.Type {
			case "money":
				if v, ok := val.(int64); ok {
					f.SetCellValue(sheet, cell, float64(v)/100)
					f.SetCellStyle(sheet, cell, cell, moneyStyle)
				} else {
					f.SetCellValue(sheet, cell, formatValue(val, col.Type))
				}
			case "date", "datetime":
				if t, ok := val.(time.Time); ok {
					f.SetCellValue(sheet, cell, t)
				} else {
					f.SetCellValue(sheet, cell, formatValue(val, col.Type))
				}
			case "percentage":
				f.SetCellValue(sheet, cell, toFloat64(val))
			default:
				f.SetCellValue(sheet, cell, formatValue(val, col.Type))
			}
		}
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("failed to write XLSX: %w", err)
	}

	return buf.Bytes(), nil
}

func colLetter(col int) string {
	name, _ := excelize.ColumnNumberToName(col)
	return name
}

// --- PDF Export ---

func (e *reportExporter) ExportPDF(columns json.RawMessage, data []map[string]interface{}, title string, summary map[string]interface{}) ([]byte, error) {
	cols, err := parseColumns(columns)
	if err != nil {
		return nil, err
	}

	pdf := fpdf.New("L", "mm", "A4", "") // Landscape for tables
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(270, 10, "HIAS Insurance")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(270, 6, "Health Insurance Administration System")
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(270, 10, title)
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 9)
	pdf.Cell(270, 5, fmt.Sprintf("Generated: %s | Rows: %d", time.Now().Format("02 Jan 2006 15:04"), len(data)))
	pdf.Ln(8)
	pdf.Line(10, pdf.GetY(), 287, pdf.GetY())
	pdf.Ln(4)

	// Summary section
	if len(summary) > 0 {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(270, 6, "Summary")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 9)
		for k, v := range summary {
			pdf.Cell(60, 5, k+":")
			pdf.Cell(80, 5, fmt.Sprintf("%v", v))
			pdf.Ln(5)
		}
		pdf.Ln(4)
	}

	// Calculate column widths
	totalWidth := 270.0 // A4 landscape usable width
	colWidth := totalWidth / float64(len(cols))
	if colWidth > 40 {
		colWidth = 40
	}

	// Table header
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(46, 134, 171) // #2E86AB
	pdf.SetTextColor(255, 255, 255)
	for _, col := range cols {
		pdf.CellFormat(colWidth, 7, col.Label, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 7)
	pdf.SetTextColor(0, 0, 0)
	for rowIdx, row := range data {
		// Alternating row colors
		if rowIdx%2 == 0 {
			pdf.SetFillColor(240, 240, 240)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		for _, col := range cols {
			val := formatValue(row[col.Name], col.Type)
			align := "L"
			if col.Type == "money" || col.Type == "number" || col.Type == "percentage" || col.Type == "decimal" {
				align = "R"
			}
			pdf.CellFormat(colWidth, 6, val, "1", 0, align, true, 0, "")
		}
		pdf.Ln(-1)
	}

	// Footer
	pdf.SetY(-20)
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(270, 5, fmt.Sprintf("Generated on %s", time.Now().Format("02 Jan 2006 15:04")))
	pdf.Ln(4)
	pdf.Cell(270, 5, "This is a system-generated report. For queries, contact support@hias.co.ke")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

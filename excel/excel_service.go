package excel

import (
	"crypto-skatt-go/exchange"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct{}

// NewExcelService creates a new instance of ExcelService
func NewExcelService() *ExcelService {
	return &ExcelService{}
}

const sheetName = "Sheet1"

var headers = []string{
	"Tidspunkt",
	"Type",
	"Inn",
	"Inn-Valuta",
	"Ut",
	"Ut-Valuta",
	"Gebyr",
	"Gebyr-Valuta",
	"Marked",
	"Notat",
}

func (es *ExcelService) ExportToFile(withdrawals []exchange.TransferHistory, deposits []exchange.TransferHistory, orders []exchange.OrderHistory) error {
	file := excelize.NewFile()
	defer file.Close()

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheetName, cell, header)
	}
	// Add withdrawals
	for i, order := range orders {
		row := i + 2 // Start from the second row
		file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), order.Timestamp)
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Handel")
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), order.BoughtAmount)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), order.BoughtCoin)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), order.SoldAmount)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), order.SoldCoin)
		file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), order.FeeAmount)
		file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), order.FeeCurrency)
		file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), order.Exchange)
		file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), "")
	}

	if err := file.SaveAs("export.xlsx"); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// Implement the logic to export data to an Excel file
	// This is a placeholder implementation
	return nil
}

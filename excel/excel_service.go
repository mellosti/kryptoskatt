package excel

import (
	"crypto-skatt-go/exchange"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct {
	OutputPath string // Directory where Excel files will be saved
}

// NewExcelService creates a new instance of ExcelService
func NewExcelService(outputPath string) *ExcelService {
	if outputPath == "" {
		outputPath = "." // Default to current directory
	}
	return &ExcelService{
		OutputPath: outputPath,
	}
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

// Transaction types
const (
	TypeTrade      = "Handel"
	TypeDeposit    = "Innskudd"
	TypeWithdrawal = "Uttak"
)

const fileName = "export.xlsx"

func (es *ExcelService) ExportToFile(withdrawals []exchange.TransferHistory, deposits []exchange.TransferHistory, orders []exchange.OrderHistory) error {
	file := excelize.NewFile()
	defer file.Close()

	// Create headers
	for i, header := range headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return fmt.Errorf("failed to convert coordinates: %w", err)
		}
		file.SetCellValue(sheetName, cell, header)
	}

	row := 2 // Start from the second row

	// Add orders
	for _, order := range orders {
		es.addOrderRow(file, row, order)
		row++
	}

	// Add deposits
	for _, deposit := range deposits {
		es.addTransferRow(file, row, deposit, TypeDeposit)
		row++
	}

	// Add withdrawals
	for _, withdrawal := range withdrawals {
		es.addTransferRow(file, row, withdrawal, TypeWithdrawal)
		row++
	}

	if err := file.SaveAs(fileName); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// addOrderRow adds a trade order row to the Excel file
func (es *ExcelService) addOrderRow(file *excelize.File, row int, order exchange.OrderHistory) {
	file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), order.Timestamp)
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), TypeTrade)
	file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), order.BoughtAmount)
	file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), order.BoughtCoin)
	file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), order.SoldAmount)
	file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), order.SoldCoin)
	file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), order.FeeAmount)
	file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), order.FeeCurrency)
	file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), order.Exchange)
	file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), "")
}

// addTransferRow adds a deposit or withdrawal row to the Excel file
func (es *ExcelService) addTransferRow(file *excelize.File, row int, transfer exchange.TransferHistory, transferType string) {
	file.SetCellValue(sheetName, fmt.Sprintf("A%d", row), transfer.Timestamp)
	file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), transferType)

	// For deposits, set the "Inn" columns
	if transferType == TypeDeposit {
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), transfer.Amount)
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), transfer.Coin)
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "")
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "")
	} else {
		// For withdrawals, set the "Ut" columns
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "")
		file.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "")
		file.SetCellValue(sheetName, fmt.Sprintf("E%d", row), transfer.Amount)
		file.SetCellValue(sheetName, fmt.Sprintf("F%d", row), transfer.Coin)
	}

	file.SetCellValue(sheetName, fmt.Sprintf("G%d", row), transfer.FeeAmount)
	file.SetCellValue(sheetName, fmt.Sprintf("H%d", row), transfer.FeeCoin)
	file.SetCellValue(sheetName, fmt.Sprintf("I%d", row), transfer.Exchange)
	file.SetCellValue(sheetName, fmt.Sprintf("J%d", row), transfer.TransactionID)
}

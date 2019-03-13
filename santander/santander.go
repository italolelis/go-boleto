// Package santander is the Santander implementation of boletos
package santander

import (
	"fmt"
	"strconv"

	boleto "github.com/italolelis/go-boleto"
)

const (
	// BarcodeNumber.BankNumbers size
	bankNumbersSize = 25
)

// Santander is the Santander bank slip implementation
// Source: (https://www.santander.com.br/document/wps/sl-tabela-de-tarifas-cobranca.pdf)
type Santander struct {
	boleto.Bank
	Company  boleto.Company
	Agency   int
	Account  int
	Carteira int
	IOS      int
}

// New creates a new instance of Santander
func New() *Santander {
	return &Santander{
		Bank: boleto.Bank{
			ID:             33,
			Aceite:         "N",
			Currency:       9,
			CurrencyName:   "R$",
			AccountMaxSize: 7,
		},
	}
}

// Barcode Get the Barcode, creating a BarcodeNumber
func (b *Santander) Barcode(d boleto.Document) (*boleto.BarcodeNumber, error) {
	// Complete the BankNumbers digits
	var bn string
	bn += "9"
	bn += fmt.Sprintf("%0"+strconv.Itoa(b.AccountMaxSize)+"d", b.Account)
	bn += fmt.Sprintf("%013d", d.OurNumber)
	bn += strconv.Itoa(b.IOS)
	bn += strconv.Itoa(b.Carteira)

	dueDateFactor, err := boleto.DateDueFactor(d.DateDue)
	if err != nil {
		return nil, fmt.Errorf("there was an error calculating the due date factor: %s", err)
	}

	// Create a new Barcode
	return &boleto.BarcodeNumber{
		BankID:        b.ID,
		CurrencyID:    b.Currency,
		DateDueFactor: dueDateFactor,
		Value:         d.Value,
		BankNumbers:   fmt.Sprintf("%0"+strconv.Itoa(bankNumbersSize)+"s", bn),
	}, nil
}

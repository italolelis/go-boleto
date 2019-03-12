// Package bradesco is the Bradesco implementation of boletos
package bradesco

import (
	"errors"
	"fmt"
	"strconv"

	boleto "github.com/italolelis/go-boleto"
)

const (
	// BarcodeNumber.BankNumbers size
	bankNumbersSize = 25
)

var (
	// ErrInvalidConvenio is used when the convenio is not supported by this bank
	ErrInvalidConvenio = errors.New("invalid convenio provided")
)

// Bradesco is the Bradesco bank slip implementation
// Source: (https://banco.bradesco/assets/pessoajuridica/pdf/4008-524-0121-08-layout-cobranca-versao-portuguesSS28785.pdf)
type Bradesco struct {
	boleto.Bank
	Company  *boleto.Company
	Agency   int
	Account  int
	Convenio int
	Carteira int
}

// New creates a new instance of BB
func New(agency int, account int, convenio int, carteira int, company *boleto.Company) *Bradesco {
	return &Bradesco{
		Bank: boleto.Bank{
			ID:             237,
			Aceite:         "N",
			Currency:       9,
			CurrencyName:   "R$",
			AgencyMaxSize:  4,
			AccountMaxSize: 7,
		},
		Agency:   agency,
		Account:  account,
		Convenio: convenio,
		Carteira: carteira,
		Company:  company,
	}
}

// Barcode Get the Barcode, creating a BarcodeNumber
func (b *Bradesco) Barcode(d *boleto.Document) (*boleto.BarcodeNumber, error) {
	// Complete the BankNumbers digits
	var bn string
	bn += fmt.Sprintf("%0"+strconv.Itoa(b.AgencyMaxSize)+"d", b.Agency)
	bn += fmt.Sprintf("%02d", b.Carteira)
	bn += fmt.Sprintf("%011d", d.OurNumber)

	if b.Carteira == 9 {
		bn += fmt.Sprintf("%0"+strconv.Itoa(b.AccountMaxSize)+"d", b.Account)
	} else {
		bn += fmt.Sprintf("%07d", b.Convenio)
	}

	bn += "0"

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

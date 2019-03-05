// Package bb is the Banco do Brasil implementation of boletos
package bb

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

// BB - Banco do Brasil
// Source: (http://www.bb.com.br/docs/pub/emp/mpe/espeboletobb.pdf)
type BB struct {
	boleto.Bank
	Agency   int
	Account  int
	Convenio int
	Carteira int
	Company  boleto.Company
}

// New creates a new instance of BB
func New() *BB {
	return &BB{
		Bank: boleto.Bank{
			ID:             001,
			Aceite:         "N",
			Currency:       9,
			CurrencyName:   "R$",
			AgencyMaxSize:  4,
			AccountMaxSize: 8,
		},
	}
}

// Barcode Get the Barcode, creating a BarcodeNumber
func (b BB) Barcode(d boleto.Document) (*boleto.BarcodeNumber, error) {
	// Complete the BankNumbers digits, by adding convenio rules according to the bank
	var bn string
	convenioSize := len(strconv.Itoa(b.Convenio))
	ourNumberSize := len(strconv.Itoa(d.OurNumber))

	switch convenioSize {
	case 4:
		// For Convenio size 4: CCCCNNNNNNN-X
		// C = Convenio number int(4)
		// N = OurNumber int(7)
		// X = DV, calculated by module11 int(1)
		if ourNumberSize > 7 {
			return nil, fmt.Errorf("document our number for this convenio has exceeded 7 digits")
		}

		bn += strconv.Itoa(b.Convenio)
		bn += fmt.Sprintf("%07d", d.OurNumber)
		bn += fmt.Sprintf("%0"+strconv.Itoa(b.AgencyMaxSize)+"d", b.Agency)
		bn += fmt.Sprintf("%0"+strconv.Itoa(b.AccountMaxSize)+"d", b.Account)
		bn += strconv.Itoa(b.Carteira)
	case 6:
		// For Convenio size 6: CCCCCCNNNNN-X
		// C = Convenio number int(6)
		// N = OurNumber int(5)
		// X = DV, calculated by module11 int(1)
		if ourNumberSize > 5 {
			return nil, fmt.Errorf("document our number for this convenio has exceeded 5 digits")
		}

		bn += strconv.Itoa(b.Convenio)
		bn += fmt.Sprintf("%05d", d.OurNumber)
		bn += fmt.Sprintf("%0"+strconv.Itoa(b.AgencyMaxSize)+"d", b.Agency)
		bn += fmt.Sprintf("%0"+strconv.Itoa(b.AccountMaxSize)+"d", b.Account)
		bn += strconv.Itoa(b.Carteira)
	case 7:
		// For Convenio size 7: CCCCCCCNNNNNNNNNN
		// C = Convenio number int(7)
		// N = OurNumber int(9)
		if ourNumberSize > 9 {
			return nil, fmt.Errorf("document our number for this convenio has exceeded 9 digits")
		}

		bn += fmt.Sprintf("%013d", b.Convenio)
		bn += fmt.Sprintf("%09d", d.OurNumber)
		bn += strconv.Itoa(b.Carteira)
	default:
		return nil, ErrInvalidConvenio
	}

	// Create a new Barcode
	n := &boleto.BarcodeNumber{
		BankID:        b.ID,
		CurrencyID:    b.Currency,
		DateDueFactor: boleto.DateDueFactor(d.DateDue),
		Value:         d.Value,
		BankNumbers:   fmt.Sprintf("%0"+strconv.Itoa(bankNumbersSize)+"s", bn),
	}

	if err := n.Verify(); err != nil {
		return nil, fmt.Errorf("the BB bank slip number could not be generated: %s", err)
	}

	return n, nil
}

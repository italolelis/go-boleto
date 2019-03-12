// Package itau is the Itau implementation of boletos
package itau

import (
	"fmt"
	"strconv"

	boleto "github.com/italolelis/go-boleto"
)

const (
	bankNumbersSize   = 25
	maxOurNumberSize  = 8
	maxDocumentIDSize = 7
	maxClientCodeSize = 5
)

// Itau is the itau bank slip implementation
// Source: (http://download.itau.com.br/bankline/cobranca_cnab240.pdf)
type Itau struct {
	boleto.Bank
	Company    *boleto.Company
	Agency     int
	Account    int
	Carteira   int
	ClientCode int
}

// New creates a new instance of BB
func New(agency int, account int, carteira int, clientCode int, company *boleto.Company) *Itau {
	return &Itau{
		Bank: boleto.Bank{
			ID:             341,
			Aceite:         "N",
			Currency:       9,
			CurrencyName:   "R$",
			AgencyMaxSize:  4,
			AccountMaxSize: 5,
		},
		Agency:     agency,
		Account:    account,
		ClientCode: clientCode,
		Carteira:   carteira,
		Company:    company,
	}
}

// Barcode Get the Barcode, creating a BarcodeNumber
func (b *Itau) Barcode(d *boleto.Document) (*boleto.BarcodeNumber, error) {
	// Complete the BankNumbers digits, by adding carteira rules according to the bank
	var bn string

	// Verify max of Document.OurNumber size
	if size := len(strconv.Itoa(d.OurNumber)); size > maxOurNumberSize {
		return nil, fmt.Errorf("our number needs to have a max of %d digits", maxOurNumberSize)
	}

	// Verify max of Document.ID size
	if len(strconv.Itoa(d.ID)) > maxDocumentIDSize {
		return nil, fmt.Errorf("document ID needs to have a max of %d digits", maxDocumentIDSize)
	}

	// Add rules to carteira equals to 107, 122, 142, 143, 196, 198
	if b.Carteira == 107 || b.Carteira == 122 || b.Carteira == 142 ||
		b.Carteira == 143 || b.Carteira == 196 || b.Carteira == 198 {
		// CCCNNNNNNNNLLLLLLLDDDDDX0
		// C = Carteira number int(3)
		// N = OurNumber int(8)
		// L = Document number int(7)
		// D = Client code int(5)
		// X = DV, calculated by module10 int(1)

		// Verify max Bank.ClientCode size
		clientCodeSize := len(strconv.Itoa(b.ClientCode))
		if clientCodeSize > maxClientCodeSize {
			return nil, fmt.Errorf("client code needs to have a max of %d digits", maxClientCodeSize)
		}

		// this code var is part of the BankNumbers,
		// we use it to generate another var with module10
		var code string
		code += strconv.Itoa(b.Carteira)
		code += fmt.Sprintf("%08d", d.OurNumber)
		code += fmt.Sprintf("%07d", d.ID)
		code += fmt.Sprintf("%05d", b.ClientCode)

		// module10 with code created above
		codeModule := strconv.Itoa(boleto.Module10(code, 2))
		bn += code + codeModule + "0"

	} else {
		// CCCNNNNNNNNXAAAATTTTTY000
		// C = Carteira number int(3)
		// N = OurNumber int(8)
		// X = DV, calculated by module10 int(1)
		// A = Agency number int(4)
		// T = Account number int(5)
		// Y = DV, calculated by module10 int(1)

		// Add rules to carteira number equals to 126, 131, 146, 150, 168
		var codeModuleCarteira int
		if b.Carteira == 126 || b.Carteira == 131 || b.Carteira == 146 ||
			b.Carteira == 150 || b.Carteira == 168 {
			// module10 with bank agency, account, carteira, and OurNumber
			m := fmt.Sprintf("%0"+strconv.Itoa(b.AgencyMaxSize)+"d", b.Agency)
			m += fmt.Sprintf("%0"+strconv.Itoa(b.AccountMaxSize)+"d", b.Account)
			m += strconv.Itoa(b.Carteira)
			m += strconv.Itoa(d.OurNumber)
			codeModuleCarteira = boleto.Module10(m, 2)
		} else {
			// module10 with carteira and OurNumber
			m := strconv.Itoa(b.Carteira)
			m += strconv.Itoa(d.OurNumber)
			codeModuleCarteira = boleto.Module10(m, 2)
		}

		// module10 with bank agency and account
		codeModuleAccount := boleto.Module10(strconv.Itoa(b.Account)+strconv.Itoa(b.Agency), 2)

		bn += strconv.Itoa(b.Carteira)
		bn += fmt.Sprintf("%08d", d.OurNumber)
		bn += strconv.Itoa(codeModuleCarteira)
		bn += fmt.Sprintf("%0"+strconv.Itoa(b.AgencyMaxSize)+"d", b.Agency)
		bn += fmt.Sprintf("%0"+strconv.Itoa(b.AccountMaxSize)+"d", b.Account)
		bn += strconv.Itoa(codeModuleAccount)
		bn += "000"
	}

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

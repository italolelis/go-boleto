package boleto

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	// The min size of a bankID
	bankMinSize = 3
	// The min size of the value formated
	valueMinSize = 10
	// The min size of a barcode
	barcodeNumberMinSize = 19
	// The max size of a barcode
	barcodeNumberMaxSize = 44
	// BarcodeNumber.BankNumbers size
	bankNumbersSize = 25
)

// Banker defines the bank barcode generation
type Banker interface {
	// Barcode returns the BarcodeNumber
	Barcode(Document) (*BarcodeNumber, error)
}

// Bank is the basic struct that defines required fileds that all the banks are required to have
type Bank struct {
	ID             int
	Aceite         string
	Currency       int
	CurrencyName   string
	AgencyMaxSize  int
	AccountMaxSize int
}

// BankSlipPrinter defines an interface to print a HTML template using a document
type BankSlipPrinter interface {
	Layout(http.ResponseWriter, Document)
}

// BarcodeNumber defines a barcode number type,
// holds numbers of the barcode
type BarcodeNumber struct {
	// Codigo do banco int(3)
	BankID int
	// Codigo da moeda int(1)
	CurrencyID int
	// Fator de vencimento int(4)
	DateDueFactor int
	// Valor formatado int(10)
	Value uint
	// Campo livre, numeros do banco com nosso numero string(25)
	BankNumbers string
	// Verification digit fomr the barcode
	Dv int
}

// Verify checks if the DV number was generated correctly
func (n *BarcodeNumber) Verify() error {
	s := n.String()

	if len(s) < barcodeNumberMinSize {
		return errors.New("there are missing values in Bank and Document structures")
	}

	if len(s) > barcodeNumberMaxSize {
		return errors.New("there are remaining values in Bank and Document structures")
	}

	n.Dv = Module11(s)

	return nil
}

// Digitable mount the barcode digitable number,
// taking all BarcodeNumber fields together:
// Field 1: AAABC.CCCCX
// A = FEBRABAN Bank identifier
// B = the currency identifier
// C = 20-24 barcode numbers
// X = DV, using module10
//
// Field 2: DDDDD.DDDDDX
// D = 25-34 barcode numbers
// X = DV, using module10
//
// Field 3: EEEEE.EEEEEX
// E = 35-44 barcode numbers
// X = DV, using module10
//
// Field 4: X
// X = DV, BarcodeNumber.Dv
//
// Field 5: UUUUVVVVVVVVVV
// U = Due date factor
// V = Value
//
// return AAABC.CCCCX DDDDD.DDDDDX EEEEE.EEEEEX X UUUUVVVVVVVVVV
func (n *BarcodeNumber) Digitable() string {
	s := n.String()

	// Field 1
	var f1 = fmt.Sprintf("%0"+strconv.Itoa(bankMinSize)+"d", n.BankID)
	f1 += strconv.Itoa(n.CurrencyID)
	f1 += string(s[19]) + "." + s[20:24]
	f1 += strconv.Itoa(Module10(f1, maxModule10))

	// Field 2
	var f2 = s[24:29] + "." + s[29:34]
	f2 += strconv.Itoa(Module10(f2, minModule10))

	// Field 3
	var f3 = s[34:39] + "." + s[39:44]
	f3 += strconv.Itoa(Module10(f3, minModule10))

	// Field 5
	var f4 = strconv.Itoa(n.Dv)

	// Field 5
	var f5 = strconv.Itoa(n.DateDueFactor)
	f5 += fmt.Sprintf("%0"+strconv.Itoa(valueMinSize)+"d", n.Value)

	// All fields together
	return fmt.Sprintf("%s %s %s %s %s", f1, f2, f3, f4, f5)
}

// toString takes BarcodeNumber, and converts to a string,
// including pad numbers and left zeros
func (n *BarcodeNumber) String() string {
	b := bytes.NewBufferString(fmt.Sprintf("%0"+strconv.Itoa(bankMinSize)+"d", n.BankID))
	b.WriteString(strconv.Itoa(n.CurrencyID))
	b.WriteString(strconv.Itoa(n.Dv))
	b.WriteString(strconv.Itoa(n.DateDueFactor))
	b.WriteString(fmt.Sprintf("%0"+strconv.Itoa(valueMinSize)+"d", n.Value))
	b.WriteString(n.BankNumbers)

	return b.String()
}

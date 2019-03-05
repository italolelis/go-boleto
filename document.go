package boleto

import "time"

// Document holds the data of the billet itself
type Document struct {
	// ID identifier of your program orders/payments
	ID      int
	Date    time.Time
	DateDue time.Time
	Value   uint

	// ValueTax is the tax for emitting the bank slip
	ValueTax uint

	// ValueDiscount if any discount is required
	ValueDiscount uint

	// ValueForfeit if any interest or debt needs to be calculated
	ValueForfeit uint

	// OurNumber literely translates to `Nosso Numero`. This is a number that is provided by the partner bank
	// where this boleto is going to be payed for.
	OurNumber int

	// FebrabanType is the document type according FEBRABAN, the default used is "DM" (Duplicata mercantil),
	// Source: (http://www.bb.com.br/docs/pub/emp/empl/dwn/011DescrCampos.pdf)
	FebrabanType string

	// Instructions are any instructions that should be printed in the bank slip
	Instructions [6]string

	// Payer is to whom this boleto is being emitted
	Payer Payer
}

// NewDocument creates a new instance of Document with the emission date set to now
func NewDocument() *Document {
	return NewDocumentFrom(time.Now())
}

// NewDocumentFrom creates a new instance of Document where you must provide the emission date
func NewDocumentFrom(date time.Time) *Document {
	return &Document{Date: date}
}

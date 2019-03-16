package boleto

import (
	"errors"
	"time"
)

const (
	// DefaultFebrabanType defines the default febraban type used in the document
	DefaultFebrabanType = "dm"
	maxInstructionsSize = 6
)

var (
	// ErrDiscountHigherThanValue is used when the discount value is higher than the value of the bank slip
	ErrDiscountHigherThanValue = errors.New("discount value cannot be higher than the value")

	// ErrForfeitHigherThanValue is used when the forfeit value is higher than the value of the bank slip
	ErrForfeitHigherThanValue = errors.New("forfeit value cannot be higher than the value")

	// ErrTooManyInstructions is used when the ammount of instructions reaches the limit
	ErrTooManyInstructions = errors.New("you cannot add more instructions to the document")

	// ErrEmptyPayer is used when the given payer is nil
	ErrEmptyPayer = errors.New("you need to pass a payer pointer")
)

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
	Instructions []string

	// Payer is to whom this boleto is being emitted
	Payer *Payer
}

// NewDocument creates a new instance of Document with the emission date set to now
func NewDocument(id int, ourNumber int, dueDate time.Time) *Document {
	return NewDocumentFrom(id, ourNumber, time.Now(), dueDate)
}

// NewDocumentFrom creates a new instance of Document where you must provide the emission date
func NewDocumentFrom(id int, ourNumber int, date time.Time, dueDate time.Time) *Document {
	return &Document{
		ID:           id,
		OurNumber:    ourNumber,
		Date:         date,
		DateDue:      dueDate,
		FebrabanType: DefaultFebrabanType,
	}
}

// DefineValue defines the values for the document
func (d *Document) DefineValue(value uint, discount uint, forfeit uint) error {
	if discount > value {
		return ErrDiscountHigherThanValue
	}

	if forfeit > value {
		return ErrForfeitHigherThanValue
	}

	d.Value = value
	d.ValueDiscount = discount
	d.ValueForfeit = forfeit

	return nil
}

// AddInstruction add a single instruction to the document
func (d *Document) AddInstruction(i string) error {
	if len(d.Instructions) > maxInstructionsSize {
		return ErrTooManyInstructions
	}

	d.Instructions = append(d.Instructions, i)

	return nil
}

// To defines to whom the document will be addressed
func (d *Document) To(p *Payer) error {
	if p == nil {
		return ErrEmptyPayer
	}

	d.Payer = p

	return nil
}

// Total calculates the total of the boleto
func (d *Document) Total() uint {
	return (d.Value + d.ValueForfeit) - d.ValueDiscount
}

package boleto

// Payer holds the data of the person to whom the bank slip should be emmited to
type Payer struct {
	Name    string
	Address string
	Contact string
}

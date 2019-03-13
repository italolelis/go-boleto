package santander

import (
	"testing"
	"time"

	boleto "github.com/italolelis/go-boleto"
)

func TestSantander_Barcode(t *testing.T) {
	bank := New()
	bank.Account = 101888
	bank.Agency = 564
	bank.Carteira = 06
	bank.IOS = 0
	bank.Company = boleto.Company{
		Name:      "ACME Corporation",
		LegalName: "ACME Corporation Inc.",
		Address:   "Setor de Clubes Esportivos Sul (SCES) - Trecho 2 - Conjunto 31 - Lotes 1A/1B, 70200-002, Brasília, DF",
		Contact:   "acme@example.com",
		Document:  "01.122.241/0001-76",
	}

	d := boleto.Document{
		ID:            24588722,
		Value:         25000,
		ValueTax:      0,
		ValueDiscount: 0,
		ValueForfeit:  0,
		OurNumber:     77000009017,
		FebrabanType:  "dm",
		Date:          time.Now(),
		DateDue:       time.Date(2015, 03, 24, 0, 0, 0, 0, time.UTC),
		Instructions: []string{
			"Não receber após o vencimento",
			"Após vencimento, receber apenas no meu banco",
		},
		Payer: &boleto.Payer{
			Name:    "John Doe",
			Address: "Setor de Clubes Esportivos Sul (SCES) - Trecho 2 - Conjunto 31 - Lotes 1A/1B, 70200-002, Brasília, DF",
			Contact: "john.doe@example.com",
		},
	}

	b, err := bank.Barcode(d)
	if err != nil {
		t.Fatalf("an error was not expected when generating a barcode: %s", err)
	}

	digitableCode, err := b.Digitable()
	if err != nil {
		t.Fatalf("an error was not expected when retrieving the digitable code: %s", err)
	}

	t.Logf("digitable number: %s", digitableCode)
	t.Logf("barcode: %s", b.String())

	expectedBarcode := "03396637700000250000090101888007700000901706"
	expectedDigitable := "03390.09018 01888.007707 00009.017064 6 63770000025000"
	if digitableCode != expectedDigitable {
		t.Fatalf("digitable code doesn't match: \nexpected: %s \nactual: %s", expectedDigitable, digitableCode)
	}

	if b.String() != expectedBarcode {
		t.Fatalf("bardcode doesn't match: \nexpected: %s \nactual: %s", expectedBarcode, b.String())
	}
}

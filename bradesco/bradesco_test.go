package bradesco

import (
	"testing"
	"time"

	boleto "github.com/italolelis/go-boleto"
)

func TestBradesco_Carteira6(t *testing.T) {
	bank := New(
		564,
		101888,
		101888,
		06,
		&boleto.Company{
			Name:      "Nome da empresa",
			LegalName: "Razao social",
			Address:   "Endereço",
			Contact:   "Email e telefone",
			Document:  "CNPJ/CPF",
		},
	)

	d := boleto.NewDocument(
		24588722,
		77000009017,
		time.Date(2015, 03, 24, 0, 0, 0, 0, time.UTC),
	)
	d.DefineValue(25000, 0, 0)
	d.AddInstruction("Não receber após o vencimento")
	d.AddInstruction("Após vencimento, receber apenas no meu banco")
	d.To(&boleto.Payer{
		Name:    "John Doe",
		Address: "Setor de Clubes Esportivos Sul (SCES) - Trecho 2 - Conjunto 31 - Lotes 1A/1B, 70200-002, Brasília, DF",
		Contact: "john.doe@example.com",
	})

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

	expectedBarcode := "23798637700000250000564067700000901701018880"
	expectedDigitable := "23790.56406 67700.000907 17010.188801 8 63770000025000"
	if digitableCode != expectedDigitable {
		t.Fatalf("digitable code doesn't match: \nexpected: %s \nactual: %s", expectedDigitable, digitableCode)
	}

	if b.String() != expectedBarcode {
		t.Fatalf("bardcode doesn't match: \nexpected: %s \nactual: %s", expectedBarcode, b.String())
	}
}

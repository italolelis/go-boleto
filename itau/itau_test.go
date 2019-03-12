package itau

import (
	"testing"
	"time"

	boleto "github.com/italolelis/go-boleto"
)

func TestItau_Carteira107(t *testing.T) {
	bank := New(
		88888888,
		4444,
		107,
		43213,
		&boleto.Company{
			Name:      "Nome da empresa",
			LegalName: "Razao social",
			Address:   "Endereço",
			Contact:   "Email e telefone",
			Document:  "CNPJ/CPF",
		},
	)

	d := boleto.NewDocument(
		1111111,
		1111111,
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

	expectedBarcode := "34197637700000250001070111111111111114321300"
	expectedDigitable := "34191.07011 11111.111114 11143.213002 7 63770000025000"
	if digitableCode != expectedDigitable {
		t.Fatalf("digitable code doesn't match: \nexpected: %s \nactual: %s", expectedDigitable, digitableCode)
	}

	if b.String() != expectedBarcode {
		t.Fatalf("bardcode doesn't match: \nexpected: %s \nactual: %s", expectedBarcode, b.String())
	}
}

package bb

import (
	"testing"
	"time"

	boleto "github.com/italolelis/go-boleto"
)

func TestBB_Carteira17Convenio4(t *testing.T) {
	bank := New(
		88888888,
		4444,
		4321,
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
		1111,
		111111,
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

	expectedBarcode := "00191637700000250004321111111144448888888817"
	expectedDigitable := "00194.32110 11111.144442 88888.888170 1 63770000025000"
	if digitableCode != expectedDigitable {
		t.Fatalf("digitable code doesn't match: \nexpected: %s \nactual: %s", expectedDigitable, digitableCode)
	}

	if b.String() != expectedBarcode {
		t.Fatalf("bardcode doesn't match: \nexpected: %s \nactual: %s", expectedBarcode, b.String())
	}
}

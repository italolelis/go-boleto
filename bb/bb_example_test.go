package bb_test

import (
	"log"
	"time"

	boleto "github.com/italolelis/go-boleto"
	"github.com/italolelis/go-boleto/bb"
)

func ExampleNew() {
	// We create an instance of BB bank
	bank := bb.New(
		88888888, // Account Number
		4444,     // Agency Number
		4321,     // Convenio Number
		06,       // Carteira Number (supported are 4, 6 and 7)
		&boleto.Company{
			Name:      "Acme",
			LegalName: "Acme Inc",
			Address:   "WarnerBros Studios, Los Angeles, CA, US",
			Contact:   "me@acme.com",
			Document:  "CNPJ/CPF",
		},
	)

	// Creates a new document that we will attach to the bank instance later
	d := boleto.NewDocument(
		1111,   // Document ID
		111111, // Our Number
		time.Date(2015, 03, 24, 0, 0, 0, 0, time.UTC), // Due Date
	)

	// Defines the document value along with any extra cost
	d.DefineValue(25000, 0, 0)

	// Add instructions to the document
	d.AddInstruction("Não receber após o vencimento")
	d.AddInstruction("Após vencimento, receber apenas no meu banco")

	// Assings to whom this bank slip will be addressed
	d.To(&boleto.Payer{
		Name:    "John Doe",
		Address: "Setor de Clubes Esportivos Sul (SCES) - Trecho 2 - Conjunto 31 - Lotes 1A/1B, 70200-002, Brasília, DF",
		Contact: "john.doe@example.com",
	})

	b, err := bank.Barcode(d)
	if err != nil {
		log.Fatalf("there was an error generating a barcode: %s", err)
	}

	digitableCode, err := b.Digitable()
	if err != nil {
		log.Fatalf("an error was not expected when retrieving the digitable code: %s", err)
	}

	log.Printf("BB digitable number is: %s", digitableCode)
}

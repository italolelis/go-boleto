package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	boleto "github.com/italolelis/go-boleto"
	"github.com/italolelis/go-boleto/bradesco"
)

var (
	fmap = template.FuncMap{
		"formatNumber": func(valueInCents uint) (string, error) {
			dollars := valueInCents / 100
			cents := valueInCents % 100
			return fmt.Sprintf("%d,%2d", dollars, cents), nil
		},
		"formatAsDate": func(t time.Time) string {
			year, month, day := t.Date()
			return fmt.Sprintf("%d/%d/%d", day, month, year)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
		"now": func() string {
			year, month, day := time.Now().Date()
			return fmt.Sprintf("%d/%d/%d", day, month, year)
		},
	}

	templates = template.Must(template.New("main").Funcs(fmap).ParseGlob("../templates/*.tmpl"))
)

type View struct {
	Bank            boleto.Banker
	Document        *boleto.Document
	DigitableNumber string
}

func main() {
	http.HandleFunc("/", mainHandler)
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("../web/images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("../web/css"))))

	//Listen on port 8080
	http.ListenAndServe(":8000", nil)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	bank := bradesco.New(
		564,
		101888,
		101888,
		9,
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
		log.Printf("there was an error generating the digitable code: %s", err)
		http.Error(w, "there was error creating the barcode", http.StatusInternalServerError)
		return
	}

	number, err := b.Digitable()
	if err != nil {
		log.Printf("there was an error generating the digitable code: %s", err)
		http.Error(w, "there was an error generating the digitable code", http.StatusInternalServerError)
		return
	}

	templates.ExecuteTemplate(w, "main", &View{
		Bank:            bank,
		Document:        d,
		DigitableNumber: number,
	})
}

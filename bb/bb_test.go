package bb

import "testing"

func TestBB_BarcodeGeneratedSuccessfully(t *testing.T) {
	bb := New()
	bb.Account = 5579
	bb.Agency = 000000

}

package validate

import "testing"

func TestJSON_OK(t *testing.T) {
	if err := JSON([]byte(`{"a":1}`)); err != nil {
		t.Fatal(err)
	}
}

func TestJSON_Invalid(t *testing.T) {
	if err := JSON([]byte(`{`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestXML_OK(t *testing.T) {
	if err := XML([]byte(`<root><a>1</a></root>`)); err != nil {
		t.Fatal(err)
	}
}

func TestSOAP_OK(t *testing.T) {
	s := `<?xml version="1.0"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/"><s:Body><x/></s:Body></s:Envelope>`
	if err := SOAP([]byte(s)); err != nil {
		t.Fatal(err)
	}
}

func TestSOAP_NoEnvelope(t *testing.T) {
	if err := SOAP([]byte(`<root/>`)); err == nil {
		t.Fatal("expected error")
	}
}

package googleauth

import (
	"html/template"
	"testing"
)

func TestEmbeddedTemplates_Parse(t *testing.T) {
	cases := []struct {
		name string
		src  string
	}{
		{name: "accounts", src: accountsTemplate},
		{name: "success_new", src: successTemplateNew},
		{name: "success", src: successTemplate},
		{name: "error", src: errorTemplate},
		{name: "cancelled", src: cancelledTemplate},
	}
	for _, tc := range cases {
		if tc.src == "" {
			t.Fatalf("%s template is empty", tc.name)
		}
		if _, err := template.New(tc.name).Parse(tc.src); err != nil {
			t.Fatalf("%s parse: %v", tc.name, err)
		}
	}
}

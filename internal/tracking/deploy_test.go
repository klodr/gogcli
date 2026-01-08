package tracking

import "testing"

func TestSanitizeWorkerName(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{input: "Test@Example.com", want: "test-example-com"},
		{input: " gog--tracker ", want: "gog-tracker"},
		{input: "___", want: ""},
	}
	for _, tc := range cases {
		if got := SanitizeWorkerName(tc.input); got != tc.want {
			t.Fatalf("SanitizeWorkerName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}

	long := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	if got := SanitizeWorkerName(long); len(got) != 63 {
		t.Fatalf("expected max length 63, got %d (%q)", len(got), got)
	}
}

func TestParseDatabaseID(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{input: `database_id = "abc-123"`, want: "abc-123"},
		{input: `database_id: abc-123`, want: "abc-123"},
		{input: `Database ID: abc-123`, want: "abc-123"},
		{input: `database_id: "xyz-789"`, want: "xyz-789"},
		{input: `Database ID: 12345`, want: "12345"},
		{input: `database_id = "with-dash"`, want: "with-dash"},
	}
	for _, tc := range cases {
		if got := parseDatabaseID(tc.input); got != tc.want {
			t.Fatalf("parseDatabaseID(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}

	if got := parseDatabaseID("nope"); got != "" {
		t.Fatalf("expected empty id, got %q", got)
	}
}

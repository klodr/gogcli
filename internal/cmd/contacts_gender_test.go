package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"google.golang.org/api/people/v1"

	"github.com/steipete/gogcli/internal/ui"
)

// ---------------------------------------------------------------------------
// contacts update --gender
// ---------------------------------------------------------------------------

func TestContactsUpdate_Gender_Set(t *testing.T) {
	var gotGetFields string
	var gotUpdateFields string
	var gotGenderValue string

	svc, closeSrv := newPeopleService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "people/c1") && r.Method == http.MethodGet && !strings.Contains(r.URL.Path, ":"):
			gotGetFields = r.URL.Query().Get("personFields")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"resourceName": "people/c1",
				"names":        []map[string]any{{"givenName": "Ada", "familyName": "Lovelace"}},
			})
			return
		case strings.Contains(r.URL.Path, ":updateContact") && (r.Method == http.MethodPatch || r.Method == http.MethodPost):
			gotUpdateFields = r.URL.Query().Get("updatePersonFields")
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if genders, ok := body["genders"].([]any); ok && len(genders) > 0 {
				if first, ok := genders[0].(map[string]any); ok {
					gotGenderValue = strings.TrimSpace(primaryValue(first, "value"))
				}
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"resourceName": "people/c1"})
			return
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(closeSrv)
	stubPeopleServices(t, svc)

	u, err := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if err != nil {
		t.Fatalf("ui.New: %v", err)
	}
	ctx := ui.WithUI(context.Background(), u)

	if err := runKong(t, &ContactsUpdateCmd{}, []string{"people/c1", "--gender", "female"}, ctx, &RootFlags{Account: "a@b.com"}); err != nil {
		t.Fatalf("runKong: %v", err)
	}

	if !strings.Contains(gotGetFields, "genders") {
		t.Fatalf("missing genders in people.get fields: %q", gotGetFields)
	}
	if !strings.Contains(gotUpdateFields, "genders") {
		t.Fatalf("missing genders in updatePersonFields: %q", gotUpdateFields)
	}
	if gotGenderValue != "female" {
		t.Fatalf("unexpected gender payload: %q, want %q", gotGenderValue, "female")
	}
}

func TestContactsUpdate_Gender_CustomValue(t *testing.T) {
	var gotGenderValue string

	svc, closeSrv := newPeopleService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "people/c1") && r.Method == http.MethodGet && !strings.Contains(r.URL.Path, ":"):
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"resourceName": "people/c1"})
			return
		case strings.Contains(r.URL.Path, ":updateContact") && (r.Method == http.MethodPatch || r.Method == http.MethodPost):
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if genders, ok := body["genders"].([]any); ok && len(genders) > 0 {
				if first, ok := genders[0].(map[string]any); ok {
					gotGenderValue = strings.TrimSpace(primaryValue(first, "value"))
				}
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"resourceName": "people/c1"})
			return
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(closeSrv)
	stubPeopleServices(t, svc)

	u, err := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if err != nil {
		t.Fatalf("ui.New: %v", err)
	}
	ctx := ui.WithUI(context.Background(), u)

	if err := runKong(t, &ContactsUpdateCmd{}, []string{"people/c1", "--gender", "nonbinary"}, ctx, &RootFlags{Account: "a@b.com"}); err != nil {
		t.Fatalf("runKong: %v", err)
	}

	if gotGenderValue != "nonbinary" {
		t.Fatalf("unexpected custom gender payload: %q, want %q", gotGenderValue, "nonbinary")
	}
}

func TestContactsUpdate_Gender_Clear(t *testing.T) {
	var gotUpdateFields string

	svc, closeSrv := newPeopleService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "people/c1") && r.Method == http.MethodGet && !strings.Contains(r.URL.Path, ":"):
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"resourceName": "people/c1",
				"genders":      []map[string]any{{"value": "male", "formattedValue": "Male"}},
			})
			return
		case strings.Contains(r.URL.Path, ":updateContact") && (r.Method == http.MethodPatch || r.Method == http.MethodPost):
			gotUpdateFields = r.URL.Query().Get("updatePersonFields")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"resourceName": "people/c1"})
			return
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(closeSrv)
	stubPeopleServices(t, svc)

	u, err := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if err != nil {
		t.Fatalf("ui.New: %v", err)
	}
	ctx := ui.WithUI(context.Background(), u)

	if err := runKong(t, &ContactsUpdateCmd{}, []string{"people/c1", "--gender", ""}, ctx, &RootFlags{Account: "a@b.com"}); err != nil {
		t.Fatalf("runKong: %v", err)
	}

	if !strings.Contains(gotUpdateFields, "genders") {
		t.Fatalf("missing genders in clear updatePersonFields: %q", gotUpdateFields)
	}
}

// ---------------------------------------------------------------------------
// contacts create --gender
// ---------------------------------------------------------------------------

func TestContactsCreate_Gender_Set(t *testing.T) {
	var gotGenderValue string

	svc, closeSrv := newPeopleService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, ":createContact") && r.Method == http.MethodPost:
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if genders, ok := body["genders"].([]any); ok && len(genders) > 0 {
				if first, ok := genders[0].(map[string]any); ok {
					gotGenderValue = strings.TrimSpace(primaryValue(first, "value"))
				}
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"resourceName": "people/c1"})
			return
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(closeSrv)
	stubPeopleServices(t, svc)

	u, err := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if err != nil {
		t.Fatalf("ui.New: %v", err)
	}
	ctx := ui.WithUI(context.Background(), u)

	if err := runKong(t, &ContactsCreateCmd{}, []string{"--given", "Ada", "--gender", "female"}, ctx, &RootFlags{Account: "a@b.com"}); err != nil {
		t.Fatalf("runKong: %v", err)
	}

	if gotGenderValue != "female" {
		t.Fatalf("unexpected gender payload: %q, want %q", gotGenderValue, "female")
	}
}

func TestContactsCreate_Gender_Omitted(t *testing.T) {
	var gotBody map[string]any

	svc, closeSrv := newPeopleService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, ":createContact") && r.Method == http.MethodPost:
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"resourceName": "people/c1"})
			return
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(closeSrv)
	stubPeopleServices(t, svc)

	u, err := ui.New(ui.Options{Stdout: io.Discard, Stderr: io.Discard, Color: "never"})
	if err != nil {
		t.Fatalf("ui.New: %v", err)
	}
	ctx := ui.WithUI(context.Background(), u)

	if err := runKong(t, &ContactsCreateCmd{}, []string{"--given", "Ada"}, ctx, &RootFlags{Account: "a@b.com"}); err != nil {
		t.Fatalf("runKong: %v", err)
	}

	if _, present := gotBody["genders"]; present {
		t.Fatalf("genders should be absent from payload when --gender not provided, got: %v", gotBody["genders"])
	}
}

// ---------------------------------------------------------------------------
// contacts get: gender displayed in text output
// ---------------------------------------------------------------------------

func TestContactsGet_Gender_DisplayedInTextOutput(t *testing.T) {
	svc, closeSrv := newPeopleService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "people/c1") && r.Method == http.MethodGet && !strings.Contains(r.URL.Path, ":"):
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"resourceName": "people/c1",
				"names":        []map[string]any{{"givenName": "Ada", "familyName": "Lovelace"}},
				"genders":      []map[string]any{{"value": "female", "formattedValue": "Female"}},
			})
			return
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(closeSrv)
	stubPeopleServices(t, svc)

	u, uiErr := ui.New(ui.Options{Stdout: nil, Stderr: io.Discard, Color: "never"})
	if uiErr != nil {
		t.Fatalf("ui.New: %v", uiErr)
	}

	out := captureStdout(t, func() {
		u2, _ := ui.New(ui.Options{Stderr: io.Discard, Color: "never"})
		ctx := ui.WithUI(context.Background(), u2)
		_ = u
		if err := runKong(t, &ContactsGetCmd{}, []string{"people/c1"}, ctx, &RootFlags{Account: "a@b.com"}); err != nil {
			t.Fatalf("runKong: %v", err)
		}
	})

	if !strings.Contains(out, "gender") || !strings.Contains(out, "Female") {
		t.Fatalf("expected 'gender' and 'Female' in get output, got: %q", out)
	}
}

func TestContactsGet_Gender_AbsentWhenEmpty(t *testing.T) {
	svc, closeSrv := newPeopleService(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "people/c1") && r.Method == http.MethodGet && !strings.Contains(r.URL.Path, ":"):
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"resourceName": "people/c1",
				"names":        []map[string]any{{"givenName": "Ada", "familyName": "Lovelace"}},
			})
			return
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(closeSrv)
	stubPeopleServices(t, svc)

	out := captureStdout(t, func() {
		u, _ := ui.New(ui.Options{Stderr: io.Discard, Color: "never"})
		ctx := ui.WithUI(context.Background(), u)
		if err := runKong(t, &ContactsGetCmd{}, []string{"people/c1"}, ctx, &RootFlags{Account: "a@b.com"}); err != nil {
			t.Fatalf("runKong: %v", err)
		}
	})

	if strings.Contains(out, "gender") {
		t.Fatalf("gender line should be absent when no gender set, got: %q", out)
	}
}

// ---------------------------------------------------------------------------
// primaryGender helper unit tests
// ---------------------------------------------------------------------------

func TestPrimaryGender_FormattedValuePreferred(t *testing.T) {
	p := &people.Person{
		Genders: []*people.Gender{
			{Value: "female", FormattedValue: "Female"},
		},
	}
	if got := primaryGender(p); got != "Female" {
		t.Fatalf("primaryGender: got %q, want %q", got, "Female")
	}
}

func TestPrimaryGender_FallbackToValue(t *testing.T) {
	p := &people.Person{
		Genders: []*people.Gender{
			{Value: "nonbinary"},
		},
	}
	if got := primaryGender(p); got != "nonbinary" {
		t.Fatalf("primaryGender: got %q, want %q", got, "nonbinary")
	}
}

func TestPrimaryGender_EmptyWhenNone(t *testing.T) {
	p := &people.Person{}
	if got := primaryGender(p); got != "" {
		t.Fatalf("primaryGender: expected empty, got %q", got)
	}
}

func TestPrimaryGender_EmptyWhenNilEntry(t *testing.T) {
	p := &people.Person{Genders: []*people.Gender{nil}}
	if got := primaryGender(p); got != "" {
		t.Fatalf("primaryGender nil entry: expected empty, got %q", got)
	}
}

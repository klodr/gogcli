package googleauth

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/steipete/gogcli/internal/secrets"
)

type fakeStore struct {
	tokens       []secrets.Token
	defaultEmail string

	setDefaultCalled string
	deleteCalled     string
}

func (s *fakeStore) Keys() ([]string, error)                { return nil, nil }
func (s *fakeStore) SetToken(string, secrets.Token) error   { return nil }
func (s *fakeStore) GetToken(string) (secrets.Token, error) { return secrets.Token{}, nil }
func (s *fakeStore) DeleteToken(email string) error {
	s.deleteCalled = email
	return nil
}

func (s *fakeStore) ListTokens() ([]secrets.Token, error) {
	return append([]secrets.Token(nil), s.tokens...), nil
}
func (s *fakeStore) GetDefaultAccount() (string, error) { return s.defaultEmail, nil }
func (s *fakeStore) SetDefaultAccount(email string) error {
	s.setDefaultCalled = email
	s.defaultEmail = email
	return nil
}

func TestManageServer_HandleAccountsPage(t *testing.T) {
	ms := &ManageServer{
		csrfToken: "csrf",
	}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ms.handleAccountsPage(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("content-type: %q", ct)
	}
	body := rr.Body.String()
	if strings.TrimSpace(body) == "" {
		tmpl, err := template.New("accounts").Parse(accountsTemplate)
		if err != nil {
			t.Fatalf("expected body, parse err=%v", err)
		}
		var buf bytes.Buffer
		execErr := tmpl.Execute(&buf, struct{ CSRFToken string }{CSRFToken: "csrf"})
		t.Fatalf("expected body; handler wrote 0 bytes; direct execute bytes=%d err=%v", buf.Len(), execErr)
	}
	if !strings.Contains(body, "csrfToken") || !strings.Contains(body, "const csrfToken") {
		t.Fatalf("expected csrf js in body")
	}
	if !strings.Contains(body, "'csrf'") && !strings.Contains(body, "\"csrf\"") {
		excerpt := body
		if len(excerpt) > 200 {
			excerpt = excerpt[:200]
		}
		t.Fatalf("expected rendered token, body excerpt=%q", excerpt)
	}
}

func TestManageServer_HandleListAccounts_DefaultFirst(t *testing.T) {
	store := &fakeStore{
		tokens: []secrets.Token{
			{Email: "a@b.com", Services: []string{"gmail"}},
			{Email: "c@d.com", Services: []string{"drive"}},
		},
	}
	ms := &ManageServer{
		csrfToken: "csrf",
		store:     store,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	ms.handleListAccounts(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status: %d", rr.Code)
	}
	var parsed struct {
		Accounts []AccountInfo `json:"accounts"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("json parse: %v", err)
	}
	if len(parsed.Accounts) != 2 || !parsed.Accounts[0].IsDefault || parsed.Accounts[1].IsDefault {
		t.Fatalf("unexpected defaults: %#v", parsed.Accounts)
	}
}

func TestManageServer_HandleListAccounts_DefaultExplicit(t *testing.T) {
	store := &fakeStore{
		tokens: []secrets.Token{
			{Email: "a@b.com", Services: []string{"gmail"}},
			{Email: "c@d.com", Services: []string{"drive"}},
		},
		defaultEmail: "c@d.com",
	}
	ms := &ManageServer{
		csrfToken: "csrf",
		store:     store,
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	ms.handleListAccounts(rr, req)

	var parsed struct {
		Accounts []AccountInfo `json:"accounts"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("json parse: %v", err)
	}
	if len(parsed.Accounts) != 2 || parsed.Accounts[0].IsDefault || !parsed.Accounts[1].IsDefault {
		t.Fatalf("unexpected defaults: %#v", parsed.Accounts)
	}
}

func TestManageServer_HandleOAuthCallback_ErrorAndValidation(t *testing.T) {
	ms := &ManageServer{
		csrfToken:  "csrf",
		oauthState: "state1",
	}
	// Need a listener for redirectURI generation even though we don't reach exchange.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { _ = ln.Close() })
	ms.listener = ln

	t.Run("cancelled", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/oauth2/callback?error=access_denied", nil)
		ms.handleOAuthCallback(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status: %d", rr.Code)
		}
	})

	t.Run("state mismatch", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/oauth2/callback?state=nope&code=abc", nil)
		ms.handleOAuthCallback(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status: %d", rr.Code)
		}
	})

	t.Run("missing code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/oauth2/callback?state=state1", nil)
		ms.handleOAuthCallback(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status: %d", rr.Code)
		}
	})
}

func TestManageServer_HandleSetDefault_AndRemove(t *testing.T) {
	store := &fakeStore{
		tokens: []secrets.Token{{Email: "a@b.com"}},
	}
	ms := &ManageServer{
		csrfToken: "csrf",
		store:     store,
	}

	t.Run("set-default csrf", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/set-default", bytes.NewReader([]byte(`{"email":"a@b.com"}`)))
		req.Header.Set("X-CSRF-Token", "nope")
		ms.handleSetDefault(rr, req)
		if rr.Code != http.StatusForbidden {
			t.Fatalf("status: %d", rr.Code)
		}
	})

	t.Run("set-default ok", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/set-default", bytes.NewReader([]byte(`{"email":"a@b.com"}`)))
		req.Header.Set("X-CSRF-Token", "csrf")
		ms.handleSetDefault(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status: %d body=%s", rr.Code, rr.Body.String())
		}
		if store.setDefaultCalled != "a@b.com" {
			t.Fatalf("expected setDefaultCalled")
		}
	})

	t.Run("remove ok", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/remove-account", bytes.NewReader([]byte(`{"email":"a@b.com"}`)))
		req.Header.Set("X-CSRF-Token", "csrf")
		ms.handleRemoveAccount(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status: %d", rr.Code)
		}
		if store.deleteCalled != "a@b.com" {
			t.Fatalf("expected deleteCalled")
		}
	})
}

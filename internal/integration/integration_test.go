//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/steipete/gogcli/internal/googleapi"
	"github.com/steipete/gogcli/internal/googleauth"
	"github.com/steipete/gogcli/internal/secrets"
)

func TestDriveSmoke(t *testing.T) {
	account := os.Getenv("GOG_IT_ACCOUNT")
	if account == "" {
		t.Skip("set GOG_IT_ACCOUNT")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	svc, err := googleapi.NewDrive(ctx, account)
	if err != nil {
		t.Fatalf("NewDrive: %v", err)
	}
	_, err = svc.Files.List().
		Q("trashed = false").
		PageSize(1).
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Fields("files(id)").
		Do()
	if err != nil {
		t.Fatalf("Drive list: %v", err)
	}
}

func TestCalendarSmoke(t *testing.T) {
	account := os.Getenv("GOG_IT_ACCOUNT")
	if account == "" {
		t.Skip("set GOG_IT_ACCOUNT")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	svc, err := googleapi.NewCalendar(ctx, account)
	if err != nil {
		t.Fatalf("NewCalendar: %v", err)
	}
	_, err = svc.CalendarList.List().MaxResults(1).Do()
	if err != nil {
		t.Fatalf("Calendar list: %v", err)
	}
}

func TestGmailSmoke(t *testing.T) {
	account := os.Getenv("GOG_IT_ACCOUNT")
	if account == "" {
		t.Skip("set GOG_IT_ACCOUNT")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	svc, err := googleapi.NewGmail(ctx, account)
	if err != nil {
		t.Fatalf("NewGmail: %v", err)
	}
	_, err = svc.Users.Labels.List("me").Do()
	if err != nil {
		t.Fatalf("Gmail labels: %v", err)
	}
}

func TestAuthRefreshTokenSmoke(t *testing.T) {
	account := os.Getenv("GOG_IT_ACCOUNT")
	if account == "" {
		t.Skip("set GOG_IT_ACCOUNT")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	store, err := secrets.OpenDefault()
	if err != nil {
		t.Fatalf("OpenDefault: %v", err)
	}
	tok, err := store.GetToken(account)
	if err != nil {
		t.Fatalf("GetToken: %v", err)
	}

	if err := googleauth.CheckRefreshToken(ctx, tok.RefreshToken, tok.Scopes, 15*time.Second); err != nil {
		t.Fatalf("CheckRefreshToken: %v", err)
	}
}

func TestContactsSmoke(t *testing.T) {
	account := os.Getenv("GOG_IT_ACCOUNT")
	if account == "" {
		t.Skip("set GOG_IT_ACCOUNT")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	svc, err := googleapi.NewPeopleContacts(ctx, account)
	if err != nil {
		t.Fatalf("NewPeople: %v", err)
	}
	_, err = svc.People.Connections.List("people/me").PersonFields("names").PageSize(1).Do()
	if err != nil {
		t.Fatalf("People connections: %v", err)
	}
}

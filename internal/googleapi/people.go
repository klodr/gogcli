package googleapi

import (
	"context"

	"google.golang.org/api/people/v1"
)

const (
	scopeContactsWrite   = "https://www.googleapis.com/auth/contacts"
	scopeContactsOtherRO = "https://www.googleapis.com/auth/contacts.other.readonly"
	scopeDirectoryRO     = "https://www.googleapis.com/auth/directory.readonly"
)

func NewPeopleContacts(ctx context.Context, email string) (*people.Service, error) {
	opts, err := optionsForAccountScopes(ctx, "contacts", email, []string{scopeContactsWrite})
	if err != nil {
		return nil, err
	}
	return people.NewService(ctx, opts...)
}

func NewPeopleOtherContacts(ctx context.Context, email string) (*people.Service, error) {
	opts, err := optionsForAccountScopes(ctx, "contacts", email, []string{scopeContactsOtherRO})
	if err != nil {
		return nil, err
	}
	return people.NewService(ctx, opts...)
}

func NewPeopleDirectory(ctx context.Context, email string) (*people.Service, error) {
	opts, err := optionsForAccountScopes(ctx, "contacts", email, []string{scopeDirectoryRO})
	if err != nil {
		return nil, err
	}
	return people.NewService(ctx, opts...)
}

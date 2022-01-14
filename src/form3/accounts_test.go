package form3

import (
	"context"
	"form3-interview-accountapi/src/models"
	"form3-interview-accountapi/src/utils"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func createTestAccount() *models.AccountData {
	account := new(models.AccountData)
	account.Type = "accounts"
	account.ID = uuid.New().String()
	account.OrganisationID = uuid.New().String()
	account.Version = utils.Int64Addr(0)
	account.Attributes = &models.AccountAttributes{
		Country:                 utils.StringAddr("GB"),
		BaseCurrency:            "GBP",
		BankID:                  "400300",
		BankIDCode:              "GBDSC",
		Bic:                     "NWBKGB22",
		Name:                    []string{"Samantha Holder"},
		AccountClassification:   utils.StringAddr("Personal"),
		JointAccount:            utils.BoolAddr(false),
		AccountMatchingOptOut:   utils.BoolAddr(false),
		SecondaryIdentification: "A1B2C3D4",
	}
	return account
}

func TestCreateAccountSuccess(t *testing.T) {
	// Create an account
	client := NewClient()

	want := createTestAccount()
	ctx := context.Background()
	got, _, err := client.CreateAccount(ctx, want)
	if err != nil {
		t.Errorf("accounts::CreateAccount returned error: %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("accounts::CreateAccount returned %+v, want %+v", got, want)
	}
}

func TestCreateAccountFail(t *testing.T) {
	// Create an account which already exists
	client := NewClient()

	want := createTestAccount()
	ctx := context.Background()
	_, _, err := client.CreateAccount(ctx, want)
	if err != nil {
		t.Errorf("accounts::CreateAccount returned error: %v", err)
	}
	// Try to create the same account again
	got, resp, err := client.CreateAccount(ctx, want)
	if err != nil {
		t.Errorf("accounts::CreateAccount returned error: %v", err)
	}
	wantStatus := http.StatusConflict
	if resp.StatusCode != wantStatus {
		t.Errorf("accounts::GetAccount with duplicate account returned status code %+v, want %+v", resp.StatusCode, wantStatus)
	}
	if got != nil {
		t.Errorf("accounts::GetAccount with duplicate account returned %+v, want nil", got)
	}
}

func TestGetAccountSuccess(t *testing.T) {
	// Create an account and confirm it exists.
	client := NewClient()

	want := createTestAccount()
	ctx := context.Background()
	// Create account
	_, _, err := client.CreateAccount(ctx, want)
	if err != nil {
		t.Errorf("accounts::CreateAccount returned error: %v", err)
	}
	// Get account
	got, _, err := client.GetAccount(ctx, want.ID)
	if err != nil {
		t.Errorf("accounts::GetAccount returned error: %v", err)
	}
	if !cmp.Equal(got, want) {
		t.Errorf("accounts::GetAccount returned %+v, want %+v", got, want)
	}
}

func TestGetAccountFail(t *testing.T) {
	// Get an account which does not exist
	client := NewClient()

	want := createTestAccount()
	ctx := context.Background()
	got, resp, err := client.GetAccount(ctx, want.ID)
	if err != nil {
		t.Errorf("accounts::GetAccount returned error: %+v", err)
	}
	wantStatus := http.StatusNotFound
	if resp.StatusCode != wantStatus {
		t.Errorf("accounts::GetAccount with invalid id returned status code %+v, want %+v", resp.StatusCode, wantStatus)
	}
	if got != nil {
		t.Errorf("accounts::GetAccount with invalid id returned %+v, want nil", got)
	}
}

func TestGetAccountList(t *testing.T) {
	// Create multiple accounts and confirm they all exist.
	client := NewClient()

	var accounts []*models.AccountData
	for i := 0; i < 5; i++ {
		want := createTestAccount()
		accounts = append(accounts, want)
	}
	ctx := context.Background()
	for _, want := range accounts {
		_, _, err := client.CreateAccount(ctx, want)
		if err != nil {
			t.Errorf("accounts::CreateAccount returned error: %v", err)
		}
	}

	accountList, _, err := client.ListAccounts(ctx)
	if err != nil {
		t.Errorf("accounts::ListAccounts returned error: %v", err)
	}
	lookup := make(map[string]*models.AccountData)
	for _, got := range accountList {
		lookup[got.ID] = got
	}
	for _, want := range accounts {
		got, ok := lookup[want.ID]
		if !ok {
			t.Errorf("accounts::GetAccountList does not contain: %v", want)
			continue
		}
		if !cmp.Equal(got, want) {
			t.Errorf("accounts::GetAccount returned %+v, want %+v", got, want)
		}
	}
}

func TestDeleteAccountSuccess(t *testing.T) {
	// Create an account, delete it, and confirm it is deleted.
	client := NewClient()

	want := createTestAccount()
	ctx := context.Background()
	// Create account
	_, _, err := client.CreateAccount(ctx, want)
	if err != nil {
		t.Errorf("accounts::CreateAccount returned error: %v", err)
	}
	// Delete account
	resp, err := client.DeleteAccount(ctx, want.ID)
	if err != nil {
		t.Errorf("accounts::DeleteAccount returned error: %v", err)
	}
	wantStatus := http.StatusNoContent
	if resp.StatusCode != wantStatus {
		t.Errorf("accounts::DeleteAccount returned status code %+v, want %+v", resp.StatusCode, wantStatus)
	}
	// Get account
	got, _, err := client.GetAccount(ctx, want.ID)
	if err != nil {
		t.Errorf("accounts::GetAccount returned error: %v", err)
	}
	if got != nil {
		t.Errorf("accounts::GetAccount returned %+v, want nil", got)
	}
}

func TestDeleteAccountFail(t *testing.T) {
	// Delete an account which does not exist.
	client := NewClient()

	want := createTestAccount()
	ctx := context.Background()
	resp, err := client.DeleteAccount(ctx, want.ID)
	if err != nil {
		t.Errorf("accounts::DeleteAccount returned error: %v", err)
	}
	wantStatus := http.StatusNotFound
	if resp.StatusCode != wantStatus {
		t.Errorf("accounts::DeleteAccount returned status code %+v, want %+v", resp.StatusCode, wantStatus)
	}
}

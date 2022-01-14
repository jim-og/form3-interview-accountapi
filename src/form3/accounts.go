package form3

import (
	"context"
	"fmt"
	"form3-interview-accountapi/src/models"
	"net/http"
)

const (
	accountsURL = "organisation/accounts"
	version     = 0
)

// CreateAccount creates a new account
func (client *Client) CreateAccount(ctx context.Context, account *models.AccountData) (*models.AccountData, *http.Response, error) {
	u := accountsURL
	accountReq := &models.AccountModel{
		Data: account,
	}
	req, err := client.NewRequest(http.MethodPost, u, accountReq)
	if err != nil {
		return nil, nil, err
	}

	createdAccount := new(models.AccountModel)
	resp, err := client.Do(ctx, req, createdAccount)
	if err != nil {
		return nil, resp, err
	}

	return createdAccount.Data, resp, nil
}

// ListAccounts returns all accounts
func (client *Client) ListAccounts(ctx context.Context) ([]*models.AccountData, *http.Response, error) {
	u := accountsURL
	req, err := client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	accountList := new(models.AccountModelList)
	resp, err := client.Do(ctx, req, accountList)
	if err != nil {
		return nil, resp, err
	}

	return accountList.Data, resp, nil
}

// GetAccount returns the account with the specified accountID
func (client *Client) GetAccount(ctx context.Context, accountID string) (*models.AccountData, *http.Response, error) {
	u := fmt.Sprintf("%v/%v", accountsURL, accountID)
	req, err := client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}
	account := new(models.AccountModel)
	resp, err := client.Do(ctx, req, account)
	if err != nil {
		return nil, resp, err
	}

	return account.Data, resp, nil
}

// DeleteAccount deletes the account with the specified accountID
func (client *Client) DeleteAccount(ctx context.Context, accountID string) (*http.Response, error) {
	u := fmt.Sprintf("%v/%v?version=%v", accountsURL, accountID, version)
	req, err := client.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(ctx, req, nil)
}

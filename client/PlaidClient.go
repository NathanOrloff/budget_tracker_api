package client

import (
	"budget_tracket/database/models"
	"context"
	"fmt"
	"os"
	"strings"

	plaid "github.com/plaid/plaid-go"
)

type PlaidClient struct {
	Client *plaid.APIClient
}

func NewPlaidClient() *PlaidClient {
	config := plaid.NewConfiguration()
	config.Host = os.Getenv("PLAID_ENV")
	config.AddDefaultHeader("PLAID-CLIENT-ID", os.Getenv("PLAID_CLIENT_ID"))
	config.AddDefaultHeader("PLAID-SECRET", os.Getenv("PLAID_SECRET"))

	client := PlaidClient{
		Client: plaid.NewAPIClient(config),
	}
	return &client
}

func (p *PlaidClient) CreateLinkToken(ctx context.Context, userID string) (string, error) {
	op := "CreateLinkToken"

	countryCodes := getCountryCodes()
	redirectUri := os.Getenv("PLAID_REDIRECT_URI")
	products := getPlaidProducts()

	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: userID,
	}

	request := plaid.NewLinkTokenCreateRequest(
		"App name",
		"en",
		countryCodes,
		user,
	)

	if redirectUri != "" {
		request.SetRedirectUri(redirectUri)
	}

	request.SetProducts(products)

	resp, _, err := p.Client.PlaidApi.LinkTokenCreate(ctx).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetLinkToken(), nil
}

func (p *PlaidClient) ExchangePublicToken(ctx context.Context, publicToken string) (string, string, error) {
	op := "ExchangePublicToken"

	resp, _, err := p.Client.PlaidApi.ItemPublicTokenExchange(ctx).ItemPublicTokenExchangeRequest(
		*plaid.NewItemPublicTokenExchangeRequest(publicToken),
	).Execute()
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetAccessToken(), resp.GetItemId(), nil
}

func (p *PlaidClient) SyncTransactions(ctx context.Context, accessToken string, cursor *string) ([]models.Transaction, []models.Transaction, []models.Transaction, *string, error) {
	op := "SyncTransactions"

	var added []models.Transaction
	var modified []models.Transaction
	var removed []models.Transaction
	hasMore := true

	for hasMore {
		request := plaid.NewTransactionsSyncRequest(accessToken)
		if cursor != nil {
			request.SetCursor(*cursor)
		}

		resp, _, err := p.Client.PlaidApi.TransactionsSync(ctx).TransactionsSyncRequest(*request).Execute()
		if err != nil {
			return []models.Transaction{}, []models.Transaction{}, []models.Transaction{}, nil, fmt.Errorf("%s: %w", op, err)
		}

		nextCursor := resp.GetNextCursor()
		cursor = &nextCursor

		added = append(added, marshallTransactions(resp.GetAdded())...)
		modified = append(modified, marshallTransactions(resp.GetModified())...)
		removed = append(removed, marshallTransactions(resp.GetRemoved())...)
		hasMore = resp.GetHasMore()
	}

	return added, modified, removed, cursor, nil
}

func marshalTransactions(plaidTransactions []plaid.Transaction) ([]models.Transaction, error) {
	op := "marshalTransactions"
	var transactions []models.Transaction
	for _, trans := range plaidTransactions {
		transaction, err := models.Transaction.Marshal(trans)
		if err != nil {
			return []models.Transaction{}, fmt.Errorf("%s: %w", op, err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}

func getCountryCodes() []plaid.CountryCode {
	countryCodes := []plaid.CountryCode{}

	countryCodeStrs := strings.Split(os.Getenv("PLAID_COUNTRY_CODES"), ",")
	for _, countryCodeStr := range countryCodeStrs {
		countryCodes = append(countryCodes, plaid.CountryCode(countryCodeStr))
	}

	return countryCodes
}

func getPlaidProducts() []plaid.Products {
	plaidProducts := []plaid.Products{}

	products := strings.Split(os.Getenv("PLAID_PRODUCTS"), ",")
	for _, product := range products {
		plaidProducts = append(plaidProducts, plaid.Product(product))
	}

	return plaidProducts
}

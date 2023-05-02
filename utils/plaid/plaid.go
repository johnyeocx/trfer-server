package my_plaid

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/plaid/plaid-go/v11/plaid"
)

func CreateClient() (*plaid.APIClient) {
	clientId := os.Getenv("PLAID_CLIENT_ID")
	secret := os.Getenv("PLAID_SANDBOX_SECRET")
	// secret := os.Getenv("PLAID_DEV_SECRET")

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientId)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)
	configuration.UseEnvironment(plaid.Sandbox)
	// configuration.UseEnvironment(plaid.Development)

	client := plaid.NewAPIClient(configuration)
	return client
}

func CreateLinkToken(plaidCli *plaid.APIClient, userId int) (string, error) {
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: fmt.Sprintf("%d", userId),
	}

	request := plaid.NewLinkTokenCreateRequest(
	  "trfer.me",
	  "en",
	  []plaid.CountryCode{plaid.COUNTRYCODE_GB},
	  user,
	)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH})
	request.SetLinkCustomizationName("default")
	// request.SetWebhook("https://google.com")

	resp, _, err := plaidCli.PlaidApi.LinkTokenCreate(context.TODO()).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}
	linkToken := resp.GetLinkToken()
	
	return linkToken, err
}

func GetAuthAccessToken(
	plaidCli *plaid.APIClient, 
	publicToken string,
) (string, error) {
	exchangePublicTokenReq := plaid.NewItemPublicTokenExchangeRequest(publicToken)
	exchangePublicTokenResp, _, err := plaidCli.PlaidApi.ItemPublicTokenExchange(context.TODO()).ItemPublicTokenExchangeRequest(
  	*exchangePublicTokenReq).Execute()

	if err != nil {
		return "", err
	}

	accessToken := exchangePublicTokenResp.GetAccessToken()
	return accessToken, nil
}

func GetBACNumbers(
	plaidCli *plaid.APIClient, 
	accessToken string,
) (*plaid.NumbersBACS, error) {
	authGetResp, _, err := plaidCli.PlaidApi.AuthGet(context.TODO()).AuthGetRequest(
		*plaid.NewAuthGetRequest(accessToken),
	).Execute()

	if err != nil {
		return nil, err
	}
	numbers := authGetResp.GetNumbers()
	if len(numbers.Bacs) == 0 {
		return nil, errors.New("No account found")
	}

	bacs := numbers.Bacs[0]
	return &bacs, nil
}

func CreatePaymentRecipient(
	plaidCli *plaid.APIClient, 
	fullName string, 
	acctNumber string, 
	sortCode string,
) (string, error) {
	request := plaid.NewPaymentInitiationRecipientCreateRequest(fullName)

	request.SetBacs(plaid.RecipientBACSNullable{
		Account:  plaid.PtrString(acctNumber),
		SortCode: plaid.PtrString(sortCode),
	})

	paymentRecipientCreateResp, _, err := plaidCli.PlaidApi.PaymentInitiationRecipientCreate(context.TODO()).PaymentInitiationRecipientCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}
	return paymentRecipientCreateResp.RecipientId, nil
}

func CreatePayment(
	plaidCli *plaid.APIClient, 
	recipientID string, 
	amount float64,
) (string, error) {
	request := plaid.NewPaymentInitiationPaymentCreateRequest(
		recipientID,
		"NIL",
		*plaid.NewPaymentAmount("GBP", amount),
	)

	response, _, err := plaidCli.PlaidApi.PaymentInitiationPaymentCreate(context.TODO()).PaymentInitiationPaymentCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}

	paymentID := response.GetPaymentId()
	return paymentID, nil
}

func CreatePaymentLinkToken(plaidCli *plaid.APIClient, userId int, paymentId string) (string, error) {
	user := plaid.LinkTokenCreateRequestUser{
		ClientUserId: fmt.Sprintf("%d", userId),
	}

	request := plaid.NewLinkTokenCreateRequest(
	  "Trfer",
	  "en",
	  []plaid.CountryCode{plaid.COUNTRYCODE_GB},
	  user,
	)
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_PAYMENT_INITIATION})

	paymentInitiationReq := plaid.NewLinkTokenCreateRequestPaymentInitiation()
	paymentInitiationReq.PaymentId = &paymentId
	request.SetPaymentInitiation(*paymentInitiationReq)
	request.SetLinkCustomizationName("default")
	request.SetWebhook("https://usual-app.com/transfer/webhook")

	resp, _, err := plaidCli.PlaidApi.LinkTokenCreate(context.TODO()).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}
	linkToken := resp.GetLinkToken()
	
	return linkToken, err
}
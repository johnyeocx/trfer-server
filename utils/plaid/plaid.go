package my_plaid

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/plaid/plaid-go/v11/plaid"
)

func CreateClient() (*plaid.APIClient) {
	clientId := os.Getenv("PLAID_CLIENT_ID")
	// secret := os.Getenv("PLAID_SANDBOX_SECRET")
	secret := os.Getenv("PLAID_DEV_SECRET")

	configuration := plaid.NewConfiguration()
	configuration.AddDefaultHeader("PLAID-CLIENT-ID", clientId)
	configuration.AddDefaultHeader("PLAID-SECRET", secret)
	// configuration.UseEnvironment(plaid.Sandbox)
	configuration.UseEnvironment(plaid.Development)

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
	request.SetProducts([]plaid.Products{plaid.PRODUCTS_AUTH, plaid.PRODUCTS_TRANSACTIONS})
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
	reference string,
) (string, error) {
	request := plaid.NewPaymentInitiationPaymentCreateRequest(
		recipientID,
		reference,
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
	request.SetWebhook("https://usual-app.com/api/payment/webhook")

	resp, _, err := plaidCli.PlaidApi.LinkTokenCreate(context.TODO()).LinkTokenCreateRequest(*request).Execute()
	if err != nil {
		return "", err
	}
	linkToken := resp.GetLinkToken()
	
	return linkToken, err
}

func GetPayment(
	plaidCli *plaid.APIClient, 
	paymentID string,
) (string, error) {
	// request := plaid.NewPaymentInitiationPaymentGetRequest(paymentID)
	// response, _, err := testClient.PlaidApi.PaymentInitiationPaymentGet(ctx).PaymentInitiationPaymentGetRequest(*request).Execute()
	// paymentID := response.GetPaymentId()

	// response, _, err := plaidCli.PlaidApi.PaymentInitiationPaymentCreate(context.TODO()).PaymentInitiationPaymentCreateRequest(*request).Execute()
	// if err != nil {
	// 	return "", err
	// }

	// paymentID := response.GetPaymentId()
	return paymentID, nil
}

func GetUserTransactions(
	plaidCli *plaid.APIClient, 
	accessToken string,
	createdDate time.Time,
) ([]plaid.Transaction, error) {
	const iso8601TimeFormat = "2006-01-02"
	startDate := createdDate.Format(iso8601TimeFormat)
	endDate := time.Now().Add(time.Minute * 60).Format(iso8601TimeFormat)

	request := plaid.NewTransactionsGetRequest(
  		accessToken,
  		startDate,
  		endDate,
	)

	options := plaid.TransactionsGetRequestOptions{
  		Count:  plaid.PtrInt32(100),
  		Offset: plaid.PtrInt32(0),
	}

	request.SetOptions(options)
	transactionsResp, _, err := plaidCli.PlaidApi.TransactionsGet(context.TODO()).TransactionsGetRequest(*request).Execute()
	if err != nil {
		return nil, err
	}

	return transactionsResp.GetTransactions(), nil
}
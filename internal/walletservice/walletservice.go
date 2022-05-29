package walletservice

import (
	"encoding/json"
	"time"
	"errors"
	"github.com/earn-alliance/wallet-commander-cli/internal/log"
	"strconv"
	"os"
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/guregu/dynamo"
	"github.com/earn-alliance/wallet-commander-cli/pkg/constants"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

type WalletService interface {
	SignMessage(roninWallet string, message string) (string, error)
	SendTransaction(roninWallet common.Address, transaction *types.Transaction) (*types.Transaction, error)
}

type WalletServiceVault struct {
	CreationGroup string
	LambdaARN string
	ClientType string
	ClientID string
	AWSSession *session.Session
	LambdaClient *lambda.Lambda
	DynamoTable dynamo.Table
}

func New() (WalletService, error) {
	creationGroup := getEnv("WALLET_SERVICE_CREATION_GROUP", "")
	lambdaARN := getEnv("WALLET_SERVICE_LAMBDA_ARN", "")
	awsRegion := getEnv("AWS_REGION", "")
	clientType := getEnv("WALLET_SERVICE_LAMBDA_TYPE", "")
	clientID := getEnv("WALLET_SERVICE_LAMBDA_TYPE_ID", "")
	dynamoAcountTable := getEnv("DYNAMODB_ACCOUNT", "")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	
	client := lambda.New(sess, &aws.Config{Region: &awsRegion})
	db := dynamo.New(sess, &aws.Config{Region: &awsRegion})
	table := db.Table(dynamoAcountTable)

	if lambdaARN == "" || clientType == "" || clientID == "" || awsRegion == "" {
		log.Logger().Errorln("[WalletService] Fail to load lambda parameters")
		return nil, errors.New("[WalletService] Fail to load lambda parameters")
	}

	log.Logger().Infof("[WalletService] Successfully loaded lambda parameters")

	return &WalletServiceVault{
		CreationGroup: creationGroup,
		LambdaARN: lambdaARN,
		ClientType: clientType,
		ClientID: clientID,
		AWSSession: sess,
		LambdaClient: client,
		DynamoTable: table,
	}, nil
}

func (w WalletServiceVault) SignMessage(roninWallet string, message string) (string, error) {
	// get the same item
	var result BlackpoolScholarAxieAccount
	err := w.DynamoTable.Get("address", roninWallet).One(&result)

	request := WalletServiceSignMessageRequest{
		Action: "wallet:sign_msg", 
		FromType: w.ClientType,
		FromID: w.ClientID,
		TargetID: result.ID,
	}

	request.Payload.Message = message

	payload, err := json.Marshal(request)
	if err != nil {
		log.Logger().Errorln("Error marshalling WalletServiceSignMessageRequest request")
		os.Exit(0)
	}

	walletServiceResponse, err := w.LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: &w.LambdaARN, Payload: payload})
	if err != nil {
		log.Logger().Errorln("Error calling WalletServiceSignMessageRequest")
		os.Exit(0)
	}

	var responsePayload WalletServiceSignMessageResponse
	err = json.Unmarshal(walletServiceResponse.Payload, &responsePayload)
	if err != nil {
		log.Logger().Errorln("Error unmarshalling WalletServiceSignMessageRequest response")
		os.Exit(0)
	}

	// If the status code is NOT 200, the call failed
	if responsePayload.StatusCode != 200 {
		log.Logger().Errorln("Error getting items, StatusCode: " + strconv.Itoa(responsePayload.StatusCode))
		os.Exit(0)
	}

	return responsePayload.Payload.Signature, nil
}

func (w WalletServiceVault) SendTransaction(roninWallet common.Address, transaction *types.Transaction) (*types.Transaction, error) {

	var result BlackpoolScholarAxieAccount
	err := w.DynamoTable.Get("address", roninWallet).One(&result)

	request := WalletServiceSignTransactionRequest{
		Action: "wallet:sign_tx", 
		FromType: w.ClientType,
		FromID: w.ClientID,
		TargetID: result.ID,
	}

	request.Payload.Chain = "ronin"
	request.Payload.Transaction = WalletServiceTransaction{
		Value: transaction.Value().String(),
		GasLimit: strconv.Itoa(int(transaction.Gas())),
		GasPrice: transaction.GasPrice().String(),
		Nonce: strconv.Itoa(int(transaction.Nonce())),
		To: transaction.To().String(),
		Data: string(transaction.Data()),
	}

	payload, err := json.Marshal(request)
	if err != nil {
		log.Logger().Errorln("Error marshalling WalletServiceSignTransactionRequest request")
		os.Exit(0)
	}

	walletServiceResponse, err := w.LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: &w.LambdaARN, Payload: payload})
	if err != nil {
		log.Logger().Errorln("Error calling WalletServiceSignTransactionRequest")
		os.Exit(0)
	}

	var responsePayload WalletServiceSignTransactionResponse
	err = json.Unmarshal(walletServiceResponse.Payload, &responsePayload)
	if err != nil {
		log.Logger().Errorln("Error unmarshalling WalletServiceSignTransactionRequest response")
		os.Exit(0)
	}

	// If the status code is NOT 200, the call failed
	if responsePayload.StatusCode != 200 {
		log.Logger().Errorln("Error getting items, StatusCode: " + strconv.Itoa(responsePayload.StatusCode))
		os.Exit(0)
	}

	txExecuted := false
	getTxRequest := WalletServiceGetTransactionRequest{
		Action: "wallet:get_tx", 
		FromType: w.ClientType,
		FromID: w.ClientID,
		TargetID: result.ID,
	}
	getTxRequest.Payload.TxID = responsePayload.Payload.Transaction.ID
	getTxPayload, err := json.Marshal(getTxRequest)
	if err != nil {
		log.Logger().Errorln("Error marshalling WalletServiceSignTransactionRequest request")
		os.Exit(0)
	}

	var getTxResponsePayload WalletServiceSignTransactionResponse

	for txExecuted == false {
		time.Sleep(15 * time.Second)

		getTxResponse, err := w.LambdaClient.Invoke(&lambda.InvokeInput{FunctionName: &w.LambdaARN, Payload: getTxPayload})
		if err != nil {
			log.Logger().Errorln("Error calling WalletServiceGetTransactionRequest")
			os.Exit(0)
		}

		err = json.Unmarshal(getTxResponse.Payload, &getTxResponsePayload)
		if err != nil {
			log.Logger().Errorln("Error unmarshalling WalletServiceSignTransactionRequest response")
			os.Exit(0)
		}

		if getTxResponsePayload.Payload.Transaction.State != "TO_SEND" {
			if getTxResponsePayload.Payload.Transaction.State != "SUCCESS" {
				return nil, errors.New("Unsuccessful tx execution")
			}
		}	
	}

	rpcClient, err := rpc.DialHTTP(constants.RONIN_PROVIDER_RPC_URI)
	rpcClient.SetHeader("context-type", "application/json")
	rpcClient.SetHeader("user-agent", "wallet-commander")	
	ethClient := ethclient.NewClient(rpcClient)	

	receipt, _, _ := ethClient.TransactionByHash(context.Background(), common.HexToHash(getTxResponsePayload.Payload.Transaction.Hash))

	return receipt, nil
}

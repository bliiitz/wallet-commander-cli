package walletservice

type WalletServiceGetTransactionRequest struct {
	Action              string              `json:"action"`    // Change name in the database
	FromType            string              `json:"from_type"`    // Change name in the database
	FromID              string              `json:"from_id"`    // Change name in the database
	TargetID            string              `json:"target_id"`    // Change name in the database
	Payload             struct {
		TxID string `json:"transaction_id"`
	} `json:"payload"`
}

type WalletServiceGetTransactionResponse struct {
	StatusCode              int              `json:"status_code"`    // Change name in the database
	Payload                 struct {
		Transaction         struct {
			ID string `json:"ID"`
			SendedFromType string `json:"sended_from_type"`
			SendedFromId string `json:"sended_from_id"`
			SendedAt string `json:"sended_at"`
			WalletID string `json:"wallet_id"`
			Chain string `json:"chain"`
			To string `json:"to"`
			Value string `json:"value"`
			Data string `json:"data"`
			Nonce string `json:"nonce"`
			GasPrice string `json:"gas_price"`
			GasLimit string `json:"gas_limit"`
			State string `json:"state"`
			Hash string `json:"hash"`
			Raw_transaction string `json:"raw_transaction"`
			Recipit string `json:"recipit"`
		} `json:"transaction"`
	} `json:"payload"`
}

type WalletServiceSignMessageRequest struct {
	Action              string              `json:"action"`    // Change name in the database
	FromType            string              `json:"from_type"`    // Change name in the database
	FromID              string              `json:"from_id"`    // Change name in the database
	TargetID            string              `json:"target_id"`    // Change name in the database
	Payload             struct {
		Message string `json:"message"`
	} `json:"payload"`
}

type WalletServiceSignMessageResponse struct {
	StatusCode              int              `json:"status_code"`    // Change name in the database
	Payload                 struct {
		Signature           string              `json:"signature"`
	} `json:"payload"`
}


type WalletServiceSignTransactionRequest struct {
	Action              string                   `json:"action"`    // Change name in the database
	FromType            string                   `json:"from_type"`    // Change name in the database
	FromID              string                   `json:"from_id"`    // Change name in the database
	TargetID            string                   `json:"target_id"`    // Change name in the database
	Payload             struct {
		Chain           string                   `json:"chain"`
		Transaction     WalletServiceTransaction `json:"transaction"`
	} `json:"payload"`
}

type WalletServiceTransaction struct {
	Value       string              `json:"value"`
	GasLimit    string              `json:"gas_limit"`
	GasPrice    string              `json:"gas_price"`
	Nonce       string              `json:"nonce"`
	To          string              `json:"to"`
	Data        string              `json:"data"`
}

type WalletServiceSignTransactionResponse struct {
	StatusCode              int              `json:"status_code"`    // Change name in the database
	Payload                 struct {
		Transaction         struct {
			ID string `json:"ID"`
			SendedFromType string `json:"sended_from_type"`
			SendedFromId string `json:"sended_from_id"`
			SendedAt string `json:"sended_at"`
			WalletID string `json:"wallet_id"`
			Chain string `json:"chain"`
			To string `json:"to"`
			Value string `json:"value"`
			Data string `json:"data"`
			Nonce string `json:"nonce"`
			GasPrice string `json:"gas_price"`
			GasLimit string `json:"gas_limit"`
			State string `json:"state"`
			Hash string `json:"hash"`
			Raw_transaction string `json:"raw_transaction"`
			Recipit string `json:"recipit"`
		} `json:"transaction"`
	} `json:"payload"`
}

type BlackpoolScholarAxieAccount struct {
	ID                  string              `dynamo:"ID"`    // Change name in the database
	Address             string              `dynamo:"address"`    // Change name in the database
	Email               string              `dynamo:"email"`    // Change name in the database
	Password            string              `dynamo:"password"`    // Change name in the database
	ScholarID           string              `dynamo:"scholar_id"`    // Change name in the database
	ManagerID           string              `dynamo:"manager_id"`    // Change name in the database
	NotDistributeBefore int                 `dynamo:"not_distribute_before"`    // Change name in the database
	NotClaimableBefore  int                 `dynamo:"not_claimable_before"`    // Change name in the database
	IsActivated         bool                `dynamo:"is_activated"`    // Change name in the database
	Sms                 bool                `dynamo:"sms"`    // Change name in the database
	CreatedAt           int                 `dynamo:"created_at"`    // Change name in the database
}

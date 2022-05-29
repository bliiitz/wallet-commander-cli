package controller

import (
	"context"
	"github.com/earn-alliance/wallet-commander-cli/internal/log"
	"github.com/earn-alliance/wallet-commander-cli/internal/query"
	"github.com/earn-alliance/wallet-commander-cli/internal/store"
	"github.com/earn-alliance/wallet-commander-cli/internal/vault"
	"github.com/earn-alliance/wallet-commander-cli/pkg/api"
	"github.com/earn-alliance/wallet-commander-cli/pkg/client"
	"github.com/earn-alliance/wallet-commander-cli/pkg/utils"
)

type Controller interface {
	ProcessWalletCommand(command query.WalletCommanderCommand)
}

type WalletCommanderController struct {
	client client.Client
	vault  vault.Vault
	store  store.Store
}

func New(vault vault.Vault, store store.Store, client client.Client) (Controller, error) {
	return &WalletCommanderController{
		client: client,
		vault:  vault,
		store:  store,
	}, nil
}

func (w *WalletCommanderController) ProcessWalletCommand(command query.WalletCommanderCommand) {
	commandIdStr := command.Id.(string)
	switch string(command.Operation) {
	case string(api.OperationClaimSLP):
		if payload, err := command.UnmarshalClaimSlpPayload(); err != nil {
			w.store.UpdateWalletCommandTransactionError(commandIdStr, api.CommandStatusInvalidPayload, err.Error())
		} else {
			log.Logger().Infof("Processing claim slp command id %s for wallet %s", command.Id, payload.AddressToClaim)
			w.claimSlp(commandIdStr, payload)
		}
	case string(api.OperationTransferSLP):
		if payload, err := command.UnmarshalTransferSlpPayload(); err != nil {
			w.store.UpdateWalletCommandTransactionError(commandIdStr, api.CommandStatusInvalidPayload, err.Error())
		} else {
			log.Logger().Infof("Processing transfer slp command id %s for"+
				"transfering %d slp from %s to %s", command.Id, payload.Amount, payload.From, payload.To)
			w.transferSlp(commandIdStr, payload)
		}
	case string(api.OperationTransferAxie):
		if payload, err := command.UnmarshalTransferAxiePayload(); err != nil {
			w.store.UpdateWalletCommandTransactionError(commandIdStr, api.CommandStatusInvalidPayload, err.Error())
		} else if err := payload.Validate(); err != nil {
			w.store.UpdateWalletCommandTransactionError(commandIdStr, api.CommandStatusInvalidPayload, err.Error())
		} else {
			log.Logger().Infof("Processing transfer axie command id %s for"+
				"transfering axie id %d from %s to %s", command.Id, payload.AxieId, payload.From, payload.To)
			w.transferAxie(commandIdStr, payload)
		}
	}
}

func (w *WalletCommanderController) processTransactionResult(commandId string, transaction string, err error) {
	if err == nil {
		log.Logger().Debugf("wrting successful result for command %s", commandId)
		w.store.UpdateWalletCommandTransactionSuccess(commandId, transaction)
	} else {
		log.Logger().Debugf("wrting an error result for command %s with err %v", commandId, err)
		w.store.UpdateWalletCommandTransactionError(commandId, api.CommandStatusError, err.Error())
	}
}

func (w *WalletCommanderController) claimSlp(commandId string, payload *api.ClaimSlpPayload) {
	tx, err := w.client.ClaimSlp(context.Background(), utils.RoninAddrToEthAddr(payload.AddressToClaim))

	if err == nil {
		logWalletCommanderSuccess(api.OperationClaimSLP, commandId, tx)
	} else {
		logWalletCommanderError(api.OperationClaimSLP, commandId, err)
	}

	w.processTransactionResult(commandId, tx, err)
	
}

func (w *WalletCommanderController) transferSlp(commandId string, payload *api.TransferSlpPayload) {
	tx, err := w.client.TransferSlp(context.Background(), utils.RoninAddrToEthAddr(payload.From), utils.RoninAddrToEthAddr(payload.To), payload.Amount)

	if err == nil {
		logWalletCommanderSuccess(api.OperationTransferSLP, commandId, tx)
	} else {
		logWalletCommanderError(api.OperationTransferSLP, commandId, err)
	}

	w.processTransactionResult(commandId, tx, err)
	
}

func (w *WalletCommanderController) transferAxie(commandId string, payload *api.TransferAxiePayload) {
	tx, err := w.client.TransferAxie(context.Background(), utils.RoninAddrToEthAddr(payload.From), utils.RoninAddrToEthAddr(payload.To), payload.AxieId)

	if err == nil {
		logWalletCommanderSuccess(api.OperationTransferAxie, commandId, tx)
	} else {
		logWalletCommanderError(api.OperationTransferAxie, commandId, err)
	}

	w.processTransactionResult(commandId, tx, err)
}

func logWalletCommanderSuccess(operation api.WalletCommanderOperation, commandId, tx string) {
	log.Logger().Infof("Successfully signed transaction to %s for command id %s. Transaction id is: %s",
		operation,
		commandId,
		tx,
	)
}

func logWalletCommanderError(operation api.WalletCommanderOperation, commandId string, err error) {
	log.Logger().Warnf("An blockchain error occurred trying to %v with command id %s with err %v"+
		"The system will attempt to retry shortly with another command request if needed", operation, commandId, err)
}

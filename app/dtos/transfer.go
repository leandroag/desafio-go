package dtos

import "github.com/leandroag/desafio/app/domain/entities"

type TransferDTO struct {
	AccountOriginID      int32   `json:"account_origin_id"`
	AccountDestinationID int32   `json:"account_destination_id"`
	Amount               float64 `json:"amount"`
}

func (dto *TransferDTO) ToTransferDomain() *entities.Transfer {
	return &entities.Transfer{
		AccountOriginID:      dto.AccountOriginID,
		AccountDestinationID: dto.AccountDestinationID,
		Amount:               dto.Amount,
	}
}

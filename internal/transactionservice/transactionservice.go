package transactionservice

import (
	"context"

	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/config"
	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/persist"
	lproto "github.com/T-Prohmpossadhorn/GolangMiniProject/pkg/Proto"
)

type TransactionService struct {
	lproto.UnimplementedFruitListServiceServer

	persist persist.Persist
	options config.Options
}

func New(persist persist.Persist, options config.Options) (*TransactionService, error) {
	s := TransactionService{
		persist: persist,
		options: options,
	}
	return &s, nil
}

func (s *TransactionService) GetFullList(ctx context.Context, in *lproto.GetFullListRequest) (*lproto.FullList, error) {
	list := s.persist.Getfulllist()

	var ret lproto.FullList
	for _, v := range list {
		ret.Fruit = append(ret.Fruit, &lproto.Fruit{
			Fruit: v,
		})
	}

	return &ret, nil
}

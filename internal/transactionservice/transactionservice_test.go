package transactionservice

import (
	"context"
	"testing"

	"github.com/T-Prohmpossadhorn/GolangMiniProject/internal/config"
	mocks "github.com/T-Prohmpossadhorn/GolangMiniProject/internal/persist/mocks"
	lproto "github.com/T-Prohmpossadhorn/GolangMiniProject/pkg/Proto"
	"github.com/stretchr/testify/assert"
)

var p = &mocks.Persist{}

func TestGetFullList(t *testing.T) {
	tests := []struct {
		name         string
		fruit        []string
		expectoutput *lproto.FullList
	}{
		{
			name:  "Test1",
			fruit: []string{"apple", "orange"},
			expectoutput: &lproto.FullList{
				Fruit: []*lproto.Fruit{
					{
						Fruit: "apple",
					},
					{
						Fruit: "orange",
					},
				},
			},
		},
		{
			name:  "Test2",
			fruit: []string{"blueberry", "pineapple"},
			expectoutput: &lproto.FullList{
				Fruit: []*lproto.Fruit{
					{
						Fruit: "blueberry",
					},
					{
						Fruit: "pineapple",
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ctx context.Context

			p.Mock.ExpectedCalls = nil
			p.On("Getfulllist").Return(test.fruit)

			options := config.Options{}

			transactionsrv, _ := New(p, options)

			res, err := transactionsrv.GetFullList(ctx, &lproto.GetFullListRequest{})

			assert.Equal(t, res, test.expectoutput)
			assert.Nil(t, err)
		})
	}
}

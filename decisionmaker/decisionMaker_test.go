package decisionmaker

import (
	"reflect"
	"testing"

	"github.com/fabiodmferreira/crypto-trading/assets"
)

type Asset = assets.Asset

type SellArgs struct {
	asset     *Asset
	sellPrice float32
}

type BuyArgs struct {
	amount   float32
	buyPrice float32
}

type TraderSpy struct {
	Buys  []BuyArgs
	Sells []SellArgs
}

func NewTraderSpy() *TraderSpy {
	return &TraderSpy{}
}

func (t *TraderSpy) Buy(amount, price float32) {
	t.Buys = append(t.Buys, BuyArgs{amount, price})
}

func (t *TraderSpy) Sell(asset *Asset, price float32) {
	t.Sells = append(t.Sells, SellArgs{asset, price})
}

func (t *TraderSpy) Reset() {
	t.Buys = []BuyArgs{}
	t.Sells = []SellArgs{}
}

func TestDecisionMaker(t *testing.T) {
	t.Run("DecideToSell should sell if order returns the gains pretended", func(t *testing.T) {
		trader := NewTraderSpy()
		decisionMaker := &DecisionMaker{trader}

		tests := []struct {
			askPrice        float32
			orders          []*Asset
			pretendedProfit float32
			sells           []SellArgs
		}{
			{
				12,
				[]*Asset{
					&Asset{Amount: 100, BuyPrice: 10},
				},
				0.1,
				[]SellArgs{
					{&Asset{Amount: 100, BuyPrice: 10}, 12},
				},
			},
			{
				10.5,
				[]*Asset{
					&Asset{Amount: 100, BuyPrice: 10},
				},
				0.1,
				[]SellArgs{},
			},
		}

		for index, tt := range tests {
			decisionMaker.DecideToSell(tt.askPrice, tt.orders, tt.pretendedProfit)

			got := trader.Sells
			want := tt.sells

			if !reflect.DeepEqual(got, want) {
				t.Errorf("#%v: got %v want %v", index+1, got, want)
			}

			trader.Reset()
		}
	})

	t.Run("DecideToBuy should buy if there is no asset or if ask price dropped", func(t *testing.T) {
		trader := NewTraderSpy()
		decisionMaker := &DecisionMaker{trader}

		tests := []struct {
			ask                  float32
			minimumAssetBuyPrice float32
			dropToBuy            float32
			buyAmount            float32
			buys                 []BuyArgs
		}{
			{
				8.9,
				10,
				0.1,
				100,
				[]BuyArgs{
					{100, 8.9},
				},
			},
			{
				10,
				0,
				0.1,
				20,
				[]BuyArgs{
					{20, 10},
				},
			},
			{
				9,
				10,
				0.1,
				20,
				[]BuyArgs{},
			},
		}

		for index, tt := range tests {
			decisionMaker.DecideToBuy(tt.ask, tt.minimumAssetBuyPrice, tt.dropToBuy, tt.buyAmount)
			got := trader.Buys
			want := tt.buys

			if !reflect.DeepEqual(got, want) {
				t.Errorf("#%v: got %v want %v", index+1, got, want)
			}

			trader.Reset()
		}
	})
}

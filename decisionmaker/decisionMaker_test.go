package decisionmaker

import (
	"reflect"
	"testing"
	"time"

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
	buyTime  time.Time
}

type TraderSpy struct {
	Buys  []BuyArgs
	Sells []SellArgs
}

func NewTraderSpy() *TraderSpy {
	return &TraderSpy{}
}

func (t *TraderSpy) Buy(amount, price float32, buyTime time.Time) {
	t.Buys = append(t.Buys, BuyArgs{amount, price, buyTime})
}

func (t *TraderSpy) Sell(asset *Asset, price float32) {
	t.Sells = append(t.Sells, SellArgs{asset, price})
}

func (t *TraderSpy) Reset() {
	t.Buys = []BuyArgs{}
	t.Sells = []SellArgs{}
}

type AccountSpy struct {
	Withdraws []float32
	Deposits  []float32
}

func (a *AccountSpy) Withdraw(amount float32) error {
	a.Withdraws = append(a.Withdraws, amount)
	return nil
}

func (a *AccountSpy) Deposit(amount float32) error {
	a.Deposits = append(a.Deposits, amount)
	return nil
}

func (a *AccountSpy) Reset() {
	a.Deposits = []float32{}
	a.Withdraws = []float32{}
}

func NewAccountSpy() *AccountSpy {
	return &AccountSpy{}
}

type AssetRepositorySpy struct {
	FindAllCalls               int
	FindCheaperAssetPriceCalls int
}

func (ar *AssetRepositorySpy) FindAll() (*[]assets.Asset, error) {
	ar.FindAllCalls++
	return &[]Asset{}, nil
}

func (ar *AssetRepositorySpy) FindCheaperAssetPrice() (float32, error) {
	ar.FindCheaperAssetPriceCalls++
	return 0, nil
}

func NewAssertRepositorySpy() *AssetRepositorySpy {
	return &AssetRepositorySpy{}
}

func TestDecisionMaker(t *testing.T) {
	decisionMakerOptionsStub := DecisionMakerOptions{0.1, 0.1, 0.1}

	t.Run("DecideToSell should sell if order returns the gains pretended", func(t *testing.T) {
		trader := NewTraderSpy()
		account := NewAccountSpy()
		assetRepository := NewAssertRepositorySpy()
		decisionMaker := &DecisionMaker{trader, account, assetRepository, decisionMakerOptionsStub}

		tests := []struct {
			askPrice        float32
			orders          *[]Asset
			pretendedProfit float32
			sells           []SellArgs
		}{
			{
				12,
				&[]Asset{
					Asset{Amount: 100, BuyPrice: 10},
				},
				0.1,
				[]SellArgs{
					{&Asset{Amount: 100, BuyPrice: 10}, 12},
				},
			},
			{
				10.5,
				&[]Asset{
					Asset{Amount: 100, BuyPrice: 10},
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
		account := NewAccountSpy()
		assetRepository := NewAssertRepositorySpy()
		decisionMaker := &DecisionMaker{trader, account, assetRepository, decisionMakerOptionsStub}

		now := time.Now()

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
					{100, 8.9, now},
				},
			},
			{
				10,
				0,
				0.1,
				20,
				[]BuyArgs{
					{20, 10, now},
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
			decisionMaker.DecideToBuy(tt.ask, tt.minimumAssetBuyPrice, tt.dropToBuy, tt.buyAmount, now)
			got := trader.Buys
			want := tt.buys

			if !reflect.DeepEqual(got, want) {
				t.Errorf("#%v: got %v want %v", index+1, got, want)
			}

			trader.Reset()
		}
	})
}

package assets

import "github.com/fabiodmferreira/crypto-trading/domain"

type BenchmarkAssetsInfo struct {
	Buys         [][]float32 `json:"buys"`
	Sells        [][]float32 `json:"sells"`
	SellsPending int         `json:"sellsPending"`
}

// GroupAssetsByState returns assets bought and sold
func GroupAssetsByState(assets *[]domain.Asset) BenchmarkAssetsInfo {
	var sells int

	Buys := [][]float32{}
	Sells := [][]float32{}

	for _, asset := range *assets {
		Buys = append(Buys, []float32{float32(asset.BuyTime.Unix()) * 1000, asset.BuyPrice})

		if asset.Sold {
			Sells = append(Sells, []float32{float32(asset.SellTime.Unix()) * 1000, asset.SellPrice})
			sells++
		}
	}

	return BenchmarkAssetsInfo{
		Buys:         Buys,
		Sells:        Sells,
		SellsPending: len(*assets) - sells,
	}
}

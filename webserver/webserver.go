package webserver

import (
	"fmt"
	"net/http"

	"github.com/fabiodmferreira/crypto-trading/benchmark"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/gorilla/mux"
)

// CryptoTradingServer provides server that handles application API
type CryptoTradingServer struct {
	http.Handler
}

// NewCryptoTradingServer returns an instance of CryptoTradingServer
func NewCryptoTradingServer(benchmark *benchmark.Service, assetsPrice domain.AssetPriceRepository) (*CryptoTradingServer, error) {
	server := new(CryptoTradingServer)

	router := mux.NewRouter()

	benchmarkController := NewBenchmarkController(benchmark)
	router.Handle("/api/benchmark", http.HandlerFunc(benchmarkController.BenchmarkHandler))
	router.Handle("/api/benchmark/data-sources", http.HandlerFunc(benchmarkController.GetBenchmarkDataSourcesHandler))
	router.HandleFunc("/api/benchmark/{id}", benchmarkController.ResourceHandler)

	assetsPricesController := NewAssetsPricesController(assetsPrice)
	router.Handle("/api/assets/{asset}/prices", http.HandlerFunc(assetsPricesController.GetAssetPrices))

	router.Handle("/", http.HandlerFunc(server.versionHandler))

	server.Handler = router

	return server, nil
}

// versionHandler responds with application version
func (c *CryptoTradingServer) versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "v0.0.0")
}

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
func NewCryptoTradingServer(
	benchmark *benchmark.Service,
	assetsPrice domain.AssetPriceRepository,
	accounts domain.AccountsRepository,
	assets domain.AssetsRepository,
	appService domain.ApplicationService,
) (*CryptoTradingServer, error) {
	server := new(CryptoTradingServer)

	router := mux.NewRouter()

	benchmarkController := NewBenchmarkController(benchmark)
	router.Handle("/api/benchmark", http.HandlerFunc(benchmarkController.BenchmarkHandler))
	router.Handle("/api/benchmark/data-sources", http.HandlerFunc(benchmarkController.GetBenchmarkDataSourcesHandler))
	router.HandleFunc("/api/benchmark/{id}", benchmarkController.ResourceHandler)
	router.HandleFunc("/api/benchmark/{id}/state", benchmarkController.GetBenchmarkExecutionStateHandler)

	assetsPricesController := NewAssetsPricesController(assetsPrice)
	router.Handle("/api/assets/{asset}/prices", http.HandlerFunc(assetsPricesController.GetAssetPrices))

	accountsController := NewAccountsController(accounts, assets)
	router.HandleFunc("/api/accounts/{id}", accountsController.GetAccountHandler)
	router.HandleFunc("/api/accounts/{id}/assets", accountsController.GetAccountAssetsHandler)

	applicationsController := NewApplicationsController(appService)
	router.HandleFunc("/api/applications", applicationsController.GetApplicationsHandler)
	router.HandleFunc("/api/applications/{id}/state/last", applicationsController.GetLastApplicationStateHandler)
	router.HandleFunc("/api/applications/{id}/log-events", applicationsController.GetApplicationLogEventsHandler)
	router.HandleFunc("/api/applications/{id}", applicationsController.ApplicationItemHandler)

	router.Handle("/", http.HandlerFunc(server.versionHandler))

	server.Handler = router

	return server, nil
}

// versionHandler responds with application version
func (c *CryptoTradingServer) versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "v0.0.0")
}

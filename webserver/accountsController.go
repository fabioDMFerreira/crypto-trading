package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fabiodmferreira/crypto-trading/assets"
	"github.com/fabiodmferreira/crypto-trading/domain"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AccountsController has the accounts routes handlers
type AccountsController struct {
	repo       domain.AccountsRepository
	assetsRepo domain.AssetsRepository
}

func NewAccountsController(repo domain.AccountsRepository, assetsRepo domain.AssetsRepository) *AccountsController {
	return &AccountsController{repo, assetsRepo}
}

func (a *AccountsController) GetAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	oid, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	account, err := a.repo.FindById(oid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func (a *AccountsController) GetAccountAssetsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	oid, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	assets, err := a.assetsRepo.FindAll(oid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(assets)
}

func (a *AccountsController) GetAccountAssetsGroupedByStateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	oid, err := primitive.ObjectIDFromHex(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	docs, err := a.assetsRepo.FindAll(oid)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	json.NewEncoder(w).Encode(assets.GroupAssetsByState(docs))
}

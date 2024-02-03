package handlers

import (
	"encoding/json"
	models "github.com/h-abranches-dev/daily-expenses-be/domain-layer"
	repositories "github.com/h-abranches-dev/daily-expenses-be/persistence-layer"
	services "github.com/h-abranches-dev/daily-expenses-be/service-layer"
	"net/http"
	"strconv"
	"strings"
)

// TransactionHandlerFunc /transactions
func TransactionHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		entityProvided := r.URL.Query().Get("entity")
		typeProvided := r.URL.Query().Get("type")
		categoriesProvided := r.URL.Query().Get("categories")

		if entityProvided != "" && !entityIsValid(entityProvided) || typeProvided != "" && !kindIsValid(typeProvided) {
			writeResponseWithError(w, http.StatusBadRequest, badRequest)
			return
		}

		if entityProvided == "" && typeProvided == "" && categoriesProvided != "" {
			categoriesProvidedSlc := strings.Split(categoriesProvided, ",")

			repos, err := repositories.GetAllRepos()
			if err != nil {
				writeResponseWithError(w, http.StatusInternalServerError, internalServerError)
				return
			}
			allTransactions, err := services.GetAllTransactions(*repos, false)
			if err != nil {
				writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
				return
			}

			transactionsFilteredByCategories, err := services.GetTransactionsFilteredByCategories(allTransactions, categoriesProvidedSlc)
			if err != nil {
				writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err = json.NewEncoder(w).Encode(transactionsFilteredByCategories); err != nil {
				writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
				return
			}
			if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, ok, transactionsFilteredByCategories); err != nil {
				logDetailedError(err)
				return
			}

			return
		}

		repo, err := repositories.GetTransRepo(typeProvided, entityProvided)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		ts, err := services.GetAllTransactionsByRepo(repo, false)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		tsDTO, err := services.NewTransactionsDTO(*ts)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(tsDTO.TransactionsDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, ok, tsDTO.TransactionsDTO); err != nil {
			logDetailedError(err)
			return
		}

	case http.MethodPost:
		entityProvided := r.URL.Query().Get("entity")
		typeProvided := r.URL.Query().Get("type")
		if !entityIsValid(entityProvided) || !kindIsValid(typeProvided) {
			writeResponseWithError(w, http.StatusBadRequest, badRequest)
			return
		}

		ntDTO := services.TransactionDTO{}

		if err := json.NewDecoder(r.Body).Decode(&ntDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		nt, err := ntDTO.NewTransaction()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		repo, err := repositories.GetTransRepo(typeProvided, entityProvided)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		newID, err := services.AddTransaction(repo, nt)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		ntDTO.ID = strconv.Itoa(newID)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(ntDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, created, ntDTO); err != nil {
			logDetailedError(err)
			return
		}

	default:
		writeResponseWithError(w, http.StatusMethodNotAllowed, methodNotAllowed)
	}
}

// TransactionsTypesHandlerFunc /transactions/types
func TransactionsTypesHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		kinds := []string{models.DebitKindTransaction, models.CreditKindTransaction}
		if err := json.NewEncoder(w).Encode(kinds); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err := logResponse(r.Method, r.URL.Path, r.URL.RawQuery, ok, kinds); err != nil {
			logDetailedError(err)
			return
		}

	default:
		writeResponseWithError(w, http.StatusMethodNotAllowed, methodNotAllowed)
	}
}

// UpdateTransactionHandlerFunc /transactions/:transaction_id
func UpdateTransactionHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPut:
		tIDStr := strings.Split(r.URL.Path, "/transactions/")[1]
		tID, err := strconv.Atoi(tIDStr)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		entityProvided := r.URL.Query().Get("entity")
		typeProvided := r.URL.Query().Get("type")
		if !entityIsValid(entityProvided) || !kindIsValid(typeProvided) {
			writeResponseWithError(w, http.StatusBadRequest, badRequest)
			return
		}

		tDTO := services.TransactionDTO{
			ID: tIDStr,
		}
		if err = json.NewDecoder(r.Body).Decode(&tDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		t, err := tDTO.NewTransaction()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		t.ID = tID

		repo, err := repositories.GetTransRepo(typeProvided, entityProvided)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		err = services.UpdateTransaction(repo, t)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(tDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, ok, tDTO); err != nil {
			logDetailedError(err)
			return
		}
		return

	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PUT")
		w.WriteHeader(http.StatusNoContent)
		if err := logResponse(r.Method, r.URL.Path, r.URL.RawQuery, noContent, ""); err != nil {
			logDetailedError(err)
			return
		}
		return

	default:
		writeResponseWithError(w, http.StatusMethodNotAllowed, methodNotAllowed)
		return
	}
}

// TransactionsCategoriesHandlerFunc /transactions/categories
func TransactionsCategoriesHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		repos, err := repositories.GetAllRepos()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		allTransactions, err := services.GetAllTransactions(*repos, false)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		m := services.GetCategoriesLabelsInTransactionsMap(allTransactions)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(m); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, ok, m); err != nil {
			logDetailedError(err)
			return
		}

	default:
		writeResponseWithError(w, http.StatusMethodNotAllowed, methodNotAllowed)
	}
}

func entityIsValid(entityProvided string) bool {
	for _, entity := range models.TransactionEntities {
		if string(entity) == entityProvided {
			return true
		}
	}
	return false
}

func kindIsValid(kindProvided string) bool {
	for _, kind := range models.TransactionKinds {
		if string(kind) == kindProvided {
			return true
		}
	}
	return false
}

package handlers

import (
	"encoding/json"
	models "github.com/h-abranches-dev/daily-expenses-be/domain-layer"
	repositories "github.com/h-abranches-dev/daily-expenses-be/persistence-layer"
	services "github.com/h-abranches-dev/daily-expenses-be/service-layer"
	"net/http"
	"strconv"
)

// TransactionsEntitiesHandlerFunc /entities
func TransactionsEntitiesHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		repo, err := repositories.GetTransEntRepo()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		tses, err := services.GetAllTransactionsEntities(repo, models.RefToTransactionsEntities)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		tsesDTO := services.NewTransactionsEntitiesDTO(*tses)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(tsesDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, ok, tsesDTO); err != nil {
			logDetailedError(err)
			return
		}

	case http.MethodPost:
		ntseDTO := services.TransactionsEntityDTO{}
		if err := json.NewDecoder(r.Body).Decode(&ntseDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		ntse, err := ntseDTO.NewTransactionsEntity()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		repo, err := repositories.GetTransEntRepo()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		newID, err := services.AddTransactionsEntity(repo, ntse)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		ntseDTO.ID = strconv.Itoa(newID)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(ntseDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, created, ntse); err != nil {
			logDetailedError(err)
			return
		}

	default:
		writeResponseWithError(w, http.StatusMethodNotAllowed, methodNotAllowed)
	}
}

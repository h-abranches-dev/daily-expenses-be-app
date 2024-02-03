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

// CategoriesHandlerFunc /categories
func CategoriesHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		repo, err := repositories.GetCategoriesRepo()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		cs, err := services.GetAllCategories(repo, models.RefToCategories)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		csDTO := services.NewCategoriesDTO(*cs)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(csDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, ok, csDTO); err != nil {
			logDetailedError(err)
			return
		}
	case http.MethodPost:
		ncDTO := services.CategoryDTO{}
		if err := json.NewDecoder(r.Body).Decode(&ncDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		nc, err := ncDTO.NewCategory()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		repo, err := repositories.GetCategoriesRepo()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		newID, err := services.AddCategory(repo, nc)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		ncDTO.ID = strconv.Itoa(newID)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		if err = json.NewEncoder(w).Encode(ncDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, created, ncDTO); err != nil {
			logDetailedError(err)
			return
		}

	default:
		writeResponseWithError(w, http.StatusMethodNotAllowed, methodNotAllowed)
	}

}

// UpdateCategoryHandlerFunc /categories/:category_id
func UpdateCategoryHandlerFunc(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPut:
		cIDStr := strings.Split(r.URL.Path, "/categories/")[1]
		cID, err := strconv.Atoi(cIDStr)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		cDTO := services.CategoryDTO{
			ID: cIDStr,
		}
		if err = json.NewDecoder(r.Body).Decode(&cDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		c, err := cDTO.NewCategory()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusBadRequest, badRequest, err)
			return
		}

		c.ID = cID

		repo, err := repositories.GetCategoriesRepo()
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		err = services.UpdateCategory(repo, c)
		if err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err = json.NewEncoder(w).Encode(cDTO); err != nil {
			writeResponseWithDetailedError(w, http.StatusInternalServerError, internalServerError, err)
			return
		}
		if err = logResponse(r.Method, r.URL.Path, r.URL.RawQuery, ok, cDTO); err != nil {
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

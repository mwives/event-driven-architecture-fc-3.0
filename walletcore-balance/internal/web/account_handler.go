package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mwives/event-driven-architecture-fc-3.0/walletcore-balance/internal/usecase/find_account_by_id"
)

type WebAccountHandler struct {
	FindAccountByIDUseCase find_account_by_id.FindAccountByIDUseCase
}

func NewWebAccountHandler(findAccountByIDUseCase find_account_by_id.FindAccountByIDUseCase) WebAccountHandler {
	return WebAccountHandler{FindAccountByIDUseCase: findAccountByIDUseCase}
}

func (h *WebAccountHandler) FindAccountByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	dto := find_account_by_id.FindAccountByIDInputDTO{
		AccountID: id,
	}

	output, err := h.FindAccountByIDUseCase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

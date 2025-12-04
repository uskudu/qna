package question_handler

import (
	"encoding/json"
	"net/http"

	"QnA/internal/shared/utils"
	"QnA/internal/usecase"
)

type QuestionHandler struct {
	uc *usecase.QuestionUsecase
}

func NewQuestionHandler(uc *usecase.QuestionUsecase) *QuestionHandler {
	return &QuestionHandler{uc: uc}
}

func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Text == "" {
		utils.WriteError(w, http.StatusBadRequest, "text is required")
		return
	}

	id, err := h.uc.CreateQuestion(r.Context(), req.Text)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"id": id})
}

func (h *QuestionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	list, err := h.uc.GetAllQuestions(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, list)
}

func (h *QuestionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rawID := r.PathValue("id")
	//id, err := uuid.Parse(rawID)
	//if err != nil {
	//	return c.Status(my_http.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	//}
	q, answers, err := h.uc.GetQuestionByID(r.Context(), rawID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"question": q,
		"answers":  answers,
	})
}

func (h *QuestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rawID := r.PathValue("id")
	//id, err := uuid.Parse(rawID)
	//if err != nil {
	//	writeError(w, http.StatusBadRequest, "invalid uuid")
	//	return
	//}

	err := h.uc.DeleteQuestion(r.Context(), rawID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

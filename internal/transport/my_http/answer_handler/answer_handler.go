package answer_handler

import (
	"QnA/internal/shared/utils"
	"QnA/internal/usecase"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

//POST   /questions/{id}/answers
//GET    /answers/{id}
//DELETE /answers/{id}

type AnswerHandler struct {
	uc *usecase.AnswerUsecase
}

func NewAnswerHandler(uc *usecase.AnswerUsecase) *AnswerHandler {
	return &AnswerHandler{uc: uc}
}

func (h *AnswerHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req struct {
		UserID     string `json:"user_id"`
		QuestionID string `json:"question_id"`
		Text       string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Text == "" {
		utils.WriteError(w, http.StatusBadRequest, "text is required")
		return
	}

	_, err := uuid.Parse(req.QuestionID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid question id")
		return
	}

	id, err := h.uc.CreateAnswer(r.Context(), req.QuestionID, req.UserID, req.Text)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"id": id})
}

func (h *AnswerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rawID := r.PathValue("id")
	id, err := uuid.Parse(rawID)
	if err != nil {
		utils.WriteError(w, http.StatusMethodNotAllowed, "invalid uuid")
	}
	a, err := h.uc.GetAnswerByID(r.Context(), id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"answer": a,
	})
}

func (h *AnswerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	rawID := r.PathValue("id")
	//id, err := uuid.Parse(rawID)
	//if err != nil {
	//	return c.Status(my_http.StatusBadRequest).JSON(fiber.Map{"error": "invalid uuid"})
	//}
	err := h.uc.DeleteAnswer(r.Context(), rawID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

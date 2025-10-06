package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"sobezos/services/core-service/internal/models"
	"sobezos/services/core-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{service: s}
}

// TaskGet godoc
// @Summary Получить задачу
// @Description Получает задачу для пользователя с учетом его тегов
// @Tags task
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /taskget [get]
func (h *Handler) TaskGet(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	result, err := h.service.TaskGet(telegramID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// TaskGetID godoc
// @Summary Получить задачу по ID
// @Description Получает задачу по её ID для пользователя
// @Tags task
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Param args query int true "ID задачи"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /taskgetid [get]
func (h *Handler) TaskGetID(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	args := r.URL.Query().Get("args")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	result, err := h.service.TaskGetID(telegramID, args)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// AnswerGet godoc
// @Summary Получить ответ на последний вопрос пользователя
// @Description Возвращает ответ на последний теоретический вопрос пользователя
// @Tags answer
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /answerget [get]
func (h *Handler) AnswerGet(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	result, err := h.service.AnswerGet(telegramID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// TagSet godoc
// @Summary Установить теги пользователю
// @Description Добавляет новые теги пользователю (требуется telegram_id в заголовке)
// @Tags tag
// @Accept json
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Param data body models.TagSetRequest true "Список тегов через запятую"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /tagset [post]
func (h *Handler) TagSet(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	var req models.TagSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid body"})
		return
	}
	result, err := h.service.TagSet(telegramID, req.Args)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// TagClear godoc
// @Summary Очистить теги пользователя
// @Description Очищает все теги пользователя
// @Tags tag
// @Accept json
// @Produce json
// @Param data body models.TagClearRequest true "Telegram ID пользователя"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid body"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /tagclear [post]
func (h *Handler) TagClear(w http.ResponseWriter, r *http.Request) {
	var req models.TagClearRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid body"})
		return
	}
	result, err := h.service.TagClear(req.TelegramID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// TagGet godoc
// @Summary Получить список тегов
// @Description Получает список всех доступных тегов
// @Tags tag
// @Produce json
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /tagget [get]
func (h *Handler) TagGet(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.TagGet()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// TaskAdd godoc
// @Summary Добавить задачу
// @Description Добавляет новую задачу (требуется telegram_id в заголовке)
// @Tags task
// @Accept json
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Param data body models.TaskAddRequest true "JSON задачи"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /taskadd [post]
func (h *Handler) TaskAdd(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}

	var req models.TaskAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid body"})
		return
	}

	result, err := h.service.TaskAdd(telegramID, req.Question, req.Answer, req.Tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// TaskEdit godoc
// @Summary Редактировать задачу
// @Description Редактирует существующую задачу (требуется telegram_id в заголовке)
// @Tags task
// @Accept json
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Param data body models.TaskEditRequest true "JSON задачи для редактирования"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /taskedit [put]
func (h *Handler) TaskEdit(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	var req models.TaskEditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid body"})
		return
	}
	result, err := h.service.TaskEdit(telegramID, req.ID, req.Question, req.Answer, req.Tags)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// UserAdd godoc
// @Summary Добавить пользователя
// @Description Добавляет нового пользователя (требуется telegram_id администратора в заголовке)
// @Tags user
// @Accept json
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Param data body models.UserAddRequest true "Аргументы для добавления пользователя"
// @Success 200 {object} models.CommonSuccessResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid admin telegram_id"
// @Failure 500 {object} models.CommonErrorResponse "internal error"
// @Router /useradd [post]
func (h *Handler) UserAdd(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid admin telegram_id"})
		return
	}
	var req models.UserAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid body"})
		return
	}
	result, err := h.service.UserAdd(telegramID, req.Args) //
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(models.CommonSuccessResponse{Result: result})
}

// UserCheck godoc
// @Summary Проверить пользователя
// @Description Проверяет существование пользователя и возвращает его роль и username
// @Tags user
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Success 200 {object} models.UserCheckResponse
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 404 {object} models.CommonErrorResponse "user not found"
// @Router /usercheck [get]
func (h *Handler) UserCheck(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	res, exists := h.service.UserCheck(telegramID)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "user not found"})
		return
	}
	json.NewEncoder(w).Encode(models.UserCheckResponse{
		Role:     res.Role,
		Username: res.Username,
	})
}

// StatsGet godoc
// @Summary Получить статистику пользователя
// @Description Получить статистику пользователя
// @Tags user
// @Produce json
// @Param telegram_id query int true "Telegram ID пользователя"
// @Success 200 {object} models.UserState
// @Failure 400 {object} models.CommonErrorResponse "invalid telegram_id"
// @Failure 404 {object} models.CommonErrorResponse "user not found"
// @Router /usercheck [get]
func (h *Handler) StatsGet(w http.ResponseWriter, r *http.Request) {
	telegramIDStr := r.URL.Query().Get("telegram_id")
	telegramID, err := strconv.Atoi(telegramIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.CommonErrorResponse{Error: "invalid telegram_id"})
		return
	}
	model, err := h.service.StatsGet(telegramID)
	if err == nil {
		json.NewEncoder(w).Encode(model)
	}

}

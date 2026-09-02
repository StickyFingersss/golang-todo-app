package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/StickyFingersss/golang-todo-app/internal/core/logger"
	core_http_request "github.com/StickyFingersss/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/StickyFingersss/golang-todo-app/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary		Получение задач
// @Description	Просмотр списка задач с опциональной фильтрацией по автору и пагинацией
// @Tags			tasks
// @Produce		json
// @Param			user_id query int false "Фильтр по ID автора задачи"
// @Param			limit query int false "Размер страницы с задачами"
// @Param			offset query int false "Смещение страницы с задачами"
// @Success		200 {object} GetTasksResponse "Успешное получение списка задач"
// @Failure		400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure		500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router			/tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID/limit/offset query params",
		)
		return
	}

	tasksDomain, err := h.taskService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userId, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}
	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userId, limit, offset, nil
}

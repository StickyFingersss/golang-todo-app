package tasks_transport_http

import (
	"net/http"

	"github.com/StickyFingersss/golang-todo-app/internal/core/domain"
	core_logger "github.com/StickyFingersss/golang-todo-app/internal/core/logger"
	core_http_request "github.com/StickyFingersss/golang-todo-app/internal/core/transport/http/request"
	core_http_response "github.com/StickyFingersss/golang-todo-app/internal/core/transport/http/response"
)

type CreateTaskRequest TaskDTOResponse

type CreateTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) CreateTask(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)
		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)

	taskDomain, err := h.taskService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

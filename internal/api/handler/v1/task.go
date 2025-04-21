package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/vadimbarashkov/workmate-test-task/internal/entity"
	"github.com/vadimbarashkov/workmate-test-task/internal/manager"
)

type TaskHandler struct {
	taskManager manager.TaskManager
}

func RegisterTaskRoutes(r chi.Router, taskManager manager.TaskManager) {
	h := &TaskHandler{taskManager: taskManager}

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", h.CreateTask)
		r.Get("/{task_id}", h.GetTask)
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateTaskRequest struct {
	TaskType string `json:"task_type"`
	Payload  string `json:"payload"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error: "invalid request",
		})
		return
	}

	taskType, err := entity.ParseTaskType(req.TaskType)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error: "missing or invalid task type",
		})
		return
	}

	ctx := context.WithoutCancel(r.Context())
	taskID, err := h.taskManager.Create(ctx, taskType, []byte(req.Payload))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error: "failed to create task",
		})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, CreateTaskResponse{
		ID: taskID.String(),
	})
}

type GetTaskResponse struct {
	TaskType string  `json:"task_type"`
	Status   string  `json:"status"`
	Result   string  `json:"result"`
	Error    *string `json:"error"`
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskIDParam := chi.URLParam(r, "task_id")
	if taskIDParam == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error: "missing task_id path param",
		})
		return
	}

	taskID, err := uuid.Parse(taskIDParam)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, ErrorResponse{
			Error: "invalid task_id path param",
		})
		return
	}

	ctx := context.WithoutCancel(r.Context())
	task, err := h.taskManager.Get(ctx, taskID)
	if err != nil {
		if errors.Is(err, manager.ErrTaskNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, ErrorResponse{
				Error: "task not found",
			})
			return
		}
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, ErrorResponse{
			Error: "failed to get task",
		})
		return
	}

	var taskErr string
	if err := task.Error(); err != nil {
		taskErr = err.Error()
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, GetTaskResponse{
		TaskType: task.TaskType().String(),
		Status:   task.Status().String(),
		Result:   string(task.Result()),
		Error:    &taskErr,
	})
}

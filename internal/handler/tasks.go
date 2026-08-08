package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	middleware "github.com/Davethompson01/rialo_hub_backend/Middleware"
	"github.com/Davethompson01/rialo_hub_backend/config"
	auth "github.com/Davethompson01/rialo_hub_backend/internal/Auth"
	"github.com/Davethompson01/rialo_hub_backend/internal/models"
	"github.com/Davethompson01/rialo_hub_backend/internal/services"
	"github.com/go-chi/chi"
)

func CreateTasks(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		task.UserID = claims.UserID

		CreateTasks, err := services.CreateTasks(api, task)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}

		RespondWithJson(
			w,
			200,
			true,
			"Tasks created Successfully",
			CreateTasks,
		)
	}
}

func ApplyForTasks(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apply models.TaskApplication
		if err := json.NewDecoder(r.Body).Decode(&apply); err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		apply.Employee_id = claims.UserID
		fmt.Println(apply.Employee_id)
		apply.Status = "Ongoing"

		tasksApplication, err := services.ApplyForTasks(api, apply)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}

		RespondWithJson(
			w,
			201,
			true,
			"Application Sent Successfully",
			tasksApplication,
		)
	}
}

func AcceptEmployee(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apply models.TaskApplication
		if err := json.NewDecoder(r.Body).Decode(&apply); err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		apply.Employer_id = claims.UserID
		apply.Status = "Accepted"

		acceptEmployee, err := services.AcceptEmployee(api, apply.ID, apply.Task_id, apply.Employer_id)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}
		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Applicant accepted successfully",
			acceptEmployee,
		)
	}
}

func RejectEmployee(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskStatus models.TaskApplication
		if err := json.NewDecoder(r.Body).Decode(&taskStatus); err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		taskStatus.Employer_id = claims.UserID

		rejectEmployee, err := services.RejectEmployee(api, taskStatus.ID, taskStatus.Task_id, taskStatus.Employer_id)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}
		RespondWithJson(
			w,
			201,
			true,
			"Applicant rejected successfully",
			rejectEmployee,
		)
	}
}

func GetTaskApplications(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
		if err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, "Invalid task ID", nil)
			return
		}
		// var taskStatus models.ApplicationResponse
		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)

		taskApplication, err := services.GetTaskApplications(api, taskID, claims.UserID)
		if err != nil {
			RespondWithJson(w, 400, false, err.Error(), nil)
			return
		}
		RespondWithJson(
			w,
			200,
			true,
			"Applicants for Tasks Fetched Successfully",
			taskApplication,
		)
	}
}

func GetMyApplications(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var taskStatus models.ApplicationResponse
		if err := json.NewDecoder(r.Body).Decode(&taskStatus); err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}
		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)

		GetMyApplications, err := services.GetMyApplications(api, claims.UserID)
		if err != nil {
			RespondWithJson(w, http.StatusUnauthorized, false, err.Error(), nil)
			return
		}
		RespondWithJson(
			w,
			201,
			true,
			"Applicants for Tasks Fetched Successfully",
			GetMyApplications,
		)

	}
}

func CancelApplication(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apply models.TaskApplication
		if err := json.NewDecoder(r.Body).Decode(&apply); err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}
		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)

		CancelApplication, err := services.CancelApplication(api, apply.ID, claims.UserID)
		if err != nil {
			RespondWithJson(w, http.StatusUnauthorized, false, err.Error(), nil)
			return
		}
		RespondWithJson(
			w,
			201,
			true,
			"Applicants Cancelled",
			CancelApplication,
		)
	}
}

func DeleteTask(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var apply models.TaskApplication
		if err := json.NewDecoder(r.Body).Decode(&apply); err != nil {
			RespondWithJson(w, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		delete, err := services.DeleteTask(api, apply.Task_id, claims.UserID)
		if err != nil {
			RespondWithJson(w, http.StatusUnauthorized, false, err.Error(), nil)
			return
		}
		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Task Deleted",
			delete,
		)
	}
}

func Taskfeed(api *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// claims := r.Context().Value(middleware.ClaimsKey).(*auth.Claims)
		feeds, err := services.TasksFeeds(api)
		if err != nil {
			RespondWithJson(w, http.StatusUnauthorized, false, err.Error(), nil)
			return
		}
		RespondWithJson(
			w,
			http.StatusOK,
			true,
			"Task Fetched",
			feeds,
		)
	}
}

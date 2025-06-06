package handler

import (
	"RobotService/internal/dto"
	"RobotService/internal/metrics"
	"RobotService/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "RobotService/cmd/docs" // импорт сгенерированных Swagger-доков

	"github.com/go-chi/chi/v5"
)

const (
	CreateRobot     = "/robots/create"
	GetRobotInfo    = "/robots/{id}"
	UpdateRobotCord = "/robots/updatecord"
	UpdateRobotName = "/robots/updatename"
	DeleteRobot     = "/robots/delete/{id}"
)

type RobotHandlers struct {
	RobotService services.RobotService
}

func (h *RobotHandlers) Register(router *chi.Mux) {
	router.Post(CreateRobot, h.RobotCreate)
	router.Get(GetRobotInfo, h.GetRobotInfo)
	router.Put(UpdateRobotCord, h.UpdateRobotCord)
	router.Put(UpdateRobotName, h.UpdateRobotName)
	router.Delete(DeleteRobot, h.DeleteRobot)
}

// @Summary Create new robot
// @Description Create a new robot with name and coordinates
// @Tags robots
// @Accept json
// @Produce json
// @Param robot body dto.CreateRobotDTO true "Robot info"
// @Success 201 {integer} int "Robot ID"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Internal error"
// @Router /robots/create [post]
func (h *RobotHandlers) RobotCreate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := "201"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("POST", "/expenses", status).Observe(duration)
	}()

	var createdto dto.CreateRobotDTO

	err := json.NewDecoder(r.Body).Decode(&createdto)
	if err != nil {
		status = "400"
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := h.RobotService.CreateRobot(createdto)
	if err != nil {
		status = "500"
		http.Error(w, fmt.Sprintf("Internal error: %s", err.Error()), http.StatusInternalServerError)
	}
	metrics.CreatedRobot.Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		status = "500"
		http.Error(w, fmt.Sprintf("Failed to encode JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

// @Summary Get robot info
// @Description Get detailed robot info by ID
// @Tags robots
// @Produce json
// @Param id path int true "Robot ID"
// @Success 200 {object} entities.Robot
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Internal server error"
// @Router /robots/{id} [get]
func (handler *RobotHandlers) GetRobotInfo(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := "200"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("GET", "/expenses/{id}", status).Observe(duration)
	}()

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		status = "400"
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	robotinfo, err := handler.RobotService.GetRobotInfo(id)
	if err != nil {
		status = "500"
		http.Error(w, fmt.Sprintf("Failed to find robot data: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	metrics.GetRobot.Inc()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(robotinfo)
}

// @Summary Update robot coordinates
// @Description Update x/y coordinates of a robot
// @Tags robots
// @Accept json
// @Param robot body dto.UpdateRobotCordDTO true "Updated coordinates"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to update robot cords"
// @Router /robots/updatecord [put]
func (h *RobotHandlers) UpdateRobotCord(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := "200"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("GET", "/expenses", status).Observe(duration)
	}()

	newRobotData := dto.UpdateRobotCordDTO{}
	err := json.NewDecoder(r.Body).Decode(&newRobotData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	err = h.RobotService.UpdateRobotCords(newRobotData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update robot cords: %s", err.Error()), http.StatusInternalServerError)
	}
	metrics.UpdateRobotCords.Inc()
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update robot name
// @Description Update robot name by ID
// @Tags robots
// @Accept json
// @Param robot body dto.UpdateRobotNameDTO true "Updated name"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to update robot name"
// @Router /robots/updatename [put]
func (h *RobotHandlers) UpdateRobotName(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	status := "200"
	defer func() {
		duration := time.Since(start).Seconds()
		metrics.RequestDuration.WithLabelValues("GET", "/expenses", status).Observe(duration)
	}()

	newRobotData := dto.UpdateRobotNameDTO{}
	err := json.NewDecoder(r.Body).Decode(&newRobotData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	err = h.RobotService.UpdateRobotName(newRobotData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update robot name: %s", err.Error()), http.StatusInternalServerError)
	}
	metrics.UpdateRobotNames.Inc()
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete robot
// @Description Delete robot by ID
// @Tags robots
// @Param id path int true "Robot ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Failed to delete robot"
// @Router /robots/delete/{id} [delete]
func (h *RobotHandlers) DeleteRobot(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.RobotService.DeleteRobot(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete robot: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	metrics.DeletedRobot.Inc()
	w.WriteHeader(http.StatusNoContent)
}

package handlers

import (
	"RobotService/internal/dto"
	"RobotService/internal/prometheusinfo"
	"RobotService/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	_ "RobotService/cmd/docs" // импорт сгенерированных Swagger-доков

	"github.com/go-chi/chi/v5"
)

type RbtHndler struct {
	Srvc services.RbtSrvic
}

func (hndler *RbtHndler) SetRoute(router *chi.Mux) {
	router.Post("/robots/create", hndler.RobotCreate)
	router.Get("/robots/{id}", hndler.GetRobotInfo)
	router.Put("/robots/updatecord", hndler.UpdateRobotCord)
	router.Put("/robots/updatename", hndler.UpdateRobotName)
	router.Put("/robots/updatetype", hndler.ChangeRobotType)
	router.Delete("/robots/delete/{id}", hndler.DeleteRobot)
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
func (hndlr *RbtHndler) RobotCreate(w http.ResponseWriter, r *http.Request) {

	var createdto dto.CreateRobotDTO

	err := json.NewDecoder(r.Body).Decode(&createdto)
	if err != nil {
		http.Error(w, "проблемы с жсоником", http.StatusBadRequest)
		return
	}

	id, err := hndlr.Srvc.CreateRobot(createdto)
	if err != nil {
		http.Error(w, "err", 500)
	}
	prometheusinfo.CreatedRobot.Inc()
	prometheusinfo.CountOfRobotType.WithLabelValues(createdto.Type).Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
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
func (hndl *RbtHndler) GetRobotInfo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "проблемы с жсоником", 400)
		return
	}

	robotinfo, err := hndl.Srvc.GetRobotInfo(id)
	if err != nil {
		http.Error(w, "Error", 500)
		return
	}

	prometheusinfo.GetRobot.Inc()
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
func (hdlr *RbtHndler) UpdateRobotCord(w http.ResponseWriter, r *http.Request) {
	newRobotData := dto.UpdateRobotCordDTO{}
	err := json.NewDecoder(r.Body).Decode(&newRobotData)
	if err != nil {
		http.Error(w, "проблемы с жсоником", 400)
	}

	_ = hdlr.Srvc.UpdateRobotCords(newRobotData)

	prometheusinfo.UpdateRobotCords.Inc()
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
func (handler *RbtHndler) UpdateRobotName(w http.ResponseWriter, r *http.Request) {
	newRobotData := dto.UpdateRobotNameDTO{}
	err := json.NewDecoder(r.Body).Decode(&newRobotData)
	if err != nil {
		http.Error(w, "проблемы с жсоником", 400)
	}

	_ = handler.Srvc.UpdateRobotName(newRobotData)
	prometheusinfo.UpdateRobotNames.Inc()
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Update robot type
// @Description Update robot name by ID
// @Tags robots
// @Accept json
// @Param robot body dto.ChangeTypeDTO true "Updated type"
// @Success 204 {string} string "No Content"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to update robot type"
// @Router /robots/updatetype [put]
func (hdler *RbtHndler) ChangeRobotType(w http.ResponseWriter, r *http.Request) {
	newRobotData := dto.ChangeTypeDTO{}
	err := json.NewDecoder(r.Body).Decode(&newRobotData)
	if err != nil {
		http.Error(w, "проблемы с жсоником", 400)
	}

	_ = hdler.Srvc.ChangeRobotType(newRobotData)

	prometheusinfo.CountOfRobotType.WithLabelValues(newRobotData.Type).Inc()
	prometheusinfo.UpdateRobotType.Inc()
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
func (hnd *RbtHndler) DeleteRobot(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "неверный айди", 400)
		return
	}

	err = hnd.Srvc.DeleteRobot(id)
	if err != nil {
		return
	}
	prometheusinfo.DeletedRobot.Inc()
	w.WriteHeader(http.StatusNoContent)
}

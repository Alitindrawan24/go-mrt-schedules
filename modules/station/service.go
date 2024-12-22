package station

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Alitindrawan24/go-mrt-schedules/common/client"
)

type Service interface {
	GetAllStation() ([]StationResponse, error)
	GetStationSchedule(string) ([]ScheduleResponse, error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *service) GetAllStation() (response []StationResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	stations := []Station{}
	err = json.Unmarshal(byteResponse, &stations)
	if err != nil {
		return nil, err
	}

	for _, station := range stations {
		response = append(response, StationResponse{
			ID:   station.ID,
			Name: station.Name,
		})
	}

	return response, nil
}

func (s *service) GetStationSchedule(id string) (response []ScheduleResponse, err error) {
	url := "https://jakartamrt.co.id/id/val/stasiuns"

	byteResponse, err := client.DoRequest(s.client, url)
	if err != nil {
		return
	}

	schedules := []Schedule{}
	err = json.Unmarshal(byteResponse, &schedules)
	if err != nil {
		return nil, err
	}

	var selectedSchedule Schedule

	for _, schedule := range schedules {
		if schedule.StationID == id {
			selectedSchedule = schedule
			break
		}
	}

	if selectedSchedule.StationID == "" {
		return nil, errors.New("Station not found")
	}

	response, err = convertDataToScheduleResponses(selectedSchedule)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func convertDataToScheduleResponses(schedule Schedule) (response []ScheduleResponse, err error) {
	var (
		lebakBulusTripName = "Stasiun Lebak Bulus Grab"
		bundaranHITripName = "Stasiun Bundaran HI Bank DKI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHI

	scheduleLebakBulusParsed, err := convertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return nil, err
	}

	scheduleBundaranHIParsed, err := convertScheduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return nil, err
	}

	for _, item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: lebakBulusTripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	for _, item := range scheduleBundaranHIParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: bundaranHITripName,
				Time:        item.Format("15:04"),
			})
		}
	}

	return response, nil
}

func convertScheduleToTimeFormat(schedule string) (response []time.Time, err error) {
	var (
		parsedTime time.Time
		schedules  = strings.Split(schedule, ",")
	)

	for _, schedule := range schedules {
		trimmedTime := strings.TrimSpace(schedule)
		if trimmedTime == "" {
			continue
		}

		parsedTime, err = time.Parse("15:04", trimmedTime)
		if err != nil {
			return nil, errors.New("Invalid time format")
		}

		response = append(response, parsedTime)
	}

	return response, nil
}

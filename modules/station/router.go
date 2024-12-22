package station

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Alitindrawan24/go-mrt-schedules/common/response"
)

func Init(router *gin.RouterGroup) {
	stationService := NewService()

	station := router.Group("/stations")

	station.GET("", func(c *gin.Context) {
		GetAllStation(c, stationService)
	})

	station.GET("/:id", func(c *gin.Context) {
		GetStationSchedule(c, stationService)
	})
}

func GetAllStation(c *gin.Context, service Service) {
	data, err := service.GetAllStation()
	if err != nil {
		c.JSON(http.StatusBadGateway, response.ApiResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.ApiResponse{
		Success: true,
		Message: "Successfully get all station",
		Data:    data,
	})
}

func GetStationSchedule(c *gin.Context, service Service) {
	id := c.Param("id")

	data, err := service.GetStationSchedule(id)
	if err != nil {
		c.JSON(http.StatusBadGateway, response.ApiResponse{
			Success: false,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, response.ApiResponse{
		Success: true,
		Message: "Successfully get schedule station",
		Data:    data,
	})
}

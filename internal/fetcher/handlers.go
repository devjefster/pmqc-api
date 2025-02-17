package fetcher

import (
	"net/http"
	"pmqc-api/internal/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FetchDataHandler handles API requests to fetch PMQC data
// @Summary Fetch PMQC Data
// @Description Fetches PMQC data from the Brazilian government API
// @Tags fetcher
// @Accept json
// @Produce json
// @Param year query int true "Year of data to fetch"
// @Param month query int true "Month of data to fetch"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /fetch [get]
func FetchDataHandler(c *gin.Context) {
	yearParam := c.Query("year")
	monthParam := c.Query("month")

	year, err := strconv.Atoi(yearParam)
	month, err2 := strconv.Atoi(monthParam)

	if err != nil || err2 != nil || year < 2016 || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year or month"})
		return
	}

	data, err := FetchPMQCData(year, month)
	if err != nil {
		config.Logger.Errorf("❌ Failed to fetch data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data"})
		return
	}

	config.Logger.Infof("✅ Successfully fetched data for %d-%02d", year, month)
	c.JSON(http.StatusOK, data)
}

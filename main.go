package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	"pmqc-api/config"
	_ "pmqc-api/docs"
	"pmqc-api/fetcher"
	"pmqc-api/models"
	"pmqc-api/storage"
	"sync"
)

// @title PMQC API
// @version 1.0
// @description API for fetching and storing PMQC data.
// @host localhost:8080
// @BasePath /
func main() {
	fmt.Println("🚀 PMQC API Running on Port 8080")

	config.InitLogger()
	config.InitDB()
	defer config.DB.Close()
	defer config.Logger.Sync()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/fetch", fetcher.FetchDataHandler)
	router.GET("/fetch/all", fetchAllHandler)

	router.POST("/store", storage.StorePMQCData)
	router.GET("/amostras", storage.GetAmostras)
	router.GET("/amostras/:id", storage.GetAmostraByID)

	router.Run(":8080")
}

func fetchAllHandler(c *gin.Context) {
	startYear := 2016
	endYear := 2024

	var wg sync.WaitGroup
	results := make(chan *fetcher.FetchResult, 12*(endYear-startYear+1))
	sem := make(chan struct{}, 5)

	for year := startYear; year <= endYear; year++ {
		for month := 1; month <= 12; month++ {
			wg.Add(1)
			go fetchParallel(year, month, &wg, results, sem)
		}
	}

	wg.Wait()
	close(results)

	for res := range results {
		if res.Error == nil {
			if res.Data != nil {
				data, ok := res.Data.(*models.PMQCData) // ✅ Fix: Type assertion
				if !ok {
					config.Logger.Errorf("🚨 Data conversion failed for %d-%02d", res.Year, res.Month)
					return
				}
				sendToStorage(res.Year, res.Month, data) // ✅ Now correctly typed
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parallel fetch complete"})
}

func fetchParallel(year, month int, wg *sync.WaitGroup, results chan<- *fetcher.FetchResult, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}

	data, err := fetcher.FetchPMQCData(year, month)
	if err != nil {
		config.Logger.Errorf("❌ Fetch failed for %d-%02d: %v", year, month, err)
		results <- &fetcher.FetchResult{Year: year, Month: month, Error: err}
		<-sem
		return
	}

	config.Logger.Infof("✅ Successfully fetched %d-%02d!", year, month)
	results <- &fetcher.FetchResult{Year: year, Month: month, Data: data}
	<-sem
}

func sendToStorage(year, month int, data *models.PMQCData) {
	url := "http://localhost:8080/store"

	payload := map[string]interface{}{
		"year":  year,
		"month": month,
		"data":  data,
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		config.Logger.Errorf("🚨 Failed to send data for %d-%02d: %v", year, month, err)
		return
	}
	defer resp.Body.Close()

	config.Logger.Infof("📤 Data for %d-%02d sent to storage. Response: %d", year, month, resp.StatusCode)
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"net/http"
	_ "pmqc-api/docs"
	config2 "pmqc-api/internal/config"
	fetcher2 "pmqc-api/internal/fetcher"
	"pmqc-api/internal/models"
	storage2 "pmqc-api/internal/storage"
	"sync"
)

// @title PMQC API
// @version 1.0
// @description API for fetching and storing PMQC data.
// @host localhost:8080
// @BasePath /
func main() {
	fmt.Println("🚀 PMQC API Running on Port 8080")

	config2.InitLogger()
	config2.InitDB()
	defer config2.DB.Close()
	defer config2.Logger.Sync()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/fetch", fetcher2.FetchDataHandler)
	router.GET("/fetch/all", fetchAllHandler)

	router.POST("/store", storage2.StorePMQCData)
	router.GET("/amostras", storage2.GetAmostras)
	router.GET("/amostras/:id", storage2.GetAmostraByID)

	router.Run(":8080")
}

func fetchAllHandler(c *gin.Context) {
	startYear := 2016
	endYear := 2024

	var wg sync.WaitGroup
	results := make(chan *fetcher2.FetchResult, 12*(endYear-startYear+1))
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
					config2.Logger.Errorf("🚨 Data conversion failed for %d-%02d", res.Year, res.Month)
					return
				}
				sendToStorage(res.Year, res.Month, data) // ✅ Now correctly typed
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Parallel fetch complete"})
}

func fetchParallel(year, month int, wg *sync.WaitGroup, results chan<- *fetcher2.FetchResult, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}

	data, err := fetcher2.FetchPMQCData(year, month)
	if err != nil {
		config2.Logger.Errorf("❌ Fetch failed for %d-%02d: %v", year, month, err)
		results <- &fetcher2.FetchResult{Year: year, Month: month, Error: err}
		<-sem
		return
	}

	config2.Logger.Infof("✅ Successfully fetched %d-%02d!", year, month)
	results <- &fetcher2.FetchResult{Year: year, Month: month, Data: data}
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
		config2.Logger.Errorf("🚨 Failed to send data for %d-%02d: %v", year, month, err)
		return
	}
	defer resp.Body.Close()

	config2.Logger.Infof("📤 Data for %d-%02d sent to storage. Response: %d", year, month, resp.StatusCode)
}

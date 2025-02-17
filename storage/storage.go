package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/jackc/pgx/v5"
	"pmqc-api/config"
	"pmqc-api/models"
)

func StorePMQCData(c *gin.Context) {
	var request struct {
		Year  int             `json:"year"`
		Month int             `json:"month"`
		Data  models.PMQCData `json:"data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	ctx := context.Background()
	tx, err := config.DB.Begin(ctx)
	if err != nil {
		config.Logger.Errorf("❌ Transaction start failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var postosBuffer bytes.Buffer
	for _, stateData := range request.Data.UF {
		for _, municipioData := range stateData.Municipios {
			for _, amostra := range municipioData.Amostras {
				postosBuffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s',%f,%f),",
					amostra.Posto.RazaoSocial, amostra.Posto.CNPJ, amostra.Posto.Distribuidora,
					amostra.Posto.Endereco, amostra.Posto.Complemento, amostra.Posto.Bairro,
					amostra.Posto.Latitude, amostra.Posto.Longitude))
			}
		}
	}

	if postosBuffer.Len() > 0 {
		query := "INSERT INTO postos (razao_social, cnpj, distribuidora, endereco, complemento, bairro, latitude, longitude) VALUES " + postosBuffer.String()[:postosBuffer.Len()-1] + " ON CONFLICT (cnpj) DO NOTHING"
		_, err := tx.Exec(ctx, query)
		if err != nil {
			config.Logger.Errorf("❌ Failed to batch insert postos: %v", err)
			tx.Rollback(ctx)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	copyData := [][]interface{}{}
	for _, stateData := range request.Data.UF {
		for _, municipioData := range stateData.Municipios {
			for _, amostra := range municipioData.Amostras {
				copyData = append(copyData, []interface{}{
					amostra.IdNumeric, amostra.DataColeta, amostra.GrupoProduto, amostra.Produto,
				})
			}
		}
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"amostras"}, []string{"id_numeric", "data_coleta", "grupo_produto", "produto"}, pgx.CopyFromRows(copyData))
	if err != nil {
		config.Logger.Errorf("❌ Failed to batch insert amostras: %v", err)
		tx.Rollback(ctx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err = tx.Commit(ctx); err != nil {
		config.Logger.Errorf("❌ Transaction commit failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	config.Logger.Infof("✅ Data for %d-%02d stored successfully!", request.Year, request.Month)
	c.JSON(http.StatusOK, gin.H{"message": "Data stored successfully"})
}

package storage

import (
	"context"
	"net/http"
	"pmqc-api/internal/config"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AmostraResponse struct {
	ID           int     `json:"id"`
	DataColeta   string  `json:"dataColeta"`
	GrupoProduto string  `json:"grupoProduto"`
	Produto      string  `json:"produto"`
	CNPJ         string  `json:"cnpj"`
	RazaoSocial  string  `json:"razaoSocial"`
	Municipio    string  `json:"municipio"`
	Estado       string  `json:"estado"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
}

// GetAmostraByID fetches a single amostra by ID
// @Summary Get a specific amostra by ID
// @Description Retrieves details of an amostra including its associated posto
// @Tags storage
// @Accept json
// @Produce json
// @Param id path int true "Amostra ID"
// @Success 200 {object} AmostraResponse
// @Failure 400 {object} map[string]string "Invalid ID format"
// @Failure 404 {object} map[string]string "Amostra not found"
// @Router /amostras/{id} [get]
func GetAmostraByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	query := `
		SELECT a.id, a.data_coleta, a.grupo_produto, a.produto, p.cnpj, p.razao_social, a.municipio, a.estado, p.latitude, p.longitude 
		FROM amostras a
		JOIN postos p ON a.posto_id = p.id
		WHERE a.id = $1
	`
	row := config.DB.QueryRow(context.Background(), query, id)

	var amostra AmostraResponse
	err = row.Scan(&amostra.ID, &amostra.DataColeta, &amostra.GrupoProduto, &amostra.Produto, &amostra.CNPJ, &amostra.RazaoSocial, &amostra.Municipio, &amostra.Estado, &amostra.Latitude, &amostra.Longitude)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Amostra not found"})
		return
	}

	c.JSON(http.StatusOK, amostra)
}

// GetAmostras fetches all amostras with pagination and optional filters
// @Summary Get all amostras with pagination and filters
// @Description Retrieves a paginated list of amostras, filtered by CNPJ, RazaoSocial, Produto, or Municipio
// @Tags storage
// @Accept json
// @Produce json
// @Param limit query int false "Number of records per page" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Param cnpj query string false "CNPJ of the posto"
// @Param razaoSocial query string false "Razao Social of the posto"
// @Param produto query string false "Product name"
// @Param municipio query string false "Municipio name"
// @Success 200 {array} AmostraResponse
// @Failure 500 {object} map[string]string "Database error"
// @Router /amostras [get]
func GetAmostras(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Dynamic WHERE conditions
	var conditions []string
	var values []interface{}
	paramIndex := 1

	if cnpj := c.Query("cnpj"); cnpj != "" {
		conditions = append(conditions, "p.cnpj = $"+strconv.Itoa(paramIndex))
		values = append(values, cnpj)
		paramIndex++
	}

	if razaoSocial := c.Query("razaoSocial"); razaoSocial != "" {
		conditions = append(conditions, "p.razao_social ILIKE $"+strconv.Itoa(paramIndex))
		values = append(values, "%"+razaoSocial+"%")
		paramIndex++
	}

	if produto := c.Query("produto"); produto != "" {
		conditions = append(conditions, "a.produto ILIKE $"+strconv.Itoa(paramIndex))
		values = append(values, "%"+produto+"%")
		paramIndex++
	}

	if municipio := c.Query("municipio"); municipio != "" {
		conditions = append(conditions, "a.municipio ILIKE $"+strconv.Itoa(paramIndex))
		values = append(values, "%"+municipio+"%")
		paramIndex++
	}

	// Build SQL query
	query := `
		SELECT a.id, a.data_coleta, a.grupo_produto, a.produto, p.cnpj, p.razao_social, a.municipio, a.estado, p.latitude, p.longitude 
		FROM amostras a
		JOIN postos p ON a.posto_id = p.id
	`
	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " ORDER BY a.id DESC LIMIT $" + strconv.Itoa(paramIndex) + " OFFSET $" + strconv.Itoa(paramIndex+1)
	values = append(values, limit, offset)

	// Execute Query
	rows, err := config.DB.Query(context.Background(), query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var amostras []AmostraResponse
	for rows.Next() {
		var amostra AmostraResponse
		err := rows.Scan(&amostra.ID, &amostra.DataColeta, &amostra.GrupoProduto, &amostra.Produto, &amostra.CNPJ, &amostra.RazaoSocial, &amostra.Municipio, &amostra.Estado, &amostra.Latitude, &amostra.Longitude)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse results"})
			return
		}
		amostras = append(amostras, amostra)
	}

	c.JSON(http.StatusOK, amostras)
}

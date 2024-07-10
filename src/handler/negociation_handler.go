package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vinaocruz/go-extractor/src/repository"
)

type NegociationHandler struct {
	Repository repository.NegociationRepository
}

type Config struct {
	Router *gin.Engine
}

func NewNegociationHandler(c *Config) error {
	handler := &NegociationHandler{
		Repository: repository.NewNegociationRepository(),
	}

	c.Router.GET("/negociations", handler.getNegociations)

	return nil
}

func (h *NegociationHandler) getNegociations(c *gin.Context) {
	ticker := c.Query("ticker")
	transactionAt := c.Query("DataNegocio")

	dto, err := h.Repository.Find(ticker, transactionAt)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	dto.Ticker = ticker
	c.JSON(http.StatusOK, dto)
}

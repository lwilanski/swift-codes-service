package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/your-github-name/swift-codes-service/internal/models"
	"github.com/your-github-name/swift-codes-service/internal/repository"
)

type Handler struct{ repo repository.SwiftRepo }

func New(r repository.SwiftRepo) Handler { return Handler{r} }

/*  GET /v1/swift-codes/:code  */
func (h Handler) GetByCode(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	item, err := h.repo.Get(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "code not found"})
		return
	}
	if item.IsHeadquarter {
		branches, _ := h.repo.GetBranches(item.Root())
		c.JSON(http.StatusOK, HQResponse(item, branches))
		return
	}
	c.JSON(http.StatusOK, BranchResponse(item))
}

/*  GET /v1/swift-codes/country/:iso2  */
func (h Handler) Country(c *gin.Context) {
	iso2 := strings.ToUpper(c.Param("iso2"))
	all, err := h.repo.CountryAll(iso2)
	if err != nil || len(all) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "country not found"})
		return
	}
	c.JSON(http.StatusOK, CountryResponse(all))
}

/*  POST /v1/swift-codes  */
func (h Handler) Create(c *gin.Context) {
	var req models.SwiftCode
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	req.CountryISO2 = strings.ToUpper(req.CountryISO2)
	req.CountryName = strings.ToUpper(req.CountryName)

	if err := h.repo.UpsertMany([]models.SwiftCode{req}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "created"})
}

/*  DELETE /v1/swift-codes/:code  */
func (h Handler) Delete(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	if err := h.repo.Delete(code); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

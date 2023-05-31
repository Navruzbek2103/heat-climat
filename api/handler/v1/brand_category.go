package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/climate.uz/api/models"
	"gitlab.com/climate.uz/internal/controller/storage/repo"
)

// @Summary Add brand_category
// @Description Through this api can add brand_category
// @Tags Brand
// @Security	BearerAuth
// @Accept json
// @Produce json
// @Param	brand_category  body models.BrandCategoryReq true "BrandCategory"
// @Success 201 {object} models.SuccessInfo
// @Failure 400 {object} models.FailureInfo
// @Failure 500 {object} models.FailureInfo
// @Router /brand-category [post]
func (h *handlerV1) CreateBrandCategory(c *gin.Context) {
	_, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Invalid access token",
			Error:   err,
		})
		h.log.Error("Error while getting claims of access token ", err)
		return
	}

	var (
		body models.BrandCategoryReq
	)

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while binding request", err)
		return
	}

	err = h.storage.BrandCategory().CreateBrandCategory(&repo.BrandCategoryReq{
		BrandId:    body.BrandId,
		CategoryId: body.CategoryId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Ooops something went wrong",
		})
		h.log.Error("Error while inserting brand_category ids: ", err)
	}
	c.JSON(http.StatusOK, models.SuccessInfo{
		Message:    "Successfully insert brand and category id",
		StatusCode: http.StatusOK,
	})
}

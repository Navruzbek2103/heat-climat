package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/climate.uz/api/models"
	"gitlab.com/climate.uz/internal/controller/storage/repo"
)

// New Brand
// @Summary Add new brand
// @Description This method add new brand
// @Tags 		Brand
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		brand 	body 		models.BrandRequest true "AddNewBrand"
// @Success 	201 	{object}	models.BrandResponse
// @Failure 	400 	{object} 	models.FailureInfo
// @Failure 	403		{object} 	models.FailureInfo
// @Failure 	500 	{object}	models.FailureInfo
// @Router  	/brand 	[post]
func (h *handlerV1) CreateBrand(c *gin.Context) {
	_, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Invalid access token",
		})
		h.log.Error("Error while getting claims of access token ", err)
		return
	}
	var (
		body models.BrandRequest
	)
	err = c.ShouldBindJSON(&body)
	fmt.Println(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while binding json", err)
		return
	}

	data, err := h.storage.Brand().CreateBrand(&repo.BrandRequst{
		BrandName: body.BrandName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while create brand info: ", err)
		return
	}
	c.JSON(http.StatusCreated, data)
}

// @Summary 	Update brand info
// @Description Through this api can update brand info
// @Tags 		Brand
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		brand 	body 		models.BrandUpdateReq true "Brand Update"
// @Success 	200	  	{object} 	models.BrandResponse
// @Failure		400   	{object}  	models.FailureInfo
// @Failure		500   	{object} 	models.FailureInfo
// @Router		/brand	[patch]
func (h *handlerV1) UpdateBrandInfo(c *gin.Context) {
	_, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Invalid access token",
		})
		h.log.Error("Error while getting claims of access token ", err)
		return
	}

	var (
		body models.BrandUpdateReq
	)
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while binding request: ", err)
		return
	}
	res, err := h.storage.Brand().UpdateBrand(&repo.BrandUpdateReq{
		Id:       body.Id,
		BandName: body.BrandName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while update brand info: ", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary Get brand info
// @Description Through this api can get brand info
// @Tags  		Brand
// @Accept 		json
// @Produce 	json
// @Param  		id 	path int true "id"
// @Success 	200 {object} models.GetBrandInfo
// @Failure 	400 {object} models.FailureInfo
// @Failure		500 {object} models.FailureInfo
// @Router  	/brand/{id} [get]
func (h *handlerV1) GetBrandId(c *gin.Context) {
	id := c.Param("id")
	brandId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while binding request: ", err)
		return
	}
	data, err := h.storage.Brand().GetBrandById(&repo.BrandId{
		Id: brandId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Ooops something went wrong",
		})
		h.log.Error("Error while getting brand info", err)
		return
	}
	info := models.GetBrandInfo{
		Id:        data.Id,
		BrandName: data.BrandName,
	}
	for _, c := range data.Categories {
		info.Categories = append(info.Categories, models.BrandCategories{
			Id:           c.Id,
			CategoryName: c.CategoryName})
	}
	c.JSON(http.StatusOK, info)
}

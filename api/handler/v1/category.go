package v1

import (
	"net/http"
	"strconv"

	"gitlab.com/climate.uz/api/models"
	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Add New category
// @Summary 		Create new categorys
// @Description 	this method create new category
// @Tags			Category
// @Security		BearerAuth
// @Accept 			json
// @Produce			json
// @Param			category  body models.CategoryReq true "CategoryRequest"
// @Success 		200 {object} repo.CategoryResponse
// @Failure 		400 {object} models.FailureInfo
// @Failure 		500 {object} models.FailureInfo
// @Router 			/category	[post]
func (h *handlerV1) CreateCategory(c *gin.Context) {
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
		body models.CategoryReq
	)

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("error while binding from request", err)
		return
	}

	res, err := h.storage.Category().CreateCategory(&repo.CategoryRequest{
		CategoryName: body.CategoryName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "error while create category",
		})
		h.log.Error("error while binding from request", err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Get category
// @Summary 		Get  category
// @Description 	this method get category by id
// @Tags			Category
// @Accept 			json
// @Produce			json
// @Param			id  path int true "category id"
// @Success 		200 {object} repo.CategoryResponse
// @Failure 		400 {object} models.FailureInfo
// @Failure 		500 {object} models.FailureInfo
// @Router 			/category/{id}	[get]
func (h *handlerV1) GetCategory(c *gin.Context) {
	id := c.Param("id")
	category_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("error while parsing int from string", err)
		return
	}
	res, err := h.storage.Category().GetCategoryById(int(category_id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "error while get category",
		})
		h.log.Error("error while binding from request", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Update category
// @Summary 		Update  category
// @Description 	this method renames the category with id
// @Tags			Category
// @Accept 			json
// @Produce			json
// @Param			page  	query 	string false "page"
// @Param 			limit 	query 	string false "limit"
// @Param 			search	query	string false "search"
// @Success 		200 		{object} repo.CategoryResponse
// @Failure 		400 		{object} models.FailureInfo
// @Router 			/categories	[get]
func (h *handlerV1) GetAllGategory(c *gin.Context) {

	queryparams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryparams)

	if errStr != nil {
		c.JSON(http.StatusBadRequest, &models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("error while parsing int from string", errStr[0])
		return
	}

	res, errs := h.storage.Category().GetAllCategories(&repo.AllCategoriesParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Search: params.Search,
	})
	if errs != nil {
		c.JSON(http.StatusInternalServerError, &models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "error while get category",
		})
		h.log.Error("error while binding from request", errs)
		return
	}
	c.JSON(http.StatusOK, res)

}

// Update category
// @Summary 		Update  category
// @Description 	this method renames the category with id
// @Tags			Category
// @Security		BearerAuth
// @Accept 			json
// @Produce			json
// @Param			category  	body models.CategoryUpdateReq true "UpdateCategoryReq"
// @Success 		200 		{object} repo.CategoryResponse
// @Failure 		400 		{object} models.FailureInfo
// @Router 			/category	[patch]
func (h *handlerV1) EditCategoryName(c *gin.Context) {
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
		body models.CategoryUpdateReq
	)
	err = c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, &models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("error while binding from request", err)
		return
	}
	res, err := h.storage.Category().UpdateCategory(&repo.CategoryUpdateReq{
		Id:           body.Id,
		CategoryName: body.CategoryName,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "error while update category",
		})
		h.log.Error("error while binding from request", err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Delete category
// @Summary 		Delete  category
// @Description 	This function is for deleting categories by Id
// @Tags			Category
// @Security		BearerAuth
// @Accept 			json
// @Produce			json
// @Param			id  	path int true "categoryId"
// @Success 		200 		{object} repo.Empty
// @Failure 		400 		{object} models.FailureInfo
// @Router 			/category/{id}	[delete]
func (h *handlerV1) DeleteCategory(c *gin.Context) {
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
	id := c.Param("id")
	category_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("error while deleting category by id", err)
		return
	}

	_, err = h.storage.Category().DeleteCategoryById(int(category_id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Error while delete category",
		})
		h.log.Error("error while binding from request ", err)
		return
	}
	c.JSON(http.StatusOK, &models.SuccessInfo{
		StatusCode: http.StatusOK,
		Message:    "Category successfully deleted",
	})
}

package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/climate.uz/api/models"
	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/utils"
)

// New news product
// @Summary Add news product
// @Description This method add news products
// @Tags 		News
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		news 	body 		models.NewProductReq true "NewsProductRequest"
// @Success 	201 	{object}	models.NewProductRes
// @Failure 	400 	{object} 	models.FailureInfo
// @Failure 	403		{object} 	models.FailureInfo
// @Failure 	500 	{object}	models.FailureInfo
// @Router  	/new-product 	[post]
func (h *handlerV1) CreateNewProduct(c *gin.Context) {
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
		body models.NewProductReq
	)
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("Error while binding request", err)
		return
	}

	char, err := json.Marshal(body.Charactestics)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Invalid characteristic info",
		})
		h.log.Error("Error while marshaling json: ", err)
		return
	}

	req := &repo.NewProductRequest{
		BrandId:       body.BrandId,
		Title:         body.Title,
		Description:   body.Description,
		NewPrice:      body.NewPrice,
		OldPrice:      body.OldPrice,
		DisplayType:   body.DisplayType,
		OsType:        body.OsType,
		Camera:        body.Camera,
		Dioganal:      body.Dioganal,
		Charactestics: string(char),
	}

	for _, m := range body.MediaLinks {
		req.MediaLinks = append(req.MediaLinks, &repo.MediaReq{
			MediaLink: m.MediaLink,
		})
	}
	res, err := h.storage.NewProduct().CreateNewProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("error while create news product", err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Get news products
// @Summary 	get news products
// @Description This method get news product info by id
// @Tags 		News
// @Accept 		json
// @Produce 	json
// @Param 		id 		path 	 int true "newproductID"
// @Success 	200 	{object} models.NewProductRes
// @Success 	400 	{object} models.FailureInfo
// @Failure 	500 	{object} models.FailureInfo
// @Router 		/new-product/{id} 	[get]
func (h *handlerV1) GetNewProduct(c *gin.Context) {
	id := c.Param("id")
	news_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("error while parsing news product id", err.Error())
		return
	}

	res, err := h.storage.NewProduct().GetNewProductById(news_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while getting news product by id", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Get All news products
// @Summary 	get all news products
// @Description Through this api can get all new products
// @Tags 		News
// @Accept 		json
// @Produce 	json
// @Param 		page 	query 		int 	false "page"
// @Param 		limit 	query 		int 	false "limit"
// @Param 		search 	query 		string 	false "search"
// @Success 	200 	{object}	models.AllNewProducts
// @Success 	400 	{object}	models.FailureInfo
// @Failure 	500 	{object}	models.FailureInfo
// @Router 		/new-products	[get]
func (h *handlerV1) GetAllNewsProducts(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("error while parsequeryparams", errStr)
		return
	}

	res, errs := h.storage.NewProduct().GetAllNewProducts(&repo.AllNewProductsParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Search: params.Search,
	})
	if errs != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Something went wrong",
		})
		h.log.Error("Error while getting all news product", errs)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Edit News product
// @Summary 	update 	news product
// @Description this 	method update news product
// @Tags 		News
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		news 	body 	 models.UpdateNewProductReq true "UpdateRequest"
// @Success 	200 	{object} models.NewProductRes
// @Failure		400 	{object} models.FailureInfo
// @Failure 	500 	{object} models.FailureInfo
// @Router 		/new-product 	[patch]
func (h *handlerV1) UpdateNewsProduct(c *gin.Context) {
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
		body models.UpdateNewProductReq
	)
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("Error while binding request", err)
		return
	}
	char, err := json.Marshal(&body.Charactestics)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while marshaling json", err.Error())
		return
	}

	req := &repo.NewProductUpdateReq{
		Id:            int64(body.Id),
		BrandId:       body.BrandId,
		Title:         body.Title,
		Description:   body.Description,
		NewPrice:      body.NewPrice,
		OldPrice:      body.OldPrice,
		DisplayType:   body.DisplayType,
		OsType:        body.OsType,
		Camera:        body.Camera,
		Dioganal:      body.Dioganal,
		Charactestics: string(char),
	}

	for _, m := range body.MediaLinks {
		req.MediaLink = append(req.MediaLink, &repo.UpdateMediaLink{
			Id:        m.Id,
			MediaLink: m.MediaLink,
		})
	}
	res, err := h.storage.NewProduct().UpdateNewProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("error while updating news product", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Delete News product
// @Summary 	delete news product info
// @Description this 	method delete news product by id
// @Tags 		News
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		id 	    	path 	int 	true "id"
// @Success 	200 		{object} models.SuccessInfo
// @Failure 	400 		{object} models.FailureInfo
// @Failure 	500 		{object} models.FailureInfo
// @Router 		/new-product/{id}	[delete]
func (h *handlerV1) DeleteNewsProduct(c *gin.Context) {
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

	newId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("error while parsing id", err)
		return
	}
	_, err = h.storage.NewProduct().DeleteNewProductById(newId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while deleting news product", err)
		return
	}
	c.JSON(http.StatusOK, models.SuccessInfo{
		StatusCode: http.StatusOK,
		Message:    "Successfully deleted news product info",
	})
}

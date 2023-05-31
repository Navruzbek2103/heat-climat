package v1

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gitlab.com/climate.uz/api/models"
	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/utils"

	"github.com/gin-gonic/gin"
)

// New add product
// @Summary 	create new product
// @Description Through this api can create product
// @Tags 		Product
// @Security	BearerAuth
// @Accept 		json
// @Produce		json
// @Param		product 	body 	 models.ProductReq true "ProductRequest"
// @Success		201 		{object} repo.ProductResponse
// @Failure 	400 		{object} models.FailureInfo
// @Failure 	500 		{object} models.FailureInfo
// @Router 		/product	[post]
func (h *handlerV1) CreateProduct(c *gin.Context) {
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
		body models.ProductReq
	)
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("Error while binding from request")
		return
	}
	char, err := json.Marshal(body.Characteristics)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Invalid characteristic info",
		})
		h.log.Error("Error while marshaling json: ", err)
		return
	}
	req := &repo.ProductRequest{
		BrandId:         body.BrandId,
		CategoryId:      body.CategoryId,
		Title:           body.Title,
		Description:     body.Description,
		Price:           body.Price,
		DisplayType:     body.DisplayType,
		OsType:          body.OsType,
		Camera:          body.Camera,
		Dioganal:        body.Dioganal,
		Characteristics: string(char),
	}

	for _, m := range body.MediaLinks {
		req.MediaLinks = append(req.MediaLinks, &repo.MediaReq{
			MediaLink: m.MediaLink,
		})
	}

	res, err := h.storage.Product().CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while create product:", err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Get product
// @Summary 	get one product
// @Description this method get product by Id
// @Tags 		Product
// @Accept 		json
// @Produce		json
// @Param		id 			path 	 int true "product_id"
// @Success		200 		{object} repo.ProductResponse
// @Failure 	400 		{object} models.FailureInfo
// @Failure 	500 		{object} models.FailureInfo
// @Router 		/product/{id}		 [get]
func (h *handlerV1) GetProductInfo(c *gin.Context) {
	id := c.Param("id")
	product_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("Error parsing product id", err)
		return
	}
	res, err := h.storage.Product().GetProductById(product_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while getting product by Id: ", err)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Update product
// @Summary 	edit product info
// @Description this method edit product info
// @Tags 		Product
// @Security	BearerAuth
// @Accept 		json
// @Produce		json
// @Param		product    	body		models.UpdateProductReq true "UpdateProductRequest"
// @Success		200 		{object} 	repo.ProductResponse
// @Failure 	400 		{object} models.FailureInfo
// @Failure 	500 		{object} 	models.FailureInfo
// @Router 		/product 	[put]
func (h *handlerV1) UpdateProduct(c *gin.Context) {
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
		body models.UpdateProductReq
	)
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("Error while binding from request ", err)
		return
	}

	charJson, err := json.Marshal(body.Characteristics)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please, Enter right info",
		})
		h.log.Error("Error while marshaling json: ", err)
		return
	}

	req := &repo.ProductUpdateReq{
		Id:              body.Id,
		Title:           body.Title,
		Description:     body.Description,
		Price:           body.Price,
		DisplayType:     body.DisplayType,
		OsType:          body.OsType,
		Camera:          body.Camera,
		Dioganal:        body.Dioganal,
		Characteristics: string(charJson),
	}

	for _, m := range body.MediaLinks {
		req.MediaLink = append(req.MediaLink, &repo.UpdateMediaLink{
			Id:        m.Id,
			MediaLink: m.MediaLink,
		})
	}

	res, err := h.storage.Product().UpdateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while updating product info", err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// Get all product
// @Summary 	get all products
// @Description this method get all products list
// @Tags 		Product
// @Accept 		json
// @Produce		json
// @Param		page 		query 	 int false "page"
// @Param 		limit 		query 	 int false "limit"
// @Param 		search 		query 	 string false "search"
// @Success		200 		{object} models.AllProducts
// @Failure 	400 		{object} models.FailureInfo
// @Failure 	500 		{object} models.FailureInfo
// @Router 		/products 	[get]
func (h *handlerV1) GetAllProduct(c *gin.Context) {
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
	res, errs := h.storage.Product().GetAllProducts(&repo.AllProductsParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Search: params.Search,
	})
	if errs != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while getting products list", errs)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Delete product
// @Summary 	delete product info
// @Description this method delete product info by productId
// @Tags 		Product
// @Security	BearerAuth
// @Accept 		json
// @Produce		json
// @Param		id    			path string		true "product_id"
// @Success		200 			{object} 		models.SuccessInfo
// @Failure 	400 			{object} 		models.FailureInfo
// @Failure 	500 			{object} 		models.FailureInfo
// @Router 		/product/{id} 	[delete]
func (h *handlerV1) DeleteProductInfo(c *gin.Context) {
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

	product_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right wrong",
		})
		h.log.Error("Error while parsing int from string", err)
		return
	}
	_, err = h.storage.Product().DeleteProductById(product_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("error while deleting product by id", err)
		return
	}
	c.JSON(http.StatusOK, models.SuccessInfo{
		Message:    "Product successfully deleted",
		StatusCode: http.StatusOK,
	})
}

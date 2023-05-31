package v1

import (
	"net/http"
	"strconv"

	"gitlab.com/climate.uz/api/models"
	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/utils"

	"github.com/gin-gonic/gin"
)

// New add user
// @Summary 	add user phone_number
// @Description This function create new user
// @Tags 		User
// @Accept 		json
// @Produce 	json
// @Param  		user	body 		models.UserReq true "UserRequest"
// @Success 	201 	{object} 	models.UserRes
// @Failure 	400 	{object} 	models.FailureInfo
// @Failure 	500 	{object} 	models.FailureInfo
// @Router 		/user	[post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body models.UserReq
	)
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("error while binding from request", err)
		return
	}
	res, err := h.storage.User().CreateUser(&repo.UserRequest{
		PhoneNumber: body.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Error while creating user",
		})
		h.log.Error("error while binding from request", err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Get user
// @Summary 	Get user
// @Description This function get user by id
// @Tags 		User
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param  		id 			path 		int true "user_id"
// @Success 	200 		{object} 	models.UserRes
// @Failure 	400 		{object} 	models.FailureInfo
// @Failure 	500 		{object} 	models.FailureInfo
// @Router 		/user/{id}	[get]
func (h *handlerV1) GetUser(c *gin.Context) {
	_, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Invalid access token",
		})
		h.log.Error("Error while getting claims of access token ", err)
		return
	}
	id := c.Param("id")
	user_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while parsing user id", err)
	}
	res, err := h.storage.User().GetUserById(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "please enter right info",
		})
		h.log.Error("error while getting user info by id")
		return
	}
	c.JSON(http.StatusOK, res)
}

// Get user
// @Summary 	Get all users
// @Description This method is to get all users
// @Tags 		User
// @Accept 		json
// @Produce 	json
// @Param  		page		query 	    string 	true "page"
// @Param		limit		query 		int 	true "limit"
// @Param		search		query 		int 	true "search"
// @Success 	200 		{object} 	models.UserRes
// @Failure 	500 		{object} 	models.FailureInfo
// @Router 		/users		[get]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
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
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("error while binding from request", errStr)
		return
	}
	res, errs := h.storage.User().GetAllUser(&repo.AllUsersParams{
		Page:   params.Page,
		Limit:  params.Limit,
		Search: params.Search,
	})
	if errs != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("error while getting users", errs)
		return
	}
	c.JSON(http.StatusOK, res)
}

// Update user
// @Summary 	Update user info
// @Description This method update user info
// @Tags 		User
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param  		user 		body 		models.UpdateUserReq true "UpdateUserReq"
// @Success 	200 		{object} 	models.UserRes
// @Failure 	500 		{object} 	models.FailureInfo
// @Router 		/user		[patch]
func (h *handlerV1) UpdateUser(c *gin.Context) {
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
		body models.UpdateUserReq
	)
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("error while binding from request ", err)
		return
	}
	res, err := h.storage.User().UpdateUser(&repo.UserUpdateReq{
		Id:          body.Id,
		PhoneNumber: body.PhoneNumber,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("error while updating user info", err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Delete user
// @Summary 	Delete user info
// @Description This method delete user info
// @Tags 		User
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param  		id 			path 		int true "user_id"
// @Success 	200 		{object} 	models.SuccessInfo
// @Failure 	500 		{object} 	models.FailureInfo
// @Router 		/user/{id}	[delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
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
	user_id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while parsing user id")
	}

	_, err = h.storage.User().DeleteUser(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please enter right info",
		})
		h.log.Error("Error while deleting user")
		return
	}

	c.JSON(http.StatusOK, models.SuccessInfo{
		Message:    "User successfully deleted",
		StatusCode: http.StatusOK,
	})
}

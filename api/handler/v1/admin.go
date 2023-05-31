package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/climate.uz/api/models"
	"gitlab.com/climate.uz/internal/controller/storage/repo"
	"gitlab.com/climate.uz/pkg/etc"
)

// admin login
// @Summary 	Login admin
// @Description Through this api admin can login
// @Tags 		Admin
// @Accept 		json
// @Produce 	json
// @Param  		username path 	  string true "username"
// @Param   	password path 	  string true "password"
// @Success 	200 	 {object} models.AdminLoginRes
// @Failure 	500 	 {object} models.FailureInfo
// @Failure 	400 	 {object} models.FailureInfo
// @Failure 	409 	 {object} models.FailureInfo
// @Router 		/admin/{username}/{password} 	[get]
func (h *handlerV1) AdminLogin(c *gin.Context) {
	var (
		username = c.Param("username")
		password = c.Param("password")
	)
	res, err := h.storage.Admin().GetAdmin(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "please enter right info",
		})
		h.log.Error("error while getting admin info", err)
		return
	}

	if !etc.CheckPasswordHash(password, res.Password) {
		c.JSON(http.StatusConflict, models.FailureInfo{
			Code:    http.StatusConflict,
			Message: "username or password error",
		})
		h.log.Error("error checking password", err)
		return
	}
	h.jwthandler.Role = "admin"
	h.jwthandler.Sub = res.Id
	h.jwthandler.Aud = []string{"cliamte.uz"}
	h.jwthandler.SigninKey = h.cfg.SigninKey
	h.jwthandler.Log = h.log
	tokens, err := h.jwthandler.GenerateAuthJWT()
	accessToken := tokens[0]
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "something went wrong, please try again",
		})
		h.log.Error("error occured while generating tokens ", err)
		return
	}
	res.AccessToken = accessToken
	res.Password = ""
	info := models.AdminLoginRes{
		Id:          res.Id,
		UserName:    res.UserName,
		Password:    res.Password,
		AccessToken: res.AccessToken,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}
	c.JSON(http.StatusOK, info)
}

// admin register
// @Summary 	Add admin
// @Description Through this api new admin can register
// @Tags 		Admin
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param  		admin 	 body models.AdminReq true "AddAdmin"
// @Success 	201 	 {object} models.AdminRes
// @Failure 	400 	 {object} models.FailureInfo
// @Failure		409		 {object} models.FailureInfo
// @Failure 	500 	 {object} models.FailureInfo
// @Router 		/admin	 [post]
func (h *handlerV1) AddAdmin(c *gin.Context) {
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
		body models.AdminReq
	)
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Something went wrong",
		})
		h.log.Error("Error while adding new admin", err)
		return
	}
	IsUserName, err := h.storage.Admin().CheckField(&repo.CheckFieldReq{
		Field: "user_name",
		Value: body.UserName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while checking username uniquness: ", err)
	}
	if IsUserName.Exists {
		c.JSON(http.StatusConflict, models.FailureInfo{
			Code:    http.StatusConflict,
			Message: "username already in use",
		})
		return
	}
	admin_password, err := etc.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Oops something went wrong",
		})
		h.log.Error("Error while hashing admin password", err)
		return
	}

	adminId := uuid.New().String()
	h.jwthandler.Sub = adminId
	h.jwthandler.Role = "admin"
	h.jwthandler.Aud = []string{"climate-frontend"}
	h.jwthandler.SigninKey = h.cfg.SigninKey
	tokens, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Ooops something went wrong",
		})
		h.log.Error("Error while creating jwt token", err)
		return
	}
	refresh_token := tokens[1]
	data, err := h.storage.Admin().AddAdmin(&repo.AdminRequest{
		Id:           adminId,
		UserName:     body.UserName,
		Password:     admin_password,
		RefreshToken: refresh_token,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Oops something went wrong",
		})
		h.log.Error("Error while creating new admin")
	}
	data.AccessToken = tokens[0]
	c.JSON(http.StatusCreated, data)
}

// @Summary 	Get admin profile
// @Description Through this api admin can get profile
// @Tags 		Admin
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Success 	200 	 {object} models.GetAdminProdile
// @Failure 	500 	 {object} models.FailureInfo
// @Router 		/admin/profile 	[get]
func (h *handlerV1) GetAdminProfile(c *gin.Context) {
	claim, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Invalid access token",
			Error:   err,
		})
		h.log.Error("Error while getting claims of access token ", err)
		return
	}
	data, err := h.storage.Admin().GetAdminInfo(claim.Sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Ooops something went wrong",
		})
		h.log.Error("Error while getting admin info: ", err)
		return
	}
	info := models.GetAdminProdile{
		Id:        data.Id,
		UserName:  data.UserName,
		Password:  data.Password,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	c.JSON(http.StatusOK, info)
}

// @Summary 	Get all admin
// @Description Through this api admin can get all admins
// @Tags 		Admin
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		keyword  path string true "search"
// @Success 	200 	 {object} models.AllAdmins
// @Failure 	400 	 {object} models.FailureInfo
// @Failure 	500 	 {object} models.FailureInfo
// @Router 		/admins  [get]
func (h *handlerV1) GetAllAdmin(c *gin.Context) {
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
	keyword := c.Param("keyword")

	data, err := h.storage.Admin().GetAllAdmins(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while get all admins: ", err)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary 	Update admin info
// @Description Through this api admin can update info
// @Tags 		Admin
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		admin    body 	  models.AdminReq true "Admin Update"
// @Success 	200 	 {object} models.GetAdminProdile
// @Failure 	400 	 {object} models.FailureInfo
// @Failure 	500 	 {object} models.FailureInfo
// @Router 		/admin   [patch]
func (h *handlerV1) UpdateAdminInfo(c *gin.Context) {
	claim, err := GetClaims(*h, c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Invalid access token",
			Error:   err,
		})
		h.log.Error("Error while getting claims of access token ", err)
		return
	}

	var body models.AdminReq
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Enter right info",
		})
		h.log.Error("Error while binding json: ", err)
		return
	}
	password, err := etc.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Something went wrong",
		})
		h.log.Error("Error while hashing password: ", err)
		return
	}
	data, err := h.storage.Admin().UpdateAdmin(&repo.UpdateAdminReq{
		Id:       claim.Sub,
		UserName: body.UserName,
		Password: password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Ooops something went wrong",
		})
		h.log.Error("Error while update admin info: ", err)
		return
	}
	c.JSON(http.StatusOK, data)
}

// @Summary 	Delete admin info
// @Description Through this api admin can delete info
// @Tags 		Admin
// @Security	BearerAuth
// @Accept 		json
// @Produce 	json
// @Param 		id 		 path     string  true "id"
// @Success 	200 	 {object} models.SuccessInfo
// @Failure 	400 	 {object} models.FailureInfo
// @Failure 	500 	 {object} models.FailureInfo
// @Router 		/admin/{id}   [delete]
func (h *handlerV1) DeleteAdminInfo(c *gin.Context) {
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
	_, err = h.storage.Admin().DeleteAdmin(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Ooops something went wrong",
		})
		h.log.Error("Error while deleting admin info: ", err)
		return
	}
	c.JSON(http.StatusOK, models.SuccessInfo{
		Message:    "Successfully deleted admin info",
		StatusCode: http.StatusOK,
	})
}

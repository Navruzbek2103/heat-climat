package v1

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/climate.uz/api/models"
)

// @Summary upload media file
// @Description Through this api can upload media file
// @Tags 		Media
// @Accept 		json
// @Produce 	json
// @Param 		file 	formData file true "image-upload"
// @Success  	201 	{object} models.FileResponse
// @Failure  	400 	{object} models.FailureInfo
// @Failure  	500 	{object} models.FailureInfo
// @Router 	 	/image-upload  	 [post]
func (h *handlerV1) UploadMedia(c *gin.Context) {
	var (
		file models.File
	)
	err := c.ShouldBind(&file)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Code:    http.StatusBadRequest,
			Message: "Please check your data",
		})
		h.log.Error("Error while binding request: ", err)
		return
	}

	ext := filepath.Ext(file.File.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".mp4" {
		c.JSON(http.StatusBadRequest, models.FailureInfo{
			Message: "Unsupported file format",
		})
		return
	}

	mediaType := "image"
	if ext == ".mp4" {
		mediaType = "video"
	}

	fileName := uuid.New().String() + filepath.Ext(file.File.Filename)
	dst, _ := os.Getwd()

	if _, err := os.Stat(dst + "/media"); os.IsNotExist(err) {
		os.Mkdir(dst+"/media", os.ModePerm)
	}

	filePath := "/media/" + fileName

	err = c.SaveUploadedFile(file.File, dst+filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FailureInfo{
			Code:    http.StatusInternalServerError,
			Message: "Ooops something went wrong",
		})
		h.log.Error("Error while save upload media file: ", err)
		return
	}

	c.JSON(http.StatusOK, models.FileResponse{
		Url:       filePath,
		MediaType: mediaType,
	})

}

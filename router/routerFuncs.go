package routeController

import (
	"context"
	"dash/database"
	"dash/helpers"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func ServeMpd(c *gin.Context) {
	c.File("./videos/2151765.mp4output.mpd")
}

func UploadVideo(c *gin.Context) {

	vdo, err := c.FormFile("video")
	title := c.PostForm("title")
	collection := database.Client.Database("videos").Collection("dash")

	if err != nil {
		println("error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	uniqueId := uuid.New()

	filePath := fmt.Sprintf("./videos/%s/%s", uniqueId.String(), vdo.Filename)

	errSave := c.SaveUploadedFile(vdo, filePath)

	if errSave != nil {
		println("error:", errSave.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	workDir := fmt.Sprintf("./videos/%s", uniqueId.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	var videoModel database.Video

	videoModel.ID = uniqueId.String()
	videoModel.Title = title
	videoModel.IsDash = false

	_, dbInsertErr := collection.InsertOne(ctx, videoModel)
	if dbInsertErr != nil {
		fmt.Println(dbInsertErr.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	go helpers.CreateSegment(vdo.Filename, "output.mpd", workDir, uniqueId.String())

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": http.StatusText(http.StatusOK),
	})

}

func GetMpdById(c *gin.Context) {
	id := c.Param("id")
	name := c.Param("name")
	var result database.Video

	collection := database.Client.Database("videos").Collection("dash")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	if !result.IsDash {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": http.StatusText(http.StatusNotFound),
		})
		return
	}

	destPath := fmt.Sprintf("./videos/%s/%s", id, name)

	c.File(destPath)
}

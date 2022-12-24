package artistHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lab1/controller/backend/query"
	"lab1/model"
	"net/http"
	"strconv"
)

func GetArtistByID(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}

	artistPage := *query.GetArtistByID(c, id)

	albumList := query.GetAlbumsByArtistId(c, id)

	c.HTML(http.StatusOK, "albums/artist.html",
		gin.H{
			"albums": albumList,
			"artist": artistPage,
			"item":   "artists",
			"ID":     id,
		})
	return
}

func GetEditArtist(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}

	artistPage := *query.GetArtistByID(c, id)

	c.HTML(http.StatusOK, "albums/newArtist.html",
		gin.H{
			"pageName":    "Добавить исполнителя",
			"answerClass": "",
			"answer":      "Успешно добавлено в базу данных",
			"artist":      artistPage,
			"error":       "",
		})
	return
}
func PostEditArtist(c *gin.Context) {
	PostArtist(c, true)

}

func DeleteArtist(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}
	var artist model.Artist

	db := c.Value("database").(*gorm.DB)
	if err = db.First(&artist, id).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	var albums []model.Album
	if err = db.Where("artist_id = ?", id).Find(&albums).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err = db.Unscoped().Delete(&albums).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err = db.Unscoped().Delete(&artist).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
}
func PostCreateArtist(c *gin.Context) {
	PostArtist(c, false)
}

func PostArtist(c *gin.Context, isEdit bool) {
	pageName := "Добавить исполнителя"
	formAction := "/artists/add"
	if isEdit {
		id := c.Param("id")
		pageName = "Изменить исполнителя"
		formAction = "/artists/" + id + "/edit"
	}
	successClass := "alert success"
	errorClass := "error"
	answerClass := successClass
	artist := &model.Artist{}
	if err := c.Bind(artist); err != nil {
		// Note: if there's a bind error, Gin will call
		// c.AbortWithError. We just need to return here.
		return
	}

	if artist.ArtistName == "" {
		answerClass = errorClass
		c.HTML(http.StatusBadRequest, "albums/newArtist.html",
			gin.H{
				"pageName":    pageName,
				"formAction":  formAction,
				"answerClass": answerClass,
				"answer":      "Ошибка валидации формы",
				"artist":      artist,
				"error":       "Заполните поле исполнитель",
			})
		return
	}
	db := c.Value("database").(*gorm.DB)
	if isEdit {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
			return
		}
		artist.ID = uint(id)
		if err = db.Save(&artist).Error; err != nil {
			answerClass = errorClass
			c.HTML(http.StatusInternalServerError, "albums/newArtist.html",
				gin.H{
					"pageName":    pageName,
					"answerClass": answerClass,
					"answer":      "Ошибка добавления в БД",
					"formAction":  formAction,
					"artist":      artist,
					"error":       err,
				})
			return
		}
	} else {
		if err := db.Create(&artist).Error; err != nil {
			answerClass = errorClass
			c.HTML(http.StatusInternalServerError, "albums/newArtist.html",
				gin.H{
					"pageName":    pageName,
					"answerClass": answerClass,
					"answer":      "Ошибка обновления данных в БД",
					"formAction":  formAction,
					"artist":      artist,
					"error":       err,
				})
			return
		}
	}
	c.HTML(http.StatusOK, "albums/newArtist.html",
		gin.H{
			"pageName":    pageName,
			"answerClass": answerClass,
			"answer":      "Успешно добавлено в базу данных",
			"formAction":  formAction,
			"artist":      artist,
			"error":       "",
		})
	return
}

func GetCreateArtist(c *gin.Context) {

	c.HTML(http.StatusOK, "albums/newArtist.html",
		gin.H{
			"formAction":  "/artists/add",
			"pageName":    "Добавить исполнителя",
			"answerClass": "",
		})
	return
}

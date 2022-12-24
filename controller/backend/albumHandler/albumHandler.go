package albumHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lab1/controller/backend/query"
	"lab1/model"
	"net/http"
	"strconv"
)

func GetAlbums(c *gin.Context) {
	albums := query.GetAllAlbums(c)
	c.HTML(http.StatusOK, "albums/index.html", gin.H{"albums": albums})

}

func GetAlbumByID(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}

	albumPage := *query.GetAlbumByID(c, id)
	songs := query.GetSongsByAlbumID(c, id)

	c.HTML(http.StatusOK, "albums/album.html",
		gin.H{
			"album": albumPage,
			"songs": songs,
			"item":  "albums",
			"ID":    id,
		})
	return
}

func GetEditAlbum(c *gin.Context) {
	pageName := "Изменить альбом"
	errorClass := "error"
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}

	album := *query.GetAlbumByID(c, id)
	var artists []model.Artist

	db := c.Value("database").(*gorm.DB)
	if err = db.Where("id <> ?", album.ArtistID).Find(&artists).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "albums/newAlbum.html",
			gin.H{
				"pageName":    pageName,
				"answerClass": errorClass,
				"answer":      "Нет исполнителей в БД",
				"error":       err,
				"album":       album,
			})
		return
	}
	var first model.Artist
	if err = db.Find(&first, album.ArtistID).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "albums/newAlbum.html",
			gin.H{
				"pageName":    "Добавление альбома",
				"answerClass": errorClass,
				"answer":      "Нет исполнителей в БД",
				"error":       err,
				"album":       album,
				"first":       first,
			})
		return
	}

	c.HTML(http.StatusOK, "albums/newAlbum.html",
		gin.H{
			"pageName": pageName,
			"artists":  artists,
			"error":    "",
			"album":    album,
			"first":    first,
		})
	return
}
func PostEditAlbum(c *gin.Context) {
	idString := c.Param("id")
	pageName := "Изменить исполнителя"
	formAction := "/albums/" + idString + "/edit"
	successClass := "alert success"
	errorClass := "error"
	answerClass := successClass
	album := &model.Album{}

	if err := c.Bind(album); err != nil {
		// Note: if there's a bind error, Gin will call
		// c.AbortWithError. We just need to return here.
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}
	album.ID = uint(id)

	first := *query.GetArtistByID(c, uint64(album.ArtistID))
	var artists []model.Artist
	db := c.Value("database").(*gorm.DB)
	if err := db.Where("id <> ?", album.ArtistID).Find(&artists).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "albums/newAlbum.html",
			gin.H{
				"pageName":    pageName,
				"answerClass": errorClass,
				"answer":      "Нет исполнителей в БД",
				"formAction":  formAction,
				"error":       err,
				"album":       album,
				"first":       first,
			})
		return
	}

	if err := db.Save(&album).Error; err != nil {
		answerClass = errorClass
		c.HTML(http.StatusInternalServerError, "albums/newAlbum.html",
			gin.H{
				"pageName":    pageName,
				"answerClass": answerClass,
				"answer":      "Ошибка обновления БД",
				"formAction":  formAction,
				"artist":      artists,
				"error":       err,
				"album":       album,
				"first":       first,
			})
		return
	}

	c.HTML(http.StatusOK, "albums/newAlbum.html",
		gin.H{
			"pageName":    pageName,
			"answerClass": answerClass,
			"answer":      "Успешно обновлено",
			"formAction":  formAction,
			"artists":     artists,
			"error":       "",
			"album":       album,
			"first":       first,
		})
	return

}
func DeleteAlbum(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}
	var album model.Album

	db := c.Value("database").(*gorm.DB)
	if err = db.First(&album, id).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if err = db.Unscoped().Delete(&album).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	c.Status(http.StatusOK)
}
func GetCreateAlbum(c *gin.Context) {
	errorClass := "error"

	var artists []model.Artist
	db := c.Value("database").(*gorm.DB)

	if err := db.Find(&artists).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "albums/newArtist.html",
			gin.H{
				"pageName":    "Добавление альбома",
				"answerClass": errorClass,
				"answer":      "Нет исполнителей в БД",
				"formAction":  "/albums/add",
				"error":       err,
			})
		return
	}

	c.HTML(http.StatusOK, "albums/newAlbum.html",
		gin.H{
			"artists":     artists,
			"formAction":  "/albums/add",
			"pageName":    "Добавить альбом",
			"answerClass": "",
		})
	return
}

func PostCreateAlbum(c *gin.Context) {
	pageName := "Добавить альбом"
	formAction := "/albums/add"
	successClass := "alert success"
	errorClass := "error"
	answerClass := successClass
	album := &model.Album{}

	if err := c.Bind(album); err != nil {
		// Note: if there's a bind error, Gin will call
		// c.AbortWithError. We just need to return here.
		return
	}
	db := c.Value("database").(*gorm.DB)
	var artists []model.Artist
	if err := db.Where("id <> ?", album.ArtistID).Find(&artists).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "albums/newAlbum.html",
			gin.H{
				"pageName":    "Добавление альбома",
				"answerClass": errorClass,
				"answer":      "Нет исполнителей в БД",
				"formAction":  "/albums/add",
				"error":       err,
				"album":       album,
			})
		return
	}
	var first model.Artist
	if err := db.Find(&first, album.ArtistID).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "albums/newAlbum.html",
			gin.H{
				"pageName":    "Добавление альбома",
				"answerClass": errorClass,
				"answer":      "Нет исполнителей в БД",
				"formAction":  "/albums/add",
				"error":       err,
				"album":       album,
				"first":       first,
			})
		return
	}
	// FIXME: There's a better way to do this validation!
	if album.Name == "" {
		answerClass = errorClass
		c.HTML(http.StatusBadRequest, "albums/newAlbum.html",
			gin.H{
				"pageName":    pageName,
				"formAction":  formAction,
				"answerClass": answerClass,
				"answer":      "Ошибка валидации формы",
				"artists":     artists,
				"error":       "Заполните поле исполнитель",
				"album":       album,
				"first":       first,
			})
		return
	}

	if err := db.Create(&album).Error; err != nil {
		answerClass = errorClass
		c.HTML(http.StatusInternalServerError, "albums/newAlbum.html",
			gin.H{
				"pageName":    pageName,
				"answerClass": answerClass,
				"answer":      "Ошибка добавления в БД",
				"formAction":  formAction,
				"artist":      artists,
				"error":       err,
				"album":       album,
				"first":       first,
			})
		return
	}

	c.HTML(http.StatusOK, "albums/newAlbum.html",
		gin.H{
			"pageName":    pageName,
			"answerClass": answerClass,
			"answer":      "Успешно добавлено в базу данных",
			"formAction":  formAction,
			"artists":     artists,
			"error":       "",
			"album":       album,
			"first":       first,
		})
	return
}

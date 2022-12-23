package artistHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lab1/model"
	"net/http"
	"strconv"
)

func GetArtistByID(c *gin.Context) {
	db := c.Value("database").(*gorm.DB)
	var artistPage model.ArtistPage

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}

	//Add 404 handling
	rez := db.Table("artists").
		Select("artists.id, artists.artist_name, artists.year").
		Where(&model.ArtistPage{ID: uint(id)}).Find(&artistPage)

	if rez.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "artist not found"})
		return
	}
	if rez.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return
	}
	var albumList []model.AlbumPage
	rez1 := db.Table("albums").
		Select("albums.id, albums.name, albums.rating, albums.image_path").
		Where(&model.Album{ArtistID: uint(id)}).Find(&albumList)
	if rez1.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "artist not found"})
		return
	}
	if rez1.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return
	}
	c.HTML(http.StatusOK, "albums/artist.html",
		gin.H{
			"albums": albumList,
			"artist": artistPage,
			"item":   "artists",
			"ID":     id,
		})
	return
}

package albumHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lab1/model"
	"net/http"
	"strconv"
)

func GetAlbums(c *gin.Context) {
	db := c.Value("database").(*gorm.DB)
	//var row []model.Album
	/*if err := db.Find(&alb).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}*/

	var albums []model.AlbumPage
	rez := db.Table("albums").
		Select("albums.id,albums.name, albums.rating, albums.image_path, albums.artist_id, artists.artist_name").
		Joins("left join artists on artists.id = albums.artist_id").
		Scan(&albums)

	if rez.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	if rez.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return
	}
	c.HTML(http.StatusOK, "albums/index.html", gin.H{"albums": albums})

}

func GetAlbumByID(c *gin.Context) {
	db := c.Value("database").(*gorm.DB)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}

	var albumPage model.AlbumPage
	rez := db.Table("albums").
		Select("albums.artist_id, albums.name, albums.rating, albums.image_path, artists.artist_name").
		Joins("left join artists on artists.id = albums.artist_id").
		Where(&model.Album{ID: uint(id)}).Find(&albumPage)
	if rez.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	if rez.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return
	}
	var songs []model.Song
	rez1 := db.Table("songs").
		Select("songs.song_name, songs.length, songs.number").
		Where(&model.Song{AlbumID: uint(id)}).Find(&songs)
	if rez1.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "artist not found"})
		return
	}
	if rez1.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return
	}

	c.HTML(http.StatusOK, "albums/album.html",
		gin.H{
			"album": albumPage,
			"songs": songs,
			"item":  "albums",
			"ID":    id,
		})
	return
}

func EditAlbum(c *gin.Context) {
	db := c.Value("database").(*gorm.DB)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Can't parse ID"})
		return
	}

	var albumPage model.AlbumPage
	rez := db.Table("albums").
		Select("albums.artist_id, albums.name, albums.rating, albums.image_path, artists.artist_name").
		Joins("left join artists on artists.id = albums.artist_id").
		Where(&model.Album{ID: uint(id)}).Find(&albumPage)
	if rez.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	if rez.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return
	}
	var songs []model.Song
	rez1 := db.Table("songs").
		Select("songs.song_name, songs.length, songs.number").
		Where(&model.Song{AlbumID: uint(id)}).Find(&songs)
	if rez1.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "artist not found"})
		return
	}
	if rez1.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return
	}

	c.HTML(http.StatusOK, "albums/album.html",
		gin.H{
			"album": albumPage,
			"songs": songs,
			"item":  "albums",
			"ID":    id,
		})
	return
}

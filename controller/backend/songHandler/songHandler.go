package songHandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lab1/model"
	"net/http"
	"strconv"
)

func GetEditSong(c *gin.Context) {
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
func PostEditSong(c *gin.Context) {}

func DeleteSong(c *gin.Context) {}

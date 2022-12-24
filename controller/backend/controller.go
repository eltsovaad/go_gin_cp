package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lab1/controller/backend/albumHandler"
	"lab1/controller/backend/artistHandler"
	"lab1/model"
	"net/http"
)

func SetupDatabase(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Artist{},
		&model.Album{},
		&model.Song{},
	)
	if err != nil {
		return fmt.Errorf("Error migrating database: %s", err)
	}

	//query.InsertIntoDb(db)

	return nil
}

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.Static("/static", "./static/")
	r.LoadHTMLGlob("templates/**/*.html")
	r.Use(connectDatabase(db))

	r.GET("/albums_list", albumHandler.GetAlbums)
	r.GET("/albums/:id", albumHandler.GetAlbumByID)
	r.GET("/albums/:id/edit", albumHandler.GetEditAlbum)
	r.GET("/albums/add", albumHandler.GetCreateAlbum)
	r.GET("/albums/:id/delete", albumHandler.DeleteAlbum)
	r.POST("/albums/:id/edit", albumHandler.PostEditAlbum)
	r.POST("/albums/add", albumHandler.PostCreateAlbum)

	r.GET("/artists/:id", artistHandler.GetArtistByID)
	r.GET("/artists/:id/edit", artistHandler.GetEditArtist)
	r.GET("/artists/add", artistHandler.GetCreateArtist)
	r.GET("/artists/:id/delete", artistHandler.DeleteArtist)
	r.POST("/artists/:id/edit", artistHandler.PostEditArtist)
	r.POST("/artists/add", artistHandler.PostCreateArtist)

	r.GET("/items/new", itemNewGetHandler)
	r.POST("/albums/new", albumNewPostHandler)
	r.GET("/exp", exp)
	//r.POST("/artists/new", artistNewPostHandler)
	//r.POST("/songs/new", songNewPostHandler)
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/albums_list/")
	})
}

// Middleware to connect the database for each request that uses this
// middleware.
func connectDatabase(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("database", db)
	}
}

func exp(c *gin.Context) {

	db := c.Value("database").(*gorm.DB)

	id := 1

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
		Select("songs.song_name, songs.length, songs.number, songs.id").
		Where(&model.Song{AlbumID: uint(id)}).Find(&songs)
	if rez1.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "artist not found"})
		return
	}
	if rez1.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return
	}

	c.HTML(http.StatusOK, "albums/exp.html",
		gin.H{
			"album": albumPage,
			"songs": songs,
		})
	return
}

func itemNewGetHandler(c *gin.Context) {
	album := &model.AlbumOld{}
	if err := c.Bind(album); err != nil {
		// Note: if there's a bind error, Gin will call
		// c.AbortWithError. We just need to return here.
		return
	}
	// FIXME: There's a better way to do this validation!
	if album.Title == "" || album.Artist == "" || album.Review == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	db := c.Value("database").(*gorm.DB)
	if err := db.Create(&album).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Redirect(http.StatusFound, "/welcome/")
}

func albumNewGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "albums/newArtist.html", gin.H{})
}

// postAlbums godoc
// @Summary      Adds album to the DB
// @Description  post json with album
// @Accept       json
// @Produce      json
// @Param        title body string true "Title of an album"
// @Param        artist body string true "Artist"
// @Param        review body float32 true "My review mark"
// @Success      200  {object}  AlbumOld
// @Failure      404  {int}  http.StatusNotFound
// @Failure      500  {int}  http.StatusInternalServerError
// @Router       /albums [post]
func albumNewPostHandler(c *gin.Context) {
	var newAlbum model.AlbumOld
	fmt.Printf("\nWe are in post Alb\n")
	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.Bind(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "can't bind"})
		return
	}
	fmt.Printf("\nAlb:%v\n", newAlbum)
	db := c.Value("database").(*gorm.DB)
	if err := db.Create(&newAlbum).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

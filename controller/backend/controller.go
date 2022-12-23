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

	/*var artists = []model.Artist{
		{ID: 1, ArtistName: "Red Hot Chili Peppers", Year: 1983},
		{ID: 2, ArtistName: "Arctic Monkeys", Year: 2002},
		{ID: 3, ArtistName: "Sting", Year: 1971}}
	db.Create(&artists)

	var albums = []model.Album{
		{ID: 1, Name: "Brand New Day", Rating: 4.5, ImagePath: "images/Sting_Brand_New_Day_album_art.jpg", ArtistID: 3},
		{ID: 2, Name: "Nothing Like the Sun", Rating: 4.7, ImagePath: "images/NLTS.jpg", ArtistID: 3},
		{ID: 3, Name: "AM", Rating: 4, ImagePath: "images/Arctic_Monkeys_AM_cover.jpg", ArtistID: 2},
		{ID: 4, Name: "Tranquility Base Hotel & Casino", Rating: 5, ImagePath: "images/Tranquility_Base_Hotel_&_Casino.jpg", ArtistID: 2},
		{ID: 5, Name: "Californication", Rating: 4.8, ImagePath: "images/Californication.jpg", ArtistID: 1},
		{ID: 6, Name: "Stadium Arcadium", Rating: 4.5, ImagePath: "images/Stadiumarcadium.jpg", ArtistID: 1},
		{ID: 7, Name: "The Getaway", Rating: 5, ImagePath: "images/Thegetawayalbum.jpg", ArtistID: 1},
	}
	db.Create(&albums)
	var songs = []model.Song{
		{ID: 1, Name: "Desert Rose", Length: "4m45s", Number: 2, AlbumID: 1},
		{ID: 2, Name: "Englishman In New York", Length: "4m45s", Number: 3, AlbumID: 2},
		{ID: 3, Name: "Arabella", Length: "3m27s", Number: 4, AlbumID: 3},
		{ID: 4, Name: "Knee Socks", Length: "4m17s", Number: 11, AlbumID: 3},
		{ID: 5, Name: "Four Out Of Five", Length: "5m12s", Number: 6, AlbumID: 4},
		{ID: 6, Name: "Californication", Length: "5m21s", Number: 6, AlbumID: 5},
		{ID: 7, Name: "Stadium Arcadium", Length: "5m15s", Number: 4, AlbumID: 6},
		{ID: 8, Name: "Go Robot", Length: "4m24s", Number: 7, AlbumID: 7},
	}
	db.Create(&songs)*/

	return nil
}

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.Static("/static", "./static/")
	r.LoadHTMLGlob("templates/**/*.html")
	r.Use(connectDatabase(db))

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/albums_list", albumHandler.GetAlbums)
	r.GET("/albums/:id", albumHandler.GetAlbumByID)
	r.GET("/albums/:id/edit", albumHandler.EditAlbum)
	r.GET("/albums/:id/delete", albumHandler.GetAlbumByID)
	r.GET("/artists/:id", artistHandler.GetArtistByID)
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
	c.HTML(http.StatusOK, "albums/new.html", gin.H{})
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

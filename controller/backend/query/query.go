package query

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"lab1/model"
	"net/http"
)

func GetAlbumByID(c *gin.Context, id uint64) *model.AlbumPage {
	db := c.Value("database").(*gorm.DB)

	var albumPage model.AlbumPage
	rez := db.Table("albums").
		Select("albums.id, albums.artist_id, albums.name, albums.rating, albums.image_path, artists.artist_name").
		Joins("left join artists on artists.id = albums.artist_id").
		Where(&model.Album{ID: uint(id)}).Find(&albumPage)
	if rez.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return &albumPage
	}
	if rez.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return &albumPage
	}

	return &albumPage
}

func GetSongsByAlbumID(c *gin.Context, id uint64) []model.Song {
	db := c.Value("database").(*gorm.DB)
	var songs []model.Song
	rez1 := db.Table("songs").
		Select("songs.song_name, songs.length, songs.number, songs.id, songs.album_id").
		Where(&model.Song{AlbumID: uint(id)}).Find(&songs)
	if rez1.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return nil
	}
	return songs
}

func GetAllAlbums(c *gin.Context) []model.AlbumPage {
	db := c.Value("database").(*gorm.DB)
	var albums []model.AlbumPage
	rez := db.Table("albums").
		Select("albums.id,albums.name, albums.rating, albums.image_path, albums.artist_id, artists.artist_name").
		Joins("left join artists on artists.id = albums.artist_id").
		Scan(&albums)

	if rez.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return nil
	}
	if rez.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return nil
	}
	return albums
}

func GetAlbumsByArtistId(c *gin.Context, id uint64) []model.AlbumPage {
	db := c.Value("database").(*gorm.DB)
	var albumList []model.AlbumPage
	rez1 := db.Table("albums").
		Select("albums.id, albums.name, albums.rating, albums.image_path").
		Where(&model.Album{ArtistID: uint(id)}).Find(&albumList)
	if rez1.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "artist not found"})
		return nil
	}
	if rez1.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred"})
		return nil
	}
	return albumList
}

func GetArtistByID(c *gin.Context, id uint64) *model.ArtistPage {
	db := c.Value("database").(*gorm.DB)
	var artistPage model.ArtistPage

	//Add 404 handling
	rez := db.Table("artists").
		Select("artists.id, artists.artist_name, artists.year").
		Where(&model.ArtistPage{ID: uint(id)}).Find(&artistPage)

	if rez.RowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "artist not found"})
		return nil
	}
	if rez.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error occurred", "error": rez.Error})
		return nil
	}
	return &artistPage
}

func AddArtist(c *gin.Context, name string, year uint) error {
	db := c.Value("database").(*gorm.DB)
	var artist = model.Artist{ArtistName: name, Year: year}
	rez := db.Create(&artist)
	if rez.Error != nil {
		return rez.Error
	}
	return nil
}

func InsertIntoDb(db *gorm.DB) {
	var artists = []model.Artist{
		{ID: 1, ArtistName: "Red Hot Chili Peppers", Year: 1983},
		{ID: 2, ArtistName: "Arctic Monkeys", Year: 2002},
		{ID: 3, ArtistName: "Sting", Year: 1971}}
	db.Create(&artists)

	var albums = []model.Album{
		{ID: 1, Name: "Brand New Day", Rating: 4.5, ImagePath: "../../static/images/Sting_Brand_New_Day_album_art.jpg", ArtistID: 3},
		{ID: 2, Name: "Nothing Like the Sun", Rating: 4.7, ImagePath: "../../static/images/NLTS.jpg", ArtistID: 3},
		{ID: 3, Name: "AM", Rating: 4, ImagePath: "../../static/images/Arctic_Monkeys_AM_cover.jpg", ArtistID: 2},
		{ID: 4, Name: "Tranquility Base Hotel & Casino", Rating: 5, ImagePath: "../../static/images/Tranquility_Base_Hotel_&_Casino.jpg", ArtistID: 2},
		{ID: 5, Name: "Californication", Rating: 4.8, ImagePath: "../../static/images/Californication.jpg", ArtistID: 1},
		{ID: 6, Name: "Stadium Arcadium", Rating: 4.5, ImagePath: "../../static/images/Stadiumarcadium.jpg", ArtistID: 1},
		{ID: 7, Name: "The Getaway", Rating: 5, ImagePath: "../../static/images/Thegetawayalbum.jpg", ArtistID: 1},
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
	db.Create(&songs)
}

package model

import (
	"gorm.io/gorm"
)

type AlbumOld struct {
	//ID in data base
	ID uint64 `swaggerignore:"true"`
	//Title of an album
	Title string `form:"title" json:"title"`
	//Artist (band's) name
	Artist string `form:"artist" json:"artist"`
	//My review for the album
	Review float32 `form:"review" json:"review"`
}

type Artist struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey;unique;autoIncrement"`
	ArtistName string `gorm:"unique"`
	Year       uint
	Albums     []Album
}

type Album struct {
	gorm.Model
	ID        uint `gorm:"primaryKey;unique;autoIncrement"`
	Name      string
	Rating    float32 `gorm:"scale:2"`
	ImagePath string
	//Year     time.Time
	ArtistID uint
	Songs    []Song
}

type Song struct {
	gorm.Model
	ID     uint   `gorm:"primaryKey;unique;autoIncrement"`
	Name   string `gorm:"column:song_name"`
	Length string
	Number uint
	//Mark 	 bool
	AlbumID uint
}

type AlbumPage struct {
	ID         uint
	Name       string
	Rating     float32
	ImagePath  string
	ArtistName string
	ArtistID   uint
}

type ArtistPage struct {
	ID         uint
	ArtistName string
	Year       uint
}

func (artist *Artist) AfterDelete(tx *gorm.DB) error {
	return tx.Model(&Album{}).Unscoped().Where("artist_id = ?", artist.ID).Delete(&Album{}).Error
}

func (album *Album) AfterDelete(tx *gorm.DB) error {
	return tx.Model(&Song{}).Unscoped().Where("album_id = ?", album.ID).Delete(&Song{}).Error
}

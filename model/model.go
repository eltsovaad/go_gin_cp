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
	ID         uint
	ArtistName string
	Year       uint
	Albums     []Album
}

type Album struct {
	gorm.Model
	ID        uint
	Name      string
	Rating    float32 `gorm:"scale:2"`
	ImagePath string
	//Year     time.Time
	ArtistID uint
	Songs    []Song
}

type Song struct {
	gorm.Model
	ID     int
	Name   string `gorm:"column:song_name"`
	Length string
	Number int
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
	Year       int
}

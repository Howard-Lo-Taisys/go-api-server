package models

import (
	"time"
)

type Service struct {
	ServiceId   string    `gorm:"size:255;primary_key;not null;unique" json:"serviceId"`
	ServiceName string    `gorm:"size:255;not null;unique" json:"serviceName"`
	Version     []Version `gorm:"foreignKey:ServiceId"`
}

type Version struct {
	VersionId  string    `gorm:"size:255;primary_key;not null;unique" json:"versionId"`
	ServiceId  string    `gorm:"size:255;not null;" json:"serviceId"`
	Version    string    `gorm:"size:255;not null" json:"version"`
	GitHash    string    `gorm:"size:255;not null" json:"gitHash"`
	CreateTime time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createTime"`
}

type CreateRequest struct {
	ServiceId string `gorm:"size:255;not null" json:"serviceId"`
	VersionId string `gorm:"size:255;not null" json:"versionId"`
	Version   string `gorm:"size:255;not null" json:"version"`
	GitHash   string `gorm:"size:255;not null" json:"gitHash"`
}

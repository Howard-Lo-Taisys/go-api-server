package controller

import (
	"go-api-server/config"
	"go-api-server/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetHealthz(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func GetServerName(c *gin.Context) {
	versions := models.Version{}

	subq := config.DB.Model(&versions).
		Select("MAX(versions.create_time) as create_time").
		Where("service_name = ?", c.Param("service-name")).
		Joins("LEFT JOIN services ON services.service_id = versions.service_id").
		Group("versions.service_id").Find(&versions)

	if err := config.DB.Model(&versions).
		Where("versions.create_time = ?", versions.CreateTime).
		Joins("LEFT JOIN (?) q ON versions.create_time = q.create_time", subq).
		Find(&versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &versions)
}

func GetServerVersion(c *gin.Context) {
	versions := models.Version{}

	if err := config.DB.Model(&versions).
		Joins("LEFT JOIN services ON services.service_id = versions.service_id").
		Where("service_name = ? AND version_number = ?", c.Param("service-name"), c.Param("version-number")).
		First(&versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &versions)
}

func PostServerName(c *gin.Context) {
	input := models.CreateRequest{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var services models.Service

	services.ServiceId = uuid.NewString()
	services.ServiceName = c.Param("service-name")

	services.Version = []models.Version{{
		VersionId:     uuid.NewString(),
		VersionNumber: input.VersionNumber,
		Version:       input.Version,
		GitHash:       input.GitHash,
		CreateTime:    time.Now()}}

	if err := config.DB.Create(&services).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"create service error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ServeiceId": &services.ServiceId, "ServiceName": &services.ServiceName,
		"Version": &services.Version})
}

func PutServerVersion(c *gin.Context) {
	input := models.CreateRequest{}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var services models.Service
	var version models.Version

	if err := config.DB.Model(&services).Where("service_name = ?", c.Param("service-name")).
		Find(&services).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&version).Where("version_number = ? AND service_id = ?", c.Param("version-number"), services.ServiceId).Find(&version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if version.VersionNumber == input.VersionNumber && version.ServiceId == input.ServiceId && services.ServiceId != "" {
		version.VersionId = uuid.NewString()
		version.VersionNumber = c.Param("version-number")
		version.ServiceId = services.ServiceId
		version.Version = input.Version
		version.GitHash = input.GitHash
		version.CreateTime = time.Now()
		config.DB.Model(&version).Create(&version)
		c.JSON(http.StatusOK, &version)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Input data error"})
	}
}

func DeleteServerVersion(c *gin.Context) {
	var versions models.Version

	if err := config.DB.Model(&versions).Joins("LEFT JOIN services ON services.service_id = versions.service_id").
		Where("service_name = ? AND version_number = ?", c.Param("service-name"), c.Param("version-number")).
		First(&versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Delete(&versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

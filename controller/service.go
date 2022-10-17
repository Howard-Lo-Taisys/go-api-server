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

	if err := config.DB.Joins("JOIN services ON services.service_id = versions.service_id").
		Where("service_name = ?", c.Param("service-name")).Order("create_time desc").
		First(&versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, &versions)
}

func GetServerVersion(c *gin.Context) {
	versions := models.Version{}

	if err := config.DB.Joins("JOIN services ON services.service_id = versions.service_id").
		Where("service_name = ? AND version_id = ?", c.Param("service-name"), c.Param("version-id")).
		First(&versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		VersionId:  input.VersionId,
		Version:    input.Version,
		GitHash:    input.GitHash,
		CreateTime: time.Now()}}

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

	if err := config.DB.Where("service_name = ?", c.Param("service-name")).
		Find(&services).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	version.VersionId = c.Param("version-id")
	version.ServiceId = services.ServiceId
	version.Version = input.Version
	version.GitHash = input.GitHash
	version.CreateTime = time.Now()

	if err := config.DB.Create(&version).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &version)
}

func DeleteServerVersion(c *gin.Context) {
	var versions models.Version

	if err := config.DB.Joins("JOIN services ON services.service_id = versions.service_id").
		Where("service_name = ? AND version_id = ?", c.Param("service-name"), c.Param("version-id")).
		First(&versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := config.DB.Delete(&versions).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

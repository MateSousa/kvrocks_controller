package handlers

import (
	"net/http"

	"github.com/KvrocksLabs/kvrocks-controller/consts"
	"github.com/KvrocksLabs/kvrocks-controller/metadata"
	"github.com/KvrocksLabs/kvrocks-controller/storage/memory"
	"github.com/gin-gonic/gin"
)

func ListNamespace(c *gin.Context) {
	storage := c.MustGet(consts.ContextKeyStorage).(*memory.MemStorage)
	namespaces, err := storage.ListNamespace()
	if err != nil {
		if metaErr, ok := err.(*metadata.Error); ok && metaErr.Code == metadata.CodeNoExists {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"namespaces": namespaces})
}

func CreateNamespace(c *gin.Context) {
	storage := c.MustGet(consts.ContextKeyStorage).(*memory.MemStorage)
	var param struct {
		Namespace string `json:"namespace"`
	}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	if len(param.Namespace) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "namespace should NOT be empty"})
		return
	}

	if err := storage.CreateNamespace(param.Namespace); err != nil {
		if metaErr, ok := err.(*metadata.Error); ok && metaErr.Code == metadata.CodeExisted {
			c.JSON(http.StatusConflict, gin.H{"err": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func RemoveNamespace(c *gin.Context) {
	storage := c.MustGet(consts.ContextKeyStorage).(*memory.MemStorage)
	namespace := c.Param("namespace")
	if err := storage.CreateNamespace(namespace); err != nil {
		if metaErr, ok := err.(*metadata.Error); ok && metaErr.Code == metadata.CodeNoExists {
			c.JSON(http.StatusNotFound, gin.H{"err": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
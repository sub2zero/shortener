package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	store "shortener/internal/struct"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

type serverAddress string

const serverAdd serverAddress = "127.0.0.1"

func main() {
	// Create Gin router
	r := gin.Default()
	r.Use(gin.Logger())
	// Instantiate recipe Handler and provide a data store
	r.Delims("{[{", "}]}")
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLFiles("./templates/index.tmpl")

	store := store.NewUrlStore()
	UrlHandler := NewUrlHandler(store)

	// Register the routes
	r.GET("/", homePage)
	r.GET("/:id", UrlHandler.RedirectUrl)
	r.GET("/url", UrlHandler.ListUrls)
	r.POST("/url", UrlHandler.CreateUrl)
	r.GET("/url/:id", UrlHandler.GetUrl)
	r.PUT("/url/:id", UrlHandler.UpdateUrl)
	r.DELETE("/url/:id", UrlHandler.DeleteUrl)

	// Start the server
	r.Run()
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func homePage(c *gin.Context) {
	// c.String(http.StatusOK, "This is my home page")
	// Serve the HTML form

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"now": time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	})
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 7

	rand.Seed(uint64(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortKey)
}

type urlstore interface {
	Add(name string, shortUrl store.ShortUrls) error
	Get(name string) (store.ShortUrls, error)
	List() (map[string]store.ShortUrls, error)
	Update(name string, shortUrl store.ShortUrls) error
	Remove(name string) error
}
type UrlHandler struct {
	store urlstore
}

func NewUrlHandler(s urlstore) *UrlHandler {
	return &UrlHandler{
		store: s,
	}
}
func (h UrlHandler) RedirectUrl(c *gin.Context) {
	id := c.Param("id")
	log.Println("getting by short id: ", id)
	url, err := h.store.Get(id)

	log.Println("redirecting url by id: ", url.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.Redirect(http.StatusMovedPermanently, url.Full)
}

func (h UrlHandler) CreateUrl(c *gin.Context) {
	var url store.ShortUrls

	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create a url friendly name
	// id := GenerateUUID()
	id := generateShortKey()
	// id := slug.Make(url.Full)
	url.Id = id
	log.Println(`saving url by id: `, url.Id)
	// add to the store
	h.store.Add(id, url)

	// return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h UrlHandler) ListUrls(c *gin.Context) {
	r, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(200, r)
}
func (h UrlHandler) GetUrl(c *gin.Context) {
	id := c.Param("id")
	println("getting by id: %s", id)
	url, err := h.store.Get(id)

	log.Println("getting url by id: ", url.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	c.JSON(200, url)
}
func (h UrlHandler) UpdateUrl(c *gin.Context) {
	// Get request body and convert it to recipes.Recipe
	var url store.ShortUrls
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	url.Id = id
	log.Println("getting by id: ", id)
	log.Println("getting by id: ", url.Id)
	err := h.store.Update(id, url)
	if err != nil {
		if err == store.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h UrlHandler) DeleteUrl(c *gin.Context) {
	id := c.Param("id")
	log.Println("getting to remove by id: ", id)
	err := h.store.Remove(id)
	if err != nil {
		if err == store.NotFoundErr {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return success payload
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

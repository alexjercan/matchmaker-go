package main

import (
	"log/slog"
	"net/http"
	"rps-multiplayer/pkg/config"
	"rps-multiplayer/pkg/database"
	"rps-multiplayer/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type HandlerV1 struct {
	db *bun.DB
}

func (this HandlerV1) CreateRoom(c *gin.Context) {
	dto := models.RoomDTO{MaxPlayers: 2, Private: false}
	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	slog.Info("DTO is {}", dto)

	// spin up game server => Address
	address := ""

	// create Room model => Generate a Code for the room
	code := "abcd"

	// save Room model in database => Fill in all fields that remain
	room := models.Room{Code: code, Address: address, Name: dto.Name, MaxPlayers: dto.MaxPlayers, Private: dto.Private}

	_, err := this.db.NewInsert().Model(&room).Exec(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, room)
}

func (this HandlerV1) ListRooms(c *gin.Context) {
    rooms := []models.Room{}

    err := this.db.NewSelect().Model(&rooms).Where("? = ?", bun.Ident("private"), "false").Scan(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (this HandlerV1) GetRoom(c *gin.Context) {
    code := c.Param("code")
    room := models.Room{}

    err := this.db.NewSelect().Model(&room).Where("? = ?", bun.Ident("code"), code).Limit(1).Scan(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, room)
}

func main() {
    cfg := config.LoadConfig()
    slog.Info("The config is {}", cfg)

    db := database.New(cfg)
	router := gin.Default()

	handler := HandlerV1{
		db,
	}

	apiV1 := router.Group("/api/v1")

	apiV1.GET("/rooms", handler.ListRooms)
	apiV1.POST("/rooms", handler.CreateRoom)
	apiV1.GET("/rooms/:code", handler.GetRoom)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

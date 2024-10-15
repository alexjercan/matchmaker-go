package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type RoomDTO struct {
	Name       string `json:"name" binding:"required"`
    MaxPlayers int    `json:"maxPlayers"`
	Private    bool   `json:"private"`
}

type Room struct {
	bun.BaseModel `bun:"table:rooms,alias:r"`

	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	Code       string    `bun:"code,type:varchar(6),notnull" json:"code"`
	Address    string    `bun:"address,type:varchar(128),notnull" json:"address"`
	Name       string    `bun:"name,type:varchar(128),notnull" json:"name"`
	MaxPlayers int       `bun:"maxPlayers,type:int,nullzero,notnull,default:2" json:"maxPlayers"`
	Private    bool      `bun:"private,notnull,default:false"`
	CreatedAt  time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `bun:"updatedAt,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

type HandlerV1 struct {
	db *bun.DB
}

func (this HandlerV1) CreateRoom(c *gin.Context) {
	dto := RoomDTO{MaxPlayers: 2, Private: false}
	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

    slog.Info("DTO is {}", dto)

	room := Room{Name: dto.Name, MaxPlayers: dto.MaxPlayers, Private: dto.Private}

	// spin up game server => Address

	// create Room model => Generate a Code for the room

	// save Room model in database => Fill in all fields that remain
	_, err := this.db.NewInsert().Model(&room).Exec(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, room)
}

func (this HandlerV1) ListRooms(c *gin.Context) {
	c.JSON(http.StatusOK, "List Rooms")
}

func (this HandlerV1) JoinRoom(c *gin.Context) {
	c.JSON(http.StatusOK, "Room Created")
}

func main() {
	router := gin.Default()

	handler := HandlerV1{
        db: nil,
	}

	apiV1 := router.Group("/api/v1")

	apiV1.GET("/rooms", handler.ListRooms)
	apiV1.POST("/rooms", handler.CreateRoom)
	apiV1.GET("/join/:code", handler.JoinRoom)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"matchmaker"
	"math/rand"
	"net"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type GameServer interface {
	Spawn(code string, maxPlayers int, name string) (string, int, int, error)
    Destroy(code string) error
}

type gameServerDocker struct {
}

func (this gameServerDocker) Spawn(code string, maxPlayers int, name string) (address string, query int, game int, err error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	query = listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	listener, err = net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	game = listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	address = "0.0.0.0"

	cmd := exec.Command(
		"docker", "run", "-d",
		"-e", fmt.Sprintf("SERVER_MAX_PLAYERS=%d", maxPlayers),
		"-p", fmt.Sprintf("%d:8080", query),
		"-p", fmt.Sprintf("%d:6969", game),
		"--name", code,
        fmt.Sprintf("%s", name),
	)
	slog.Info("Running command: {}", cmd)
	if err = cmd.Run(); err != nil {
		return
	}

	return
}

func (this gameServerDocker)Destroy(code string) (err error) {
	cmd := exec.Command("docker", "kill", code)
	slog.Info("Running command: {}", cmd)
	if err = cmd.Run(); err != nil {
		return
	}

	return
}

func NewGameServer(cfg matchmaker.Config) GameServer {
	return gameServerDocker{}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateCode(n int) string {
	bs := make([]rune, n)
	for i := range bs {
		bs[i] = letters[rand.Intn(len(letters))]
	}
	return string(bs)
}

type HandlerV1 struct {
    cfg        matchmaker.Config
	db         *bun.DB
	gameServer GameServer
}

func (this HandlerV1) CreateRoom(c *gin.Context) {
	dto := matchmaker.RoomDTO{MaxPlayers: 2, Private: false}
	if err := c.ShouldBind(&dto); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	slog.Info("DTO is {}", dto)

	code := generateCode(6)

	address, queryPort, gamePort, err := this.gameServer.Spawn(code, dto.MaxPlayers, this.cfg.Game.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	room := matchmaker.Room{
		Code:       code,
		Address:    address,
		QueryPort:  queryPort,
		GamePort:   gamePort,
		Name:       dto.Name,
		MaxPlayers: dto.MaxPlayers,
		Private:    dto.Private,
	}

	_, err = this.db.NewInsert().Model(&room).Exec(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, room)
}

func (this HandlerV1) ListRooms(c *gin.Context) {
	rooms := []matchmaker.Room{}

	err := this.db.NewSelect().Model(&rooms).Where("? = ?", bun.Ident("private"), "false").Scan(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (this HandlerV1) GetRoom(c *gin.Context) {
	code := c.Param("code")
	room := matchmaker.Room{}

	err := this.db.NewSelect().Model(&room).Where("? = ?", bun.Ident("code"), code).Limit(1).Scan(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, room)
}

func mainHTTP(cfg matchmaker.Config, db *bun.DB, server GameServer) {
	handler := HandlerV1{
        cfg,
		db,
		server,
	}

	router := gin.Default()

	apiV1 := router.Group("/api/v1")

	apiV1.GET("/rooms", handler.ListRooms)
	apiV1.POST("/rooms", handler.CreateRoom)
	apiV1.GET("/rooms/:code", handler.GetRoom)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func healthCheck(db *bun.DB, server GameServer) (error) {
    c := context.Background()
	rooms := []matchmaker.Room{}

    err := db.NewSelect().Model(&rooms).Scan(c)
	if err != nil {
		return err
	}

    for _, room := range rooms {
        status, err := room.GetStatus()
        if err != nil {
            return err
        }

        if status.Players <= 0 {
            slog.Info("Found an empty room ", room.Code, " will attempt to destroy it")
            server.Destroy(room.Code)
            _, err := db.NewDelete().Model(&room).WherePK().Exec(c)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func main() {
	cfg := matchmaker.LoadConfig()
	slog.Info("The config is {}", cfg)

	db := matchmaker.NewDB(cfg)
	server := NewGameServer(cfg)

    go mainHTTP(cfg, db, server)

	for {
        err := healthCheck(db, server)
        if err != nil {
            slog.Error("Error while performing health check {}", err)
        }
        time.Sleep(60 * time.Second)
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"golang.org/x/net/netutil"
)

type Config struct {
	Server struct {
		Host      string `env:"SERVER_HOST" env-default:"0.0.0.0"`
		QueryPort int    `env:"SERVER_QUERY_PORT" env-default:"8080"`
		GamePort  int    `env:"SERVER_GAME_PORT" env-default:"6969"`
		Code      string `env:"SERVER_CODE"`
		MaxPlayer int    `env:"SERVER_MAX_PLAYERS"`
	}
}

func LoadConfig() Config {
	cfg := Config{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal("Error read env variables")
	}

	return cfg
}

type GameClient struct {
	server *GameServer
	conn   net.Conn
}

func (this GameClient) Read(b []byte) (n int, err error) {
	return this.conn.Read(b)
}

func (this GameClient) Write(b []byte) (n int, err error) {
	return this.conn.Write(b)
}

func (this GameClient) Close() error {
	delete(this.server.Players, this)
	return this.conn.Close()
}

type GameServer struct {
	Players    map[GameClient]bool
	MaxPlayers int
	listener   net.Listener
}

func NewGameServer(address string, port int, maxPlayers int) (s GameServer, err error) {
	address = fmt.Sprintf("%s:%d", address, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	listener = netutil.LimitListener(listener, maxPlayers)
	return GameServer{Players: map[GameClient]bool{}, MaxPlayers: maxPlayers, listener: listener}, nil
}

func (this GameServer) Accept() (c GameClient, err error) {
	conn, err := this.listener.Accept()
	if err != nil {
		return
	}

	c = GameClient{conn: conn, server: &this}
	this.Players[c] = true

	return
}

func (this GameServer) Write(b []byte) (n int, err error) {
	for c := range this.Players {
		n, err = c.Write(b)
		if err != nil {
			return
		}
	}

	return
}

func (this GameServer) Close() error {
	for c := range this.Players {
		c.Close()
	}
	return this.listener.Close()
}

func handleConnection(conn GameClient) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn.server)

	for {
		bytes, err := reader.ReadBytes(byte('\n'))
		if err != nil {
			if err == io.EOF {
				return
			}
			slog.Error("Failed to read: {err}", err)
			continue
		}

		slog.Info("Read bytes: {bytes}", bytes)

		_, err = writer.Write(bytes)
		if err != nil {
			slog.Error("Failed to write: {err}", err)
			continue
		}

		err = writer.Flush()
		if err != nil {
			slog.Error("Failed to send data: {err}", err)
			continue
		}
	}
}

type HandlerV1 struct {
	gameServer GameServer
}

func NewHandler(gameServer GameServer) HandlerV1 {
	return HandlerV1{gameServer: gameServer}
}

type StatusResponse struct {
	Players int `json:"players"`
}

func (this HandlerV1) Status(c *gin.Context) {
	players := len(this.gameServer.Players)

	c.JSON(http.StatusOK, StatusResponse{Players: players})
}

func mainHTTP(cfg Config, server GameServer) {
	handler := NewHandler(server)

	router := gin.Default()

	apiV1 := router.Group("/api/v1")

	apiV1.GET("/status", handler.Status)

	router.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.QueryPort))
}

func main() {
	cfg := LoadConfig()
	slog.Info("The config is {cfg}", cfg)

	server, err := NewGameServer(cfg.Server.Host, cfg.Server.GamePort, cfg.Server.MaxPlayer)
	if err != nil {
		log.Fatal("Error creating game server: {err}", err)
	}
	defer server.Close()

	go mainHTTP(cfg, server)

	for {
		client, err := server.Accept()
		if err != nil {
			slog.Error("Failed to accept connection: {err}", err)
			continue
		}

		slog.Info("Accepted connection: {conn}", client)

		go handleConnection(client)
	}
}

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Host      string `env:"SERVER_HOST" env-default:"0.0.0.0"`
		QueryPort int    `env:"SERVER_QUERY_PORT" env-default:"8080"`
		GamePort  int    `env:"SERVER_GAME_PORT" env-default:"6969"`
		Code      string `env:"SERVER_CODE"`
		MaxPlayer string `env:"SERVER_MAX_PLAYERS"`
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

func handleConnection(conn net.Conn) {
    defer conn.Close()
    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

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

func main() {
	cfg := LoadConfig()
	slog.Info("The config is {cfg}", cfg)

    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GamePort))
    if err != nil {
        log.Fatal("Error creating tcp listener: {err}", err)
    }
    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            slog.Error("Failed to accept connection: {err}", err)
            continue
        }

        slog.Info("Accepted connection: {conn}", conn)

        go handleConnection(conn)
    }
}

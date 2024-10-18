package matchmaker

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Config struct {
	Database struct {
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Database string `env:"DB_DATABASE"`
		Insecure bool   `env:"DB_INSECURE"`
	}
}

func LoadConfig() Config {
	cfg := Config{}

	env := os.Getenv("ENV")
	if "" == env {
		env = "local"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal("Error read env variables")
	}

	return cfg
}

func NewDB(cfg Config) *bun.DB {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Database.Port)),
		pgdriver.WithUser(cfg.Database.User),
		pgdriver.WithPassword(cfg.Database.Password),
		pgdriver.WithDatabase(cfg.Database.Database),
		pgdriver.WithInsecure(cfg.Database.Insecure),
	)
	sqldb := sql.OpenDB(pgconn)
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db
}

type RoomDTO struct {
	Name       string `json:"name" binding:"required"`
	MaxPlayers int    `json:"maxPlayers"`
	Private    bool   `json:"private"`
}

type Room struct {
	bun.BaseModel `bun:"table:rooms,alias:r"`

	ID         int64     `bun:"id,pk,autoincrement" json:"id"`
	Code       string    `bun:"code,type:varchar(6),unique" json:"code"`
	Address    string    `bun:"address,type:varchar(128),notnull" json:"address"`
	QueryPort  int       `bun:"queryPort,type:int,notnull" json:"queryPort"`
	GamePort   int       `bun:"gamePort,type:int,notnull" json:"gamePort"`
	Name       string    `bun:"name,type:varchar(128),notnull" json:"name"`
	Players    int       `bun:"players,type:int,nullzero,notnull,default:0" json:"players"`
	MaxPlayers int       `bun:"maxPlayers,type:int,nullzero,notnull,default:2" json:"maxPlayers"`
	Private    bool      `bun:"private,notnull,default:false"`
	CreatedAt  time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `bun:"updatedAt,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

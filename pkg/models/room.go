package models

import (
	"time"

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
	Players    int       `bun:"players,type:int,nullzero,notnull,default:0" json:"players"`
	MaxPlayers int       `bun:"maxPlayers,type:int,nullzero,notnull,default:2" json:"maxPlayers"`
	Private    bool      `bun:"private,notnull,default:false"`
	CreatedAt  time.Time `bun:"createdAt,nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time `bun:"updatedAt,nullzero,notnull,default:current_timestamp" json:"updatedAt"`
}

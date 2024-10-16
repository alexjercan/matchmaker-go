package migrations

import (
	"context"
	"fmt"
	"rps-multiplayer/pkg/models"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")

		_, err := db.NewCreateTable().
			Model((*models.Room)(nil)).
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")

		_, err := db.NewDropTable().
			Model((*models.Room)(nil)).
			IfExists().
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	})
}

package migrations

import (
	"context"
	"fmt"
	"matchmaker"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [up migration] ")

		_, err := db.NewCreateTable().
			Model((*matchmaker.Room)(nil)).
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")

		_, err := db.NewDropTable().
			Model((*matchmaker.Room)(nil)).
			IfExists().
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	})
}

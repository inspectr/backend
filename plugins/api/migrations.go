package inspectr

import (
	"fmt"

	log "github.com/codeamp/logger"
	resolvers "github.com/inspectr/backend/plugins/api/resolvers"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	gormigrate "gopkg.in/gormigrate.v1"
)

func (x *API) Migrate() {
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		viper.GetString("plugins.api.postgres.host"),
		viper.GetString("plugins.api.postgres.port"),
		viper.GetString("plugins.api.postgres.user"),
		viper.GetString("plugins.api.postgres.dbname"),
		viper.GetString("plugins.api.postgres.sslmode"),
		viper.GetString("plugins.api.postgres.password"),
	))
	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(false)
	db.Set("gorm:auto_preload", true)

	db.AutoMigrate(
		&resolvers.User{},
		&resolvers.UserPermission{},
		&resolvers.Trail{},
	)

	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create users
		{
			ID: "201803021521",
			Migrate: func(tx *gorm.DB) error {
				emails := []string{
					"kilgore@kilgore.trout",
				}

				for _, email := range emails {
					user := resolvers.User{
						Email: email,
					}
					db.Save(&user)

					userPermission := resolvers.UserPermission{
						UserId: user.Model.ID,
						Value:  "admin",
					}
					db.Save(&userPermission)
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return nil
			},
		},
	})

	if err = m.Migrate(); err != nil {
		log.Fatal("Could not migrate: %v", err)
	}

	log.Info("Migration did run successfully")

	defer db.Close()
}

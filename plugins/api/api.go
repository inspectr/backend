package inspectr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"

	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
	redis "github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/inspectr/backend/assets"
	"github.com/inspectr/backend/plugins"
	resolvers "github.com/inspectr/backend/plugins/api/resolvers"
	"github.com/inspectr/backend/plugins/api/utils"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func init() {
	transistor.RegisterPlugin("api", func() transistor.Plugin {
		return NewAPI()
	},
		plugins.Trail{},
		plugins.HeartBeat{})
}

type API struct {
	ServiceAddress string `mapstructure:"service_address"`
	Events         chan transistor.Event
	DB             *gorm.DB
	Redis          *redis.Client
	Schema         *graphql.Schema
}

func NewAPI() *API {
	return &API{}
}

func (x *API) Listen() {
	_, filename, _, _ := runtime.Caller(0)
	fs := http.FileServer(http.Dir(path.Join(path.Dir(filename), "static/")))
	http.Handle("/", fs)
	http.Handle("/query", utils.CorsMiddleware(utils.AuthMiddleware(&relay.Handler{Schema: x.Schema}, x.DB, x.Redis)))

	log.Info(fmt.Sprintf("running API GraphQL server on %v", x.ServiceAddress))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", x.ServiceAddress), handlers.LoggingHandler(os.Stdout, http.DefaultServeMux)))
}

func (x *API) Start(events chan transistor.Event) error {
	var err error

	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s",
		viper.GetString("plugins.api.postgres.host"),
		viper.GetString("plugins.api.postgres.port"),
		viper.GetString("plugins.api.postgres.user"),
		viper.GetString("plugins.api.postgres.dbname"),
		viper.GetString("plugins.api.postgres.sslmode"),
		viper.GetString("plugins.api.postgres.password"),
	))
	//defer x.DB.Close()

	db.AutoMigrate(
		&resolvers.User{},
		&resolvers.UserPermission{},
	)

	schema, err := assets.Asset("plugins/api/schema.graphql")
	if err != nil {
		log.Panic(err)
	}

	parsedSchema, err := graphql.ParseSchema(string(schema), &resolvers.Resolver{DB: db})
	if err != nil {
		log.Panic(err)
	}

	redisDb, err := strconv.Atoi(viper.GetString("redis.database"))
	if err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.server"),
		Password: viper.GetString("redis.password"),
		DB:       redisDb,
	})

	if _, err := redisClient.Ping().Result(); err != nil {
		log.Fatal(err)
	}

	x.Events = events
	x.Schema = parsedSchema
	x.Redis = redisClient

	// DEBUG
	db.LogMode(false)

	x.DB = db

	go x.Listen()

	return nil
}

func (x *API) Stop() {
	log.Info("stopping API service")
}

func (x *API) Subscribe() []string {
	return []string{
		"heartbeat",
		"trail:create",
	}
}

func (x *API) Process(e transistor.Event) error {
	log.DebugWithFields("process API event", log.Fields{
		"event_name": e.Name,
	})

	if e.Name == "trail" {
		payload := e.Payload.(plugins.Trail)
		if e.Action == "create" {
			eventMetadataMarshaled, err := json.Marshal(payload.EventMetadata)
			if err != nil {
				log.Error(err)
			}
			eventMetadataJsonb := postgres.Jsonb{eventMetadataMarshaled}

			actorMetadataMarshaled, err := json.Marshal(payload.ActorMetadata)
			if err != nil {
				log.Error(err)
			}
			actorMetadataJsonb := postgres.Jsonb{actorMetadataMarshaled}

			targetMetadataMarshaled, err := json.Marshal(payload.TargetMetadata)
			if err != nil {
				log.Error(err)
			}
			targetMetadataJsonb := postgres.Jsonb{targetMetadataMarshaled}

			originMetadataMarshaled, err := json.Marshal(payload.OriginMetadata)
			if err != nil {
				log.Error(err)
			}
			originMetadataJsonb := postgres.Jsonb{originMetadataMarshaled}

			trail := resolvers.Trail{
				Timestamp:      payload.Timestamp,
				Event:          payload.Event,
				EventMetadata:  eventMetadataJsonb,
				Actor:          payload.Actor,
				ActorMetadata:  actorMetadataJsonb,
				Target:         payload.Target,
				TargetMetadata: targetMetadataJsonb,
				Origin:         payload.Origin,
				OriginMetadata: originMetadataJsonb,
			}

			if x.DB.Create(&trail).Error != nil {
				log.Error(err)
				x.Events <- e.NewEvent(transistor.GetAction("status"), transistor.GetState("failed"), "ack")
			} else {
				x.Events <- e.NewEvent(transistor.GetAction("status"), transistor.GetState("complete"), "ack")
			}
		}
	}

	return nil
}

package inspectr

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/Inspectr/backend/plugins/api/utils"
	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
	redis "github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/inspectr/backend/assets"
	resolvers "github.com/inspectr/backend/plugins/api/resolvers"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func init() {
	transistor.RegisterPlugin("api", func() transistor.Plugin { return NewAPI() })
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

func (x *API) Migrate() {

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

	parsedSchema, err := graphql.ParseSchema(string(schema), &resolvers.Resolver{DB: x.DB})
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
		"plugins.HeartBeat",
		"plugins.Tick:create",
	}
}

func (x *API) Process(e transistor.Event) error {
	log.InfoWithFields("process API event", log.Fields{
		"event_name": e.Name,
	})

	return nil
}

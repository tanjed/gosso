package apiservice

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/tanjed/go-sso/internal/config"
	"github.com/tanjed/go-sso/internal/db/mongodb"
	"github.com/tanjed/go-sso/internal/db/redisdb"
)


type ApiService struct {
	DB *mongodb.DB
	Redis *redisdb.Redis
	Config *config.Config
	server *http.Server
	router *mux.Router
}

var (
	apiServiceContainer *ApiService
	once             sync.Once
)

func GetApp() *ApiService {
	once.Do(func() {
		apiServiceContainer = NewApiService()
		apiServiceContainer.Boot()
	})
	return apiServiceContainer
}

func (a *ApiService) Boot() {
	a.Config = config.NewConfig()
	a.DB = mongodb.NewDB(a.Config)
	a.Redis = redisdb.NewRedis(a.Config)
}

func (a *ApiService) Destroy() {
	
	a.DB.Close()
	a.Redis.Close()
	a.Config.Close()
}

func NewApiService() *ApiService{
	return &ApiService{}
}
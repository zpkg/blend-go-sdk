package main

import (
	"sync"

	"github.com/blend/go-sdk/logger"
	web "github.com/blendlabs/go-web"
)

// Any is a type alias for interface{}
type Any = interface{}

// APIController implements a simple controller.
type APIController struct {
	db     map[string]Any
	dbLock sync.Mutex
}

// Register adds routes for the controller to the app.
func (ac *APIController) Register(app *web.App) {
	app.GET("/", ac.index)
	app.GET("/api", ac.all)
	app.GET("/api/:key", ac.get)
	app.POST("/api/:key", ac.post)
	app.PUT("/api/:key", ac.put)
	app.DELETE("/api/:key", ac.delete)
}

func (ac *APIController) ensureDB() {
	if ac.db == nil {
		ac.db = map[string]Any{}
	}
}

func (ac *APIController) index(r *web.Ctx) web.Result {
	return r.JSON().OK()
}

func (ac *APIController) all(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	return r.JSON().Result(ac.db)
}

func (ac *APIController) get(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	key, err := r.Param("key")
	if err != nil {
		return r.JSON().BadRequest(err)
	}

	value, hasValue := ac.db[key]
	if !hasValue {
		return r.JSON().NotFound()
	}
	return r.JSON().Result(value)
}

func (ac *APIController) post(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	body, err := r.PostBody()
	if err != nil {
		return r.JSON().InternalError(err)
	}

	key, err := r.Param("key")
	if err != nil {
		return r.JSON().BadRequest(err)
	}
	ac.db[key] = string(body)
	return r.JSON().OK()
}

func (ac *APIController) put(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	key, err := r.Param("key")
	if err != nil {
		return r.JSON().BadRequest(err)
	}

	_, hasValue := ac.db[key]
	if !hasValue {
		return r.JSON().NotFound()
	}

	body, err := r.PostBody()
	if err != nil {
		return r.JSON().InternalError(err)
	}
	ac.db[key] = string(body)

	return r.JSON().OK()
}

func (ac *APIController) delete(r *web.Ctx) web.Result {
	ac.dbLock.Lock()
	defer ac.dbLock.Unlock()
	ac.ensureDB()

	key, err := r.Param("key")
	if err != nil {
		return r.JSON().BadRequest(err)
	}

	_, hasValue := ac.db[key]
	if !hasValue {
		return r.JSON().NotFound()
	}
	delete(ac.db, key)
	return r.JSON().OK()
}

func main() {
	app := web.New().WithLogger(logger.NewFromEnv())
	app.Register(new(APIController))
	app.Start()
}

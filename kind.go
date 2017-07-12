package handlerstore

import (
	"cloud.google.com/go/datastore"
	"context"

	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

type Kind struct {
	client  *datastore.Client
	context context.Context
	Kind    string
}

func (c *Kind) Get(key string, entity interface{}) error {
	return c.client.Get(c.context, datastore.NameKey(c.Kind, key, nil), entity)
}

func (c *Kind) GetAll(q *datastore.Query, entities interface{}) ([]*datastore.Key, error) {
	return c.client.GetAll(c.context, q, entities)
}

func (c *Kind) Put(key string, entity interface{}) (*datastore.Key, error) {
	var datastoreKey *datastore.Key
	if key == "" {
		datastoreKey = datastore.IncompleteKey(c.Kind, nil)
	} else {
		datastoreKey = datastore.NameKey(c.Kind, key, nil)
	}
	return c.client.Put(c.context, datastoreKey, entity)
}

func NewEntity(client *datastore.Client, context context.Context, kind string) *Kind {
	return &Kind{
		client:  client,
		context: context,
		Kind:    kind,
	}
}

type QueryHandler func( map[string]string, *datastore.Query) *datastore.Query

type HandleError func(http.ResponseWriter, error)

func GetEntitiesHandler(e *Kind, fn QueryHandler, entities interface{}, errfn HandleError) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := fn(mux.Vars(r), datastore.NewQuery(e.Kind))
		if _, err := e.GetAll(q, entities); err != nil {
			errfn(w, err)
		}
		json.NewEncoder(w).Encode(entities)
	}
}

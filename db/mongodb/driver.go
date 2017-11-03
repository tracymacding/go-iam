package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"strings"
)

type MongoDriver struct{}

func (mgd *MongoDriver) Open(args ...interface{}) (db.Service, error) {
	mgServer := ""
	for _, s := range args {
		mgServer = mgServer + s.(string) + ","
	}
	mgServer = strings.Trim(mgServer, ",")

	session, err := mgo.Dial(mgServer)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("user")
	index := mgo.Index{
		Key:        []string{"account", "name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	c = session.DB("go_iam").C("key")
	index = mgo.Index{
		Key:        []string{"entity", "entitype"},
		Unique:     false,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	c = session.DB("go_iam").C("group")
	index = mgo.Index{
		Key:        []string{"account", "name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	c = session.DB("go_iam").C("policy")
	index = mgo.Index{
		Key:        []string{"account", "name", "type"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	c = session.DB("go_iam").C("account")
	index = mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		return nil, err
	}

	return &mongoService{mgServer}, nil
}

func init() {
	db.RegisterDriver("mongodb", &MongoDriver{})
}

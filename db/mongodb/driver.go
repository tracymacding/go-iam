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

	return &mongoService{mgServer}, nil
}

//func Init() error {
//	session, err := mgo.Dial("192.168.100.100")
//	if err != nil {
//		return err
//	}
//
//	defer session.Close()
//	session.SetMode(mgo.Monotonic, true)
//
//	// 为objects表创建索引("bucket" + "object_name" + "version")
//	c := session.DB("galaxy_s3_gateway").C("objects")
//	index := mgo.Index{
//		Key:        []string{"bucket", "object_name", "version"},
//		Unique:     true,
//		DropDups:   true,
//		Background: true,
//		Sparse:     true,
//	}
//	err = c.EnsureIndex(index)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func init() {
	db.RegisterDriver("mongodb", &MongoDriver{})
}

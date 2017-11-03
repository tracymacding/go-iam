package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ms *mongoService) CreateKey(key *db.KeyBean) (*db.KeyBean, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("key")
	key.AccessKeyId = bson.NewObjectId()
	key.AccessKeySecret = bson.NewObjectId().Hex()
	err = c.Insert(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (ms *mongoService) GetKey(id string, key *db.KeyBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("key")
	err = c.FindId(bson.ObjectIdHex(id)).One(key)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.KeyNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) DeleteKey(key string) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("key")
	err = c.RemoveId(bson.ObjectIdHex(key))
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.KeyNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) UpdateKey(id string, key *db.KeyBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("key")
	err = c.UpdateId(bson.ObjectIdHex(id), key)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.KeyNotExistError
		}
		if mgo.IsDup(err) {
			return db.KeyExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) ListKey(entity string, entitype int, keys *[]*db.KeyBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("key")
	iter := c.Find(bson.M{"entity": entity, "type": entitype}).Iter()
	return iter.All(keys)
}

func (ms *mongoService) KeyCountOfEntity(entity string, entitype int) (int, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return 0, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("key")
	return c.Find(bson.M{"entity": entity, "entitype": entitype}).Count()
}

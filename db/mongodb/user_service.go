package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ms *mongoService) CreateIamUser(usr *db.UserBean) (*db.UserBean, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("user")
	usr.UserId = bson.NewObjectId()
	err = c.Insert(usr)
	// TODO: check if is user already exist error
	if err != nil {
		if mgo.IsDup(err) {
			return nil, db.UserExistError
		}
		return nil, err
	}
	return usr, nil
}

func (ms *mongoService) UserCountOfAccount(accountId string) (int, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return 0, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	collection := session.DB("go_iam").C("user")
	return collection.Find(bson.M{"account": accountId}).Count()
}

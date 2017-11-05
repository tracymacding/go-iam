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
	if err != nil {
		if mgo.IsDup(err) {
			return nil, db.UserExistError
		}
		return nil, err
	}
	return usr, nil
}

func (ms *mongoService) GetIamUser(account, user string, usr *db.UserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("user")
	err = c.Find(bson.M{"name": user, "account": account}).One(usr)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.UserNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) GetIamUserById(userId string, usr *db.UserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("user")
	err = c.FindId(userId).One(usr)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.UserNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) DeleteIamUser(account, user string) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("user")
	err = c.Remove(bson.M{"name": user, "account": account})
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.UserNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) UpdateIamUser(account, user string, usr *db.UserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("user")
	err = c.Update(bson.M{"name": user, "account": account}, usr)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.UserNotExistError
		}
		if mgo.IsDup(err) {
			return db.UserExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) ListIamUser(account, marker string, max int, usrs *[]*db.UserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("user")
	for len(*usrs) < max {
		query := c.Find(bson.M{"account": account, "name": bson.M{"$gt": marker}}).Sort("name").Limit(max - len(*usrs))
		cnt, err := query.Count()
		if err != nil {
			return err
		}
		if cnt == 0 {
			break
		}
		iter := query.Iter()

		var user db.UserBean
		for iter.Next(&user) {
			*usrs = append(*usrs, &user)
			marker = user.UserName
		}
	}
	return nil
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

package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ms *mongoService) CreateAccount(account *db.AccountBean) (*db.AccountBean, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("account")
	account.AccountId = bson.NewObjectId()
	err = c.Insert(account)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, db.AccountExistError
		}
		return nil, err
	}
	return account, nil
}

func (ms *mongoService) GetAccount(accountId string, account *db.AccountBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("account")
	err = c.FindId(bson.ObjectIdHex(accountId)).One(account)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.AccountNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) DeleteAccount(accountId string) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("account")
	err = c.RemoveId(bson.ObjectIdHex(accountId))
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.AccountNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) UpdateAccount(accountId string, account *db.AccountBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("account")
	account.AccountId = bson.ObjectIdHex(accountId)
	err = c.UpdateId(account.AccountId, account)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.AccountNotExistError
		}
		if mgo.IsDup(err) {
			return db.AccountExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) ListAccount(accountType int, accounts *[]*db.AccountBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("account")
	var iter *mgo.Iter
	if accountType != 0 {
		iter = c.Find(bson.M{"type": accountType}).Iter()
	} else {
		iter = c.Find(nil).Iter()
	}
	err = iter.All(accounts)

	if err != nil {
		return err
	}
	return nil
}

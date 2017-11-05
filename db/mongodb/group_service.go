package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ms *mongoService) CreateGroup(group *db.GroupBean) (*db.GroupBean, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group")
	group.GroupId = bson.NewObjectId()
	err = c.Insert(group)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, db.GroupExistError
		}
		return nil, err
	}
	return group, nil
}

func (ms *mongoService) GetGroup(account, group string, grp *db.GroupBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group")
	err = c.Find(bson.M{"name": group, "account": account}).One(grp)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.GroupNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) GetGroupById(groupId string, grp *db.GroupBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group")
	err = c.FindId(groupId).One(grp)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.GroupNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) DeleteGroup(account, group string) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group")
	err = c.Remove(bson.M{"name": group, "account": account})
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.GroupNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) UpdateGroup(account, group string, grp *db.GroupBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group")
	err = c.Update(bson.M{"name": group, "account": account}, grp)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.GroupNotExistError
		}
		if mgo.IsDup(err) {
			return db.GroupExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) ListGroup(account, marker string, max int, groups *[]*db.GroupBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group")
	for len(*groups) < max {
		query := c.Find(bson.M{"account": account, "name": bson.M{"$gt": marker}}).Sort("name").Limit(max - len(*groups))
		cnt, err := query.Count()
		if err != nil {
			return err
		}
		if cnt == 0 {
			break
		}
		iter := query.Iter()

		var group db.GroupBean
		for iter.Next(&group) {
			*groups = append(*groups, &group)
			marker = group.GroupName
		}
	}
	return nil
}

func (ms *mongoService) GroupCountOfAccount(accountId string) (int, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return 0, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group")
	return c.Find(bson.M{"account": accountId}).Count()
}

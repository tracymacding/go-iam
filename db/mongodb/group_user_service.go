package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ms *mongoService) GroupAddUser(bean *db.GroupUserBean) (*db.GroupUserBean, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group_user")
	bean.JoinId = bson.NewObjectId()
	err = c.Insert(bean)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, db.UserJoinedGroupError
		}
		return nil, err
	}
	return bean, nil
}

func (ms *mongoService) GroupRemoveUser(bean *db.GroupUserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group_user")
	err = c.Remove(bson.M{"group": bean.GroupId, "user": bean.UserId})
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.UserNotJoinedGroupError
		}
		return err
	}
	return nil
}

func (ms *mongoService) ListUserGroup(userId string, beans *[]*db.GroupUserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group_user")
	iter := c.Find(bson.M{"user": userId}).Iter()
	if err = iter.All(beans); err != nil {
		return err
	}
	return nil
}

func (ms *mongoService) ListGroupUser(groupId string, beans *[]*db.GroupUserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group_user")
	iter := c.Find(bson.M{"group": groupId}).Iter()
	if err = iter.All(beans); err != nil {
		return err
	}
	return nil
}

func (ms *mongoService) UserJoinedGroupsNum(userId string) (int, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return 0, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("group_user")
	return c.Find(bson.M{"user": userId}).Count()
}

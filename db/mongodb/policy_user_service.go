package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ms *mongoService) UserAttachPolicy(bean *db.PolicyUserBean) (*db.PolicyUserBean, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_user")
	bean.AttachId = bson.NewObjectId()
	err = c.Insert(bean)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, db.PolicyAttachedUserError
		}
		return nil, err
	}
	return bean, nil
}

func (ms *mongoService) UserDetachPolicy(bean *db.PolicyUserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_user")
	err = c.Remove(bson.M{"policy": bean.PolicyId, "user": bean.UserId})
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.PolicyNotAttachedUserError
		}
		return err
	}
	return nil
}

// func (ms *mongoService) PolicyRemoveUser(bean *db.PolicyUserBean) error {
// 	session, err := mgo.Dial(ms.servers)
// 	if err != nil {
// 		return err
// 	}
//
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
//
// 	c := session.DB("go_iam").C("group_user")
// 	err = c.Remove(bson.M{"group": bean.PolicyId, "user": bean.UserId})
// 	if err != nil {
// 		if err == mgo.ErrNotFound {
// 			return db.UserNotJoinedPolicyError
// 		}
// 		return err
// 	}
// 	return nil
// }
//
// func (ms *mongoService) ListUserPolicy(userId string, beans *[]*db.PolicyUserBean) error {
// 	session, err := mgo.Dial(ms.servers)
// 	if err != nil {
// 		return err
// 	}
//
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
//
// 	c := session.DB("go_iam").C("group_user")
// 	iter := c.Find(bson.M{"user": userId}).Iter()
// 	if err = iter.All(beans); err != nil {
// 		return err
// 	}
// 	return nil
// }
//

func (ms *mongoService) ListUserPolicy(userId string, beans *[]*db.PolicyUserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_user")
	iter := c.Find(bson.M{"user": userId}).Iter()
	if err = iter.All(beans); err != nil {
		return err
	}
	return nil
}

func (ms *mongoService) ListPolicyUser(policyId string, beans *[]*db.PolicyUserBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_user")
	iter := c.Find(bson.M{"policy": policyId}).Iter()
	if err = iter.All(beans); err != nil {
		return err
	}
	return nil
}

func (ms *mongoService) UserAttachedPolicyNum(userId string) (int, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return 0, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_user")
	return c.Find(bson.M{"user": userId}).Count()
}

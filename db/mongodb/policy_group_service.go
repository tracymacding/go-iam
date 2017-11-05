package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ms *mongoService) GroupAttachPolicy(bean *db.PolicyGroupBean) (*db.PolicyGroupBean, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_group")
	bean.AttachId = bson.NewObjectId()
	err = c.Insert(bean)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, db.PolicyAttachedGroupError
		}
		return nil, err
	}
	return bean, nil
}

func (ms *mongoService) GroupDetachPolicy(bean *db.PolicyGroupBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_group")
	err = c.Remove(bson.M{"policy": bean.PolicyId, "group": bean.GroupId})
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.PolicyNotAttachedGroupError
		}
		return err
	}
	return nil
}

// func (ms *mongoService) PolicyRemoveGroup(bean *db.PolicyGroupBean) error {
// 	session, err := mgo.Dial(ms.servers)
// 	if err != nil {
// 		return err
// 	}
//
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
//
// 	c := session.DB("go_iam").C("group_group")
// 	err = c.Remove(bson.M{"group": bean.PolicyId, "group": bean.GroupId})
// 	if err != nil {
// 		if err == mgo.ErrNotFound {
// 			return db.GroupNotJoinedPolicyError
// 		}
// 		return err
// 	}
// 	return nil
// }
//
// func (ms *mongoService) ListGroupPolicy(groupId string, beans *[]*db.PolicyGroupBean) error {
// 	session, err := mgo.Dial(ms.servers)
// 	if err != nil {
// 		return err
// 	}
//
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
//
// 	c := session.DB("go_iam").C("group_group")
// 	iter := c.Find(bson.M{"group": groupId}).Iter()
// 	if err = iter.All(beans); err != nil {
// 		return err
// 	}
// 	return nil
// }
//

func (ms *mongoService) ListGroupPolicy(groupId string, beans *[]*db.PolicyGroupBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_group")
	iter := c.Find(bson.M{"group": groupId}).Iter()
	if err = iter.All(beans); err != nil {
		return err
	}
	return nil
}

func (ms *mongoService) ListPolicyGroup(policyId string, beans *[]*db.PolicyGroupBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_group")
	iter := c.Find(bson.M{"policy": policyId}).Iter()
	if err = iter.All(beans); err != nil {
		return err
	}
	return nil
}

func (ms *mongoService) GroupAttachedPolicyNum(groupId string) (int, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return 0, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy_group")
	return c.Find(bson.M{"group": groupId}).Count()
}

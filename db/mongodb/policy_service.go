package mongodb

import (
	"github.com/go-iam/db"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (ms *mongoService) CreatePolicy(policy *db.PolicyBean) (*db.PolicyBean, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return nil, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy")
	policy.PolicyId = bson.NewObjectId()
	err = c.Insert(policy)
	if err != nil {
		if mgo.IsDup(err) {
			return nil, db.PolicyExistError
		}
		return nil, err
	}
	return policy, nil
}

func (ms *mongoService) GetPolicy(account, policy string, bean *db.PolicyBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy")
	err = c.Find(bson.M{"name": policy, "account": account}).One(bean)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.PolicyNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) GetPolicyById(policyId string, bean *db.PolicyBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy")
	err = c.FindId(policyId).One(bean)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.PolicyNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) DeletePolicy(account, policy string) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy")
	err = c.Remove(bson.M{"name": policy, "account": account})
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.PolicyNotExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) UpdatePolicy(account, policy string, bean *db.PolicyBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy")
	err = c.Update(bson.M{"name": policy, "account": account}, bean)
	if err != nil {
		if err == mgo.ErrNotFound {
			return db.PolicyNotExistError
		}
		if mgo.IsDup(err) {
			return db.PolicyExistError
		}
		return err
	}
	return nil
}

func (ms *mongoService) ListPolicy(account, marker string, ptype, max int, policys *[]*db.PolicyBean) error {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy")
	for len(*policys) < max {
		query := c.Find(bson.M{"account": account, "type": ptype, "name": bson.M{"$gt": marker}}).Sort("name").Limit(max - len(*policys))
		cnt, err := query.Count()
		if err != nil {
			return err
		}
		if cnt == 0 {
			break
		}
		iter := query.Iter()

		var policy db.PolicyBean
		for iter.Next(&policy) {
			*policys = append(*policys, &policy)
			marker = policy.PolicyName
		}
	}
	return nil
}

func (ms *mongoService) PolicyCountOfAccount(accountId string) (int, error) {
	session, err := mgo.Dial(ms.servers)
	if err != nil {
		return 0, err
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("go_iam").C("policy")
	return c.Find(bson.M{"account": accountId}).Count()
}

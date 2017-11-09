package key

import (
	"github.com/bitly/go-simplejson"
	"github.com/go-iam/db"
)

type Key struct {
	accessKeyId     string
	accessKeySecret string
	owner           string
	ownerType       KeyOwnerType
	status          KeyStatus
	createDate      string
}

func FromBean(bean *db.KeyBean) Key {
	return Key{
		accessKeyId:     bean.AccessKeyId.Hex(),
		accessKeySecret: bean.AccessKeySecret,
		owner:           bean.Entity,
		ownerType:       KeyOwnerType(bean.Entitype),
		status:          KeyStatus(bean.Status),
		createDate:      bean.CreateDate,
	}
}

func (k *Key) ToBean() db.KeyBean {
	return db.KeyBean{
		AccessKeySecret: k.accessKeySecret,
		Entity:          k.owner,
		Entitype:        int(k.ownerType),
		Status:          int(k.status),
		CreateDate:      k.createDate,
	}
}

func (k *Key) Json() *simplejson.Json {
	j := simplejson.New()
	j.Set("AccessKeyId", k.accessKeyId)
	j.Set("AccessKeySecret", k.accessKeySecret)
	j.Set("Status", k.status.String())
	j.Set("CreateDate", k.createDate)
	return j
}

type KeyOwnerType int

const (
	Account KeyOwnerType = 1
	IAMUser KeyOwnerType = 2
)

type KeyStatus int

const (
	ErrStatus KeyStatus = 0
	Active    KeyStatus = 1
	Inactive  KeyStatus = 2
)

func (ks KeyStatus) String() string {
	switch ks {
	case Active:
		return "Active"
	case Inactive:
		return "Inactive"
	default:
		panic("invvalid key status")
	}
}

func (ks KeyStatus) IsValid() bool {
	return ks == Active || ks == Inactive
}

func ParseKeyStatus(status string) KeyStatus {
	if status == "Active" {
		return Active
	}
	if status == "Inactive" {
		return Inactive
	}
	return ErrStatus
}

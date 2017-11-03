package key

import (
	"github.com/bitly/go-simplejson"
)

type Key struct {
	accessKeyId     string
	accessKeySecret string
	owner           string
	ownerType       KeyOwnerType
	status          KeyStatus
	createDate      string
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
	if status == "InActive" {
		return Inactive
	}
	return ErrStatus
}

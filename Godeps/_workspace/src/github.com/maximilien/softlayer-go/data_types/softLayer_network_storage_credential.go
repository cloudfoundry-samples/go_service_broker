package data_types

import (
	"time"
)

type SoftLayer_Network_Storage_Credential struct {
	AccountId           int       `json:"accountId"`
	CreateDate          time.Time `json:"createDate"`
	Id                  int       `json:"Id"`
	ModifyDate          time.Time `json:"modifyDate"`
	NasCredentialTypeId int       `json:"nasCredentialTypeId"`
	Password            string    `json:"password"`
	Username            string    `json:"username"`
}

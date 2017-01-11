package model

import (
	"beaconCloud/rpcServer/user-srv/proto/account"
)

type Account struct {
	id       string
	nodeId   string
	username string
	email    string
	salt     string
	password string
	created  int64
	updated  int64
}

func (a *Account) Create(req *account.CreateRequest, rsp *account.CreateResponse) error {
	return nil
}

func (a *Account) Delete(req *account.DeleteRequest, rsp *account.DeleteResponse) error {
	return nil
}

func (a *Account) Update(req *account.UpdateRequest, rsp *account.UpdateResponse) error {
	return nil
}

func (a *Account) Read(req *account.ReadRequest, rsp *account.ReadResponse) error {
	return nil
}

func (a *Account) Search(req *account.SearchRequest, rsp *account.SearchResponse) error {
	return nil
}

func (a *Account) UpdatePassword(req *account.UpdatePasswordRequest, rsp *account.UpdatePasswordResponse) error {
	return nil
}

func (a *Account) SaltAndPassword(req *account.UpdatePasswordRequest, rsp *account.UpdatePasswordResponse) error {
	return nil
}

func (a *Account) ReadUserByUsernameAndEmail(req *account.UpdatePasswordRequest, rsp *account.UpdatePasswordResponse) error {
	return nil
}

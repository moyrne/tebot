package data

import (
	"context"

	"github.com/moyrne/tebot/internal/biz/cqhttp"
	"github.com/moyrne/tractor/dbx"
)

var _ cqhttp.EventRepo = eventRepo{}

func NewEventRepo() cqhttp.EventRepo {
	return eventRepo{
		group:   groupRepo{},
		message: messageRepo{},
		signIn:  signInRepo{},
		user:    userRepo{},
		log:     logRepo{},
	}
}

type eventRepo struct {
	group   cqhttp.GroupRepo
	message cqhttp.MessageRepo
	signIn  cqhttp.SignInRepo
	user    cqhttp.UserRepo
	log     cqhttp.LogRepo
}

func (e eventRepo) SaveMessage(ctx context.Context, tx dbx.Transaction, message *cqhttp.Message) error {
	return e.message.Save(ctx, tx, message)
}

func (e eventRepo) SetMessageReply(ctx context.Context, tx dbx.Transaction, id int64, reply string) error {
	return e.message.SetReply(ctx, tx, id, reply)
}

func (e eventRepo) GetGroupByID(ctx context.Context, tx dbx.Transaction, groupID int64) (*cqhttp.Group, error) {
	return e.group.GetByID(ctx, tx, groupID)
}

func (e eventRepo) GetSignInByUserID(ctx context.Context, tx dbx.Transaction, userID int64) (cqhttp.SignIn, error) {
	return e.signIn.GetByUserID(ctx, tx, userID)
}

func (e eventRepo) SaveSignIn(ctx context.Context, tx dbx.Transaction, signIn *cqhttp.SignIn) error {
	return e.signIn.Save(ctx, tx, signIn)
}

func (e eventRepo) GetUserByUserID(ctx context.Context, tx dbx.Transaction, id int64) (*cqhttp.User, error) {
	return e.user.GetByUserID(ctx, tx, id)
}

func (e eventRepo) SaveUser(ctx context.Context, tx dbx.Transaction, u *cqhttp.User) error {
	return e.user.Save(ctx, tx, u)
}

func (e eventRepo) UpdateUserArea(ctx context.Context, tx dbx.Transaction, userID int64, area string) error {
	return e.user.UpdateArea(ctx, tx, userID, area)
}

func (e eventRepo) Log(ctx context.Context, tx dbx.Transaction, log *cqhttp.Log) error {
	return e.log.Save(ctx, tx, log)
}

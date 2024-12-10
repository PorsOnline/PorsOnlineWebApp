package context

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type appContext struct {
	context.Context
	db           *gorm.DB
	secretDB     *gorm.DB
	shouldCommit bool
}

type AppContextOpt func(*appContext) *appContext // option pattern

func WithDB(db *gorm.DB, shouldCommit bool) AppContextOpt {
	return func(ac *appContext) *appContext {
		ac.db = db
		ac.shouldCommit = shouldCommit
		return ac
	}
}

func NewAppContext(parent context.Context, opts ...AppContextOpt) context.Context {
	ctx := &appContext{Context: parent}
	for _, opt := range opts {
		ctx = opt(ctx)
	}

	return ctx
}

func SetDB(ctx context.Context, db *gorm.DB, secretDB *gorm.DB, shouldCommit bool) {
	appCtx, ok := ctx.(*appContext)
	if !ok {
		return
	}

	appCtx.db = db
	appCtx.secretDB = secretDB
	appCtx.shouldCommit = shouldCommit
}

func GetDB(ctx context.Context) *gorm.DB {
	appCtx, ok := ctx.(*appContext)
	if !ok {
		return nil
	}

	return appCtx.db
}
func GetSecretDB(ctx context.Context) *gorm.DB {
	appCtx, ok := ctx.(*appContext)
	if !ok {
		return nil
	}

	return appCtx.secretDB
}

func Commit(ctx context.Context) error {
	appCtx, ok := ctx.(*appContext)
	if !ok || !appCtx.shouldCommit {
		return nil
	}
	er := appCtx.db.Commit().Error
	if er != nil {
		return er
	}
	er = appCtx.secretDB.Commit().Error
	if er != nil {
		return er
	}
	return nil
}

func Rollback(ctx context.Context) error {
	appCtx, ok := ctx.(*appContext)
	if !ok || !appCtx.shouldCommit {
		return nil
	}
	er := appCtx.db.Rollback().Error
	if er != nil {
		return er
	}
	er = appCtx.secretDB.Rollback().Error
	if er != nil {
		return er
	}
	return nil
}

func CommitOrRollback(ctx context.Context, shouldLog bool) error {
	commitErr := Commit(ctx)
	if commitErr == nil {
		return nil
	}

	if shouldLog {
		log.Println("error on committing transaction, err :", commitErr.Error())
	}

	if err := Rollback(ctx); err != nil {
		log.Println("error on rollback transaction, err :", err.Error())
	}

	return commitErr
}

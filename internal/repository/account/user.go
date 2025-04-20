package account

import (
	"context"
	"time"

	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/HasanNugroho/golang-starter/internal/model/account"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserRepository struct {
	coll *mongo.Collection
	db   *mongo.Database
}

func NewUserRepository(mongoDB *mongo.Database, logger *zerolog.Logger) *UserRepository {
	return &UserRepository{
		coll: mongoDB.Collection("users"),
		db:   mongoDB,
	}
}

func (u *UserRepository) Create(ctx context.Context, user *account.User) error {
	_, err := u.coll.InsertOne(ctx, &user)
	return err
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*account.User, error) {
	var user account.User

	filter := bson.M{"email": email}
	err := u.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &account.User{}, errs.NotFound("not found", err)
		}

		return &account.User{}, errs.Internal("failed to find user", err)
	}

	return &user, nil
}

func (u *UserRepository) FindById(ctx context.Context, id string) (*account.User, error) {
	var user account.User
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return &account.User{}, errs.BadRequest("invalid ID format", err)
	}

	filter := bson.M{"_id": objectID}
	err = u.coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &account.User{}, errs.NotFound("not found", err)
		}

		return &account.User{}, errs.Internal("failed to find user", err)
	}

	return &user, nil
}

func (u *UserRepository) FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]account.User, int, error) {
	var users []account.User
	var totalItems int64

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.Limit)).
		SetLimit(int64(filter.Limit))

	cursor, err := u.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, errs.Internal("failed to fetch user", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, errs.Internal("failed to decode users", err)
	}

	totalItems, err = u.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, errs.Internal("failed to count users", err)
	}

	return &users, int(totalItems), nil
}

func (u *UserRepository) Update(ctx context.Context, id string, user *account.User) error {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errs.BadRequest("invalid ID format", err)
	}

	filter := bson.M{"_id": objectId}
	err = u.coll.FindOneAndUpdate(ctx, filter, bson.M{
		"$set": bson.M{
			"name":       user.Name,
			"email":      user.Email,
			"password":   user.Password,
			"updated_at": time.Now(),
		}}).Err()

	if err != nil {
		return errs.Internal("failed to update data", err)
	}

	return nil
}

func (u *UserRepository) Delete(ctx context.Context, id string) error {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errs.BadRequest("invalid ID format", err)
	}

	filter := bson.M{"_id": objectId}

	err = u.coll.FindOneAndDelete(ctx, filter).Err()
	if err != nil {
		return errs.Internal("failed to delete data", err)
	}

	return nil
}

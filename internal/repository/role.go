package repository

import (
	"context"
	"time"

	"github.com/HasanNugroho/golang-starter/internal/errs"
	"github.com/HasanNugroho/golang-starter/internal/model"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type RoleRepository struct {
	coll *mongo.Collection
	db   *mongo.Database
}

func NewRoleRepository(mongoDB *mongo.Database, logger *zerolog.Logger) *RoleRepository {
	return &RoleRepository{
		coll: mongoDB.Collection("roles"),
		db:   mongoDB,
	}
}

func (r *RoleRepository) Create(ctx context.Context, role *model.Role) error {
	_, err := r.coll.InsertOne(ctx, role)
	return err
}

func (r *RoleRepository) FindById(ctx context.Context, id string) (*model.Role, error) {
	var role model.Role
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return &model.Role{}, errs.BadRequest("invalid ID format", err)
	}

	filter := bson.M{"_id": objectID}
	err = r.coll.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &model.Role{}, errs.NotFound("data not found", err)
		}

		return &model.Role{}, errs.Internal("failed to find data", err)
	}

	return &role, nil
}

func (r *RoleRepository) FindManyByID(ctx context.Context, ids []string) (*[]model.Role, error) {
	panic("")
}

func (r *RoleRepository) FindAll(ctx context.Context, filter *model.PaginationFilter) (*[]model.Role, int, error) {
	var roles *[]model.Role
	var totalItems int64

	opts := options.Find().
		SetSkip(int64((filter.Page - 1) * filter.Limit)).
		SetLimit(int64(filter.Limit))
	// opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, errs.Internal("failed to fetch data", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &roles); err != nil {
		return nil, 0, errs.Internal("failed to decode data", err)
	}

	totalItems, err = r.coll.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, errs.Internal("failed to count data", err)
	}

	return roles, int(totalItems), nil
}

func (r *RoleRepository) Update(ctx context.Context, id string, role *model.Role) error {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errs.BadRequest("invalid ID format", err)
	}

	filter := bson.M{"_id": objectId}
	err = r.coll.FindOneAndUpdate(ctx, filter, bson.M{
		"$set": bson.M{
			"name":        role.Name,
			"permissions": role.Permissions,
			"updated_at":  time.Now(),
		}}).Err()

	if err != nil {
		return errs.Internal("failed to update data", err)
	}

	return nil
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errs.BadRequest("invalid ID format", err)
	}

	filter := bson.M{"_id": objectId}

	err = r.coll.FindOneAndDelete(ctx, filter).Err()
	if err != nil {
		return errs.Internal("failed to update data", err)
	}

	return nil
}

func (r *RoleRepository) AssignUser(ctx context.Context, userId string, roleId string) error {
	userCollection := r.db.Collection("users")
	objectUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return errs.BadRequest("invalid userID format", err)
	}

	objectRoleID, err := bson.ObjectIDFromHex(roleId)
	if err != nil {
		return errs.BadRequest("invalid roleID format", err)
	}

	filter := bson.M{"_id": objectUserID}
	update := bson.M{
		"$addToSet": bson.M{"roles": objectRoleID},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.Internal("failed to update data", err)
	}

	return nil
}

func (r *RoleRepository) UnassignUser(ctx context.Context, userId string, roleId string) error {
	userCollection := r.db.Collection("users")
	objectUserID, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return errs.BadRequest("invalid userID format", err)
	}

	objectRoleID, err := bson.ObjectIDFromHex(roleId)
	if err != nil {
		return errs.BadRequest("invalid roleID format", err)
	}
	filter := bson.M{"_id": objectUserID}
	update := bson.M{
		"$pull": bson.M{"roles": objectRoleID},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errs.Internal("failed to update data", err)
	}

	return nil
}

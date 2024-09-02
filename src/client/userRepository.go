package client

import (
	"context"
	"users-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	ReadAll(ctx context.Context, filter bson.M) ([]models.User, error)
	ReadOne(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, id string, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Client) UserRepository {
	return &mongoUserRepository{
		collection: db.Database("your_db_name").Collection("users"),
	}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *mongoUserRepository) ReadAll(ctx context.Context, filter bson.M) ([]models.User, error) {
	var users []models.User
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var user models.User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *mongoUserRepository) ReadOne(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *mongoUserRepository) Update(ctx context.Context, id string, user *models.User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"lastname":  user.Lastname,
			"birthdate": user.Birthdate,
			"role":      user.Role,
			"email":     user.Email,
			"password":  user.Password,
			"avatar":    user.Avatar,
		},
	}

	_, err = r.collection.UpdateByID(ctx, objectID, update)
	return err
}

func (r *mongoUserRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

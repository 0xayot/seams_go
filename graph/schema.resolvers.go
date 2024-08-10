package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"errors"
	"fmt"
	"seams_go/graph/model"
	"seams_go/models"
	. "seams_go/utils"

	"gorm.io/gorm"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUser) (*model.User, error) {
	fmt.Printf("Input user: %+v\n", input)

	// Check if required fields are present
	if input.Email == nil {
		return nil, fmt.Errorf("email is required")
	}

	// Check if user already exists
	var existingUser models.User
	result := DB.Where("email = ?", *input.Email).First(&existingUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// User doesn't exist, create a new one
			newUser := models.User{
				Email:    *input.Email,
				Name:     input.Name, // Assuming Name is not a pointer in your models.User
				Provider: input.Provider,
			}

			// Safely assign optional fields
			if input.Username != nil {
				newUser.Username = *input.Username
			}
			if input.Avi != nil {
				newUser.Avi = *input.Avi
			}

			if err := DB.Create(&newUser).Error; err != nil {
				return nil, fmt.Errorf("error creating user: %w", err)
			}

			fmt.Printf("Created user: %+v\n", newUser)

			jwt := "string"
			response := model.User{
				ID:       newUser.ID.String(),
				Name:     newUser.Name,
				Provider: newUser.Provider,
				Avi:      &newUser.Avi,
				Username: &newUser.Username,
				Jwt:      &jwt,
			}
			return &response, nil
		}
		return nil, fmt.Errorf("error checking for existing user: %w", result.Error)
	} else {
		jwt := "string"
		response := model.User{
			ID:       existingUser.ID.String(),
			Email:    &existingUser.Email,
			Name:     existingUser.Name,
			Avi:      &existingUser.Avi,
			Username: &existingUser.Username,
			Jwt:      &jwt,
		}
		return &response, nil
	}
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}

// HelloWorld is the resolver for the helloWorld field.
func (r *queryResolver) HelloWorld(ctx context.Context) (*string, error) {
	var res = "hello world"
	return &res, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

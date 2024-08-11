package graph

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
	if &input.Email == nil {
		return nil, fmt.Errorf("email is required")
	}
	validEmail := ValidateGoogleIdToken(input.Email, input.Token)

	fmt.Print("\n Speechless \n \n ", validEmail)

	// Check if user already exists
	var existingUser models.User
	result := DB.Where("email = ?", input.Email).First(&existingUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			newUser := models.User{
				Email:    input.Email,
				Name:     input.Name,
				Provider: input.Provider,
			}

			if input.Username != nil {
				newUser.Username = *input.Username
			}
			if input.Avi != nil {
				newUser.Avi = *input.Avi
			}

			if err := DB.Create(&newUser).Error; err != nil {
				return nil, fmt.Errorf("error creating user: %w", err)
			}

			jwt := "string"
			response := model.User{
				ID:       newUser.ID.String(),
				Name:     newUser.Name,
				Provider: &newUser.Provider,
				Avi:      &newUser.Avi,
				Username: &newUser.Username,
				Email:    &newUser.Email,
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

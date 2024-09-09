package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"seams_go/graph/model"
	"seams_go/models"
	. "seams_go/utils"

	"gorm.io/datatypes"
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

	if !validEmail {
		return nil, fmt.Errorf("Invalid credential")
	}

	// fmt.Print("\n Speechless \n \n ", validEmail)

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

			jwt, err := GenerateJWT(newUser.ID.String())

			if err != nil {
				fmt.Printf("Error signing token: %v\n", err)
			}

			response := model.User{
				ID:       newUser.ID.String(),
				Name:     newUser.Name,
				Provider: &newUser.Provider,
				Avi:      &newUser.Avi,
				Username: &newUser.Username,
				Email:    &newUser.Email,
				Jwt:      &jwt,
			}
			return &response, err
		}
		return nil, fmt.Errorf("error checking for existing user: %w", result.Error)
	} else {
		jwt, err := GenerateJWT(existingUser.ID.String())

		if err != nil {
			fmt.Printf("Error signing token: %v\n", err)
		}

		stuff, err := EnsureAuthurised(jwt)
		fmt.Printf("\n \n stuff: %+v\n ", stuff)

		response := model.User{
			ID:       existingUser.ID.String(),
			Email:    &existingUser.Email,
			Name:     existingUser.Name,
			Avi:      &existingUser.Avi,
			Username: &existingUser.Username,
			Jwt:      &jwt,
		}
		return &response, err
	}
}

// EditUser is the resolver for the editUser field.
func (r *mutationResolver) EditUser(ctx context.Context, input model.EditUser) (*model.User, error) {
	user := UseGQLContext(ctx)

	if user == nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	// Update fields conditionally
	if input.Avi != nil {
		user.Avi = *input.Avi
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Type != nil {
		user.Type = *input.Type
	}

	// Save the updated record
	if err := DB.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	response := model.User{
		ID:       user.ID.String(),
		Email:    &user.Email,
		Name:     user.Name,
		Avi:      &user.Avi,
		Username: &user.Username,
		Type:     &user.Type,
	}
	return &response, nil
}

// CreateMeasurement is the resolver for the createMeasurement field.
func (r *mutationResolver) CreateMeasurement(ctx context.Context, input model.MeasurementInput) (*model.Measurement, error) {
	user := UseGQLContext(ctx)

	if user == nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	var measurementsJSON datatypes.JSON
	if input.Measurements != nil {
		b, err := json.Marshal(input.Measurements)
		if err != nil {
			return nil, fmt.Errorf("error marshaling measurements: %v", err)
		}
		measurementsJSON = datatypes.JSON(b)
	}
	// add validation
	measurement := models.Measurement{
		Name:         *input.Name,
		MeasuredBy:   input.MeasuredBy,
		Measurements: measurementsJSON, // JSON type
		ShoeSize:     input.ShoeSize,
		Active:       input.Active,
		UserID:       user.ID,
	}

	if err := DB.Create(&measurement).Error; err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	response := model.Measurement{
		ID:           measurement.ID.String(),
		Name:         input.Name,
		MeasuredBy:   input.MeasuredBy,
		Measurements: input.Measurements, // JSON type
		ShoeSize:     input.ShoeSize,
		Active:       input.Active,
	}

	return &response, nil
}

// HelloWorld is the resolver for the helloWorld field.
func (r *queryResolver) HelloWorld(ctx context.Context) (*string, error) {
	var res = "hello world"
	return &res, nil
}

// GetCurrentUser is the resolver for the getCurrentUser field.
func (r *queryResolver) GetCurrentUser(ctx context.Context) (*model.User, error) {
	user := UseGQLContext(ctx)

	if user == nil {
		return nil, fmt.Errorf("Unauthorized")
	}

	response := model.User{
		ID:       user.ID.String(),
		Email:    &user.Email,
		Name:     user.Name,
		Avi:      &user.Avi,
		Username: &user.Username,
	}
	return &response, nil
}

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.PublicUser, error) {
	var existingUser models.User
	result := DB.Where("id = ?", id).First(&existingUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User not found")
		}
		return nil, result.Error
	}

	var latestMeasurement models.Measurement
	var currentMeasurement model.Measurement

	err := DB.Where("user_id = ?", id).Where("active = ?", true).
		Order("created_at desc").  // Order by created_at descending
		First(&latestMeasurement). // Get the first result (latest)
		Error

	var measurementsMap map[string]interface{}
	if err == nil {

		err = json.Unmarshal(latestMeasurement.Measurements, &measurementsMap)

		if err != nil {
			log.Fatalf("Error converting JSON to map: %v", err)
		}

		currentMeasurement = model.Measurement{
			ID:           latestMeasurement.ID.String(),
			MeasuredBy:   latestMeasurement.MeasuredBy,
			Measurements: measurementsMap,
			ShoeSize:     latestMeasurement.ShoeSize,
		}

		response := model.PublicUser{
			ID:                 existingUser.ID.String(),
			Name:               existingUser.Name,
			Avi:                &existingUser.Avi,
			Username:           &existingUser.Username,
			CurrentMeasurement: &currentMeasurement,
		}

		return &response, nil

	} else {
		fmt.Println("measurement not found")
		response := model.PublicUser{
			ID:       existingUser.ID.String(),
			Name:     existingUser.Name,
			Avi:      &existingUser.Avi,
			Username: &existingUser.Username,
		}

		return &response, nil
	}
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

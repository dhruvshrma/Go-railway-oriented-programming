package rop

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type UserInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}

type User struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Age       int
	Verified  bool
}

type EnrichedUser struct {
	User
	FullName    string
	IsAdult     bool
	AccountType string
	CreatedAt   time.Time
}

type FormattedUserProfile struct {
	DisplayName string
	Contact     string
	Status      string
	JoinDate    string
}

func parseUserJson(jsonData string) (UserInput, error) {
	var input UserInput
	err := json.Unmarshal([]byte(jsonData), &input)
	return input, err
}

func validateUser(input UserInput) (UserInput, error) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(input.Email) {
		return UserInput{}, errors.New("invalid email format")
	}

	if len(input.FirstName) < 2 || len(input.LastName) < 2 {
		return UserInput{}, errors.New("name too short")
	}

	if input.Age < 13 {
		return UserInput{}, errors.New("users must be at least 13 years old")
	}

	return input, nil
}

func createUser(input UserInput) User {
	return User{
		ID:        fmt.Sprintf("user_%d", time.Now().UnixNano()),
		Email:     strings.ToLower(input.Email),
		FirstName: strings.Title(strings.ToLower(input.FirstName)),
		LastName:  strings.Title(strings.ToLower(input.LastName)),
		Age:       input.Age,
		Verified:  false,
	}
}

func enrichUser(user User) EnrichedUser {
	return EnrichedUser{
		User:        user,
		FullName:    fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		IsAdult:     user.Age >= 18,
		AccountType: determineAccountType(user.Age),
		CreatedAt:   time.Now(),
	}
}

func determineAccountType(age int) string {
	if age < 18 {
		return "Junior"
	} else if age < 65 {
		return "Standard"
	}
	return "Senior"
}

func formatUserProfile(user EnrichedUser) FormattedUserProfile {
	status := "Pending Verification"
	if user.Verified {
		status = "Verified"
	}

	return FormattedUserProfile{
		DisplayName: user.FullName,
		Contact:     user.Email,
		Status:      fmt.Sprintf("%s (%s Account)", status, user.AccountType),
		JoinDate:    user.CreatedAt.Format("Jan 2, 2006"),
	}
}

func prettyPrint(prefix string, v interface{}) {
	fmt.Printf("\n--- %s ---\n", prefix)
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func MapExample() {

	validJSON := `{"email":"john.doe@example.com","first_name":"John","last_name":"Doe","age":25}`

	invalidJSON := `{"email":"not-an-email","first_name":"J","last_name":"D","age":10}`

	fmt.Println("Valid User using Map and Bind")

	validResult := Pipe(
		Ok(validJSON), parseUserJson)

	validatedResult := Pipe(validResult, validateUser)
	userResult := Map(validatedResult, createUser)          // Map: UserInput -> User (cannot fail)
	enrichedResult := Map(userResult, enrichUser)           // Map: User -> EnrichedUser (cannot fail)
	profileResult := Map(enrichedResult, formatUserProfile) // Map: EnrichedUser -> FormattedUserProfile (cannot fail)

	// Process and display results
	validResult.OnSuccess(func(input UserInput) {
		prettyPrint("PARSED INPUT", input)
	})

	validatedResult.OnSuccess(func(input UserInput) {
		prettyPrint("VALIDATED INPUT", input)
	})

	userResult.OnSuccess(func(user User) {
		prettyPrint("CREATED USER", user)
	})

	enrichedResult.OnSuccess(func(user EnrichedUser) {
		prettyPrint("ENRICHED USER", user)
	})

	profileResult.OnSuccess(func(profile FormattedUserProfile) {
		prettyPrint("FINAL PROFILE", profile)
	}).OnError(func(err error) {
		fmt.Printf("Error: %v\n", err)
	})

	invalidResult := Pipe(Ok(invalidJSON), parseUserJson)
	userResult.OnSuccess(func(user User) {
		prettyPrint("CREATED USER", user)
	})

	enrichedResult.OnSuccess(func(user EnrichedUser) {
		prettyPrint("ENRICHED USER", user)
	})

	profileResult.OnSuccess(func(profile FormattedUserProfile) {
		prettyPrint("FINAL PROFILE", profile)
	}).OnError(func(err error) {
		fmt.Printf("Error: %v\n", err)
	})

	invalidValidatedResult := Pipe(invalidResult, validateUser)
	invalidUserResult := Map(invalidValidatedResult, createUser)

	// This will show the error from validation
	invalidUserResult.OnSuccess(func(user User) {
		prettyPrint("CREATED USER (Invalid)", user)
	}).OnError(func(err error) {
		fmt.Printf("Error: %v\n", err)
	})
}

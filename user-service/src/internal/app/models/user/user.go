package user

import "time"

type Domain struct {
	ID               string
	Username         string
	Email            string
	RegistrationDate time.Time
	Version          string
}

type Profile struct {
	Domain      string
	FirstName   string
	LastName    string
	PhoneNumber string
	DOB         time.Time
	Address     string
	AboutMe     string
	ProfPicURL  string
}

type Auth struct {
	Domain         string
	Password       string
	ActivationLink string
	Activated      bool
}

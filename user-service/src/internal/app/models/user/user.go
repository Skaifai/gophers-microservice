package user

import "time"

type User struct {
	ID               string
	Role             string
	Username         string
	Email            string
	Password         string
	RegistrationDate time.Time
	FirstName        string
	LastName         string
	PhoneNumber      string
	DOB              time.Time
	Address          string
	AboutMe          string
	ProfPicUrl       string
	Activated        bool
	Version          string
}

func Assemble(d Domain, a Auth, p Profile) *User {
	return &User{
		ID:               d.ID,
		Role:             a.Role,
		Username:         d.Username,
		Email:            d.Email,
		Password:         a.Password,
		RegistrationDate: d.RegistrationDate,
		FirstName:        p.FirstName,
		LastName:         p.LastName,
		PhoneNumber:      p.PhoneNumber,
		DOB:              p.DOB,
		Address:          p.Address,
		AboutMe:          p.AboutMe,
		ProfPicUrl:       p.ProfPicURL,
		Activated:        a.Activated,
		Version:          d.Version,
	}
}

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
	Role           string
	Password       string
	ActivationLink string
	Activated      bool
}

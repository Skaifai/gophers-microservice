package mail

type Mailer interface {
	Send(string, string, any) error
}

type service struct {
	mailer Mailer
}

func New(mailer Mailer) *service {
	return &service{
		mailer: mailer,
	}
}

func (svc *service) SendActivationMail(firstname, email, activationUUID string) error {
	data := map[string]any{
		"name":  firstname,
		"email": email,
		"uuid":  activationUUID,
	}

	err := svc.mailer.Send(email, "user_welcome.tmpl", data)
	if err != nil {
		return err
	}

	return nil
}

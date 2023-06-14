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

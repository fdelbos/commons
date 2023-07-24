package kratos

type (
	Flows interface {
		Login() string
		Register() string
		Recovery() string
		Logout() string
	}

	flows struct {
		publicKratosFrontEnd string
	}
)

func NewFlows(publicKratosFrontEnd string) Flows {
	return &flows{
		publicKratosFrontEnd: publicKratosFrontEnd,
	}
}

func (f flows) Login() string {
	return f.publicKratosFrontEnd + "/self-service/login/browser"
}

func (f flows) Register() string {
	return f.publicKratosFrontEnd + "/self-service/registration/browser"
}

func (f flows) Recovery() string {
	return f.publicKratosFrontEnd + "/self-service/recovery/browser"
}

func (f flows) Logout() string {
	return f.publicKratosFrontEnd + "/self-service/logout/browser"
}

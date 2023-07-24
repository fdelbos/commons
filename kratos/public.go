package kratos

import client "github.com/ory/kratos-client-go"

type (
	Public interface {
		Client() *client.APIClient
	}

	public struct {
		publicApiURL string
	}
)

func NewPublic(publicApiURL string) Public {
	return &public{
		publicApiURL: publicApiURL,
	}
}

func (p public) Client() *client.APIClient {
	conf := client.NewConfiguration()
	conf.Servers = []client.ServerConfiguration{
		{
			URL: p.publicApiURL,
		},
	}

	client := client.NewAPIClient(conf)
	return client
}

package plugin

import (
	"fmt"
	"github.com/openbao/openbao/api"
	"github.com/openbao/openbao/sdk/framework"
	"github.com/openbao/openbao/sdk/logical"
	"net/http"
)

const AUTH_USER_KEY = "id"
const AUTH_TOKEN_KEY = "token"

const CONFIG_EXCHANGE_URL_KEY = "url"
const CONFIG_TOKEN_KEY = "token"
const CONFIG_AGBOT_RENEWAL_KEY = "renewal"
const CONFIG_VAULT_API_KEY = "apiurl"

type ohAuthPlugin struct {

	// The vault auth plugin framework.
	*framework.Backend

	// An HTTP client used to call the openhorizon exchange.
	httpClient *http.Client

	// A vault client used to interact with the system.
	vc *api.Client
}

func OHAuthPlugin(c *logical.BackendConfig) *ohAuthPlugin {
	var b ohAuthPlugin
	var err error

	b.httpClient, err = NewHTTPClient()
	if err != nil {
		panic(ohlog(fmt.Sprintf("could not establish HTTP client, error: %v", err)))
	}

	b.vc, err = api.NewClient(nil)
	if err != nil {
		panic(ohlog(fmt.Sprintf("could not create vault client, error: %v", err)))
	}

	b.Backend = &framework.Backend{
		BackendType: logical.TypeCredential,
		AuthRenew:   b.pathAuthRenew,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{"login"},
			SealWrapStorage: []string{"config"},
		},
		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "login",
				Fields: map[string]*framework.FieldSchema{
					AUTH_USER_KEY: &framework.FieldSchema{
						Type: framework.TypeString,
					},
					AUTH_TOKEN_KEY: &framework.FieldSchema{
						Type: framework.TypeString,
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathAuthLogin,
				},
			},
			&framework.Path{
				Pattern: "config",
				Fields: map[string]*framework.FieldSchema{
					CONFIG_EXCHANGE_URL_KEY: &framework.FieldSchema{
						Type: framework.TypeString,
					},
					CONFIG_TOKEN_KEY: &framework.FieldSchema{
						Type: framework.TypeString,
					},
					CONFIG_AGBOT_RENEWAL_KEY: &framework.FieldSchema{
						Type: framework.TypeInt,
					},
					CONFIG_VAULT_API_KEY: &framework.FieldSchema{
						Type: framework.TypeString,
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathConfig,
				},
			},
		},
	}

	return &b
}

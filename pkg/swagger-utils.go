package pkg

import (
	"embed"
	"encoding/json"
	"io"
	"reflect"

	sgroupsv1 "github.com/PRO-Robotech/sgroups-proto/pkg/api/sgroups/v1"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

var (
	//go:embed api/*/*.swagger.json api/*/*/*.swagger.json
	swaggerStore embed.FS
)

var (
	//ErrSwaggerNotExist when document is not found
	ErrSwaggerNotExist = errors.New("swagger doc is no exist")
)

// SwaggerUtil ...
type SwaggerUtil[T any] struct{}

func (u SwaggerUtil[T]) reg(p string) {
	swaggerPaths[reflect.TypeOf((*T)(nil)).Elem()] = p
}

// GetSpec ...
func (u SwaggerUtil[T]) GetSpec() (*spec.Swagger, error) {
	res := new(spec.Swagger)
	ty := reflect.TypeOf((*T)(nil)).Elem()
	p := swaggerPaths[ty]
	err := whenFindSwagger(p, func(r io.Reader) error {
		return json.NewDecoder(r).Decode(res)
	})
	return res, err
}

// GetRaw ...
func (u SwaggerUtil[T]) GetRaw() (json.RawMessage, error) {
	var ret json.RawMessage
	ty := reflect.TypeOf((*T)(nil)).Elem()
	p := swaggerPaths[ty]
	err := whenFindSwagger(p, func(reader io.Reader) error {
		data, e := io.ReadAll(reader)
		ret = data
		return e
	})
	return ret, err
}

var (
	swaggerPaths = make(map[reflect.Type]string)
)

func whenFindSwagger(p string, f func(reader io.Reader) error) error {
	data, e := swaggerStore.Open(p)
	if e != nil {
		return ErrSwaggerNotExist
	}
	return f(data)
}

func init() {
	const (
		apiSGroups = "api/sgroups/v1/services.swagger.json"
	)

	apis := [...]interface{ reg(string) }{
		SwaggerUtil[sgroupsv1.SGroupsNamespaceAPIServer]{},
		SwaggerUtil[sgroupsv1.SGroupsNamespaceAPIClient]{},

		SwaggerUtil[sgroupsv1.SGroupsAddressGroupsAPIServer]{},
		SwaggerUtil[sgroupsv1.SGroupsAddressGroupsAPIClient]{},

		SwaggerUtil[sgroupsv1.SGroupsHostsAPIServer]{},
		SwaggerUtil[sgroupsv1.SGroupsHostsAPIClient]{},

		SwaggerUtil[sgroupsv1.SGroupsHostBindingAPIServer]{},
		SwaggerUtil[sgroupsv1.SGroupsHostBindingAPIClient]{},

		SwaggerUtil[sgroupsv1.SGroupsNetworksAPIServer]{},
		SwaggerUtil[sgroupsv1.SGroupsNetworksAPIClient]{},

		SwaggerUtil[sgroupsv1.SGroupsNetworkBindingAPIServer]{},
		SwaggerUtil[sgroupsv1.SGroupsNetworkBindingAPIClient]{},
	}

	for _, api := range apis {
		api.reg(apiSGroups)
	}
}

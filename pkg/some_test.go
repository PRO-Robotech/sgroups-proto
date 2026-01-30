package pkg

import (
	"testing"

	sgroupsv1 "github.com/PRO-Robotech/sgroups-proto/pkg/api/sgroups/v1"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestSwaggerUtil(t *testing.T) {
	var u SwaggerUtil[sgroupsv1.SgroupsAPIServer]
	s, e := u.GetSpec()
	require.NoError(t, e)
	require.NotNil(t, s)
}

func TestClosableClient(t *testing.T) {
	var conn *grpc.ClientConn
	var c ClosableClient[sgroupsv1.SgroupsAPIClient]
	err := c.Init(conn)
	require.NoError(t, err)
	require.NotNil(t, c.C)
}

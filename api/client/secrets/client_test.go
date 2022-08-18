// Copyright 2021 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package secrets_test

import (
	"time"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api/base/testing"
	apisecrets "github.com/juju/juju/api/client/secrets"
	"github.com/juju/juju/core/secrets"
	"github.com/juju/juju/rpc/params"
	coretesting "github.com/juju/juju/testing"
)

var _ = gc.Suite(&SecretsSuite{})

type SecretsSuite struct {
	coretesting.BaseSuite
}

func (s *SecretsSuite) TestNewClient(c *gc.C) {
	apiCaller := testing.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		return nil
	})
	client := apisecrets.NewClient(apiCaller)
	c.Assert(client, gc.NotNil)
}

func ptr[T any](v T) *T {
	return &v
}

func (s *SecretsSuite) TestListSecrets(c *gc.C) {
	data := map[string]string{"foo": "bar"}
	now := time.Now()
	uri := secrets.NewURI()
	apiCaller := testing.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		c.Check(objType, gc.Equals, "Secrets")
		c.Check(version, gc.Equals, 0)
		c.Check(id, gc.Equals, "")
		c.Check(request, gc.Equals, "ListSecrets")
		c.Check(arg, gc.DeepEquals, params.ListSecretsArgs{
			ShowSecrets: true,
		})
		c.Assert(result, gc.FitsTypeOf, &params.ListSecretResults{})
		*(result.(*params.ListSecretResults)) = params.ListSecretResults{
			[]params.ListSecretResult{{
				URI:            uri.String(),
				Version:        1,
				OwnerTag:       "application-mysql",
				ScopeTag:       "unit-mysql-0",
				Provider:       "juju",
				ProviderID:     "provider-id",
				RotatePolicy:   string(secrets.RotateHourly),
				ExpireTime:     ptr(now),
				NextRotateTime: ptr(now.Add(time.Hour)),
				Description:    "shhh",
				Label:          "foobar",
				Revision:       2,
				CreateTime:     now,
				UpdateTime:     now.Add(time.Second),
				Value:          &params.SecretValueResult{Data: data},
			}},
		}
		return nil
	})
	client := apisecrets.NewClient(apiCaller)
	result, err := client.ListSecrets(true)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, jc.DeepEquals, []apisecrets.SecretDetails{{
		Metadata: secrets.SecretMetadata{
			URI:            uri,
			Version:        1,
			OwnerTag:       "application-mysql",
			ScopeTag:       "unit-mysql-0",
			Provider:       "juju",
			ProviderID:     "provider-id",
			RotatePolicy:   secrets.RotateHourly,
			ExpireTime:     ptr(now),
			NextRotateTime: ptr(now.Add(time.Hour)),
			Description:    "shhh",
			Label:          "foobar",
			Revision:       2,
			CreateTime:     now,
			UpdateTime:     now.Add(time.Second),
		},
		Value: secrets.NewSecretValue(data),
	}})
}

func (s *SecretsSuite) TestListSecretsError(c *gc.C) {
	apiCaller := testing.APICallerFunc(func(objType string, version int, id, request string, arg, result interface{}) error {
		*(result.(*params.ListSecretResults)) = params.ListSecretResults{
			[]params.ListSecretResult{{
				URI: "secret:9m4e2mr0ui3e8a215n4g",
				Value: &params.SecretValueResult{
					Error: &params.Error{Message: "boom"},
				},
			}},
		}
		return nil
	})
	client := apisecrets.NewClient(apiCaller)
	result, err := client.ListSecrets(true)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(result, gc.HasLen, 1)
	c.Assert(result[0].Error, gc.Equals, "boom")
}

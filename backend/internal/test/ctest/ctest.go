package ctest

import (
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/model"
	arepository "github.com/mickamy/sampay/internal/domain/auth/repository"
	"github.com/mickamy/sampay/internal/lib/either"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func ConnErr(t *testing.T, err error) *connect.Error {
	t.Helper()

	require.Error(t, err)
	var connErr *connect.Error
	require.ErrorAs(t, err, &connErr)

	return connErr
}

func AssertCode(t *testing.T, expected connect.Code, connErr *connect.Error) {
	t.Helper()

	assert.Equalf(t, expected, connect.CodeOf(connErr), "code=%s", connect.CodeOf(connErr).String())
}

func LocalizedMessage(t *testing.T, connErr *connect.Error) string {
	t.Helper()

	details := connErr.Details()

	for _, detail := range details {
		if localized, ok := either.Must(detail.Value()).(*errdetails.LocalizedMessage); ok {
			return localized.GetMessage()
		}
	}

	t.Fatal("localized message not found in error details")
	return ""
}

func FieldViolation(t *testing.T, connErr *connect.Error) *errdetails.BadRequest_FieldViolation {
	t.Helper()

	violations := FieldViolations(t, connErr)
	require.Len(t, violations, 1)
	return violations[0]
}

func FieldViolations(t *testing.T, connErr *connect.Error) []*errdetails.BadRequest_FieldViolation {
	t.Helper()

	var violations []*errdetails.BadRequest_FieldViolation
	for _, detail := range connErr.Details() {
		if br, ok := either.Must(detail.Value()).(*errdetails.BadRequest); ok {
			violations = append(violations, br.GetFieldViolations()...)
		}
	}

	if len(violations) == 0 {
		t.Fatal("field violations not found in error details")
	}

	return violations
}

func AuthorizationHeader(t *testing.T, infra *di.Infra) (string, string) {
	t.Helper()

	endUser := tseed.EndUser(t, infra.WriterDB)
	session := model.MustNewSession(endUser.UserID)
	require.NoError(t, arepository.NewSession(infra.KVS).Create(t.Context(), session))
	return endUser.UserID, "Bearer " + session.Tokens.Access.Value
}

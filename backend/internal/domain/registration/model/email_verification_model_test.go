package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/domain/registration/fixture"
	"mickamy.com/sampay/internal/domain/registration/model"
)

func TestEmailVerification_IsRequested(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T) model.EmailVerification
		assert  func(t *testing.T, got bool)
	}{
		{
			name: "not requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerification(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, false)
			},
		},
		{
			name: "requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationRequested(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, true)
			},
		},
		{
			name: "verified",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationVerified(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, true)
			},
		},
		{
			name: "consumed",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationConsumed(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, true)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			m := tc.arrange(t)
			t.Logf("m.Requested: %+v", m.Requested)
			t.Logf("m.Verified: %+v", m.Verified)
			t.Logf("m.Consumed: %+v", m.Consumed)

			// act
			got := m.IsRequested()

			// assert
			tc.assert(t, got)
		})
	}
}

func TestEmailVerification_IsVerified(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T) model.EmailVerification
		assert  func(t *testing.T, got bool)
	}{
		{
			name: "not requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerification(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, false)
			},
		},
		{
			name: "requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationRequested(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, false)
			},
		},
		{
			name: "verified",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationVerified(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, true)
			},
		},
		{
			name: "consumed",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationConsumed(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, true)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			m := tc.arrange(t)
			t.Logf("m.Requested: %+v", m.Requested)
			t.Logf("m.Verified: %+v", m.Verified)
			t.Logf("m.Consumed: %+v", m.Consumed)

			// act
			got := m.IsVerified()

			// assert
			tc.assert(t, got)
		})
	}
}

func TestEmailVerification_IsConsumed(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T) model.EmailVerification
		assert  func(t *testing.T, got bool)
	}{
		{
			name: "not requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerification(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, false)
			},
		},
		{
			name: "requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationRequested(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, false)
			},
		},
		{
			name: "verified",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationVerified(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, false)
			},
		},
		{
			name: "consumed",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationConsumed(nil)
			},
			assert: func(t *testing.T, got bool) {
				assert.Equal(t, got, true)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			m := tc.arrange(t)
			t.Logf("m.Requested: %+v", m.Requested)
			t.Logf("m.Verified: %+v", m.Verified)
			t.Logf("m.Consumed: %+v", m.Consumed)

			// act
			got := m.IsConsumed()

			// assert
			tc.assert(t, got)
		})
	}
}

func TestEmailVerification_Request(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T) model.EmailVerification
		assert  func(t *testing.T, err error)
	}{
		{
			name: "not requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerification(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationRequested(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "verified",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationVerified(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, model.ErrEmailVerificationVerified)
			},
		},
		{
			name: "consumed",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationConsumed(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, model.ErrEmailVerificationConsumed)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			m := tc.arrange(t)
			t.Logf("m.Requested: %+v", m.Requested)
			t.Logf("m.Verified: %+v", m.Verified)
			t.Logf("m.Consumed: %+v", m.Consumed)

			// act
			err := m.Request(config.Auth().EmailVerificationExpiresInDuration())

			// assert
			tc.assert(t, err)
		})
	}
}

func TestEmailVerification_Verify(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T) model.EmailVerification
		assert  func(t *testing.T, err error)
	}{
		{
			name: "not requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerification(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, model.ErrEmailVerificationNotRequested)
			},
		},
		{
			name: "requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationRequested(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "verified",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationVerified(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "consumed",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationConsumed(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, model.ErrEmailVerificationConsumed)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			m := tc.arrange(t)
			t.Logf("m.Requested: %+v", m.Requested)
			t.Logf("m.Verified: %+v", m.Verified)
			t.Logf("m.Consumed: %+v", m.Consumed)

			// act
			err := m.Verify()

			// assert
			tc.assert(t, err)
		})
	}
}

func TestEmailVerification_Consume(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(t *testing.T) model.EmailVerification
		assert  func(t *testing.T, err error)
	}{
		{
			name: "not requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerification(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, model.ErrEmailVerificationNotRequested)
			},
		},
		{
			name: "requested",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationRequested(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, model.ErrEmailVerificationNotVerified)
			},
		},
		{
			name: "verified",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationVerified(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "consumed",
			arrange: func(t *testing.T) model.EmailVerification {
				return fixture.EmailVerificationConsumed(nil)
			},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			m := tc.arrange(t)
			t.Logf("m.Requested: %+v", m.Requested)
			t.Logf("m.Verified: %+v", m.Verified)
			t.Logf("m.Consumed: %+v", m.Consumed)

			// act
			err := m.Consume()

			// assert
			tc.assert(t, err)
		})
	}
}

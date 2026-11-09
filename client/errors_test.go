package client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestErrorErrorReturnsMessage(t *testing.T) {
	t.Parallel()
	err := NewRequestError(400, "Bad Request")
	assert.Equal(t, "(400) Bad Request", err.Error())
}

func TestRequestErrorContainsStatusCodeAndMessage(t *testing.T) {
	t.Parallel()
	err := NewRequestError(500, "Server Error")
	assert.Equal(t, 500, err.StatusCode)
	assert.Equal(t, "Server Error", err.Message)
}

func TestNotFoundErrorReturnsFormattedMessage(t *testing.T) {
	t.Parallel()
	err := NewNotFoundError("User not found")
	assert.Equal(t, "NotFoundError (404) User not found", err.Error())
}

func TestNotFoundErrorHasStatusCode404(t *testing.T) {
	t.Parallel()
	err := NewNotFoundError("")
	assert.Equal(t, 404, err.StatusCode)
}

func TestNotFoundErrorErrorsIs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		err         error
		errorsIsFun func(assert.TestingT, error, error, ...interface{}) bool
	}{
		{
			"TestNotFoundErrorDifferentMessage",
			NewNotFoundError("not found"),
			assert.ErrorIs,
		},
		{
			"TestWrappedNotFoundError",
			fmt.Errorf("wrapper of %w", NewNotFoundError("wrapped")),
			assert.ErrorIs,
		},
		{
			"TestDeeplyWrappedNotFoundError",
			fmt.Errorf("wrapper of wrapper %w", fmt.Errorf("wrapper of %w", NewNotFoundError("wrapped"))),
			assert.ErrorIs,
		},
		{
			"TestNonRequestError",
			errors.New(""),
			assert.NotErrorIs,
		},
		{
			"TestWrappedOtherError",
			fmt.Errorf("wrapper of %w", NewRequestError(400, "Bad request")),
			assert.NotErrorIs,
		},
		{
			"TestNilError",
			nil,
			assert.NotErrorIs,
		},
	}
	notFoundError := NewNotFoundError("")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.errorsIsFun(t, tt.err, notFoundError)
		})
	}
}

// Framework resources check errors.Is(err, client.NotFoundError{}), with the
// zero value as the target.
func TestNotFoundErrorErrorsIsZeroValueTarget(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		err         error
		errorsIsFun func(assert.TestingT, error, error, ...interface{}) bool
	}{
		{
			"TestNotFoundErrorFromAPI",
			NewNotFoundError(`{"errors":[{"title":"Not found or unauthorized","status":"404"}]}`),
			assert.ErrorIs,
		},
		{
			"TestWrappedNotFoundError",
			fmt.Errorf("wrapper of %w", NewNotFoundError("wrapped")),
			assert.ErrorIs,
		},
		{
			"TestZeroValueNotFoundError",
			NotFoundError{},
			assert.ErrorIs,
		},
		{
			"TestRequestError404",
			NewRequestError(404, "Not Found"),
			assert.NotErrorIs,
		},
		{
			"TestNonRequestError",
			errors.New(""),
			assert.NotErrorIs,
		},
		{
			"TestNilError",
			nil,
			assert.NotErrorIs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.errorsIsFun(t, tt.err, NotFoundError{})
		})
	}
}

// A zero-value NotFoundError must still match a NewNotFoundError target.
func TestZeroValueNotFoundErrorMatchesNewNotFoundErrorTarget(t *testing.T) {
	t.Parallel()
	assert.ErrorIs(t, NotFoundError{}, NewNotFoundError(""))
}

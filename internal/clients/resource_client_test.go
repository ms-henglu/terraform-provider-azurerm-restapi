package clients_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/stretchr/testify/assert"
)

func TestRetryClient(t *testing.T) {
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(5, 30, 2, 1.5, []string{"retry error"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps)
	_, err := retryClient.Get(context.Background(), "", "")
	assert.NoError(t, err)
	assert.Equal(t, 3, mock.requestCount)
}

func TestRetryClientRegexp(t *testing.T) {
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(1, 5, 1.5, 1.5, []string{"^retry"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps)
	_, err := retryClient.Get(context.Background(), "", "")
	assert.NoError(t, err)
	assert.Equal(t, 3, mock.RequestCount())
}

func TestRetryClientMultiRegexp(t *testing.T) {
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(1, 5, 1.5, 1.5, []string{"nomatch", "^retry"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps)
	_, err := retryClient.Get(context.Background(), "", "")
	assert.NoError(t, err)
	assert.Equal(t, 3, mock.RequestCount())
}

func TestRetryClientContextDeadline(t *testing.T) {
	mock := NewMockResourceClient(t, nil, nil, 3, errors.New("retry error"))
	bkof, errRegExps := clients.NewRetryableErrors(60, 60, 1.5, 1.5, []string{"^retry"})
	retryClient := clients.NewResourceClientRetryableErrors(mock, bkof, errRegExps)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()
	_, err := retryClient.Get(ctx, "", "")
	elapsed := time.Since(start)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.True(t, elapsed < 15*time.Second)
	// Test that the context was cancelled
	_, ok := <-ctx.Done()
	assert.False(t, ok)
}

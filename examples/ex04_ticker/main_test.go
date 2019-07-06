package main

import (
	"github.com/elliotchance/pepper/peppertest"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClock_Now(t *testing.T) {
	c := &Clock{}

	doc, err := peppertest.RenderToDocument(c)
	require.NoError(t, err)

	expectedTime := time.Now().Format(time.RFC1123)
	require.Contains(t, doc.Text(), "The time now is " + expectedTime)
}

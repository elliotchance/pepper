package main

import (
	"github.com/elliotchance/pepper/peppertest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCounter_AddOne(t *testing.T) {
	c := &Counter{}

	doc, err := peppertest.RenderToDocument(c)
	require.NoError(t, err)
	require.Contains(t, doc.Text(), "Counter: 0")

	c.AddOne()
	c.AddOne()

	doc, err = peppertest.RenderToDocument(c)
	require.NoError(t, err)
	require.Contains(t, doc.Text(), "Counter: 2")
}

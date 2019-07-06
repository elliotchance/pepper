package main

import (
	"github.com/elliotchance/pepper/peppertest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCounters_Total(t *testing.T) {
	c := &Counters{
		Counters: []*Counter{
			{Number: 3}, {}, {Number: 2},
		},
	}

	doc, err := peppertest.RenderToDocument(c)
	require.NoError(t, err)

	require.Contains(t, doc.Find("tr").Last().Text(), "Total: 5")
}

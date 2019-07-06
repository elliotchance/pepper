package main

import (
	"github.com/elliotchance/pepper/peppertest"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPeople_Render(t *testing.T) {
	c := &People{
		Names: []string{"Jack", "Jill"},
	}

	doc, err := peppertest.RenderToDocument(c)
	require.NoError(t, err)

	rows := doc.Find("tr")
	require.Equal(t, 2, rows.Length())
	require.Contains(t, rows.Eq(0).Text(), "Jack")
	require.Contains(t, rows.Eq(1).Text(), "Jill")
}

func TestPeople_Add(t *testing.T) {
	c := &People{
		Names: []string{"Jack", "Jill"},
	}

	c.Name = "Bob"
	c.Add()

	doc, err := peppertest.RenderToDocument(c)
	require.NoError(t, err)

	rows := doc.Find("tr")
	require.Equal(t, 3, rows.Length())
	require.Contains(t, rows.Eq(0).Text(), "Jack")
	require.Contains(t, rows.Eq(1).Text(), "Jill")
	require.Contains(t, rows.Eq(2).Text(), "Bob")
}

func TestPeople_Delete(t *testing.T) {
	c := &People{
		Names: []string{"Jack", "Jill"},
	}

	c.Delete("Jack")

	doc, err := peppertest.RenderToDocument(c)
	require.NoError(t, err)

	rows := doc.Find("tr")
	require.Equal(t, 1, rows.Length())
	require.Contains(t, rows.Eq(0).Text(), "Jill")
}

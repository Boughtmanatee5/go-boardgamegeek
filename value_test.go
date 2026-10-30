package boardgamegeek

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValueTypesReadValueAttribute(t *testing.T) {
	var v struct {
		Int     *IntValue   `xml:"int"`
		Empty   *IntValue   `xml:"empty"`
		Float   *FloatValue `xml:"float"`
		Text    *ValueAttr  `xml:"text"`
		Missing *IntValue   `xml:"missing"`
	}
	err := xml.Unmarshal([]byte(`<x><int value="-2200"/><empty value=""/><float value="7.09012"/><text value="Not Ranked"/></x>`), &v)
	require.NoError(t, err)

	require.NotNil(t, v.Int)
	assert.Equal(t, IntValue(-2200), *v.Int)
	require.NotNil(t, v.Empty)
	assert.Equal(t, IntValue(0), *v.Empty)
	require.NotNil(t, v.Float)
	assert.Equal(t, FloatValue(7.09012), *v.Float)
	require.NotNil(t, v.Text)
	assert.Equal(t, ValueAttr("Not Ranked"), *v.Text)
	assert.Nil(t, v.Missing)
}

func TestValueTypesRejectNonNumbers(t *testing.T) {
	var i struct {
		V IntValue `xml:"v"`
	}
	assert.ErrorContains(t, xml.Unmarshal([]byte(`<x><v value="Not Ranked"/></x>`), &i), `<v value="Not Ranked">`)

	var f struct {
		V FloatValue `xml:"v"`
	}
	assert.Error(t, xml.Unmarshal([]byte(`<x><v value="n/a"/></x>`), &f))
}

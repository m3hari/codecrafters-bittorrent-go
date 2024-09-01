package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeString(t *testing.T) {
	result, _, err := decodeString("4:spam")
	assert.Equal(t, "spam", result)
	assert.Nil(t, err)

	result, remaining, err := decodeString("2:hii100elove")
	assert.Equal(t, "hi", result)
	assert.Equal(t, "i100elove", remaining)
	assert.Nil(t, err)
}

func TestDecodeString_Invalid(t *testing.T) {
	_, _, err := decodeString("4:sp")
	assert.NotNil(t, err)
	_, _, err = decodeString("10:")
	assert.NotNil(t, err)
	_, _, err = decodeString("1x:ha")
	assert.NotNil(t, err)

}

func TestDecodeInteger(t *testing.T) {
	result, _, err := decodeInteger("i52e")
	assert.Equal(t, 52, result)
	assert.Nil(t, err)

	result, _, err = decodeInteger("i-100ethisdoesnotmatter")
	assert.Equal(t, -100, result)
	assert.Nil(t, err)
}
func TestDecodeInteger_Invalid(t *testing.T) {
	_, _, err := decodeInteger("i35d2e")
	assert.NotNil(t, err)
	_, _, err = decodeInteger("i2")
	assert.NotNil(t, err)
}

func TestDecodeList(t *testing.T) {
	result, remaining, err := decodeList("l4:spam4:eggse")
	assert.Equal(t, []any{"spam", "eggs"}, result)
	assert.Equal(t, remaining, "")
	assert.Nil(t, err)

	result, remaining, err = decodeList("l4:spami52ee")
	assert.Equal(t, []any{"spam", 52}, result)
	assert.Equal(t, remaining, "")
	assert.Nil(t, err)

	result, remaining, err = decodeList("le")
	assert.Equal(t, []any{}, result)
	assert.Equal(t, remaining, "")
	assert.Nil(t, err)
}

func TestDecodeList_Invalid(t *testing.T) {
	_, _, err := decodeList("l4:")
	assert.NotNil(t, err)

	_, _, err = decodeList("lhi")
	assert.NotNil(t, err)

	_, _, err = decodeList("l2:hi")
	assert.NotNil(t, err)

}

func TestMain(t *testing.T) {
	result, err := run([]string{"decode", "5:hello"})
	assert.Equal(t, result, "\"hello\"")
	assert.Nil(t, err)

	_, err = run([]string{"decode", "invalid"})
	assert.NotNil(t, err)
	_, err = run([]string{"decode", "-5:hmm"})
	assert.NotNil(t, err)
	_, err = run([]string{"decode", "hi:"})
	assert.NotNil(t, err)

	result, err = run([]string{"decode", "i52e"})
	assert.Equal(t, result, "52")
	assert.Nil(t, err)

	result, err = run([]string{"decode", "l5:helloi52ee"})
	assert.Equal(t, result, "[\"hello\",52]")
	assert.Nil(t, err)
}

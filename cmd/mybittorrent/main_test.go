package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	result, err := run([]string{"decode", "5:hello"})
	assert.Equal(t, result, "\"hello\"")
	assert.Nil(t, err)

	result, err = run([]string{"decode", "i52e"})
	assert.Equal(t, result, "52")
	assert.Nil(t, err)

	result, err = run([]string{"decode", "l5:helloi52ee"})
	assert.Equal(t, result, "[\"hello\",52]")
	assert.Nil(t, err)

	result, err = run([]string{"decode", "d3:foo3:bar5:helloi52ee"})
	assert.Equal(t, result, "{\"foo\":\"bar\",\"hello\":52}")
	assert.Nil(t, err)

	result, err = run([]string{"decode", "d4:spaml1:a1:bee"})
	assert.Equal(t, result, "{\"spam\":[\"a\",\"b\"]}")
	assert.Nil(t, err)
}

func TestRun_Invalid(t *testing.T) {
	_, err := run([]string{})
	assert.NotNil(t, err)

	_, err = run([]string{"unknown_command"})
	assert.NotNil(t, err)

	_, err = run([]string{"decode", "invalid"})
	assert.NotNil(t, err)

	_, err = run([]string{"decode", "-5:hmm"})
	assert.NotNil(t, err)

	_, err = run([]string{"decode", "hi:"})
	assert.NotNil(t, err)

}

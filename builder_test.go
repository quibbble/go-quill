package go_quill

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Builder(t *testing.T) {
	b := Builder{}
	assert.True(t, b.Info() != nil)
}

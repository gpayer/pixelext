package nodes

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnimUpdate(t *testing.T) {
	assert := assert.New(t)
	a := NewAnimation("")
	assert.NotNil(a)
	var values []float64

	a.Add(0, 1, 0.5, 1, func(v float64) {
		values = append(values, v)
	})

	a.Update(0.1)
	assert.Len(values, 0)
	a.Update(0.9)
	a.Update(0.4)
	a.Update(0.3)

	assert.Len(values, 3)
	assert.Equal(0.5, values[0])
	assert.Equal(9.0, math.Round(values[1]*10))
	assert.Equal(1.0, values[2])
}

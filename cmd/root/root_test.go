package root

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*type promptMock struct {
	msg string
	err error
}

func (p promptMock) Run() (string, error) {
	// return expected result
	return p.msg, p.err
}

type selectMock struct {
	msg string
	err error
}

func (p selectMock) Run() (int, string, error) {
	// return expected result
	return 1, p.msg, p.err
}*/

func TestExecute(t *testing.T) {
	opts.VerboseLog = true
	err := rootCmd.Execute()

	assert.Nil(t, err)
	opts.SetZeroValues()
}

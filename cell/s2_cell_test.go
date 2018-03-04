package cell

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func InitTest() {
	os.Setenv("GO_ENV", "test")
}

func TestFetchCellToken(t *testing.T) {
	cellToken := FetchCellToken(12.34, 77.67, 15)
	assert.Equal(t, "3bae8e1dc", cellToken, "Cell token is not matching")
}

func TestFetchCell(t *testing.T) {
  cellToken := FetchCell(12.34, 77.67, 15)
	assert.Equal(t, "3bae8e1dc", cellToken.Token, "Cell token is not matching")
  assert.Equal(t, 12.340884614266653, cellToken.Center[0], "Cell token is not matching")
  assert.Equal(t, 77.66936079375242, cellToken.Center[1], "Cell token is not matching")
}

func TestMain(m *testing.M) {
	InitTest()
	os.Exit(m.Run())
}

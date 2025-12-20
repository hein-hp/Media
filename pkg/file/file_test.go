package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountFiles(t *testing.T) {

	files, err := CountFiles("/Users/xxx/xxx")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 20, files)

}

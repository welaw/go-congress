package congress

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	tedCruz         = "C001098"
	tedCruzFullName = "Rafael Edward (Ted) CRUZ"
)

func TestGoQuery(t *testing.T) {
	//url := GetActionsUrl("115", "house", "1")
	//color.Red("url %v", url)
	votes, err := GetRollCalls("115", "house", "1", "1")
	for k, v := range votes {
		//color.Red("k=%v, v=%v", k, v)
	}
	require.NoError(t, err)
}

package filesystem

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

const (
	testEnvPath = "../.env"
)

func newTestLogger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	return logger
}

func TestListBranchVersions(t *testing.T) {
	//err := godotenv.Load(testEnvPath)
	//require.NoError(t, err)

	_, err := OpenFile("filesystem.go")
	require.NoError(t, err)

	//fmt.Printf("result count: %v\n", len(sets))
	//for _, set := range sets {
	//fmt.Printf("set: %+v\n", set)
	//}

}

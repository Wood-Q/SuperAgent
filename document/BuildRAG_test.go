package document

import (
	"SuperAgent/initialize"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildRAG(t *testing.T) {
	initialize.InitConfig("..")
	initialize.InitClient()
	initialize.InitEmbedder()
	initialize.InitIndexer()
	err := BuildRAG("mulyse", `缪尔赛思`, map[string]any{"test": "test"}, context.Background())
	println(err)
	require.NoError(t, err)
}

package document

import (
	"SuperAgent/initialize"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchRAG(t *testing.T) {
	initialize.InitConfig("..")
	initialize.InitClient()
	initialize.InitEmbedder()
	initialize.InitRetriever()
	documents := SearchRAG("缪缪", context.Background())
	fmt.Println(documents)
	require.NotEmpty(t, documents)
}

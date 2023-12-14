package socket

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSocket(t *testing.T) {
	_, err := Socket("google.com", 80, 3)
	require.NoError(t, err, "连接socket失败")
}

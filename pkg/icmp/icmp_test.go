package icmp

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIcmp(t *testing.T) {
	r, err := Icmp("baidu.com", 4)
	require.NoError(t, err, "执行Icmp失败")
	fmt.Println(r)
}

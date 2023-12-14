package whois

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWhois(t *testing.T) {
	resp, err := Whois(context.Background(), 5, "startops.com")
	require.NoError(t, err, "查看Whois失败")
	fmt.Println(resp)
}

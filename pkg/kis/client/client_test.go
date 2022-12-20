package pkgkisclient

import (
	pkgconfig "github.com/daeungkim/kis-go/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_Authenticate(t *testing.T) {
	t.Parallel()

	cred := Credential{}
	require.NoError(t, pkgconfig.LoadConfig("KIS", &cred))

	tests := []struct {
		name  string
		given Credential
		want  bool
	}{
		{
			name: "AppKey 가 잘못되면 에러를 리턴한다.",
			given: Credential{
				APISecret: cred.APISecret,
			},
			want: true,
		},
		{
			name: "AppSecret 이 잘못되면 에러를 리턴한다.",
			given: Credential{
				APIKey: cred.APIKey,
			},
			want: true,
		},
		{
			name:  "Token 을 리턴한다.",
			given: cred,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewClient(tt.given, true)

			got, err := c.authenticate()

			if tt.want {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.NotEmpty(t, got)
		})
	}
}

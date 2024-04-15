package main

import (
	"go/token"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gnolang/gno/gnovm/pkg/integration"
)

func Test_ScriptsTranspile(t *testing.T) {
	p := testscript.Params{
		Dir: "testdata/gno_transpile",
	}

	if coverdir, ok := integration.ResolveCoverageDir(); ok {
		err := integration.SetupTestscriptsCoverage(&p, coverdir)
		require.NoError(t, err)
	}

	err := integration.SetupGno(&p, t.TempDir())
	require.NoError(t, err)

	testscript.Run(t, p)
}

func Test_parseGoBuildErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		output        string
		expectedError error
	}{
		{
			name:          "empty output",
			output:        "",
			expectedError: nil,
		},
		{
			name:   "random output",
			output: "xxx",
			expectedError: goscanner.ErrorList{
				&goscanner.Error{
					Msg: "Additional go build errors:\nxxx",
				},
			},
		},
		{
			name: "some errors",
			output: `xxx
main.gno:6:2: nasty error
pkg/file.gno:60:20: ugly error`,
			expectedError: goscanner.ErrorList{
				&goscanner.Error{
					Pos: token.Position{
						Filename: "main.gno",
						Line:     6,
						Column:   2,
					},
					Msg: "nasty error",
				},
				&goscanner.Error{
					Pos: token.Position{
						Filename: "pkg/file.gno",
						Line:     60,
						Column:   20,
					},
					Msg: "ugly error",
				},
				&goscanner.Error{
					Msg: "Additional go build errors:\nxxx",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseGoBuildErrors(tt.output)

			assert.Equal(t, tt.expectedError, err)
		})
	}
}

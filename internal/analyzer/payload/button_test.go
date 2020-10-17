package payload

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBadPayload(t *testing.T) {
	var testCases = map[string]struct {
		payload string
		taskId  int64
		err     bool
	}{
		`12345 with quotes`: {
			payload: "{\"button\":\"12345\"}",
			taskId:  12345,
			err:     false,
		},
		`12345`: {
			payload: `{"button":"12345"}`,
			taskId:  12345,
			err:     false,
		},
		`not a number`: {
			payload: `{"button":"what?"}`,
			taskId:  0,
			err:     true,
		},
		`broken JSON`: {
			payload: `{"button":}`,
			taskId:  0,
			err:     true,
		},
		`empty JSON`: {
			payload: ``,
			taskId:  0,
			err:     true,
		},
	}

	for testName, testCase := range testCases {
		testError := fmt.Sprintf(`test "%s" error`, testName)
		actual, err := GetTaskId(testCase.payload)
		if (err != nil) != testCase.err {
			t.Errorf("Parse() error = %v, wantErr %v", err, testCase.err)
			return
		}
		assert.Equal(t, testCase.taskId, actual, testError)
	}
}

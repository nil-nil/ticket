package ticketeventbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	table := []testCase{
		{
			description:      "L1 wildcard",
			topic:            "*.test.example",
			expectL1:         "*",
			expectL2:         "test",
			expectL3:         "example",
			expectFuncOutput: "foo",
		},
		{
			description:      "L2 wildcard",
			topic:            "test.*.example",
			expectL1:         "test",
			expectL2:         "*",
			expectL3:         "example",
			expectFuncOutput: "bar",
		},
		{
			description:      "L2 wildcard",
			topic:            "test.example.*",
			expectL1:         "test",
			expectL2:         "example",
			expectL3:         "*",
			expectFuncOutput: "baz",
		},
	}

	bus, _ := NewBus()

	for _, tc := range table {
		t.Run(tc.description, testCaseFunc(bus, tc))
	}
}

func testCaseFunc(bus *ticketEventBus, tc testCase) func(t *testing.T) {
	return func(t *testing.T) {
		var got string
		bus.sub(tc.topic, func(eventKey string, data interface{}) {
			got = tc.expectFuncOutput
		})

		_, ok := bus.subs[tc.expectL1]
		assert.True(t, ok, "expected first sub layer [%s]", tc.expectL1)

		_, ok = bus.subs[tc.expectL1][tc.expectL2]
		assert.True(t, ok, "expected first sub layer [%s][%s]", tc.expectL1, tc.expectL2)

		_, ok = bus.subs[tc.expectL1][tc.expectL2][tc.expectL3]
		assert.True(t, ok, "expected first sub layer [%s][%s][%s]", tc.expectL1, tc.expectL2, tc.expectL3)

		subs := bus.subs[tc.expectL1][tc.expectL2][tc.expectL3]
		for _, sub := range subs {
			sub("", nil)
		}
		assert.Equal(t, tc.expectFuncOutput, got, "expected original func")
	}
}

type testCase struct {
	description      string
	topic            string
	expectL1         string
	expectL2         string
	expectL3         string
	expectFuncOutput string
}

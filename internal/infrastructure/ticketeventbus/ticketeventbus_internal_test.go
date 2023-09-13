package ticketeventbus

import (
	"testing"

	"github.com/nil-nil/ticket/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	var partialWcOutput, exactOutput int
	var wcOutput bool

	bus, _ := NewBus(".")
	bus.sub("*.test.example", func(eventKey string, data interface{}) {
		partialWcOutput = 1
	})
	bus.sub("test.*.example", func(eventKey string, data interface{}) {
		partialWcOutput = 2
	})
	bus.sub("test.example.*", func(eventKey string, data interface{}) {
		partialWcOutput = 3
	})
	bus.sub("1.test.example", func(eventKey string, data interface{}) {
		exactOutput = 1
	})
	bus.sub("test.2.example", func(eventKey string, data interface{}) {
		exactOutput = 2
	})
	bus.sub("test.example.3", func(eventKey string, data interface{}) {
		exactOutput = 3
	})
	bus.sub("*.*.*", func(eventKey string, data interface{}) {
		wcOutput = true
	})

	table := []struct {
		description           string
		expectPartialWcOutput int
		expectExactOutput     int
	}{
		{description: "1.test.example", expectPartialWcOutput: 1, expectExactOutput: 1},
		{description: "2.test.example", expectPartialWcOutput: 1, expectExactOutput: 0},
		{description: "test.2.example", expectPartialWcOutput: 2, expectExactOutput: 2},
		{description: "test.3.example", expectPartialWcOutput: 2, expectExactOutput: 0},
		{description: "test.example.3", expectPartialWcOutput: 3, expectExactOutput: 3},
		{description: "test.example.4", expectPartialWcOutput: 3, expectExactOutput: 0},
	}

	for _, tc := range table {
		t.Run(tc.description, func(t *testing.T) {
			partialWcOutput, exactOutput = 0, 0
			wcOutput = false

			var funcs []func(eventKey string, data interface{}) = bus.match(tc.description)
			for _, f := range funcs {
				f("", nil)
			}

			assert.Equal(t, tc.expectExactOutput, exactOutput, "expected exact topic match")
			assert.Equal(t, tc.expectPartialWcOutput, partialWcOutput, "expected partial wildcard topic match")
			assert.True(t, wcOutput, "always expects full wildcards to match")
		})
	}
}

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

	bus, _ := NewBus(".")

	err := bus.sub("notvalid")
	assert.EqualError(t, err, domain.ErrEventKeyInvalid.Error(), "expect meaningful error on invalid topic")

	for _, tc := range table {
		t.Run(tc.description, testCaseFunc(bus, tc))
	}
}

func testCaseFunc(bus *ticketEventBus, tc testCase) func(t *testing.T) {
	return func(t *testing.T) {
		var got string
		err := bus.sub(tc.topic, func(eventKey string, data interface{}) {
			got = tc.expectFuncOutput
		})
		assert.NoError(t, err, "valid sub shouldn't error")

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

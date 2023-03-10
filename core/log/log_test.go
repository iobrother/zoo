package log

import (
	"context"
	"fmt"
	"testing"
)

func suiteSetUp(suiteName string) func() {
	fmt.Printf("\tSetUp fixture for suite %s\n", suiteName)
	return func() {
		fmt.Printf("\tTearDown fixture for suite %s\n", suiteName)
	}
}

func testWithName(t *testing.T) {
	l := WithName("zoo")
	l.Info("testWithName")
}

func testWithField(t *testing.T) {
	l := WithField(String("key1", "foo"), Int("key2", 100))
	l.Info("testWithField")
}

func testWithContext(t *testing.T) {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, "requestId", "req-1") //nolint:staticcheck,revive
	l := WithContext(ctx, "requestId")
	l.Info("testWithContext")
}

func testWithMap(t *testing.T) {
	l := WithMap(map[string]any{"key3": "bar", "key4": 200})
	l.Info("testWithMap")
}

func TestWith(t *testing.T) {
	t.Cleanup(suiteSetUp(t.Name()))
	t.Run("testWithName", testWithName)
	t.Run("testWithField", testWithField)
	t.Run("testWithContext", testWithContext)
	t.Run("testWithMap", testWithMap)
}

func testInfo(t *testing.T) {
	Info("testInfo")
}

func testSetLevel(t *testing.T) {
	Info("testSetLevel")
	Debug("this message is not visible")
	SetLevel(DebugLevel)
	Debug("this message is visible")
}

func TestLog(t *testing.T) {
	t.Cleanup(suiteSetUp(t.Name()))
	t.Run("testInfo", testInfo)
	t.Run("testSetLevel", testSetLevel)
}

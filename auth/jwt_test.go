package auth

import "testing"

func TestToken(t *testing.T) {
	t.Run("it runs without errors", func(t *testing.T) {
		_, err := Token(0)
		if err != nil {
			t.Errorf("got %q error", err.Error())
		}
	})

	t.Run("it returns non empty tokens", func(t *testing.T) {
		details, _ := Token(0)
		if len(details.Access) == 0 || len(details.Refresh) == 0 {
			t.Error("got empty tokens")
		}
	})
}

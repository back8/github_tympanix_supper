package app

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppWebsiteIndex(t *testing.T) {
	for _, uri := range []string{
		"http://localhost",
		"http://localhost/blablabla",
		"http://localhost/doesnotexists.php",
		"http://localhost/should_always_show_index.html",
	} {
		req := httptest.NewRequest("GET", uri, nil)
		w := httptest.NewRecorder()

		handler := WebAppHandler("../web")

		handler.ServeHTTP(w, req)

		resp := w.Result()

		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		file, err := os.Open("../web/index.html")
		require.NoError(t, err)

		golden, err := ioutil.ReadAll(file)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, golden, body)
	}
}

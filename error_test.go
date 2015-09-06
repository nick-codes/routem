package routem

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHTTPError(t *testing.T) {
	err := fmt.Errorf("Error Message")
	httpErr := NewHTTPError(http.StatusBadRequest, err)

	assert.Equal(t, err.Error(), httpErr.Error())
	assert.Equal(t, httpErr.Code(), http.StatusBadRequest)
}

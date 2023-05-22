package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddZerotierToken(c *gin.Context) {
	// Read the port number from the file
	w := c.Writer
	r := c.Request
	port, err := ioutil.ReadFile("/var/lib/zerotier-one/zerotier-one.port")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the request path and remove "/zt"
	path := strings.TrimPrefix(r.URL.Path, "/v1/zt")
	fmt.Println(path)

	// Build the target URL
	targetURL := fmt.Sprintf("http://localhost:%s%s", strings.TrimSpace(string(port)), path)

	// Create a new request
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add the X-ZT1-AUTH header
	authToken, err := ioutil.ReadFile("/var/lib/zerotier-one/authtoken.secret")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("X-ZT1-AUTH", strings.TrimSpace(string(authToken)))

	copyHeaders(req.Header, r.Header)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	copyHeaders(w.Header(), resp.Header)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the response to the client
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

func copyHeaders(destination, source http.Header) {
	for key, values := range source {
		for _, value := range values {
			destination.Add(key, value)
		}
	}
}

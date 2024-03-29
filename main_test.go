package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createTestFile(name, data string) error {
	return os.WriteFile(name, []byte(data), 0644)
}
func deleteTestFile(name string) error {
	return os.Remove(name)
}

func Test_httpRequest(t *testing.T) {
	testConfig := &Config{KeyDir: "./"}
	handler := requestHandler(testConfig)
	filename := "test-file"
	fileContent := "test-data"
	err := createTestFile(testConfig.KeyDir+filename, fileContent)
	defer deleteTestFile(testConfig.KeyDir + filename)
	if err != nil {
		t.Errorf("error creating testfile: %v", err)
		return
	}

	t.Run("Normal request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+filename, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		res := w.Result()
		if res.StatusCode != 200 {
			t.Errorf("request, expected statuscode to be 200 got %v", res.StatusCode)
		}
		defer res.Body.Close()
		reply := KVPair{}
		json.NewDecoder(res.Body).Decode(&reply)
		if err != nil {
			t.Errorf("reply parsing json, expected error to be nil got %v", err)
		}
		t.Logf("content: %+v", reply)
		if reply.Key != filename {
			t.Errorf("reply, expected key to be %v got %v", filename, reply.Key)
		}
		if reply.Key != filename {
			t.Errorf("reply, expected value to be %v got %v", fileContent, reply.Value)
		}
	})
	filename = "hello-world.crt"
	fileContent = `-----BEGIN CERTIFICATE-----
MIICszCCAZsCFFliWPpolbzZOFb2Vf/KOcU/914/MA0GCSqGSIb3DQEBCwUAMBYx
FDASBgNVBAMMC2hlbGxvLXdvcmxkMB4XDTI0MDMyOTE0NTkyNloXDTM0MDMyNzE0
NTkyNlowFjEUMBIGA1UEAwwLaGVsbG8td29ybGQwggEiMA0GCSqGSIb3DQEBAQUA
A4IBDwAwggEKAoIBAQC5S+W+51QSssEoUYLt4xGgWI5q4i+SauWf3ZcZ8RIbemUU
rmMMvVxqrbCpSbWn4jWErY2jQM8dr93eSdAL3x+aS9IGDYn+0Mo2Xwafa5XEyDIC
+PhChE9Sx1Oh9mbLykEU2iLH9cec5DCVoWcQ2NFrW3JDvB8U7fywdzoFrKh4Llr7
d+YqU1jNOO5ZKnRwpQ6/Pgr88QfHnHNtwBdSUsuCCr+ukz3vPgGiQopu9p7AhtIJ
jWRvB12lXTn6qqdHXM3kZ6dqKq3ekpGeYy7RldXjmPNeyfYe7Qkz/25q+RBPvVkW
wrz3WPuWaojXrtjsCCamLQuiYJQsXWpEy1GllxAjAgMBAAEwDQYJKoZIhvcNAQEL
BQADggEBAAGiFg5ghqJmndaqpUPEsQbd07/Dq0xd6JN7HhVfvTazY2NDOW9hAOdN
kjrD2ScyUmpn8JUF93GmF41F1XtzHOrWupPp8W1xiS3oxxTMggAyuJuUyfHKg2xS
T88y9iUWgu89PjUdtNmVlFySDRQEjsWch4tlCDmGzeYhAmoMkadv7JzfYcLr2Uik
VNkF9ct2FzvGPLNCOi+qeuj2CWc2/9p/F/1musJFi3lYcG3iuN+E3rVxobvv52r6
wzwra++In6XiwTH5iupipc9k+5xfWD2MH60j1KiAaH2u8jUbFWm56aLWN20kteeM
n3Ov8z9WHplkIV3S+bQTRSliIFbYIZ0=
-----END CERTIFICATE-----`
	err = createTestFile(testConfig.KeyDir+filename, fileContent)
	defer deleteTestFile(testConfig.KeyDir + filename)

	testConfig = &Config{KeyDir: "./", LogJson: true}
	handler = requestHandler(testConfig)
	t.Run("complicated request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+filename, nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		res := w.Result()
		if res.StatusCode != 200 {
			t.Errorf("request, expected statuscode to be 200 got %v", res.StatusCode)
		}
		defer res.Body.Close()
		reply := KVPair{}
		json.NewDecoder(res.Body).Decode(&reply)
		if err != nil {
			t.Errorf("reply parsing json, expected error to be nil got %v", err)
		}
		t.Logf("content: %+v", reply)
		if reply.Key != filename {
			t.Errorf("reply, expected key to be %v got %v", filename, reply.Key)
		}
		if reply.Key != filename {
			t.Errorf("reply, expected value to be %v got %v", fileContent, reply.Value)
		}
	})
}

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type podDetailsStruct struct {
	Name             string `json:"Name"`
	ApplicationGroup string `json:"ApplicationGroup"`
	RunningPodsCount int    `json:"RunningPodsCount"`
}

type podDetailsList []podDetailsStruct

func TestReturnPodsPerService(t *testing.T) {
	req, err := http.NewRequest("GET", "/services", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Creates a new recorder to record the response received by the services endpoint
	respRec := httptest.NewRecorder()
	handler := http.HandlerFunc(returnPodsPerService)
	handler.ServeHTTP(respRec, req)
	if status := respRec.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var podDetails []podDetailsStruct

	// Convert JSON to struct
	err = json.Unmarshal([]byte(respRec.Body.String()), &podDetails)
	if err != nil {
		panic(err)
	}
	// Need to do custom sorting based on "Name", because map returns elements in a random order
	sort.Sort(podDetailsList(podDetails))

	// Convert back from struct to JSON for comparison
	podDetailsString, err := json.Marshal(podDetails)
	if err != nil {
		panic(err)
	}

	// convert byte to string
	got := string(podDetailsString)
	expected := `[{"Name":"blissful-goodall","ApplicationGroup":"beta","RunningPodsCount":1},{"Name":"confident-cartwright","ApplicationGroup":"beta","RunningPodsCount":1},{"Name":"happy-colden","ApplicationGroup":"none","RunningPodsCount":1},{"Name":"quirky-raman","ApplicationGroup":"gamma","RunningPodsCount":1},{"Name":"stoic-sammet","ApplicationGroup":"alpha","RunningPodsCount":2}]`

	if got != expected {
		t.Errorf("Handler returned unexpected body: got %v expected %v", got, expected)
	}
}

func TestReturnPodsPerAppGroup(t *testing.T) {

	req, err := http.NewRequest("GET", "/services/{applicationGroup}", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{"applicationGroup": "alpha"})
	respRec := httptest.NewRecorder()
	handler := http.HandlerFunc(returnPodsPerAppGroup)
	handler.ServeHTTP(respRec, req)
	if status := respRec.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"Name":"stoic-sammet","ApplicationGroup":"alpha","RunningPodsCount":2}]`
	got := strings.TrimSuffix(respRec.Body.String(), "\n")

	if got != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", got, expected)
	}
}

func (pod podDetailsList) Len() int {
	return len(pod)
}

func (pod podDetailsList) Less(i, j int) bool {
	return pod[i].Name < pod[j].Name
}

func (pod podDetailsList) Swap(i, j int) {
	pod[i], pod[j] = pod[j], pod[i]
}

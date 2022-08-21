package main

import (
	"challenge/pkg/client"
	"challenge/pkg/controller"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	zapcr "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var clientset *kubernetes.Clientset
var log = logf.Log.WithName("platform-test")

func main() {

	// setting up custom zap logger to have more flexibilty
	log := zapcr.New()
	ctrl.SetLogger(log)

	var err error
	clientset, err = client.PrepareClientSet()
	if err != nil {
		log.Error(err, "Unable to create client set")
	}

	handleRequests()
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// register rest endpoints
	myRouter.HandleFunc("/services", returnPodsPerService)
	myRouter.HandleFunc("/services/{applicationGroup}", returnPodsPerAppGroup)

	// pass in newly created mux router as the second argument
	log.Error(http.ListenAndServe(":9090", myRouter), "Failed to bootstrap the server")
}

// Handles rest endpoint "/services"
func returnPodsPerService(w http.ResponseWriter, r *http.Request) {

	log.Info("Endpoint Hit: returnPodsPerServiceAndAppGroup")

	// encoding our podDetails array into a JSON string and then writing as part of our response.
	json.NewEncoder(w).Encode(controller.PodsPerService(controller.PrepareMap(clientset)))
}

// Handles rest endpoint "/services/{applicationGroup}"
func returnPodsPerAppGroup(w http.ResponseWriter, r *http.Request) {

	log.Info("Endpoint Hit: returnPodsPerAppGroup")

	appGroupId := r.URL.Query().Get("applicationGroup")

	// encoding our podDetails array into a JSON string and then writing as part of our response.
	json.NewEncoder(w).Encode(controller.PodsPerAppGroup(controller.PrepareMap(clientset), appGroupId))
}

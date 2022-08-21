package controller

import (
	"challenge/constants"
	"challenge/pkg/client"
	"context"
	"fmt"

	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

type PodDetails struct {
	Name             string `json:"Name"`
	ApplicationGroup string `json:"ApplicationGroup"`
	RunningPodsCount int    `json:"RunningPodsCount"`
}

var log = logf.Log.WithName("controller")

// Prepares a map with key as "service" and value as "applicationGroup", hold details of all pods in the cluster
func PrepareMap() map[string]string {

	var podMap = map[string]string{}

	clientset, err := client.PrepareClientSet()
	if err != nil {
		log.Error(err, "Unable to create client set")
	}

	// access the API to list deployments
	// Need to access deployment, because label applictionGroup is only part of deployment and not pods
	deploy, err := clientset.AppsV1().Deployments(constants.Namespace).List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Error(err, "Unable to fetch deployments present in the cluster")
	}

	for i, v := range deploy.Items {
		applicationGroup := v.ObjectMeta.Labels["applicationGroup"]
		if len(applicationGroup) == 0 {
			applicationGroup = constants.STR_NONE
		}
		log.Info(fmt.Sprintf("Deployment: %d\t, Name: %v\t, SelectorLabel: %v\t, AppGroup: %v\t  \n", i+1, v.ObjectMeta.Name, v.Spec.Selector.MatchLabels["service"], applicationGroup))
		if selector, exists := podMap[v.Spec.Selector.MatchLabels["service"]]; !exists {
			podMap[v.Spec.Selector.MatchLabels["service"]] = applicationGroup
		} else {
			log.Error(errors.Errorf("Not recommended to have more than one same selector label %v", selector), constants.EMPTY_STR)
			return nil
		}
	}
	return podMap
}

// Returns pods per service
func PodsPerService(podMap map[string]string) []PodDetails {

	var podDetails []PodDetails
	for selector, appGroup := range podMap {
		podDetails = PopulatePodDetailsMap(selector, podDetails, appGroup)
	}
	return podDetails
}

// Returns pods per applicationGroup
func PodsPerAppGroup(podMap map[string]string, appGroupId string) []PodDetails {
	var podDetails []PodDetails
	for selector, appGroup := range podMap {

		if appGroup == appGroupId {
			podDetails = PopulatePodDetailsMap(selector, podDetails, appGroup)
		}
	}
	return podDetails
}

func PopulatePodDetailsMap(selector string, podDetails []PodDetails, appGroup string) []PodDetails {
	// setup list options
	listOptions := v1.ListOptions{
		LabelSelector: constants.LABEL_SELECTOR + selector,
	}

	// access the API to list pods
	clientSet, _ := client.PrepareClientSet()
	pods, err := clientSet.CoreV1().Pods(constants.Namespace).List(context.Background(), listOptions)
	if err != nil {
		log.Error(err, "Unable to fetch pods")
	}
	podDetails = append(podDetails, PodDetails{selector, appGroup, len(pods.Items)})
	return podDetails
}

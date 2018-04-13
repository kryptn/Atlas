package main

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
)

type Client struct {
	namespace string
	clientset *kubernetes.Clientset
	services  *v1.ServiceList
}

func makeClient() *Client {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	client := Client{clientset: clientset}
	return &client
}

func (c *Client) GetServices() {
	// todo: figure out how to determine which namespace we're in
	services, err := c.clientset.CoreV1().Services(c.namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	c.services = services
}

func (c *Client) httpHandler(handler ClientHttpHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	}
}

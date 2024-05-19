package controllers

import (
	"context"
	"fmt"

	webappv1 "assignment/Ingress-Controller/api/v1"
	v1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ingressRouterEval struct {
	Client       client.Client
	routingTable map[string]string
	ingRoute     map[string]string
}

type IngressRouterEval interface {
	RouteNewSimpleIngress(simpleIngress webappv1.SimpleIngress, ctx context.Context) error
	RunNewRoute(route string) (*http.Response, error)
	DeleteRoute(route string)
}

// NewIngressRouterEval creates a new instance of IngressRouterEval.
func NewIngressRouterEval(client client.Client) IngressRouterEval {
	return &ingressRouterEval{
		Client:       client,
		routingTable: make(map[string]string),
		ingRoute:     make(map[string]string),
	}
}

// RouteNewSimpleIngress adds a new route for the given SimpleIngress.
func (e *ingressRouterEval) RouteNewSimpleIngress(simpleIngress webappv1.SimpleIngress, ctx context.Context) error {
	reqSvc := v1.Service{}
	err := e.Client.Get(ctx, client.ObjectKey{Namespace: simpleIngress.Namespace, Name: simpleIngress.Spec.SvcName}, &reqSvc)
	if err != nil {
		return fmt.Errorf("failed to get service %s/%s: %w", simpleIngress.Namespace, simpleIngress.Spec.SvcName, err)
	}

	// Assuming there is at least one port defined
	if len(reqSvc.Spec.Ports) == 0 {
		return fmt.Errorf("no ports found for service %s/%s", reqSvc.Namespace, reqSvc.Name)
	}

	route := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", reqSvc.Name, reqSvc.Namespace, reqSvc.Spec.Ports[0].Port)
	e.routingTable[simpleIngress.Spec.Host] = route
	e.ingRoute[simpleIngress.Name] = simpleIngress.Spec.Host
	return nil
}

// RunNewRoute sends a GET request to the specified route.
func (e *ingressRouterEval) RunNewRoute(route string) (*http.Response, error) {
	if val, ok := e.routingTable[route]; ok {
		return http.Get(val)
	}
	return nil, fmt.Errorf("route %s not found", route)
}

// DeleteRoute removes the route from the routing tables.
func (e *ingressRouterEval) DeleteRoute(route string) {
	if val, ok := e.ingRoute[route]; ok {
		delete(e.routingTable, val)
		delete(e.ingRoute, route)
	}
}

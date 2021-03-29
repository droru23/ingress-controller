package controllers

import (
	webappv1 "assignment/Ingress-Controller/api/v1"
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IngressRouterEval struct {
	Client       client.Client
	RoutingTable map[string]string
}

func NewIngressRouterEval(client client.Client) *IngressRouterEval {
	ingressEval := IngressRouterEval{}
	ingressEval.Client = client
	ingressEval.RoutingTable = map[string]string{}

	return &ingressEval
}

func (e *IngressRouterEval) routeNewSimpleIngress(simpleIngress webappv1.SimpleIngress, ctx context.Context) error {
	reqSvc := v1.Service{}
	err := e.Client.Get(ctx, client.ObjectKey{Namespace: simpleIngress.Namespace, Name: simpleIngress.Spec.SvcName}, &reqSvc)
	if err != nil {
		return fmt.Errorf("service not found with err: %s", err)
	}

	route := fmt.Sprintf("http://%s.%s.svc.cluster.local:%v", reqSvc.Name, reqSvc.Namespace, reqSvc.Spec.Ports[0].Port)
	e.RoutingTable[simpleIngress.Spec.Host] = route
	return nil
}

func (e *IngressRouterEval) RunNewRoute(route string) (*http.Response, error) {
	val, ok := e.RoutingTable[route]
	if ok {
		res, err := http.Get(val)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, fmt.Errorf("not found")
}

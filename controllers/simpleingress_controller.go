/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"assignment/Ingress-Controller/controllers/utils"
	"context"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	webappv1 "assignment/Ingress-Controller/api/v1"
)

// SimpleIngressReconciler reconciles a SimpleIngress object
type SimpleIngressReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	IngressRouter *IngressRouterEval
}

// +kubebuilder:rbac:groups=webapp.my.domain,resources=simpleingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=simpleingresses/status,verbs=get;update;patch

func (r *SimpleIngressReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("simpleingress", req.NamespacedName)
	r.Log.Info("new simple ingress resource asked to be deployed")
	ctx := utils.GenerateNewContext(req)

	simIngress := webappv1.SimpleIngress{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, &simIngress)

	if err != nil {
		r.Log.Info("err getting simpleIngress resource")
		return ctrl.Result{}, nil
	}

	err = r.IngressRouter.routeNewSimpleIngress(simIngress, ctx)
	if err != nil {
		r.Log.Info("error creating new routing rule with err", "err: ", err)
		return ctrl.Result{}, nil
	}
	r.Log.Info("successfully created new routing rule")

	return ctrl.Result{}, nil
}

func (r *SimpleIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.SimpleIngress{}).
		Complete(r)
}

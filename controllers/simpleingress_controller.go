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
	"context"
	"github.com/go-logr/logr"

	webappv1 "assignment/Ingress-Controller/api/v1"
	"assignment/Ingress-Controller/controllers/utils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SimpleIngressReconciler reconciles a SimpleIngress object
type SimpleIngressReconciler struct {
	client.Client
	Log           logr.Logger
	Scheme        *runtime.Scheme
	IngressRouter IngressRouterEval
}

// +kubebuilder:rbac:groups=webapp.my.domain,resources=simpleingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=simpleingresses/status,verbs=get;update;patch

// Reconcile reconciles the SimpleIngress resources
func (r *SimpleIngressReconciler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx := utils.GenerateNewContext(req)
	log := r.Log.WithValues("simpleingress", req.NamespacedName)
	log.Info("Reconciling SimpleIngress resource")

	simIngress := webappv1.SimpleIngress{}
	err := r.Client.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, &simIngress)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("SimpleIngress resource not found. Deleting route if exists.")
			r.IngressRouter.DeleteRoute(req.Name)
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get SimpleIngress resource")
		return ctrl.Result{}, err
	}

	err = r.IngressRouter.RouteNewSimpleIngress(simIngress, ctx)
	if err != nil {
		log.Error(err, "Failed to create new routing rule")
		return ctrl.Result{}, err
	}
	log.Info("Successfully created new routing rule")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SimpleIngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.SimpleIngress{}).
		Complete(r)
}

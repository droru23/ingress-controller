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

package main

import (
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"net/http"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"strings"

	webappv1 "assignment/Ingress-Controller/api/v1"
	"assignment/Ingress-Controller/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = webappv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     false,
		LeaderElectionID:   "056bc37b.my.domain",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// creates the clientset

	ingRouter := controllers.NewIngressRouterEval(mgr.GetClient())

	if err = (&controllers.SimpleIngressReconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName("controllers").WithName("SimpleIngress"),
		Scheme:        mgr.GetScheme(),
		IngressRouter: ingRouter,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "SimpleIngress")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	go func() {
		servErr := http.ListenAndServe(":3000", nil)
		if servErr != nil {
			panic(servErr)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctrl.Log.Info("incoming http request")
		reqHost := strings.Split(r.Host, ".")[0]
		if reqHost == "" {
			fmt.Fprintf(w, "Host not found")
			return
		}
		ctrl.Log.Info("search route: ", "route", reqHost)

		res, err := ingRouter.RunNewRoute(reqHost)
		if err != nil {
			fmt.Fprintf(w, "Host not found")
			return
		}
		err = res.Write(w)
		if err != nil {
			fmt.Fprintf(w, "faild to return answer")
		}

	})

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

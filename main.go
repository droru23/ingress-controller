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
	webappv1 "assignment/Ingress-Controller/api/v1"
	"assignment/Ingress-Controller/controllers"
	"assignment/Ingress-Controller/controllers/router"
	"flag"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
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

	// Parse command-line flags
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	// Set up logger
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// Create a new manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "056bc37b.my.domain",
	})
	if err != nil {
		setupLog.Error(err, "Unable to start manager")
		os.Exit(1)
	}

	// Create the IngressRouterEval
	ingRouter := controllers.NewIngressRouterEval(mgr.GetClient())

	// Set up the SimpleIngressReconciler with the manager
	if err = (&controllers.SimpleIngressReconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName("controllers").WithName("SimpleIngress"),
		Scheme:        mgr.GetScheme(),
		IngressRouter: ingRouter,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "Unable to create controller", "controller", "SimpleIngress")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	// Start a separate goroutine for HTTP server
	handler := router.NewHandler(ingRouter)
	healthRouter := router.SetupHealth(handler)
	ingeressRouter := router.SetupRouter(handler)

	go func() {
		if err := healthRouter.Run(":8081"); err != nil {
			setupLog.Error(err, "Failed to start health server")
			os.Exit(1)
		}
	}()

	go func() {
		if err := ingeressRouter.Run(":8080"); err != nil {
			setupLog.Error(err, "Failed to start ingress server")
			os.Exit(1)
		}
	}()

	setupLog.Info("Starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "Problem running manager")
		os.Exit(1)
	}
}

/*
Copyright 2026.

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

package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	gpuv1alpha1 "github.com/lkhun9311/gpu-node-readiness-operator/api/v1alpha1"
)

// GpuNodeCheckReconciler reconciles a GpuNodeCheck object
type GpuNodeCheckReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gpu.infra.ai,resources=gpunodechecks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gpu.infra.ai,resources=gpunodechecks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gpu.infra.ai,resources=gpunodechecks/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GpuNodeCheck object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.23.3/pkg/reconcile
func (r *GpuNodeCheckReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	var check gpuv1alpha1.GpuNodeCheck
	if err := r.Get(ctx, req.NamespacedName, &check); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var nodeList corev1.NodeList
	if err := r.List(ctx, &nodeList); err != nil {
		return ctrl.Result{}, err
	}

	var total int32
	var ready int32
	var missingLabelNodes []string

	for _, node := range nodeList.Items {
		if !matchLabels(node.Labels, check.Spec.NodeSelector) {
			continue
		}

		total++

		if isNodeReady(node) {
			ready++
		}

		for _, requiredLabel := range check.Spec.RequiredLabels {
			if _, ok := node.Labels[requiredLabel]; !ok {
				missingLabelNodes = append(missingLabelNodes, node.Name)
				break
			}
		}
	}

	now := metav1.Now()
	check.Status.TotalNodes = total
	check.Status.ReadyNodes = ready
	check.Status.MissingLabelNodes = missingLabelNodes
	check.Status.LastCheckedTime = &now

	conditionStatus := metav1.ConditionFalse
	reason := "GpuNodesNotReady"
	message := fmt.Sprintf("ready GPU-like nodes: %d%d", ready, total)

	if total > 0 && ready == total && len(missingLabelNodes) == 0 {
		conditionStatus = metav1.ConditionTrue
		reason = "GpuNodesQualified"
		message = fmt.Sprintf("all selected GPU-like nodes are qualified: %d%d", ready, total)
	}

	meta.SetStatusCondition(&check.Status.Conditions, metav1.Condition{
		Type:               "Qualified",
		Status:             conditionStatus,
		Reason:             reason,
		Message:            message,
		ObservedGeneration: check.Generation,
		LastTransitionTime: now,
	})

	if err := r.Status().Update(ctx, &check); err != nil {
		return ctrl.Result{}, err
	}

	log.Info("reconciled GpuNodeCheck",
		"name", check.Name,
		"totalNodes", total,
		"readyNodes", ready,
		"missingLabelNodes", missingLabelNodes,
	)

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func matchLabels(labels map[string]string, selector map[string]string) bool {
	for key, expected := range selector {
		actual, ok := labels[key]
		if !ok || actual != expected {
			return false
		}
	}
	return true
}

func isNodeReady(node corev1.Node) bool {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *GpuNodeCheckReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpuv1alpha1.GpuNodeCheck{}).
		Named("gpunodecheck").
		Complete(r)
}

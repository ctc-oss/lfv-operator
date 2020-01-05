package datavolume

import (
	"context"
	comv1alpha1 "github.com/jw3/example-operator/pkg/apis/com/v1alpha1"
	batch1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_datavolume")

// Add creates a new DataVolume Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDataVolume{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("datavolume-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource DataVolume
	err = c.Watch(&source.Kind{Type: &comv1alpha1.DataVolume{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &batch1.Job{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &comv1alpha1.DataVolume{},
	})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner DataVolume
	//err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
	//	IsController: true,
	//	OwnerType:    &comv1alpha1.DataVolume{},
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}

// blank assignment to verify that ReconcileDataVolume implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDataVolume{}

// ReconcileDataVolume reconciles a DataVolume object
type ReconcileDataVolume struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a DataVolume object and makes changes based on the state read
// and what is in the DataVolume.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDataVolume) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling DataVolume")

	// Fetch the DataVolume instance
	instance := &comv1alpha1.DataVolume{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new pvc
	pvc := newPvForCR(instance)
	if err := controllerutil.SetControllerReference(instance, pvc, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pvc already exists
	found := &corev1.PersistentVolumeClaim{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pvc.Name, Namespace: pvc.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pvc", "Pvc.Namespace", pvc.Namespace, "Pvc.Name", pvc.Name)
		err = r.client.Create(context.TODO(), pvc)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Define a new cloner job
	cloner := clonerPodForCR(instance)
	if err := controllerutil.SetControllerReference(instance, cloner, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Job already exists
	foundJob := &batch1.Job{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: cloner.Name, Namespace: cloner.Namespace}, foundJob)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Cloner", "Namespace", pvc.Namespace, "Name", pvc.Name)
		err = r.client.Create(context.TODO(), cloner)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// update status
	if !instance.Status.Ready {
		reqLogger.Info("Checking instance status")
		if foundJob.Status.Succeeded > 0 {
			reqLogger.Info("Updated status")
			instance.Status = comv1alpha1.DataVolumeStatus{Ready: true}
			err = r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	reqLogger.Info("Skip reconcile: Pvc already exists", "Pvc.Namespace", found.Namespace, "Pvc.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPvForCR returns a new Pvc
func newPvForCR(cr *comv1alpha1.DataVolume) *corev1.PersistentVolumeClaim {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: cr.Spec.Size,
				},
			},
		},
	}
}

func clonerPodForCR(cr *comv1alpha1.DataVolume) *batch1.Job {
	var jobTTL int32 = 90
	labels := map[string]string{
		"app": cr.Name,
	}
	return &batch1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-cloner",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batch1.JobSpec{
			TTLSecondsAfterFinished: &jobTTL,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Volumes: []corev1.Volume{
						{
							Name: cr.Name,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: cr.Name,
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:    "cloner",
							Image:   "centos/python-36-centos7",
							Command: []string{"git", "clone", "-v", cr.Spec.Uri, "/opt/app-root/mount"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      cr.Name,
									MountPath: "/opt/app-root/mount",
								},
							},
						},
					},
				},
			},
		},
	}
}

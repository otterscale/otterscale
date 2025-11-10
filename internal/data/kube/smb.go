package kube

import (
	"context"
	"fmt"
	"math"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	oscore "github.com/otterscale/otterscale/internal/core"
)

type smb struct {
	kube *Kube
}

func NewSMB(kube *Kube) oscore.KubeSMBRepo {
	return &smb{
		kube: kube,
	}
}

var _ oscore.KubeSMBRepo = (*smb)(nil)

const (
	smbGroup   = "samba-operator.samba.org"
	smbVersion = "v1alpha1"
	guestOkYes = "yes"
	guestOkNo  = "no"
)

var (
	smbShareGVR = schema.GroupVersionResource{
		Group:    smbGroup,
		Version:  smbVersion,
		Resource: "smbshares",
	}

	smbCommonConfigGVR = schema.GroupVersionResource{
		Group:    smbGroup,
		Version:  smbVersion,
		Resource: "smbcommonconfigs",
	}

	smbSecurityConfigGVR = schema.GroupVersionResource{
		Group:    smbGroup,
		Version:  smbVersion,
		Resource: "smbsecurityconfigs",
	}
)

func (r *smb) ListSMBShares(ctx context.Context, config *rest.Config, namespace string) ([]oscore.SMBShare, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	opts := metav1.ListOptions{}
	unstructuredList, err := dynamicClient.Resource(smbShareGVR).Namespace(namespace).List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list SmbShares: %w", err)
	}

	shares := make([]oscore.SMBShare, 0, len(unstructuredList.Items))
	for _, item := range unstructuredList.Items {
		share := r.convertUnstructuredToSMBShare(ctx, dynamicClient, clientset, namespace, &item)
		shares = append(shares, share)
	}

	return shares, nil
}

func (r *smb) GetSMBShare(ctx context.Context, config *rest.Config, namespace, name string) (*oscore.SMBShare, error) {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Get clientset for deployment access
	clientset, err := r.kube.clientset(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	opts := metav1.GetOptions{}
	unstructuredObj, err := dynamicClient.Resource(smbShareGVR).Namespace(namespace).Get(ctx, name, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get SmbShare: %w", err)
	}

	share := r.convertUnstructuredToSMBShare(ctx, dynamicClient, clientset, namespace, unstructuredObj)

	return &share, nil
}

func (r *smb) UpdateSMBCommonConfig(ctx context.Context, config *rest.Config, namespace, name string, shareConfig *oscore.SMBShareConfig) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Get SmbCommonConfig
	getOpts := metav1.GetOptions{}
	commonConfig, err := dynamicClient.Resource(smbCommonConfigGVR).Namespace(namespace).Get(ctx, name, getOpts)
	if err != nil {
		return fmt.Errorf("failed to get SmbCommonConfig: %w", err)
	}

	spec := commonConfig.Object["spec"].(map[string]interface{})
	if shareConfig.MapToGuest != "" {
		spec["customGlobalConfig"] = map[string]interface{}{
			"useUnsafeCustomConfig": true,
			"configs": map[string]string{
				"map to guest": shareConfig.MapToGuest,
			},
		}
	} else {
		delete(spec, "customGlobalConfig")
	}

	// Update the resource
	updateOpts := metav1.UpdateOptions{}
	_, err = dynamicClient.Resource(smbCommonConfigGVR).Namespace(namespace).Update(ctx, commonConfig, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update SmbCommonConfig CRD: %w", err)
	}

	return nil
}

func (r *smb) DeleteSMBCommonConfig(ctx context.Context, config *rest.Config, namespace, name string) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	deleteOpts := metav1.DeleteOptions{}
	if err := dynamicClient.Resource(smbCommonConfigGVR).Namespace(namespace).Delete(ctx, name, deleteOpts); err != nil {
		return fmt.Errorf("failed to delete SmbCommonConfig: %w", err)
	}

	return nil
}

func (r *smb) DeleteSMBSecurityConfig(ctx context.Context, config *rest.Config, namespace, name string) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	deleteOpts := metav1.DeleteOptions{}
	if err := dynamicClient.Resource(smbSecurityConfigGVR).Namespace(namespace).Delete(ctx, name, deleteOpts); err != nil {
		return fmt.Errorf("failed to delete SmbSecurityConfig: %w", err)
	}

	return nil
}

func (r *smb) UpdateSMBSecurityConfig(ctx context.Context, config *rest.Config, namespace, name, mode, realm, secretName string) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Get SmbSecurityConfig
	getOpts := metav1.GetOptions{}
	securityConfig, err := dynamicClient.Resource(smbSecurityConfigGVR).Namespace(namespace).Get(ctx, name, getOpts)
	if err != nil {
		return fmt.Errorf("failed to get SmbSecurityConfig: %w", err)
	}

	spec := securityConfig.Object["spec"].(map[string]interface{})
	spec["mode"] = mode

	switch mode {
	case "active-directory":
		spec["realm"] = realm
		if secretName != "" {
			spec["joinSources"] = []interface{}{
				map[string]interface{}{
					"userJoin": map[string]interface{}{
						"secret": secretName,
						"key":    "join",
					},
				},
			}
		}
		delete(spec, "users")
	case "user":
		spec["users"] = map[string]interface{}{
			"secret": secretName,
			"key":    "users",
		}
		delete(spec, "realm")
		delete(spec, "joinSources")
	}

	// Update the resource
	updateOpts := metav1.UpdateOptions{}
	_, err = dynamicClient.Resource(smbSecurityConfigGVR).Namespace(namespace).Update(ctx, securityConfig, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update SmbSecurityConfig CRD: %w", err)
	}

	return nil
}

func (r *smb) updatePVC(ctx context.Context, clientset *kubernetes.Clientset, namespace, pvcName string, sizeBytes uint64) error {
	opts := metav1.GetOptions{}
	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(ctx, pvcName, opts)
	if err != nil {
		return fmt.Errorf("failed to get PVC %s: %w", pvcName, err)
	}

	if pvc.Spec.Resources.Requests == nil {
		pvc.Spec.Resources.Requests = corev1.ResourceList{}
	}

	if sizeBytes > math.MaxInt64 {
		return fmt.Errorf("size too large: %d bytes exceeds maximum allowed", sizeBytes)
	}
	pvc.Spec.Resources.Requests[corev1.ResourceStorage] = *resource.NewQuantity(int64(sizeBytes), resource.BinarySI)

	updateOpts := metav1.UpdateOptions{}
	_, err = clientset.CoreV1().PersistentVolumeClaims(namespace).Update(ctx, pvc, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update PVC %s: %w", pvcName, err)
	}
	return nil
}

func (r *smb) updateSMBShareSpecStorage(spec map[string]interface{}, sizeBytes uint64) error {
	storage, ok := spec["storage"].(map[string]interface{})
	if !ok {
		storage = make(map[string]interface{})
		spec["storage"] = storage
	}

	pvcMap, ok := storage["pvc"].(map[string]interface{})
	if !ok {
		pvcMap = make(map[string]interface{})
		storage["pvc"] = pvcMap
	}

	pvcSpec, ok := pvcMap["spec"].(map[string]interface{})
	if !ok {
		pvcSpec = make(map[string]interface{})
		pvcMap["spec"] = pvcSpec
	}

	resources, ok := pvcSpec["resources"].(map[string]interface{})
	if !ok {
		resources = make(map[string]interface{})
		pvcSpec["resources"] = resources
	}

	requests, ok := resources["requests"].(map[string]interface{})
	if !ok {
		requests = make(map[string]interface{})
		resources["requests"] = requests
	}

	if sizeBytes > math.MaxInt64 {
		return fmt.Errorf("size too large: %d bytes exceeds maximum allowed", sizeBytes)
	}
	requests["storage"] = resource.NewQuantity(int64(sizeBytes), resource.BinarySI).String()
	return nil
}

func (r *smb) UpdateSMBShare(ctx context.Context, config *rest.Config, namespace, name string, sizeBytes uint64, browseable, readOnly, guestOk bool, validUsers string) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	clientset, err := r.kube.clientset(config)
	if err != nil {
		return fmt.Errorf("failed to create clientset: %w", err)
	}

	opts := metav1.GetOptions{}
	smbShare, err := dynamicClient.Resource(smbShareGVR).Namespace(namespace).Get(ctx, name, opts)
	if err != nil {
		return fmt.Errorf("failed to get SmbShare: %w", err)
	}

	spec := smbShare.Object["spec"].(map[string]interface{})
	spec["browseable"] = browseable
	spec["readOnly"] = readOnly

	if sizeBytes > 0 {
		// Get PVC name and update it
		pvcName := name
		if storage, found, _ := unstructured.NestedMap(smbShare.Object, "spec", "storage", "pvc"); found {
			if n, ok := storage["name"].(string); ok {
				pvcName = n
			}
		}

		if err := r.updatePVC(ctx, clientset, namespace, pvcName, sizeBytes); err != nil {
			return err
		}

		if err := r.updateSMBShareSpecStorage(spec, sizeBytes); err != nil {
			return err
		}
	}

	// Get existing customShareConfig
	existingConfigs := make(map[string]string)
	if configs, found, _ := unstructured.NestedStringMap(smbShare.Object, "spec", "customShareConfig", "configs"); found {
		existingConfigs = configs
	}

	if guestOk {
		existingConfigs["guest ok"] = guestOkYes
	} else {
		existingConfigs["guest ok"] = guestOkNo
	}

	if validUsers != "" {
		existingConfigs["valid users"] = validUsers
	} else {
		delete(existingConfigs, "valid users")
	}

	if len(existingConfigs) > 0 {
		spec["customShareConfig"] = map[string]interface{}{
			"useUnsafeCustomConfig": true,
			"configs":               existingConfigs,
		}
	} else {
		delete(spec, "customShareConfig")
	}

	// Update the SmbShare resource
	updateOpts := metav1.UpdateOptions{}
	_, err = dynamicClient.Resource(smbShareGVR).Namespace(namespace).Update(ctx, smbShare, updateOpts)
	if err != nil {
		return fmt.Errorf("failed to update SmbShare CRD: %w", err)
	}

	return nil
}

func (r *smb) CreateSMBCommonConfig(ctx context.Context, config *rest.Config, namespace, name string, shareConfig *oscore.SMBShareConfig) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	commonConfig := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", smbGroup, smbVersion),
			"kind":       "SmbCommonConfig",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{},
		},
	}

	if shareConfig.MapToGuest != "" {
		spec := commonConfig.Object["spec"].(map[string]interface{})
		spec["customGlobalConfig"] = map[string]interface{}{
			"useUnsafeCustomConfig": true,
			"configs": map[string]string{
				"map to guest": shareConfig.MapToGuest,
			},
		}
	}

	opts := metav1.CreateOptions{}
	_, err = dynamicClient.Resource(smbCommonConfigGVR).Namespace(namespace).Create(ctx, commonConfig, opts)
	if err != nil {
		return fmt.Errorf("failed to create SmbCommonConfig CRD: %w", err)
	}

	return nil
}

func (r *smb) CreateSMBSecurityConfig(ctx context.Context, config *rest.Config, namespace, name, mode, realm, secretName string) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	securityConfig := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", smbGroup, smbVersion),
			"kind":       "SmbSecurityConfig",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"dns": map[string]interface{}{
					"register": "cluster-ip",
				},
				"mode": mode,
			},
		},
	}

	spec := securityConfig.Object["spec"].(map[string]interface{})

	switch mode {
	case "active-directory":
		spec["realm"] = realm
		if secretName != "" {
			spec["joinSources"] = []interface{}{
				map[string]interface{}{
					"userJoin": map[string]interface{}{
						"secret": secretName,
						"key":    "join",
					},
				},
			}
		}
	case "user":
		spec["users"] = map[string]interface{}{
			"secret": secretName,
			"key":    "users",
		}
	}

	opts := metav1.CreateOptions{}
	_, err = dynamicClient.Resource(smbSecurityConfigGVR).Namespace(namespace).Create(ctx, securityConfig, opts)
	if err != nil {
		return fmt.Errorf("failed to create SmbSecurityConfig CRD: %w", err)
	}

	return nil
}

func (r *smb) CreateSMBShare(ctx context.Context, config *rest.Config, namespace, name string, sizeBytes uint64, browseable, readOnly, guestOk bool, validUsers, commonConfig, securityConfig, _ string) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create dynamic client: %w", err)
	}

	if sizeBytes > math.MaxInt64 {
		return fmt.Errorf("size too large: %d bytes exceeds maximum allowed", sizeBytes)
	}
	storageQuantity := resource.NewQuantity(int64(sizeBytes), resource.BinarySI).String()

	smbShare := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": fmt.Sprintf("%s/%s", smbGroup, smbVersion),
			"kind":       "SmbShare",
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"shareName":  name,
				"browseable": browseable,
				"readOnly":   readOnly,
				"storage": map[string]interface{}{
					"pvc": map[string]interface{}{
						"spec": map[string]interface{}{
							"accessModes": []interface{}{corev1.ReadWriteMany},
							"resources": map[string]interface{}{
								"requests": map[string]interface{}{
									"storage": storageQuantity,
								},
							},
							"storageClassName": "cephfs",
						},
					},
				},
				"commonConfig": commonConfig,
			},
		},
	}

	spec := smbShare.Object["spec"].(map[string]interface{})

	customShareConfigs := make(map[string]string)
	if guestOk {
		customShareConfigs["guest ok"] = guestOkYes
	}
	if validUsers != "" {
		customShareConfigs["valid users"] = validUsers
	}
	if len(customShareConfigs) > 0 {
		spec["customShareConfig"] = map[string]interface{}{
			"useUnsafeCustomConfig": true,
			"configs":               customShareConfigs,
		}
	}

	if securityConfig != "" {
		spec["securityConfig"] = securityConfig
	}

	opts := metav1.CreateOptions{}
	_, err = dynamicClient.Resource(smbShareGVR).Namespace(namespace).Create(ctx, smbShare, opts)
	if err != nil {
		return fmt.Errorf("failed to create SmbShare CRD: %w", err)
	}

	return nil
}

func (r *smb) convertUnstructuredToSMBShare(ctx context.Context, dynamicClient dynamic.Interface, clientset *kubernetes.Clientset, namespace string, obj *unstructured.Unstructured) oscore.SMBShare {
	share := oscore.SMBShare{
		Name:       obj.GetName(),
		Namespace:  namespace,
		Browseable: true,
		ReadOnly:   true,
		Status:     "Progressing",
		CreatedAt:  obj.GetCreationTimestamp().Time,
	}

	// Get deployment status to determine share status
	deploymentName := obj.GetName()
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err == nil {
		if deployment.Status.ReadyReplicas > 0 && deployment.Status.ReadyReplicas == deployment.Status.Replicas {
			share.Status = "Ready"
		} else {
			share.Status = "Progressing"
		}
	}

	if storage, found, _ := unstructured.NestedString(obj.Object, "spec", "storage", "pvc", "spec", "resources", "requests", "storage"); found {
		if quantity, err := resource.ParseQuantity(storage); err == nil {
			value := quantity.Value()
			if value >= 0 {
				share.SizeBytes = uint64(value)
			}
		}
	}

	if browsable, found, _ := unstructured.NestedBool(obj.Object, "spec", "browseable"); found {
		share.Browseable = browsable
	}
	if readOnly, found, _ := unstructured.NestedBool(obj.Object, "spec", "readOnly"); found {
		share.ReadOnly = readOnly
	}

	if configs, found, _ := unstructured.NestedStringMap(obj.Object, "spec", "customShareConfig", "configs"); found {
		if guestOk, ok := configs["guest ok"]; ok {
			share.GuestOk = (guestOk == "yes")
		}
		if validUsers, ok := configs["valid users"]; ok {
			share.ValidUsers = strings.Split(validUsers, " ")
		}
	}

	commonConfigName, _, _ := unstructured.NestedString(obj.Object, "spec", "commonConfig")
	if commonConfigName != "" {
		commonConfigObj, err := dynamicClient.Resource(smbCommonConfigGVR).Namespace(namespace).Get(ctx, commonConfigName, metav1.GetOptions{})
		if err == nil {
			if configs, found, _ := unstructured.NestedStringMap(commonConfigObj.Object, "spec", "customGlobalConfig", "configs"); found {
				if mapToGuest, ok := configs["map to guest"]; ok {
					share.MapToGuest = mapToGuest
				}
			}
		}
	}

	// Get securityConfig name
	securityConfigName, _, _ := unstructured.NestedString(obj.Object, "spec", "securityConfig")
	if securityConfigName != "" {
		securityConfigObj, err := dynamicClient.Resource(smbSecurityConfigGVR).Namespace(namespace).Get(ctx, securityConfigName, metav1.GetOptions{})
		if err == nil {
			if mode, found, _ := unstructured.NestedString(securityConfigObj.Object, "spec", "mode"); found {
				share.SecurityMode = mode
			}
			if realm, found, _ := unstructured.NestedString(securityConfigObj.Object, "spec", "realm"); found {
				share.Realm = realm
				share.ADAuth = &oscore.ADAuth{
					Realm: realm,
				}
			}
		}
	}

	return share
}

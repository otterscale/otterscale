package smb

import (
	"encoding/json"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/samba-in-kubernetes/samba-operator/api/v1alpha1"

	"github.com/otterscale/otterscale/internal/core/application/config"
	"github.com/otterscale/otterscale/internal/core/application/persistent"
	"github.com/otterscale/otterscale/internal/core/application/service"
)

func (uc *UseCase) buildUsersSecret(namespace, name string, users []User) *config.Secret {
	userEntries := []userEntry{}

	for _, user := range users {
		userEntries = append(userEntries, userEntry{
			Name:     user.Username,
			Password: user.Password,
		})
	}

	data, _ := json.Marshal(&sambaContainerConfig{
		SCCVersion: "v0",
		Users: map[string][]userEntry{
			"all_entries": userEntries,
		},
	})

	return &config.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name": "samba",
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"users": data,
		},
	}
}

func (uc *UseCase) buildJoinSecret(namespace, name string, user *User) *config.Secret {
	data, _ := json.Marshal(&User{
		Username: user.Username,
		Password: user.Password,
	})

	return &config.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name": "samba",
			},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"join": data,
		},
	}
}

func (uc *UseCase) buildCommonConfig(namespace, name, mapToGuest string) *CommonConfig {
	return &CommonConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.SmbCommonConfigSpec{
			CustomGlobalConfig: &v1alpha1.SmbCommonConfigGlobalConfig{
				UseUnsafeCustomConfig: true,
				Configs: map[string]string{
					MapToGuestKey: mapToGuest,
				},
			},
		},
	}
}

func (uc *UseCase) buildSecurityConfig(namespace, name, securityMode, usersSecretName, realm, joinSecretName string) *SecurityConfig {
	return &SecurityConfig{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.SmbSecurityConfigSpec{
			Mode: securityMode,
			Users: &v1alpha1.SmbSecurityUsersSpec{
				Secret: usersSecretName,
				Key:    "users",
			},
			DNS: &v1alpha1.SmbSecurityDNSSpec{
				Register: "cluster-ip",
			},
			Realm: realm,
			JoinSources: []v1alpha1.SmbSecurityJoinSpec{
				{
					UserJoin: &v1alpha1.SmbSecurityUserJoinSpec{
						Secret: joinSecretName,
						Key:    "join",
					},
				},
			},
		},
	}
}

func (uc *UseCase) buildShareConfig(guestOK bool, validUsers []string) *v1alpha1.SmbShareConfig {
	configs := map[string]string{}

	if guestOK {
		configs[GuestOKkey] = "yes"
	}

	if len(validUsers) > 0 {
		configs[ValidUsersKey] = strings.Join(validUsers, " ")
	}

	return &v1alpha1.SmbShareConfig{
		UseUnsafeCustomConfig: true,
		Configs:               configs,
	}
}

func (uc *UseCase) buildShare(namespace, name string, browsable, readOnly bool, sizeBytes uint64, securityConfigName, commonConfigName string, customShareConfig *v1alpha1.SmbShareConfig) *v1alpha1.SmbShare {
	pvc := uc.buildPersistentVolumeClaim(namespace, name, sizeBytes)

	return &Share{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1alpha1.SmbShareSpec{
			ShareName:  name,
			Browseable: browsable,
			ReadOnly:   readOnly,
			Storage: v1alpha1.SmbShareStorageSpec{
				Pvc: &v1alpha1.SmbSharePvcSpec{
					Spec: &pvc.Spec,
				},
			},
			SecurityConfig:    securityConfigName,
			CommonConfig:      commonConfigName,
			CustomShareConfig: customShareConfig,
		},
	}
}

func (uc *UseCase) buildPersistentVolumeClaim(namespace, name string, sizeBytes uint64) *persistent.PersistentVolumeClaim {
	storageClassName := "cephfs"

	return &persistent.PersistentVolumeClaim{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": *resource.NewQuantity(int64(sizeBytes), resource.BinarySI), //nolint:gosec // ignore
				},
			},
			StorageClassName: &storageClassName,
		},
	}
}

func (uc *UseCase) buildService(namespace, name string, nodePort int32) *service.Service {
	return &service.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:     "smb",
					Port:     445,
					NodePort: nodePort,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
}

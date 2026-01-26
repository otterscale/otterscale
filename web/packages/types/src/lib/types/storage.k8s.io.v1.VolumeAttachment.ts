/** Generated from Remote JSON Schema for storage.k8s.io.v1.VolumeAttachment */

/**
 * VolumeAttachment captures the intent to attach or detach the specified volume to/from the specified node.
 *
 * VolumeAttachment objects are non-namespaced.
 */
export interface StorageK8SIoV1VolumeAttachment {
  /**
   * APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
   */
  apiVersion?: string;
  /**
   * Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
   */
  kind?: string;
  /**
   * ObjectMeta is metadata that all persisted resources must have, which includes all objects users must create.
   */
  metadata?: {
    /**
     * An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.
     *
     * Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
     */
    resourceVersion?: string;
    /**
     * Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
     */
    annotations?: {
      [k: string]: string;
    };
    /**
     * Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
     */
    labels?: {
      [k: string]: string;
    };
    /**
     * ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field. A workflow can be the user's name, a controller's name, or the name of a specific apply path like "ci-cd". The set of fields is always in the version that the workflow used when modifying the object.
     */
    managedFields?: {
      /**
       * Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
       */
      subresource?: string;
      /**
       * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
       */
      time?: string;
      /**
       * APIVersion defines the version of this resource that this field set applies to. The format is "group/version" just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
       */
      apiVersion?: string;
      /**
       * FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: "FieldsV1"
       */
      fieldsType?: string;
      /**
       * FieldsV1 stores a set of fields in a data structure like a Trie, in JSON format.
       *
       * Each key is either a '.' representing the field itself, and will always map to an empty set, or a string representing a sub-field or item. The string will follow one of these four formats: 'f:<name>', where <name> is the name of a field in a struct, or key in a map 'v:<value>', where <value> is the exact json formatted value of a list item 'i:<index>', where <index> is position of a item in a list 'k:<keys>', where <keys> is a map of  a list item's key fields to their unique values If a key maps to an empty Fields value, the field that key represents is part of the set.
       *
       * The exact format is defined in sigs.k8s.io/structured-merge-diff
       */
      fieldsV1?: {
        [k: string]: unknown;
      };
      /**
       * Manager is an identifier of the workflow managing these fields.
       */
      manager?: string;
      /**
       * Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are 'Apply' and 'Update'.
       */
      operation?: string;
      [k: string]: unknown;
    }[];
    /**
     * Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.
     */
    selfLink?: string;
    /**
     * Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
     */
    deletionGracePeriodSeconds?: number;
    /**
     * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
     */
    deletionTimestamp?: string;
    /**
     * A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.
     */
    generation?: number;
    /**
     * Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
     */
    name?: string;
    /**
     * List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller.
     */
    ownerReferences?: {
      /**
       * API version of the referent.
       */
      apiVersion: string;
      /**
       * If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. See https://kubernetes.io/docs/concepts/architecture/garbage-collection/#foreground-deletion for how the garbage collector interacts with this field and enforces the foreground deletion. Defaults to false. To set this field, a user needs "delete" permission of the owner, otherwise 422 (Unprocessable Entity) will be returned.
       */
      blockOwnerDeletion?: boolean;
      /**
       * If true, this reference points to the managing controller.
       */
      controller?: boolean;
      /**
       * Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
       */
      kind: string;
      /**
       * Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
       */
      name: string;
      /**
       * UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
       */
      uid: string;
      [k: string]: unknown;
    }[];
    /**
     * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
     */
    creationTimestamp?: string;
    /**
     * Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order.  Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list.
     */
    finalizers?: string[];
    /**
     * GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
     *
     * If this field is specified and the generated name exists, the server will return a 409.
     *
     * Applied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency
     */
    generateName?: string;
    /**
     * Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
     *
     * Must be a DNS_LABEL. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces
     */
    namespace?: string;
    /**
     * UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.
     *
     * Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
     */
    uid?: string;
    [k: string]: unknown;
  };
  /**
   * VolumeAttachmentSpec is the specification of a VolumeAttachment request.
   */
  spec: {
    /**
     * attacher indicates the name of the volume driver that MUST handle this request. This is the name returned by GetPluginName().
     */
    attacher: string;
    /**
     * nodeName represents the node that the volume should be attached to.
     */
    nodeName: string;
    /**
     * VolumeAttachmentSource represents a volume that should be attached. Right now only PersistentVolumes can be attached via external attacher, in the future we may allow also inline volumes in pods. Exactly one member can be set.
     */
    source: {
      /**
       * PersistentVolumeSpec is the specification of a persistent volume.
       */
      inlineVolumeSpec?: {
        /**
         * storageClassName is the name of StorageClass to which this persistent volume belongs. Empty value means that this volume does not belong to any StorageClass.
         */
        storageClassName?: string;
        /**
         * Represents a vSphere volume resource.
         */
        vsphereVolume?: {
          /**
           * storagePolicyName is the storage Policy Based Management (SPBM) profile name.
           */
          storagePolicyName?: string;
          /**
           * volumePath is the path that identifies vSphere volume vmdk
           */
          volumePath: string;
          /**
           * fsType is filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
           */
          fsType?: string;
          /**
           * storagePolicyID is the storage Policy Based Management (SPBM) profile ID associated with the StoragePolicyName.
           */
          storagePolicyID?: string;
          [k: string]: unknown;
        };
        /**
         * ObjectReference contains enough information to let you inspect or modify the referred object.
         */
        claimRef?: {
          /**
           * Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
           */
          namespace?: string;
          /**
           * Specific resourceVersion to which this reference is made, if any. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
           */
          resourceVersion?: string;
          /**
           * UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids
           */
          uid?: string;
          /**
           * API version of the referent.
           */
          apiVersion?: string;
          /**
           * If referring to a piece of an object instead of an entire object, this string should contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2]. For example, if the object reference is to a container within a pod, this would take on a value like: "spec.containers{name}" (where "name" refers to the name of the container that triggered the event) or if no container name is specified "spec.containers[2]" (container with index 2 in this pod). This syntax is chosen only to have some well-defined way of referencing a part of an object.
           */
          fieldPath?: string;
          /**
           * Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
           */
          kind?: string;
          /**
           * Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
           */
          name?: string;
          [k: string]: unknown;
        };
        /**
         * Represents a Flocker volume mounted by the Flocker agent. One and only one of datasetName and datasetUUID should be set. Flocker volumes do not support ownership management or SELinux relabeling.
         */
        flocker?: {
          /**
           * datasetName is Name of the dataset stored as metadata -> name on the dataset for Flocker should be considered as deprecated
           */
          datasetName?: string;
          /**
           * datasetUUID is the UUID of the dataset. This is unique identifier of a Flocker dataset
           */
          datasetUUID?: string;
          [k: string]: unknown;
        };
        /**
         * ISCSIPersistentVolumeSource represents an ISCSI disk. ISCSI volumes can only be mounted as read/write once. ISCSI volumes support ownership management and SELinux relabeling.
         */
        iscsi?: {
          /**
           * iscsiInterface is the interface Name that uses an iSCSI transport. Defaults to 'default' (tcp).
           */
          iscsiInterface?: string;
          /**
           * lun is iSCSI Target Lun number.
           */
          lun: number;
          /**
           * readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false.
           */
          readOnly?: boolean;
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          secretRef?: {
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            [k: string]: unknown;
          };
          /**
           * fsType is the filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#iscsi
           */
          fsType?: string;
          /**
           * initiatorName is the custom iSCSI Initiator Name. If initiatorName is specified with iscsiInterface simultaneously, new iSCSI interface <target portal>:<volume name> will be created for the connection.
           */
          initiatorName?: string;
          /**
           * portals is the iSCSI Target Portal List. The Portal is either an IP or ip_addr:port if the port is other than default (typically TCP ports 860 and 3260).
           */
          portals?: string[];
          /**
           * targetPortal is iSCSI Target Portal. The Portal is either an IP or ip_addr:port if the port is other than default (typically TCP ports 860 and 3260).
           */
          targetPortal: string;
          /**
           * chapAuthDiscovery defines whether support iSCSI Discovery CHAP authentication
           */
          chapAuthDiscovery?: boolean;
          /**
           * chapAuthSession defines whether support iSCSI Session CHAP authentication
           */
          chapAuthSession?: boolean;
          /**
           * iqn is Target iSCSI Qualified Name.
           */
          iqn: string;
          [k: string]: unknown;
        };
        /**
         * Represents an NFS mount that lasts the lifetime of a pod. NFS volumes do not support ownership management or SELinux relabeling.
         */
        nfs?: {
          /**
           * readOnly here will force the NFS export to be mounted with read-only permissions. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
           */
          readOnly?: boolean;
          /**
           * server is the hostname or IP address of the NFS server. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
           */
          server: string;
          /**
           * path that is exported by the NFS server. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
           */
          path: string;
          [k: string]: unknown;
        };
        /**
         * Represents a Persistent Disk resource in Google Compute Engine.
         *
         * A GCE PD must exist before mounting to a container. The disk must also be in the same GCE project and zone as the kubelet. A GCE PD can only be mounted as read/write once or read-only many times. GCE PDs support ownership management and SELinux relabeling.
         */
        gcePersistentDisk?: {
          /**
           * fsType is filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
           */
          fsType?: string;
          /**
           * partition is the partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as "1". Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty). More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
           */
          partition?: number;
          /**
           * pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
           */
          pdName: string;
          /**
           * readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
           */
          readOnly?: boolean;
          [k: string]: unknown;
        };
        /**
         * Local represents directly-attached storage with node affinity
         */
        local?: {
          /**
           * fsType is the filesystem type to mount. It applies only when the Path is a block device. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". The default value is to auto-select a filesystem if unspecified.
           */
          fsType?: string;
          /**
           * path of the full path to the volume on the node. It can be either a directory or block device (disk, partition, ...).
           */
          path: string;
          [k: string]: unknown;
        };
        /**
         * mountOptions is the list of mount options, e.g. ["ro", "soft"]. Not validated - mount will simply fail if one is invalid. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#mount-options
         */
        mountOptions?: string[];
        /**
         * persistentVolumeReclaimPolicy defines what happens to a persistent volume when released from its claim. Valid options are Retain (default for manually created PersistentVolumes), Delete (default for dynamically provisioned PersistentVolumes), and Recycle (deprecated). Recycle must be supported by the volume plugin underlying this PersistentVolume. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#reclaiming
         *
         * Possible enum values:
         *  - `"Delete"` means the volume will be deleted from Kubernetes on release from its claim. The volume plugin must support Deletion.
         *  - `"Recycle"` means the volume will be recycled back into the pool of unbound persistent volumes on release from its claim. The volume plugin must support Recycling.
         *  - `"Retain"` means the volume will be left in its current phase (Released) for manual reclamation by the administrator. The default policy is Retain.
         */
        persistentVolumeReclaimPolicy?: 'Delete' | 'Recycle' | 'Retain';
        /**
         * Represents a Photon Controller persistent disk resource.
         */
        photonPersistentDisk?: {
          /**
           * fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
           */
          fsType?: string;
          /**
           * pdID is the ID that identifies Photon Controller persistent disk
           */
          pdID: string;
          [k: string]: unknown;
        };
        /**
         * Represents a Rados Block Device mount that lasts the lifetime of a pod. RBD volumes support ownership management and SELinux relabeling.
         */
        rbd?: {
          /**
           * keyring is the path to key ring for RBDUser. Default is /etc/ceph/keyring. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
           */
          keyring?: string;
          /**
           * monitors is a collection of Ceph monitors. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
           */
          monitors: string[];
          /**
           * pool is the rados pool name. Default is rbd. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
           */
          pool?: string;
          /**
           * readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
           */
          readOnly?: boolean;
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          secretRef?: {
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            [k: string]: unknown;
          };
          /**
           * user is the rados user name. Default is admin. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
           */
          user?: string;
          /**
           * fsType is the filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#rbd
           */
          fsType?: string;
          /**
           * image is the rados image name. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
           */
          image: string;
          [k: string]: unknown;
        };
        /**
         * Represents a Fibre Channel volume. Fibre Channel volumes can only be mounted as read/write once. Fibre Channel volumes support ownership management and SELinux relabeling.
         */
        fc?: {
          /**
           * wwids Optional: FC volume world wide identifiers (wwids) Either wwids or combination of targetWWNs and lun must be set, but not both simultaneously.
           */
          wwids?: string[];
          /**
           * fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
           */
          fsType?: string;
          /**
           * lun is Optional: FC target lun number
           */
          lun?: number;
          /**
           * readOnly is Optional: Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
           */
          readOnly?: boolean;
          /**
           * targetWWNs is Optional: FC target worldwide names (WWNs)
           */
          targetWWNs?: string[];
          [k: string]: unknown;
        };
        /**
         * Represents a Persistent Disk resource in AWS.
         *
         * An AWS EBS disk must exist before mounting to a container. The disk must also be in the same AWS zone as the kubelet. An AWS EBS disk can only be mounted as read/write once. AWS EBS volumes support ownership management and SELinux relabeling.
         */
        awsElasticBlockStore?: {
          /**
           * partition is the partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as "1". Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty).
           */
          partition?: number;
          /**
           * readOnly value true will force the readOnly setting in VolumeMounts. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
           */
          readOnly?: boolean;
          /**
           * volumeID is unique ID of the persistent disk resource in AWS (Amazon EBS volume). More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
           */
          volumeID: string;
          /**
           * fsType is the filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
           */
          fsType?: string;
          [k: string]: unknown;
        };
        /**
         * AzureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
         */
        azureDisk?: {
          /**
           * cachingMode is the Host Caching mode: None, Read Only, Read Write.
           *
           * Possible enum values:
           *  - `"None"`
           *  - `"ReadOnly"`
           *  - `"ReadWrite"`
           */
          cachingMode?: 'None' | 'ReadOnly' | 'ReadWrite';
          /**
           * diskName is the Name of the data disk in the blob storage
           */
          diskName: string;
          /**
           * diskURI is the URI of data disk in the blob storage
           */
          diskURI: string;
          /**
           * fsType is Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
           */
          fsType?: string;
          /**
           * kind expected values are Shared: multiple blob disks per storage account  Dedicated: single blob disk per storage account  Managed: azure managed data disk (only in managed availability set). defaults to shared
           *
           * Possible enum values:
           *  - `"Dedicated"`
           *  - `"Managed"`
           *  - `"Shared"`
           */
          kind?: 'Dedicated' | 'Managed' | 'Shared';
          /**
           * readOnly Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
           */
          readOnly?: boolean;
          [k: string]: unknown;
        };
        /**
         * FlexPersistentVolumeSource represents a generic persistent volume resource that is provisioned/attached using an exec based plugin.
         */
        flexVolume?: {
          /**
           * options is Optional: this field holds extra command options if any.
           */
          options?: {
            [k: string]: string;
          };
          /**
           * readOnly is Optional: defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
           */
          readOnly?: boolean;
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          secretRef?: {
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            [k: string]: unknown;
          };
          /**
           * driver is the name of the driver to use for this volume.
           */
          driver: string;
          /**
           * fsType is the Filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". The default filesystem depends on FlexVolume script.
           */
          fsType?: string;
          [k: string]: unknown;
        };
        /**
         * capacity is the description of the persistent volume's resources and capacity. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity
         */
        capacity?: {
          /**
           * Quantity is a fixed-point representation of a number. It provides convenient marshaling/unmarshaling in JSON and YAML, in addition to String() and AsInt64() accessors.
           *
           * The serialization format is:
           *
           * ``` <quantity>        ::= <signedNumber><suffix>
           *
           * 	(Note that <suffix> may be empty, from the "" case in <decimalSI>.)
           *
           * <digit>           ::= 0 | 1 | ... | 9 <digits>          ::= <digit> | <digit><digits> <number>          ::= <digits> | <digits>.<digits> | <digits>. | .<digits> <sign>            ::= "+" | "-" <signedNumber>    ::= <number> | <sign><number> <suffix>          ::= <binarySI> | <decimalExponent> | <decimalSI> <binarySI>        ::= Ki | Mi | Gi | Ti | Pi | Ei
           *
           * 	(International System of units; See: http://physics.nist.gov/cuu/Units/binary.html)
           *
           * <decimalSI>       ::= m | "" | k | M | G | T | P | E
           *
           * 	(Note that 1024 = 1Ki but 1000 = 1k; I didn't choose the capitalization.)
           *
           * <decimalExponent> ::= "e" <signedNumber> | "E" <signedNumber> ```
           *
           * No matter which of the three exponent forms is used, no quantity may represent a number greater than 2^63-1 in magnitude, nor may it have more than 3 decimal places. Numbers larger or more precise will be capped or rounded up. (E.g.: 0.1m will rounded up to 1m.) This may be extended in the future if we require larger or smaller quantities.
           *
           * When a Quantity is parsed from a string, it will remember the type of suffix it had, and will use the same type again when it is serialized.
           *
           * Before serializing, Quantity will be put in "canonical form". This means that Exponent/suffix will be adjusted up or down (with a corresponding increase or decrease in Mantissa) such that:
           *
           * - No precision is lost - No fractional digits will be emitted - The exponent (or suffix) is as large as possible.
           *
           * The sign will be omitted unless the number is negative.
           *
           * Examples:
           *
           * - 1.5 will be serialized as "1500m" - 1.5Gi will be serialized as "1536Mi"
           *
           * Note that the quantity will NEVER be internally represented by a floating point number. That is the whole point of this exercise.
           *
           * Non-canonical values will still parse as long as they are well formed, but will be re-emitted in their canonical form. (So always use canonical form, or don't diff.)
           *
           * This format is intended to make it difficult to use these numbers without writing some sort of special handling code in the hopes that that will cause implementors to also use a fixed point implementation.
           */
          [k: string]: string | number;
        };
        /**
         * VolumeNodeAffinity defines constraints that limit what nodes this volume can be accessed from.
         */
        nodeAffinity?: {
          /**
           * A node selector represents the union of the results of one or more label queries over a set of nodes; that is, it represents the OR of the selectors represented by the node selector terms.
           */
          required?: {
            /**
             * Required. A list of node selector terms. The terms are ORed.
             */
            nodeSelectorTerms: {
              /**
               * A list of node selector requirements by node's fields.
               */
              matchFields?: {
                /**
                 * The label key that the selector applies to.
                 */
                key: string;
                /**
                 * Represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
                 *
                 * Possible enum values:
                 *  - `"DoesNotExist"`
                 *  - `"Exists"`
                 *  - `"Gt"`
                 *  - `"In"`
                 *  - `"Lt"`
                 *  - `"NotIn"`
                 */
                operator: 'DoesNotExist' | 'Exists' | 'Gt' | 'In' | 'Lt' | 'NotIn';
                /**
                 * An array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. If the operator is Gt or Lt, the values array must have a single element, which will be interpreted as an integer. This array is replaced during a strategic merge patch.
                 */
                values?: string[];
                [k: string]: unknown;
              }[];
              /**
               * A list of node selector requirements by node's labels.
               */
              matchExpressions?: {
                /**
                 * An array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. If the operator is Gt or Lt, the values array must have a single element, which will be interpreted as an integer. This array is replaced during a strategic merge patch.
                 */
                values?: string[];
                /**
                 * The label key that the selector applies to.
                 */
                key: string;
                /**
                 * Represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
                 *
                 * Possible enum values:
                 *  - `"DoesNotExist"`
                 *  - `"Exists"`
                 *  - `"Gt"`
                 *  - `"In"`
                 *  - `"Lt"`
                 *  - `"NotIn"`
                 */
                operator: 'DoesNotExist' | 'Exists' | 'Gt' | 'In' | 'Lt' | 'NotIn';
                [k: string]: unknown;
              }[];
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          [k: string]: unknown;
        };
        /**
         * PortworxVolumeSource represents a Portworx volume resource.
         */
        portworxVolume?: {
          /**
           * volumeID uniquely identifies a Portworx volume
           */
          volumeID: string;
          /**
           * fSType represents the filesystem type to mount Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs". Implicitly inferred to be "ext4" if unspecified.
           */
          fsType?: string;
          /**
           * readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
           */
          readOnly?: boolean;
          [k: string]: unknown;
        };
        /**
         * ScaleIOPersistentVolumeSource represents a persistent ScaleIO volume
         */
        scaleIO?: {
          /**
           * sslEnabled is the flag to enable/disable SSL communication with Gateway, default false
           */
          sslEnabled?: boolean;
          /**
           * volumeName is the name of a volume already created in the ScaleIO system that is associated with this volume source.
           */
          volumeName?: string;
          /**
           * fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Default is "xfs"
           */
          fsType?: string;
          /**
           * gateway is the host address of the ScaleIO API Gateway.
           */
          gateway: string;
          /**
           * protectionDomain is the name of the ScaleIO Protection Domain for the configured storage.
           */
          protectionDomain?: string;
          /**
           * readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
           */
          readOnly?: boolean;
          /**
           * storageMode indicates whether the storage for a volume should be ThickProvisioned or ThinProvisioned. Default is ThinProvisioned.
           */
          storageMode?: string;
          /**
           * storagePool is the ScaleIO Storage Pool associated with the protection domain.
           */
          storagePool?: string;
          /**
           * system is the name of the storage system as configured in ScaleIO.
           */
          system: string;
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          secretRef: {
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            [k: string]: unknown;
          };
          [k: string]: unknown;
        };
        /**
         * Represents a StorageOS persistent volume resource.
         */
        storageos?: {
          /**
           * fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
           */
          fsType?: string;
          /**
           * readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
           */
          readOnly?: boolean;
          /**
           * ObjectReference contains enough information to let you inspect or modify the referred object.
           */
          secretRef?: {
            /**
             * If referring to a piece of an object instead of an entire object, this string should contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2]. For example, if the object reference is to a container within a pod, this would take on a value like: "spec.containers{name}" (where "name" refers to the name of the container that triggered the event) or if no container name is specified "spec.containers[2]" (container with index 2 in this pod). This syntax is chosen only to have some well-defined way of referencing a part of an object.
             */
            fieldPath?: string;
            /**
             * Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
             */
            kind?: string;
            /**
             * Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
             */
            name?: string;
            /**
             * Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
             */
            namespace?: string;
            /**
             * Specific resourceVersion to which this reference is made, if any. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
             */
            resourceVersion?: string;
            /**
             * UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids
             */
            uid?: string;
            /**
             * API version of the referent.
             */
            apiVersion?: string;
            [k: string]: unknown;
          };
          /**
           * volumeName is the human-readable name of the StorageOS volume.  Volume names are only unique within a namespace.
           */
          volumeName?: string;
          /**
           * volumeNamespace specifies the scope of the volume within StorageOS.  If no namespace is specified then the Pod's namespace will be used.  This allows the Kubernetes name scoping to be mirrored within StorageOS for tighter integration. Set VolumeName to any name to override the default behaviour. Set to "default" if you are not using namespaces within StorageOS. Namespaces that do not pre-exist within StorageOS will be created.
           */
          volumeNamespace?: string;
          [k: string]: unknown;
        };
        /**
         * Name of VolumeAttributesClass to which this persistent volume belongs. Empty value is not allowed. When this field is not set, it indicates that this volume does not belong to any VolumeAttributesClass. This field is mutable and can be changed by the CSI driver after a volume has been updated successfully to a new class. For an unbound PersistentVolume, the volumeAttributesClassName will be matched with unbound PersistentVolumeClaims during the binding process. This is a beta field and requires enabling VolumeAttributesClass feature (off by default).
         */
        volumeAttributesClassName?: string;
        /**
         * AzureFile represents an Azure File Service mount on the host and bind mount to the pod.
         */
        azureFile?: {
          /**
           * readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
           */
          readOnly?: boolean;
          /**
           * secretName is the name of secret that contains Azure Storage Account Name and Key
           */
          secretName: string;
          /**
           * secretNamespace is the namespace of the secret that contains Azure Storage Account Name and Key default is the same as the Pod
           */
          secretNamespace?: string;
          /**
           * shareName is the azure Share Name
           */
          shareName: string;
          [k: string]: unknown;
        };
        /**
         * accessModes contains all ways the volume can be mounted. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes
         */
        accessModes?: ('ReadOnlyMany' | 'ReadWriteMany' | 'ReadWriteOnce' | 'ReadWriteOncePod')[];
        /**
         * Represents a cinder volume resource in Openstack. A Cinder volume must exist before mounting to a container. The volume must also be in the same region as the kubelet. Cinder volumes support ownership management and SELinux relabeling.
         */
        cinder?: {
          /**
           * fsType Filesystem type to mount. Must be a filesystem type supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
           */
          fsType?: string;
          /**
           * readOnly is Optional: Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
           */
          readOnly?: boolean;
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          secretRef?: {
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            [k: string]: unknown;
          };
          /**
           * volumeID used to identify the volume in cinder. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
           */
          volumeID: string;
          [k: string]: unknown;
        };
        /**
         * Represents storage that is managed by an external CSI volume driver
         */
        csi?: {
          /**
           * volumeAttributes of the volume to publish.
           */
          volumeAttributes?: {
            [k: string]: string;
          };
          /**
           * driver is the name of the driver to use for this volume. Required.
           */
          driver: string;
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          nodePublishSecretRef?: {
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            [k: string]: unknown;
          };
          /**
           * readOnly value to pass to ControllerPublishVolumeRequest. Defaults to false (read/write).
           */
          readOnly?: boolean;
          /**
           * volumeHandle is the unique volume name returned by the CSI volume plugin’s CreateVolume to refer to the volume on all subsequent calls. Required.
           */
          volumeHandle: string;
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          controllerExpandSecretRef?: {
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            [k: string]: unknown;
          };
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          controllerPublishSecretRef?: {
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            [k: string]: unknown;
          };
          /**
           * fsType to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs".
           */
          fsType?: string;
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          nodeExpandSecretRef?: {
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            [k: string]: unknown;
          };
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          nodeStageSecretRef?: {
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            [k: string]: unknown;
          };
          [k: string]: unknown;
        };
        /**
         * Represents a Quobyte mount that lasts the lifetime of a pod. Quobyte volumes do not support ownership management or SELinux relabeling.
         */
        quobyte?: {
          /**
           * tenant owning the given Quobyte volume in the Backend Used with dynamically provisioned Quobyte volumes, value is set by the plugin
           */
          tenant?: string;
          /**
           * user to map volume access to Defaults to serivceaccount user
           */
          user?: string;
          /**
           * volume is a string that references an already created Quobyte volume by name.
           */
          volume: string;
          /**
           * group to map volume access to Default is no group
           */
          group?: string;
          /**
           * readOnly here will force the Quobyte volume to be mounted with read-only permissions. Defaults to false.
           */
          readOnly?: boolean;
          /**
           * registry represents a single or multiple Quobyte Registry services specified as a string as host:port pair (multiple entries are separated with commas) which acts as the central registry for volumes
           */
          registry: string;
          [k: string]: unknown;
        };
        /**
         * volumeMode defines if a volume is intended to be used with a formatted filesystem or to remain in raw block state. Value of Filesystem is implied when not included in spec.
         *
         * Possible enum values:
         *  - `"Block"` means the volume will not be formatted with a filesystem and will remain a raw block device.
         *  - `"Filesystem"` means the volume will be or is formatted with a filesystem.
         */
        volumeMode?: 'Block' | 'Filesystem';
        /**
         * Represents a Ceph Filesystem mount that lasts the lifetime of a pod Cephfs volumes do not support ownership management or SELinux relabeling.
         */
        cephfs?: {
          /**
           * SecretReference represents a Secret Reference. It has enough information to retrieve secret in any namespace
           */
          secretRef?: {
            /**
             * name is unique within a namespace to reference a secret resource.
             */
            name?: string;
            /**
             * namespace defines the space within which the secret name must be unique.
             */
            namespace?: string;
            [k: string]: unknown;
          };
          /**
           * user is Optional: User is the rados user name, default is admin More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
           */
          user?: string;
          /**
           * monitors is Required: Monitors is a collection of Ceph monitors More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
           */
          monitors: string[];
          /**
           * path is Optional: Used as the mounted root, rather than the full Ceph tree, default is /
           */
          path?: string;
          /**
           * readOnly is Optional: Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts. More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
           */
          readOnly?: boolean;
          /**
           * secretFile is Optional: SecretFile is the path to key ring for User, default is /etc/ceph/user.secret More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
           */
          secretFile?: string;
          [k: string]: unknown;
        };
        /**
         * Represents a Glusterfs mount that lasts the lifetime of a pod. Glusterfs volumes do not support ownership management or SELinux relabeling.
         */
        glusterfs?: {
          /**
           * path is the Glusterfs volume path. More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
           */
          path: string;
          /**
           * readOnly here will force the Glusterfs volume to be mounted with read-only permissions. Defaults to false. More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
           */
          readOnly?: boolean;
          /**
           * endpoints is the endpoint name that details Glusterfs topology. More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
           */
          endpoints: string;
          /**
           * endpointsNamespace is the namespace that contains Glusterfs endpoint. If this field is empty, the EndpointNamespace defaults to the same namespace as the bound PVC. More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
           */
          endpointsNamespace?: string;
          [k: string]: unknown;
        };
        /**
         * Represents a host path mapped into a pod. Host path volumes do not support ownership management or SELinux relabeling.
         */
        hostPath?: {
          /**
           * path of the directory on the host. If the path is a symlink, it will follow the link to the real path. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
           */
          path: string;
          /**
           * type for HostPath Volume Defaults to "" More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
           *
           * Possible enum values:
           *  - `""` For backwards compatible, leave it empty if unset
           *  - `"BlockDevice"` A block device must exist at the given path
           *  - `"CharDevice"` A character device must exist at the given path
           *  - `"Directory"` A directory must exist at the given path
           *  - `"DirectoryOrCreate"` If nothing exists at the given path, an empty directory will be created there as needed with file mode 0755, having the same group and ownership with Kubelet.
           *  - `"File"` A file must exist at the given path
           *  - `"FileOrCreate"` If nothing exists at the given path, an empty file will be created there as needed with file mode 0644, having the same group and ownership with Kubelet.
           *  - `"Socket"` A UNIX socket must exist at the given path
           */
          type?:
            | ''
            | 'BlockDevice'
            | 'CharDevice'
            | 'Directory'
            | 'DirectoryOrCreate'
            | 'File'
            | 'FileOrCreate'
            | 'Socket';
          [k: string]: unknown;
        };
        [k: string]: unknown;
      };
      /**
       * persistentVolumeName represents the name of the persistent volume to attach.
       */
      persistentVolumeName?: string;
      [k: string]: unknown;
    };
    [k: string]: unknown;
  };
  /**
   * VolumeAttachmentStatus is the status of a VolumeAttachment request.
   */
  status?: {
    /**
     * VolumeError captures an error encountered during a volume operation.
     */
    attachError?: {
      /**
       * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
       */
      time?: string;
      /**
       * errorCode is a numeric gRPC code representing the error encountered during Attach or Detach operations.
       *
       * This is an optional, alpha field that requires the MutableCSINodeAllocatableCount feature gate being enabled to be set.
       */
      errorCode?: number;
      /**
       * message represents the error encountered during Attach or Detach operation. This string may be logged, so it should not contain sensitive information.
       */
      message?: string;
      [k: string]: unknown;
    };
    /**
     * attached indicates the volume is successfully attached. This field must only be set by the entity completing the attach operation, i.e. the external-attacher.
     */
    attached: boolean;
    /**
     * attachmentMetadata is populated with any information returned by the attach operation, upon successful attach, that must be passed into subsequent WaitForAttach or Mount calls. This field must only be set by the entity completing the attach operation, i.e. the external-attacher.
     */
    attachmentMetadata?: {
      [k: string]: string;
    };
    /**
     * VolumeError captures an error encountered during a volume operation.
     */
    detachError?: {
      /**
       * errorCode is a numeric gRPC code representing the error encountered during Attach or Detach operations.
       *
       * This is an optional, alpha field that requires the MutableCSINodeAllocatableCount feature gate being enabled to be set.
       */
      errorCode?: number;
      /**
       * message represents the error encountered during Attach or Detach operation. This string may be logged, so it should not contain sensitive information.
       */
      message?: string;
      /**
       * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
       */
      time?: string;
      [k: string]: unknown;
    };
    [k: string]: unknown;
  };
  [k: string]: unknown;
}

/** Generated from Remote JSON Schema for apps.v1.StatefulSet */

/**
 * StatefulSet represents a set of pods with consistent identities. Identities are defined as:
 *   - Network: A single stable DNS and hostname.
 *   - Storage: As many VolumeClaims as requested.
 *
 * The StatefulSet guarantees that a given network identity will always map to the same storage identity.
 */
export interface AppsV1StatefulSet {
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
     * Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
     */
    annotations?: {
      [k: string]: string;
    };
    /**
     * Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order.  Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list.
     */
    finalizers?: string[];
    /**
     * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
     */
    deletionTimestamp?: string;
    /**
     * GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
     *
     * If this field is specified and the generated name exists, the server will return a 409.
     *
     * Applied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency
     */
    generateName?: string;
    /**
     * An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.
     *
     * Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
     */
    resourceVersion?: string;
    /**
     * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
     */
    creationTimestamp?: string;
    /**
     * Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
     */
    deletionGracePeriodSeconds?: number;
    /**
     * A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.
     */
    generation?: number;
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
      /**
       * Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
       */
      subresource?: string;
      [k: string]: unknown;
    }[];
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
     * Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.
     */
    selfLink?: string;
    [k: string]: unknown;
  };
  /**
   * A StatefulSetSpec is the specification of a StatefulSet.
   */
  spec?: {
    /**
     * Minimum number of seconds for which a newly created pod should be ready without any of its container crashing for it to be considered available. Defaults to 0 (pod will be considered available as soon as it is ready)
     */
    minReadySeconds?: number;
    /**
     * StatefulSetPersistentVolumeClaimRetentionPolicy describes the policy used for PVCs created from the StatefulSet VolumeClaimTemplates.
     */
    persistentVolumeClaimRetentionPolicy?: {
      /**
       * WhenDeleted specifies what happens to PVCs created from StatefulSet VolumeClaimTemplates when the StatefulSet is deleted. The default policy of `Retain` causes PVCs to not be affected by StatefulSet deletion. The `Delete` policy causes those PVCs to be deleted.
       */
      whenDeleted?: string;
      /**
       * WhenScaled specifies what happens to PVCs created from StatefulSet VolumeClaimTemplates when the StatefulSet is scaled down. The default policy of `Retain` causes PVCs to not be affected by a scaledown. The `Delete` policy causes the associated PVCs for any excess pods above the replica count to be deleted.
       */
      whenScaled?: string;
      [k: string]: unknown;
    };
    /**
     * podManagementPolicy controls how pods are created during initial scale up, when replacing pods on nodes, or when scaling down. The default policy is `OrderedReady`, where pods are created in increasing order (pod-0, then pod-1, etc) and the controller will wait until each pod is ready before continuing. When scaling down, the pods are removed in the opposite order. The alternative policy is `Parallel` which will create pods in parallel to match the desired scale without waiting, and on scale down will delete all pods at once.
     *
     * Possible enum values:
     *  - `"OrderedReady"` will create pods in strictly increasing order on scale up and strictly decreasing order on scale down, progressing only when the previous pod is ready or terminated. At most one pod will be changed at any time.
     *  - `"Parallel"` will create and delete pods as soon as the stateful set replica count is changed, and will not wait for pods to be ready or complete termination.
     */
    podManagementPolicy?: 'OrderedReady' | 'Parallel';
    /**
     * replicas is the desired number of replicas of the given Template. These are replicas in the sense that they are instantiations of the same Template, but individual replicas also have a consistent identity. If unspecified, defaults to 1.
     */
    replicas?: number;
    /**
     * revisionHistoryLimit is the maximum number of revisions that will be maintained in the StatefulSet's revision history. The revision history consists of all revisions not represented by a currently applied StatefulSetSpec version. The default value is 10.
     */
    revisionHistoryLimit?: number;
    /**
     * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
     */
    selector: {
      /**
       * matchExpressions is a list of label selector requirements. The requirements are ANDed.
       */
      matchExpressions?: {
        /**
         * key is the label key that the selector applies to.
         */
        key: string;
        /**
         * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
         */
        operator: string;
        /**
         * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
         */
        values?: string[];
        [k: string]: unknown;
      }[];
      /**
       * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
       */
      matchLabels?: {
        [k: string]: string;
      };
      [k: string]: unknown;
    };
    /**
     * PodTemplateSpec describes the data a pod should have when created from a template
     */
    template: {
      /**
       * ObjectMeta is metadata that all persisted resources must have, which includes all objects users must create.
       */
      metadata?: {
        /**
         * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
         */
        creationTimestamp?: string;
        /**
         * Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order.  Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list.
         */
        finalizers?: string[];
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
          /**
           * Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
           */
          subresource?: string;
          /**
           * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
           */
          time?: string;
          [k: string]: unknown;
        }[];
        /**
         * An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.
         *
         * Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
         */
        resourceVersion?: string;
        /**
         * Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.
         */
        selfLink?: string;
        /**
         * Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
         */
        deletionGracePeriodSeconds?: number;
        /**
         * A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.
         */
        generation?: number;
        /**
         * UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.
         *
         * Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
         */
        uid?: string;
        /**
         * Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
         */
        annotations?: {
          [k: string]: string;
        };
        /**
         * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
         */
        deletionTimestamp?: string;
        /**
         * GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
         *
         * If this field is specified and the generated name exists, the server will return a 409.
         *
         * Applied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency
         */
        generateName?: string;
        /**
         * Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
         */
        name?: string;
        /**
         * Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
         *
         * Must be a DNS_LABEL. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces
         */
        namespace?: string;
        /**
         * List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller.
         */
        ownerReferences?: {
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
          /**
           * API version of the referent.
           */
          apiVersion: string;
          /**
           * If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. See https://kubernetes.io/docs/concepts/architecture/garbage-collection/#foreground-deletion for how the garbage collector interacts with this field and enforces the foreground deletion. Defaults to false. To set this field, a user needs "delete" permission of the owner, otherwise 422 (Unprocessable Entity) will be returned.
           */
          blockOwnerDeletion?: boolean;
          [k: string]: unknown;
        }[];
        [k: string]: unknown;
      };
      /**
       * PodSpec is a description of a pod.
       */
      spec?: {
        /**
         * Affinity is a group of affinity scheduling rules.
         */
        affinity?: {
          /**
           * Node affinity is a group of node affinity scheduling rules.
           */
          nodeAffinity?: {
            /**
             * A node selector represents the union of the results of one or more label queries over a set of nodes; that is, it represents the OR of the selectors represented by the node selector terms.
             */
            requiredDuringSchedulingIgnoredDuringExecution?: {
              /**
               * Required. A list of node selector terms. The terms are ORed.
               */
              nodeSelectorTerms: {
                /**
                 * A list of node selector requirements by node's labels.
                 */
                matchExpressions?: {
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
                [k: string]: unknown;
              }[];
              [k: string]: unknown;
            };
            /**
             * The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling affinity expressions, etc.), compute a sum by iterating through the elements of this field and adding "weight" to the sum if the node matches the corresponding matchExpressions; the node(s) with the highest sum are the most preferred.
             */
            preferredDuringSchedulingIgnoredDuringExecution?: {
              /**
               * A null or empty node selector term matches no objects. The requirements of them are ANDed. The TopologySelectorTerm type implements a subset of the NodeSelectorTerm.
               */
              preference: {
                /**
                 * A list of node selector requirements by node's labels.
                 */
                matchExpressions?: {
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
                [k: string]: unknown;
              };
              /**
               * Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.
               */
              weight: number;
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          /**
           * Pod affinity is a group of inter pod affinity scheduling rules.
           */
          podAffinity?: {
            /**
             * The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling affinity expressions, etc.), compute a sum by iterating through the elements of this field and adding "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the node(s) with the highest sum are the most preferred.
             */
            preferredDuringSchedulingIgnoredDuringExecution?: {
              /**
               * weight associated with matching the corresponding podAffinityTerm, in the range 1-100.
               */
              weight: number;
              /**
               * Defines a set of pods (namely those matching the labelSelector relative to the given namespace(s)) that this pod should be co-located (affinity) or not co-located (anti-affinity) with, where co-located is defined as running on a node whose value of the label with key <topologyKey> matches that of any node on which a pod of the set of pods is running
               */
              podAffinityTerm: {
                /**
                 * namespaces specifies a static list of namespace names that the term applies to. The term is applied to the union of the namespaces listed in this field and the ones selected by namespaceSelector. null or empty namespaces list and null namespaceSelector means "this pod's namespace".
                 */
                namespaces?: string[];
                /**
                 * This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching the labelSelector in the specified namespaces, where co-located is defined as running on a node whose value of the label with key topologyKey matches that of any node on which any of the selected pods is running. Empty topologyKey is not allowed.
                 */
                topologyKey: string;
                /**
                 * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
                 */
                labelSelector?: {
                  /**
                   * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                   */
                  matchExpressions?: {
                    /**
                     * key is the label key that the selector applies to.
                     */
                    key: string;
                    /**
                     * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                     */
                    operator: string;
                    /**
                     * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                     */
                    values?: string[];
                    [k: string]: unknown;
                  }[];
                  /**
                   * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                   */
                  matchLabels?: {
                    [k: string]: string;
                  };
                  [k: string]: unknown;
                };
                /**
                 * MatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key in (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both matchLabelKeys and labelSelector. Also, matchLabelKeys cannot be set when labelSelector isn't set.
                 */
                matchLabelKeys?: string[];
                /**
                 * MismatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key notin (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both mismatchLabelKeys and labelSelector. Also, mismatchLabelKeys cannot be set when labelSelector isn't set.
                 */
                mismatchLabelKeys?: string[];
                /**
                 * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
                 */
                namespaceSelector?: {
                  /**
                   * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                   */
                  matchExpressions?: {
                    /**
                     * key is the label key that the selector applies to.
                     */
                    key: string;
                    /**
                     * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                     */
                    operator: string;
                    /**
                     * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                     */
                    values?: string[];
                    [k: string]: unknown;
                  }[];
                  /**
                   * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                   */
                  matchLabels?: {
                    [k: string]: string;
                  };
                  [k: string]: unknown;
                };
                [k: string]: unknown;
              };
              [k: string]: unknown;
            }[];
            /**
             * If the affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node. If the affinity requirements specified by this field cease to be met at some point during pod execution (e.g. due to a pod label update), the system may or may not try to eventually evict the pod from its node. When there are multiple elements, the lists of nodes corresponding to each podAffinityTerm are intersected, i.e. all terms must be satisfied.
             */
            requiredDuringSchedulingIgnoredDuringExecution?: {
              /**
               * MismatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key notin (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both mismatchLabelKeys and labelSelector. Also, mismatchLabelKeys cannot be set when labelSelector isn't set.
               */
              mismatchLabelKeys?: string[];
              /**
               * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
               */
              namespaceSelector?: {
                /**
                 * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                 */
                matchLabels?: {
                  [k: string]: string;
                };
                /**
                 * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                 */
                matchExpressions?: {
                  /**
                   * key is the label key that the selector applies to.
                   */
                  key: string;
                  /**
                   * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                   */
                  operator: string;
                  /**
                   * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                   */
                  values?: string[];
                  [k: string]: unknown;
                }[];
                [k: string]: unknown;
              };
              /**
               * namespaces specifies a static list of namespace names that the term applies to. The term is applied to the union of the namespaces listed in this field and the ones selected by namespaceSelector. null or empty namespaces list and null namespaceSelector means "this pod's namespace".
               */
              namespaces?: string[];
              /**
               * This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching the labelSelector in the specified namespaces, where co-located is defined as running on a node whose value of the label with key topologyKey matches that of any node on which any of the selected pods is running. Empty topologyKey is not allowed.
               */
              topologyKey: string;
              /**
               * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
               */
              labelSelector?: {
                /**
                 * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                 */
                matchExpressions?: {
                  /**
                   * key is the label key that the selector applies to.
                   */
                  key: string;
                  /**
                   * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                   */
                  operator: string;
                  /**
                   * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                   */
                  values?: string[];
                  [k: string]: unknown;
                }[];
                /**
                 * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                 */
                matchLabels?: {
                  [k: string]: string;
                };
                [k: string]: unknown;
              };
              /**
               * MatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key in (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both matchLabelKeys and labelSelector. Also, matchLabelKeys cannot be set when labelSelector isn't set.
               */
              matchLabelKeys?: string[];
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          /**
           * Pod anti affinity is a group of inter pod anti affinity scheduling rules.
           */
          podAntiAffinity?: {
            /**
             * The scheduler will prefer to schedule pods to nodes that satisfy the anti-affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling anti-affinity expressions, etc.), compute a sum by iterating through the elements of this field and adding "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the node(s) with the highest sum are the most preferred.
             */
            preferredDuringSchedulingIgnoredDuringExecution?: {
              /**
               * Defines a set of pods (namely those matching the labelSelector relative to the given namespace(s)) that this pod should be co-located (affinity) or not co-located (anti-affinity) with, where co-located is defined as running on a node whose value of the label with key <topologyKey> matches that of any node on which a pod of the set of pods is running
               */
              podAffinityTerm: {
                /**
                 * namespaces specifies a static list of namespace names that the term applies to. The term is applied to the union of the namespaces listed in this field and the ones selected by namespaceSelector. null or empty namespaces list and null namespaceSelector means "this pod's namespace".
                 */
                namespaces?: string[];
                /**
                 * This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching the labelSelector in the specified namespaces, where co-located is defined as running on a node whose value of the label with key topologyKey matches that of any node on which any of the selected pods is running. Empty topologyKey is not allowed.
                 */
                topologyKey: string;
                /**
                 * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
                 */
                labelSelector?: {
                  /**
                   * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                   */
                  matchExpressions?: {
                    /**
                     * key is the label key that the selector applies to.
                     */
                    key: string;
                    /**
                     * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                     */
                    operator: string;
                    /**
                     * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                     */
                    values?: string[];
                    [k: string]: unknown;
                  }[];
                  /**
                   * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                   */
                  matchLabels?: {
                    [k: string]: string;
                  };
                  [k: string]: unknown;
                };
                /**
                 * MatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key in (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both matchLabelKeys and labelSelector. Also, matchLabelKeys cannot be set when labelSelector isn't set.
                 */
                matchLabelKeys?: string[];
                /**
                 * MismatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key notin (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both mismatchLabelKeys and labelSelector. Also, mismatchLabelKeys cannot be set when labelSelector isn't set.
                 */
                mismatchLabelKeys?: string[];
                /**
                 * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
                 */
                namespaceSelector?: {
                  /**
                   * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                   */
                  matchExpressions?: {
                    /**
                     * key is the label key that the selector applies to.
                     */
                    key: string;
                    /**
                     * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                     */
                    operator: string;
                    /**
                     * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                     */
                    values?: string[];
                    [k: string]: unknown;
                  }[];
                  /**
                   * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                   */
                  matchLabels?: {
                    [k: string]: string;
                  };
                  [k: string]: unknown;
                };
                [k: string]: unknown;
              };
              /**
               * weight associated with matching the corresponding podAffinityTerm, in the range 1-100.
               */
              weight: number;
              [k: string]: unknown;
            }[];
            /**
             * If the anti-affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node. If the anti-affinity requirements specified by this field cease to be met at some point during pod execution (e.g. due to a pod label update), the system may or may not try to eventually evict the pod from its node. When there are multiple elements, the lists of nodes corresponding to each podAffinityTerm are intersected, i.e. all terms must be satisfied.
             */
            requiredDuringSchedulingIgnoredDuringExecution?: {
              /**
               * MismatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key notin (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both mismatchLabelKeys and labelSelector. Also, mismatchLabelKeys cannot be set when labelSelector isn't set.
               */
              mismatchLabelKeys?: string[];
              /**
               * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
               */
              namespaceSelector?: {
                /**
                 * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                 */
                matchLabels?: {
                  [k: string]: string;
                };
                /**
                 * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                 */
                matchExpressions?: {
                  /**
                   * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                   */
                  values?: string[];
                  /**
                   * key is the label key that the selector applies to.
                   */
                  key: string;
                  /**
                   * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                   */
                  operator: string;
                  [k: string]: unknown;
                }[];
                [k: string]: unknown;
              };
              /**
               * namespaces specifies a static list of namespace names that the term applies to. The term is applied to the union of the namespaces listed in this field and the ones selected by namespaceSelector. null or empty namespaces list and null namespaceSelector means "this pod's namespace".
               */
              namespaces?: string[];
              /**
               * This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching the labelSelector in the specified namespaces, where co-located is defined as running on a node whose value of the label with key topologyKey matches that of any node on which any of the selected pods is running. Empty topologyKey is not allowed.
               */
              topologyKey: string;
              /**
               * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
               */
              labelSelector?: {
                /**
                 * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                 */
                matchExpressions?: {
                  /**
                   * key is the label key that the selector applies to.
                   */
                  key: string;
                  /**
                   * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                   */
                  operator: string;
                  /**
                   * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                   */
                  values?: string[];
                  [k: string]: unknown;
                }[];
                /**
                 * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                 */
                matchLabels?: {
                  [k: string]: string;
                };
                [k: string]: unknown;
              };
              /**
               * MatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key in (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both matchLabelKeys and labelSelector. Also, matchLabelKeys cannot be set when labelSelector isn't set.
               */
              matchLabelKeys?: string[];
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          [k: string]: unknown;
        };
        /**
         * Host networking requested for this pod. Use the host's network namespace. If this option is set, the ports that will be used must be specified. Default to false.
         */
        hostNetwork?: boolean;
        /**
         * ResourceClaims defines which ResourceClaims must be allocated and reserved before the Pod is allowed to start. The resources will be made available to those containers which consume them by name.
         *
         * This is an alpha field and requires enabling the DynamicResourceAllocation feature gate.
         *
         * This field is immutable.
         */
        resourceClaims?: {
          /**
           * Name uniquely identifies this resource claim inside the pod. This must be a DNS_LABEL.
           */
          name: string;
          /**
           * ResourceClaimName is the name of a ResourceClaim object in the same namespace as this pod.
           *
           * Exactly one of ResourceClaimName and ResourceClaimTemplateName must be set.
           */
          resourceClaimName?: string;
          /**
           * ResourceClaimTemplateName is the name of a ResourceClaimTemplate object in the same namespace as this pod.
           *
           * The template will be used to create a new ResourceClaim, which will be bound to this pod. When this pod is deleted, the ResourceClaim will also be deleted. The pod name and resource name, along with a generated component, will be used to form a unique name for the ResourceClaim, which will be recorded in pod.status.resourceClaimStatuses.
           *
           * This field is immutable and no changes will be made to the corresponding ResourceClaim by the control plane after creating the ResourceClaim.
           *
           * Exactly one of ResourceClaimName and ResourceClaimTemplateName must be set.
           */
          resourceClaimTemplateName?: string;
          [k: string]: unknown;
        }[];
        /**
         * ResourceRequirements describes the compute resource requirements.
         */
        resources?: {
          /**
           * Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container.
           *
           * This is an alpha field and requires enabling the DynamicResourceAllocation feature gate.
           *
           * This field is immutable. It can only be set for containers.
           */
          claims?: {
            /**
             * Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
             */
            name: string;
            /**
             * Request is the name chosen for a request in the referenced claim. If empty, everything from the claim is made available, otherwise only the result of this request.
             */
            request?: string;
            [k: string]: unknown;
          }[];
          /**
           * Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
           */
          limits?: {
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
           * Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
           */
          requests?: {
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
          [k: string]: unknown;
        };
        /**
         * ServiceAccountName is the name of the ServiceAccount to use to run this pod. More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
         */
        serviceAccountName?: string;
        /**
         * Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). If this value is nil, the default grace period will be used instead. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. Defaults to 30 seconds.
         */
        terminationGracePeriodSeconds?: number;
        /**
         * Set DNS policy for the pod. Defaults to "ClusterFirst". Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'. DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy. To have DNS options set along with hostNetwork, you have to specify DNS policy explicitly to 'ClusterFirstWithHostNet'.
         *
         * Possible enum values:
         *  - `"ClusterFirst"` indicates that the pod should use cluster DNS first unless hostNetwork is true, if it is available, then fall back on the default (as determined by kubelet) DNS settings.
         *  - `"ClusterFirstWithHostNet"` indicates that the pod should use cluster DNS first, if it is available, then fall back on the default (as determined by kubelet) DNS settings.
         *  - `"Default"` indicates that the pod should use the default (as determined by kubelet) DNS settings.
         *  - `"None"` indicates that the pod should use empty DNS settings. DNS parameters such as nameservers and search paths should be defined via DNSConfig.
         */
        dnsPolicy?: 'ClusterFirst' | 'ClusterFirstWithHostNet' | 'Default' | 'None';
        /**
         * HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts file if specified.
         */
        hostAliases?: {
          /**
           * Hostnames for the above IP address.
           */
          hostnames?: string[];
          /**
           * IP address of the host file entry.
           */
          ip: string;
          [k: string]: unknown;
        }[];
        /**
         * Specifies the hostname of the Pod If not specified, the pod's hostname will be set to a system-defined value.
         */
        hostname?: string;
        /**
         * The priority value. Various system components use this field to find the priority of the pod. When Priority Admission Controller is enabled, it prevents users from setting this field. The admission controller populates this field from PriorityClassName. The higher the value, the higher the priority.
         */
        priority?: number;
        /**
         * PodDNSConfig defines the DNS parameters of a pod in addition to those generated from DNSPolicy.
         */
        dnsConfig?: {
          /**
           * A list of DNS name server IP addresses. This will be appended to the base nameservers generated from DNSPolicy. Duplicated nameservers will be removed.
           */
          nameservers?: string[];
          /**
           * A list of DNS resolver options. This will be merged with the base options generated from DNSPolicy. Duplicated entries will be removed. Resolution options given in Options will override those that appear in the base DNSPolicy.
           */
          options?: {
            /**
             * Name is this DNS resolver option's name. Required.
             */
            name?: string;
            /**
             * Value is this DNS resolver option's value.
             */
            value?: string;
            [k: string]: unknown;
          }[];
          /**
           * A list of DNS search domains for host-name lookup. This will be appended to the base search paths generated from DNSPolicy. Duplicated search paths will be removed.
           */
          searches?: string[];
          [k: string]: unknown;
        };
        /**
         * EnableServiceLinks indicates whether information about services should be injected into pod's environment variables, matching the syntax of Docker links. Optional: Defaults to true.
         */
        enableServiceLinks?: boolean;
        /**
         * List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing pod to perform user-initiated actions such as debugging. This list cannot be specified when creating a pod, and it cannot be modified by updating the pod spec. In order to add an ephemeral container to an existing pod, use the pod's ephemeralcontainers subresource.
         */
        ephemeralContainers?: {
          /**
           * Arguments to the entrypoint. The image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
           */
          args?: string[];
          /**
           * Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
           *
           * Possible enum values:
           *  - `"Always"` means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.
           *  - `"IfNotPresent"` means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
           *  - `"Never"` means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn't present
           */
          imagePullPolicy?: 'Always' | 'IfNotPresent' | 'Never';
          /**
           * Lifecycle describes actions that the management system should take in response to container lifecycle events. For the PostStart and PreStop lifecycle handlers, management of the container blocks until the action is complete, unless the container process fails, in which case the handler is aborted.
           */
          lifecycle?: {
            /**
             * LifecycleHandler defines a specific action that should be taken in a lifecycle hook. One and only one of the fields, except TCPSocket must be specified.
             */
            postStart?: {
              /**
               * ExecAction describes a "run in container" action.
               */
              exec?: {
                /**
                 * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
                 */
                command?: string[];
                [k: string]: unknown;
              };
              /**
               * HTTPGetAction describes an action based on HTTP Get requests.
               */
              httpGet?: {
                /**
                 * Path to access on the HTTP server.
                 */
                path?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                /**
                 * Scheme to use for connecting to the host. Defaults to HTTP.
                 *
                 * Possible enum values:
                 *  - `"HTTP"` means that the scheme used will be http://
                 *  - `"HTTPS"` means that the scheme used will be https://
                 */
                scheme?: 'HTTP' | 'HTTPS';
                /**
                 * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
                 */
                host?: string;
                /**
                 * Custom headers to set in the request. HTTP allows repeated headers.
                 */
                httpHeaders?: {
                  /**
                   * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                   */
                  name: string;
                  /**
                   * The header field value
                   */
                  value: string;
                  [k: string]: unknown;
                }[];
                [k: string]: unknown;
              };
              /**
               * SleepAction describes a "sleep" action.
               */
              sleep?: {
                /**
                 * Seconds is the number of seconds to sleep.
                 */
                seconds: number;
                [k: string]: unknown;
              };
              /**
               * TCPSocketAction describes an action based on opening a socket
               */
              tcpSocket?: {
                /**
                 * Optional: Host name to connect to, defaults to the pod IP.
                 */
                host?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            /**
             * LifecycleHandler defines a specific action that should be taken in a lifecycle hook. One and only one of the fields, except TCPSocket must be specified.
             */
            preStop?: {
              /**
               * ExecAction describes a "run in container" action.
               */
              exec?: {
                /**
                 * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
                 */
                command?: string[];
                [k: string]: unknown;
              };
              /**
               * HTTPGetAction describes an action based on HTTP Get requests.
               */
              httpGet?: {
                /**
                 * Custom headers to set in the request. HTTP allows repeated headers.
                 */
                httpHeaders?: {
                  /**
                   * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                   */
                  name: string;
                  /**
                   * The header field value
                   */
                  value: string;
                  [k: string]: unknown;
                }[];
                /**
                 * Path to access on the HTTP server.
                 */
                path?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                /**
                 * Scheme to use for connecting to the host. Defaults to HTTP.
                 *
                 * Possible enum values:
                 *  - `"HTTP"` means that the scheme used will be http://
                 *  - `"HTTPS"` means that the scheme used will be https://
                 */
                scheme?: 'HTTP' | 'HTTPS';
                /**
                 * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
                 */
                host?: string;
                [k: string]: unknown;
              };
              /**
               * SleepAction describes a "sleep" action.
               */
              sleep?: {
                /**
                 * Seconds is the number of seconds to sleep.
                 */
                seconds: number;
                [k: string]: unknown;
              };
              /**
               * TCPSocketAction describes an action based on opening a socket
               */
              tcpSocket?: {
                /**
                 * Optional: Host name to connect to, defaults to the pod IP.
                 */
                host?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            /**
             * StopSignal defines which signal will be sent to a container when it is being stopped. If not specified, the default is defined by the container runtime in use. StopSignal can only be set for Pods with a non-empty .spec.os.name
             *
             * Possible enum values:
             *  - `"SIGABRT"`
             *  - `"SIGALRM"`
             *  - `"SIGBUS"`
             *  - `"SIGCHLD"`
             *  - `"SIGCLD"`
             *  - `"SIGCONT"`
             *  - `"SIGFPE"`
             *  - `"SIGHUP"`
             *  - `"SIGILL"`
             *  - `"SIGINT"`
             *  - `"SIGIO"`
             *  - `"SIGIOT"`
             *  - `"SIGKILL"`
             *  - `"SIGPIPE"`
             *  - `"SIGPOLL"`
             *  - `"SIGPROF"`
             *  - `"SIGPWR"`
             *  - `"SIGQUIT"`
             *  - `"SIGRTMAX"`
             *  - `"SIGRTMAX-1"`
             *  - `"SIGRTMAX-10"`
             *  - `"SIGRTMAX-11"`
             *  - `"SIGRTMAX-12"`
             *  - `"SIGRTMAX-13"`
             *  - `"SIGRTMAX-14"`
             *  - `"SIGRTMAX-2"`
             *  - `"SIGRTMAX-3"`
             *  - `"SIGRTMAX-4"`
             *  - `"SIGRTMAX-5"`
             *  - `"SIGRTMAX-6"`
             *  - `"SIGRTMAX-7"`
             *  - `"SIGRTMAX-8"`
             *  - `"SIGRTMAX-9"`
             *  - `"SIGRTMIN"`
             *  - `"SIGRTMIN+1"`
             *  - `"SIGRTMIN+10"`
             *  - `"SIGRTMIN+11"`
             *  - `"SIGRTMIN+12"`
             *  - `"SIGRTMIN+13"`
             *  - `"SIGRTMIN+14"`
             *  - `"SIGRTMIN+15"`
             *  - `"SIGRTMIN+2"`
             *  - `"SIGRTMIN+3"`
             *  - `"SIGRTMIN+4"`
             *  - `"SIGRTMIN+5"`
             *  - `"SIGRTMIN+6"`
             *  - `"SIGRTMIN+7"`
             *  - `"SIGRTMIN+8"`
             *  - `"SIGRTMIN+9"`
             *  - `"SIGSEGV"`
             *  - `"SIGSTKFLT"`
             *  - `"SIGSTOP"`
             *  - `"SIGSYS"`
             *  - `"SIGTERM"`
             *  - `"SIGTRAP"`
             *  - `"SIGTSTP"`
             *  - `"SIGTTIN"`
             *  - `"SIGTTOU"`
             *  - `"SIGURG"`
             *  - `"SIGUSR1"`
             *  - `"SIGUSR2"`
             *  - `"SIGVTALRM"`
             *  - `"SIGWINCH"`
             *  - `"SIGXCPU"`
             *  - `"SIGXFSZ"`
             */
            stopSignal?:
              | 'SIGABRT'
              | 'SIGALRM'
              | 'SIGBUS'
              | 'SIGCHLD'
              | 'SIGCLD'
              | 'SIGCONT'
              | 'SIGFPE'
              | 'SIGHUP'
              | 'SIGILL'
              | 'SIGINT'
              | 'SIGIO'
              | 'SIGIOT'
              | 'SIGKILL'
              | 'SIGPIPE'
              | 'SIGPOLL'
              | 'SIGPROF'
              | 'SIGPWR'
              | 'SIGQUIT'
              | 'SIGRTMAX'
              | 'SIGRTMAX-1'
              | 'SIGRTMAX-10'
              | 'SIGRTMAX-11'
              | 'SIGRTMAX-12'
              | 'SIGRTMAX-13'
              | 'SIGRTMAX-14'
              | 'SIGRTMAX-2'
              | 'SIGRTMAX-3'
              | 'SIGRTMAX-4'
              | 'SIGRTMAX-5'
              | 'SIGRTMAX-6'
              | 'SIGRTMAX-7'
              | 'SIGRTMAX-8'
              | 'SIGRTMAX-9'
              | 'SIGRTMIN'
              | 'SIGRTMIN+1'
              | 'SIGRTMIN+10'
              | 'SIGRTMIN+11'
              | 'SIGRTMIN+12'
              | 'SIGRTMIN+13'
              | 'SIGRTMIN+14'
              | 'SIGRTMIN+15'
              | 'SIGRTMIN+2'
              | 'SIGRTMIN+3'
              | 'SIGRTMIN+4'
              | 'SIGRTMIN+5'
              | 'SIGRTMIN+6'
              | 'SIGRTMIN+7'
              | 'SIGRTMIN+8'
              | 'SIGRTMIN+9'
              | 'SIGSEGV'
              | 'SIGSTKFLT'
              | 'SIGSTOP'
              | 'SIGSYS'
              | 'SIGTERM'
              | 'SIGTRAP'
              | 'SIGTSTP'
              | 'SIGTTIN'
              | 'SIGTTOU'
              | 'SIGURG'
              | 'SIGUSR1'
              | 'SIGUSR2'
              | 'SIGVTALRM'
              | 'SIGWINCH'
              | 'SIGXCPU'
              | 'SIGXFSZ';
            [k: string]: unknown;
          };
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          livenessProbe?: {
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              [k: string]: unknown;
            };
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                /**
                 * The header field value
                 */
                value: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            [k: string]: unknown;
          };
          /**
           * Restart policy for the container to manage the restart behavior of each container within a pod. This may only be set for init containers. You cannot set this field on ephemeral containers.
           */
          restartPolicy?: string;
          /**
           * SecurityContext holds security configuration that will be applied to a container. Some fields are present in both SecurityContext and PodSecurityContext.  When both are set, the values in SecurityContext take precedence.
           */
          securityContext?: {
            /**
             * SELinuxOptions are the labels to be applied to the container
             */
            seLinuxOptions?: {
              /**
               * Level is SELinux level label that applies to the container.
               */
              level?: string;
              /**
               * Role is a SELinux role label that applies to the container.
               */
              role?: string;
              /**
               * Type is a SELinux type label that applies to the container.
               */
              type?: string;
              /**
               * User is a SELinux user label that applies to the container.
               */
              user?: string;
              [k: string]: unknown;
            };
            /**
             * WindowsSecurityContextOptions contain Windows-specific options and credentials.
             */
            windowsOptions?: {
              /**
               * GMSACredentialSpec is where the GMSA admission webhook (https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the GMSA credential spec named by the GMSACredentialSpecName field.
               */
              gmsaCredentialSpec?: string;
              /**
               * GMSACredentialSpecName is the name of the GMSA credential spec to use.
               */
              gmsaCredentialSpecName?: string;
              /**
               * HostProcess determines if a container should be run as a 'Host Process' container. All of a Pod's containers must have the same effective HostProcess value (it is not allowed to have a mix of HostProcess containers and non-HostProcess containers). In addition, if HostProcess is true then HostNetwork must also be set to true.
               */
              hostProcess?: boolean;
              /**
               * The UserName in Windows to run the entrypoint of the container process. Defaults to the user specified in image metadata if unspecified. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
               */
              runAsUserName?: string;
              [k: string]: unknown;
            };
            /**
             * AllowPrivilegeEscalation controls whether a process can gain more privileges than its parent process. This bool directly controls if the no_new_privs flag will be set on the container process. AllowPrivilegeEscalation is true always when the container is: 1) run as Privileged 2) has CAP_SYS_ADMIN Note that this field cannot be set when spec.os.name is windows.
             */
            allowPrivilegeEscalation?: boolean;
            /**
             * AppArmorProfile defines a pod or container's AppArmor settings.
             */
            appArmorProfile?: {
              /**
               * type indicates which kind of AppArmor profile will be applied. Valid options are:
               *   Localhost - a profile pre-loaded on the node.
               *   RuntimeDefault - the container runtime's default profile.
               *   Unconfined - no AppArmor enforcement.
               *
               * Possible enum values:
               *  - `"Localhost"` indicates that a profile pre-loaded on the node should be used.
               *  - `"RuntimeDefault"` indicates that the container runtime's default AppArmor profile should be used.
               *  - `"Unconfined"` indicates that no AppArmor profile should be enforced.
               */
              type: 'Localhost' | 'RuntimeDefault' | 'Unconfined';
              /**
               * localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is "Localhost".
               */
              localhostProfile?: string;
              [k: string]: unknown;
            };
            /**
             * Adds and removes POSIX capabilities from running containers.
             */
            capabilities?: {
              /**
               * Added capabilities
               */
              add?: string[];
              /**
               * Removed capabilities
               */
              drop?: string[];
              [k: string]: unknown;
            };
            /**
             * procMount denotes the type of proc mount to use for the containers. The default value is Default which uses the container runtime defaults for readonly paths and masked paths. This requires the ProcMountType feature flag to be enabled. Note that this field cannot be set when spec.os.name is windows.
             *
             * Possible enum values:
             *  - `"Default"` uses the container runtime defaults for readonly and masked paths for /proc. Most container runtimes mask certain paths in /proc to avoid accidental security exposure of special devices or information.
             *  - `"Unmasked"` bypasses the default masking behavior of the container runtime and ensures the newly created /proc the container stays in tact with no modifications.
             */
            procMount?: 'Default' | 'Unmasked';
            /**
             * Whether this container has a read-only root filesystem. Default is false. Note that this field cannot be set when spec.os.name is windows.
             */
            readOnlyRootFilesystem?: boolean;
            /**
             * Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
             */
            runAsNonRoot?: boolean;
            /**
             * The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
             */
            runAsUser?: number;
            /**
             * SeccompProfile defines a pod/container's seccomp profile settings. Only one profile source may be set.
             */
            seccompProfile?: {
              /**
               * localhostProfile indicates a profile defined in a file on the node should be used. The profile must be preconfigured on the node to work. Must be a descending path, relative to the kubelet's configured seccomp profile location. Must be set if type is "Localhost". Must NOT be set for any other type.
               */
              localhostProfile?: string;
              /**
               * type indicates which kind of seccomp profile will be applied. Valid options are:
               *
               * Localhost - a profile defined in a file on the node should be used. RuntimeDefault - the container runtime default profile should be used. Unconfined - no profile should be applied.
               *
               * Possible enum values:
               *  - `"Localhost"` indicates a profile defined in a file on the node should be used. The file's location relative to <kubelet-root-dir>/seccomp.
               *  - `"RuntimeDefault"` represents the default container runtime seccomp profile.
               *  - `"Unconfined"` indicates no seccomp profile is applied (A.K.A. unconfined).
               */
              type: 'Localhost' | 'RuntimeDefault' | 'Unconfined';
              [k: string]: unknown;
            };
            /**
             * Run container in privileged mode. Processes in privileged containers are essentially equivalent to root on the host. Defaults to false. Note that this field cannot be set when spec.os.name is windows.
             */
            privileged?: boolean;
            /**
             * The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
             */
            runAsGroup?: number;
            [k: string]: unknown;
          };
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          startupProbe?: {
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                /**
                 * The header field value
                 */
                value: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              [k: string]: unknown;
            };
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            [k: string]: unknown;
          };
          /**
           * Whether this container should allocate a TTY for itself, also requires 'stdin' to be true. Default is false.
           */
          tty?: boolean;
          /**
           * Entrypoint array. Not executed within a shell. The image's ENTRYPOINT is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
           */
          command?: string[];
          /**
           * List of sources to populate environment variables in the container. The keys defined within a source must be a C_IDENTIFIER. All invalid keys will be reported as an event when the container is starting. When a key exists in multiple sources, the value associated with the last source will take precedence. Values defined by an Env with a duplicate key will take precedence. Cannot be updated.
           */
          envFrom?: {
            /**
             * ConfigMapEnvSource selects a ConfigMap to populate the environment variables with.
             *
             * The contents of the target ConfigMap's Data field will represent the key-value pairs as environment variables.
             */
            configMapRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              /**
               * Specify whether the ConfigMap must be defined
               */
              optional?: boolean;
              [k: string]: unknown;
            };
            /**
             * Optional text to prepend to the name of each environment variable. Must be a C_IDENTIFIER.
             */
            prefix?: string;
            /**
             * SecretEnvSource selects a Secret to populate the environment variables with.
             *
             * The contents of the target Secret's Data field will represent the key-value pairs as environment variables.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              /**
               * Specify whether the Secret must be defined
               */
              optional?: boolean;
              [k: string]: unknown;
            };
            [k: string]: unknown;
          }[];
          /**
           * Container image name. More info: https://kubernetes.io/docs/concepts/containers/images
           */
          image?: string;
          /**
           * Ports are not allowed for ephemeral containers.
           */
          ports?: {
            /**
             * Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.
             */
            containerPort: number;
            /**
             * What host IP to bind the external port to.
             */
            hostIP?: string;
            /**
             * Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this.
             */
            hostPort?: number;
            /**
             * If specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services.
             */
            name?: string;
            /**
             * Protocol for port. Must be UDP, TCP, or SCTP. Defaults to "TCP".
             *
             * Possible enum values:
             *  - `"SCTP"` is the SCTP protocol.
             *  - `"TCP"` is the TCP protocol.
             *  - `"UDP"` is the UDP protocol.
             */
            protocol?: 'SCTP' | 'TCP' | 'UDP';
            [k: string]: unknown;
          }[];
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          readinessProbe?: {
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                /**
                 * The header field value
                 */
                value: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              [k: string]: unknown;
            };
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            [k: string]: unknown;
          };
          /**
           * Whether the container runtime should close the stdin channel after it has been opened by a single attach. When stdin is true the stdin stream will remain open across multiple attach sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the first client attaches to stdin, and then remains open and accepts data until the client disconnects, at which time stdin is closed and remains closed until the container is restarted. If this flag is false, a container processes that reads from stdin will never receive an EOF. Default is false
           */
          stdinOnce?: boolean;
          /**
           * Name of the ephemeral container specified as a DNS_LABEL. This name must be unique among all containers, init containers and ephemeral containers.
           */
          name: string;
          /**
           * Resources resize policy for the container.
           */
          resizePolicy?: {
            /**
             * Name of the resource to which this resource resize policy applies. Supported values: cpu, memory.
             */
            resourceName: string;
            /**
             * Restart policy to apply when specified resource is resized. If not specified, it defaults to NotRequired.
             */
            restartPolicy: string;
            [k: string]: unknown;
          }[];
          /**
           * ResourceRequirements describes the compute resource requirements.
           */
          resources?: {
            /**
             * Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container.
             *
             * This is an alpha field and requires enabling the DynamicResourceAllocation feature gate.
             *
             * This field is immutable. It can only be set for containers.
             */
            claims?: {
              /**
               * Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
               */
              name: string;
              /**
               * Request is the name chosen for a request in the referenced claim. If empty, everything from the claim is made available, otherwise only the result of this request.
               */
              request?: string;
              [k: string]: unknown;
            }[];
            /**
             * Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
             */
            limits?: {
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
             * Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
             */
            requests?: {
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
            [k: string]: unknown;
          };
          /**
           * Indicate how the termination message should be populated. File will use the contents of terminationMessagePath to populate the container status message on both success and failure. FallbackToLogsOnError will use the last chunk of container log output if the termination message file is empty and the container exited with an error. The log output is limited to 2048 bytes or 80 lines, whichever is smaller. Defaults to File. Cannot be updated.
           *
           * Possible enum values:
           *  - `"FallbackToLogsOnError"` will read the most recent contents of the container logs for the container status message when the container exits with an error and the terminationMessagePath has no contents.
           *  - `"File"` is the default behavior and will set the container status message to the contents of the container's terminationMessagePath when the container exits.
           */
          terminationMessagePolicy?: 'FallbackToLogsOnError' | 'File';
          /**
           * volumeDevices is the list of block devices to be used by the container.
           */
          volumeDevices?: {
            /**
             * devicePath is the path inside of the container that the device will be mapped to.
             */
            devicePath: string;
            /**
             * name must match the name of a persistentVolumeClaim in the pod
             */
            name: string;
            [k: string]: unknown;
          }[];
          /**
           * List of environment variables to set in the container. Cannot be updated.
           */
          env?: {
            /**
             * Name of the environment variable. Must be a C_IDENTIFIER.
             */
            name: string;
            /**
             * Variable references $(VAR_NAME) are expanded using the previously defined environment variables in the container and any service environment variables. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Defaults to "".
             */
            value?: string;
            /**
             * EnvVarSource represents a source for the value of an EnvVar.
             */
            valueFrom?: {
              /**
               * SecretKeySelector selects a key of a Secret.
               */
              secretKeyRef?: {
                /**
                 * The key of the secret to select from.  Must be a valid secret key.
                 */
                key: string;
                /**
                 * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                 */
                name?: string;
                /**
                 * Specify whether the Secret or its key must be defined
                 */
                optional?: boolean;
                [k: string]: unknown;
              };
              /**
               * Selects a key from a ConfigMap.
               */
              configMapKeyRef?: {
                /**
                 * The key to select.
                 */
                key: string;
                /**
                 * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                 */
                name?: string;
                /**
                 * Specify whether the ConfigMap or its key must be defined
                 */
                optional?: boolean;
                [k: string]: unknown;
              };
              /**
               * ObjectFieldSelector selects an APIVersioned field of an object.
               */
              fieldRef?: {
                /**
                 * Version of the schema the FieldPath is written in terms of, defaults to "v1".
                 */
                apiVersion?: string;
                /**
                 * Path of the field to select in the specified API version.
                 */
                fieldPath: string;
                [k: string]: unknown;
              };
              /**
               * ResourceFieldSelector represents container resources (cpu, memory) and their output format
               */
              resourceFieldRef?: {
                /**
                 * Required: resource to select
                 */
                resource: string;
                /**
                 * Container name: required for volumes, optional for env vars
                 */
                containerName?: string;
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
                divisor?: string | number;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            [k: string]: unknown;
          }[];
          /**
           * Whether this container should allocate a buffer for stdin in the container runtime. If this is not set, reads from stdin in the container will always result in EOF. Default is false.
           */
          stdin?: boolean;
          /**
           * If set, the name of the container from PodSpec that this ephemeral container targets. The ephemeral container will be run in the namespaces (IPC, PID, etc) of this container. If not set then the ephemeral container uses the namespaces configured in the Pod spec.
           *
           * The container runtime must implement support for this feature. If the runtime does not support namespace targeting then the result of setting this field is undefined.
           */
          targetContainerName?: string;
          /**
           * Optional: Path at which the file to which the container's termination message will be written is mounted into the container's filesystem. Message written is intended to be brief final status, such as an assertion failure message. Will be truncated by the node if greater than 4096 bytes. The total message length across all containers will be limited to 12kb. Defaults to /dev/termination-log. Cannot be updated.
           */
          terminationMessagePath?: string;
          /**
           * Pod volumes to mount into the container's filesystem. Subpath mounts are not allowed for ephemeral containers. Cannot be updated.
           */
          volumeMounts?: {
            /**
             * Path within the container at which the volume should be mounted.  Must not contain ':'.
             */
            mountPath: string;
            /**
             * mountPropagation determines how mounts are propagated from the host to container and the other way around. When not set, MountPropagationNone is used. This field is beta in 1.10. When RecursiveReadOnly is set to IfPossible or to Enabled, MountPropagation must be None or unspecified (which defaults to None).
             *
             * Possible enum values:
             *  - `"Bidirectional"` means that the volume in a container will receive new mounts from the host or other containers, and its own mounts will be propagated from the container to the host or other containers. Note that this mode is recursively applied to all mounts in the volume ("rshared" in Linux terminology).
             *  - `"HostToContainer"` means that the volume in a container will receive new mounts from the host or other containers, but filesystems mounted inside the container won't be propagated to the host or other containers. Note that this mode is recursively applied to all mounts in the volume ("rslave" in Linux terminology).
             *  - `"None"` means that the volume in a container will not receive new mounts from the host or other containers, and filesystems mounted inside the container won't be propagated to the host or other containers. Note that this mode corresponds to "private" in Linux terminology.
             */
            mountPropagation?: 'Bidirectional' | 'HostToContainer' | 'None';
            /**
             * This must match the Name of a Volume.
             */
            name: string;
            /**
             * Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false.
             */
            readOnly?: boolean;
            /**
             * RecursiveReadOnly specifies whether read-only mounts should be handled recursively.
             *
             * If ReadOnly is false, this field has no meaning and must be unspecified.
             *
             * If ReadOnly is true, and this field is set to Disabled, the mount is not made recursively read-only.  If this field is set to IfPossible, the mount is made recursively read-only, if it is supported by the container runtime.  If this field is set to Enabled, the mount is made recursively read-only if it is supported by the container runtime, otherwise the pod will not be started and an error will be generated to indicate the reason.
             *
             * If this field is set to IfPossible or Enabled, MountPropagation must be set to None (or be unspecified, which defaults to None).
             *
             * If this field is not specified, it is treated as an equivalent of Disabled.
             */
            recursiveReadOnly?: string;
            /**
             * Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root).
             */
            subPath?: string;
            /**
             * Expanded path within the volume from which the container's volume should be mounted. Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment. Defaults to "" (volume's root). SubPathExpr and SubPath are mutually exclusive.
             */
            subPathExpr?: string;
            [k: string]: unknown;
          }[];
          /**
           * Container's working directory. If not specified, the container runtime's default will be used, which might be configured in the container image. Cannot be updated.
           */
          workingDir?: string;
          [k: string]: unknown;
        }[];
        /**
         * NodeName indicates in which node this pod is scheduled. If empty, this pod is a candidate for scheduling by the scheduler defined in schedulerName. Once this field is set, the kubelet for this node becomes responsible for the lifecycle of this pod. This field should not be used to express a desire for the pod to be scheduled on a specific node. https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#nodename
         */
        nodeName?: string;
        /**
         * Overhead represents the resource overhead associated with running a pod for a given RuntimeClass. This field will be autopopulated at admission time by the RuntimeClass admission controller. If the RuntimeClass admission controller is enabled, overhead must not be set in Pod create requests. The RuntimeClass admission controller will reject Pod create requests which have the overhead already set. If RuntimeClass is configured and selected in the PodSpec, Overhead will be set to the value defined in the corresponding RuntimeClass, otherwise it will remain unset and treated as zero. More info: https://git.k8s.io/enhancements/keps/sig-node/688-pod-overhead/README.md
         */
        overhead?: {
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
         * Restart policy for all containers within the pod. One of Always, OnFailure, Never. In some contexts, only a subset of those values may be permitted. Default to Always. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy
         *
         * Possible enum values:
         *  - `"Always"`
         *  - `"Never"`
         *  - `"OnFailure"`
         */
        restartPolicy?: 'Always' | 'Never' | 'OnFailure';
        /**
         * DeprecatedServiceAccount is a deprecated alias for ServiceAccountName. Deprecated: Use serviceAccountName instead.
         */
        serviceAccount?: string;
        /**
         * Share a single process namespace between all of the containers in a pod. When this is set containers will be able to view and signal processes from other containers in the same pod, and the first process in each container will not be assigned PID 1. HostPID and ShareProcessNamespace cannot both be set. Optional: Default to false.
         */
        shareProcessNamespace?: boolean;
        /**
         * List of containers belonging to the pod. Containers cannot currently be added or removed. There must be at least one container in a Pod. Cannot be updated.
         */
        containers: {
          /**
           * RestartPolicy defines the restart behavior of individual containers in a pod. This field may only be set for init containers, and the only allowed value is "Always". For non-init containers or when this field is not specified, the restart behavior is defined by the Pod's restart policy and the container type. Setting the RestartPolicy as "Always" for the init container will have the following effect: this init container will be continually restarted on exit until all regular containers have terminated. Once all regular containers have completed, all init containers with restartPolicy "Always" will be shut down. This lifecycle differs from normal init containers and is often referred to as a "sidecar" container. Although this init container still starts in the init container sequence, it does not wait for the container to complete before proceeding to the next init container. Instead, the next init container starts immediately after this init container is started, or after any startupProbe has successfully completed.
           */
          restartPolicy?: string;
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          startupProbe?: {
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              [k: string]: unknown;
            };
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field value
                 */
                value: string;
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              [k: string]: unknown;
            };
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            [k: string]: unknown;
          };
          /**
           * Pod volumes to mount into the container's filesystem. Cannot be updated.
           */
          volumeMounts?: {
            /**
             * This must match the Name of a Volume.
             */
            name: string;
            /**
             * Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false.
             */
            readOnly?: boolean;
            /**
             * RecursiveReadOnly specifies whether read-only mounts should be handled recursively.
             *
             * If ReadOnly is false, this field has no meaning and must be unspecified.
             *
             * If ReadOnly is true, and this field is set to Disabled, the mount is not made recursively read-only.  If this field is set to IfPossible, the mount is made recursively read-only, if it is supported by the container runtime.  If this field is set to Enabled, the mount is made recursively read-only if it is supported by the container runtime, otherwise the pod will not be started and an error will be generated to indicate the reason.
             *
             * If this field is set to IfPossible or Enabled, MountPropagation must be set to None (or be unspecified, which defaults to None).
             *
             * If this field is not specified, it is treated as an equivalent of Disabled.
             */
            recursiveReadOnly?: string;
            /**
             * Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root).
             */
            subPath?: string;
            /**
             * Expanded path within the volume from which the container's volume should be mounted. Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment. Defaults to "" (volume's root). SubPathExpr and SubPath are mutually exclusive.
             */
            subPathExpr?: string;
            /**
             * Path within the container at which the volume should be mounted.  Must not contain ':'.
             */
            mountPath: string;
            /**
             * mountPropagation determines how mounts are propagated from the host to container and the other way around. When not set, MountPropagationNone is used. This field is beta in 1.10. When RecursiveReadOnly is set to IfPossible or to Enabled, MountPropagation must be None or unspecified (which defaults to None).
             *
             * Possible enum values:
             *  - `"Bidirectional"` means that the volume in a container will receive new mounts from the host or other containers, and its own mounts will be propagated from the container to the host or other containers. Note that this mode is recursively applied to all mounts in the volume ("rshared" in Linux terminology).
             *  - `"HostToContainer"` means that the volume in a container will receive new mounts from the host or other containers, but filesystems mounted inside the container won't be propagated to the host or other containers. Note that this mode is recursively applied to all mounts in the volume ("rslave" in Linux terminology).
             *  - `"None"` means that the volume in a container will not receive new mounts from the host or other containers, and filesystems mounted inside the container won't be propagated to the host or other containers. Note that this mode corresponds to "private" in Linux terminology.
             */
            mountPropagation?: 'Bidirectional' | 'HostToContainer' | 'None';
            [k: string]: unknown;
          }[];
          /**
           * List of environment variables to set in the container. Cannot be updated.
           */
          env?: {
            /**
             * Name of the environment variable. Must be a C_IDENTIFIER.
             */
            name: string;
            /**
             * Variable references $(VAR_NAME) are expanded using the previously defined environment variables in the container and any service environment variables. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Defaults to "".
             */
            value?: string;
            /**
             * EnvVarSource represents a source for the value of an EnvVar.
             */
            valueFrom?: {
              /**
               * Selects a key from a ConfigMap.
               */
              configMapKeyRef?: {
                /**
                 * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                 */
                name?: string;
                /**
                 * Specify whether the ConfigMap or its key must be defined
                 */
                optional?: boolean;
                /**
                 * The key to select.
                 */
                key: string;
                [k: string]: unknown;
              };
              /**
               * ObjectFieldSelector selects an APIVersioned field of an object.
               */
              fieldRef?: {
                /**
                 * Version of the schema the FieldPath is written in terms of, defaults to "v1".
                 */
                apiVersion?: string;
                /**
                 * Path of the field to select in the specified API version.
                 */
                fieldPath: string;
                [k: string]: unknown;
              };
              /**
               * ResourceFieldSelector represents container resources (cpu, memory) and their output format
               */
              resourceFieldRef?: {
                /**
                 * Container name: required for volumes, optional for env vars
                 */
                containerName?: string;
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
                divisor?: string | number;
                /**
                 * Required: resource to select
                 */
                resource: string;
                [k: string]: unknown;
              };
              /**
               * SecretKeySelector selects a key of a Secret.
               */
              secretKeyRef?: {
                /**
                 * The key of the secret to select from.  Must be a valid secret key.
                 */
                key: string;
                /**
                 * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                 */
                name?: string;
                /**
                 * Specify whether the Secret or its key must be defined
                 */
                optional?: boolean;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            [k: string]: unknown;
          }[];
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          livenessProbe?: {
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                /**
                 * The header field value
                 */
                value: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              [k: string]: unknown;
            };
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            [k: string]: unknown;
          };
          /**
           * Resources resize policy for the container.
           */
          resizePolicy?: {
            /**
             * Restart policy to apply when specified resource is resized. If not specified, it defaults to NotRequired.
             */
            restartPolicy: string;
            /**
             * Name of the resource to which this resource resize policy applies. Supported values: cpu, memory.
             */
            resourceName: string;
            [k: string]: unknown;
          }[];
          /**
           * SecurityContext holds security configuration that will be applied to a container. Some fields are present in both SecurityContext and PodSecurityContext.  When both are set, the values in SecurityContext take precedence.
           */
          securityContext?: {
            /**
             * WindowsSecurityContextOptions contain Windows-specific options and credentials.
             */
            windowsOptions?: {
              /**
               * GMSACredentialSpec is where the GMSA admission webhook (https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the GMSA credential spec named by the GMSACredentialSpecName field.
               */
              gmsaCredentialSpec?: string;
              /**
               * GMSACredentialSpecName is the name of the GMSA credential spec to use.
               */
              gmsaCredentialSpecName?: string;
              /**
               * HostProcess determines if a container should be run as a 'Host Process' container. All of a Pod's containers must have the same effective HostProcess value (it is not allowed to have a mix of HostProcess containers and non-HostProcess containers). In addition, if HostProcess is true then HostNetwork must also be set to true.
               */
              hostProcess?: boolean;
              /**
               * The UserName in Windows to run the entrypoint of the container process. Defaults to the user specified in image metadata if unspecified. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
               */
              runAsUserName?: string;
              [k: string]: unknown;
            };
            /**
             * Adds and removes POSIX capabilities from running containers.
             */
            capabilities?: {
              /**
               * Added capabilities
               */
              add?: string[];
              /**
               * Removed capabilities
               */
              drop?: string[];
              [k: string]: unknown;
            };
            /**
             * procMount denotes the type of proc mount to use for the containers. The default value is Default which uses the container runtime defaults for readonly paths and masked paths. This requires the ProcMountType feature flag to be enabled. Note that this field cannot be set when spec.os.name is windows.
             *
             * Possible enum values:
             *  - `"Default"` uses the container runtime defaults for readonly and masked paths for /proc. Most container runtimes mask certain paths in /proc to avoid accidental security exposure of special devices or information.
             *  - `"Unmasked"` bypasses the default masking behavior of the container runtime and ensures the newly created /proc the container stays in tact with no modifications.
             */
            procMount?: 'Default' | 'Unmasked';
            /**
             * Whether this container has a read-only root filesystem. Default is false. Note that this field cannot be set when spec.os.name is windows.
             */
            readOnlyRootFilesystem?: boolean;
            /**
             * The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
             */
            runAsGroup?: number;
            /**
             * Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
             */
            runAsNonRoot?: boolean;
            /**
             * The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
             */
            runAsUser?: number;
            /**
             * SELinuxOptions are the labels to be applied to the container
             */
            seLinuxOptions?: {
              /**
               * Level is SELinux level label that applies to the container.
               */
              level?: string;
              /**
               * Role is a SELinux role label that applies to the container.
               */
              role?: string;
              /**
               * Type is a SELinux type label that applies to the container.
               */
              type?: string;
              /**
               * User is a SELinux user label that applies to the container.
               */
              user?: string;
              [k: string]: unknown;
            };
            /**
             * SeccompProfile defines a pod/container's seccomp profile settings. Only one profile source may be set.
             */
            seccompProfile?: {
              /**
               * localhostProfile indicates a profile defined in a file on the node should be used. The profile must be preconfigured on the node to work. Must be a descending path, relative to the kubelet's configured seccomp profile location. Must be set if type is "Localhost". Must NOT be set for any other type.
               */
              localhostProfile?: string;
              /**
               * type indicates which kind of seccomp profile will be applied. Valid options are:
               *
               * Localhost - a profile defined in a file on the node should be used. RuntimeDefault - the container runtime default profile should be used. Unconfined - no profile should be applied.
               *
               * Possible enum values:
               *  - `"Localhost"` indicates a profile defined in a file on the node should be used. The file's location relative to <kubelet-root-dir>/seccomp.
               *  - `"RuntimeDefault"` represents the default container runtime seccomp profile.
               *  - `"Unconfined"` indicates no seccomp profile is applied (A.K.A. unconfined).
               */
              type: 'Localhost' | 'RuntimeDefault' | 'Unconfined';
              [k: string]: unknown;
            };
            /**
             * AllowPrivilegeEscalation controls whether a process can gain more privileges than its parent process. This bool directly controls if the no_new_privs flag will be set on the container process. AllowPrivilegeEscalation is true always when the container is: 1) run as Privileged 2) has CAP_SYS_ADMIN Note that this field cannot be set when spec.os.name is windows.
             */
            allowPrivilegeEscalation?: boolean;
            /**
             * AppArmorProfile defines a pod or container's AppArmor settings.
             */
            appArmorProfile?: {
              /**
               * localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is "Localhost".
               */
              localhostProfile?: string;
              /**
               * type indicates which kind of AppArmor profile will be applied. Valid options are:
               *   Localhost - a profile pre-loaded on the node.
               *   RuntimeDefault - the container runtime's default profile.
               *   Unconfined - no AppArmor enforcement.
               *
               * Possible enum values:
               *  - `"Localhost"` indicates that a profile pre-loaded on the node should be used.
               *  - `"RuntimeDefault"` indicates that the container runtime's default AppArmor profile should be used.
               *  - `"Unconfined"` indicates that no AppArmor profile should be enforced.
               */
              type: 'Localhost' | 'RuntimeDefault' | 'Unconfined';
              [k: string]: unknown;
            };
            /**
             * Run container in privileged mode. Processes in privileged containers are essentially equivalent to root on the host. Defaults to false. Note that this field cannot be set when spec.os.name is windows.
             */
            privileged?: boolean;
            [k: string]: unknown;
          };
          /**
           * Whether the container runtime should close the stdin channel after it has been opened by a single attach. When stdin is true the stdin stream will remain open across multiple attach sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the first client attaches to stdin, and then remains open and accepts data until the client disconnects, at which time stdin is closed and remains closed until the container is restarted. If this flag is false, a container processes that reads from stdin will never receive an EOF. Default is false
           */
          stdinOnce?: boolean;
          /**
           * Whether this container should allocate a TTY for itself, also requires 'stdin' to be true. Default is false.
           */
          tty?: boolean;
          /**
           * volumeDevices is the list of block devices to be used by the container.
           */
          volumeDevices?: {
            /**
             * devicePath is the path inside of the container that the device will be mapped to.
             */
            devicePath: string;
            /**
             * name must match the name of a persistentVolumeClaim in the pod
             */
            name: string;
            [k: string]: unknown;
          }[];
          /**
           * Container's working directory. If not specified, the container runtime's default will be used, which might be configured in the container image. Cannot be updated.
           */
          workingDir?: string;
          /**
           * Entrypoint array. Not executed within a shell. The container image's ENTRYPOINT is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
           */
          command?: string[];
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          readinessProbe?: {
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              [k: string]: unknown;
            };
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                /**
                 * The header field value
                 */
                value: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              [k: string]: unknown;
            };
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            [k: string]: unknown;
          };
          /**
           * Indicate how the termination message should be populated. File will use the contents of terminationMessagePath to populate the container status message on both success and failure. FallbackToLogsOnError will use the last chunk of container log output if the termination message file is empty and the container exited with an error. The log output is limited to 2048 bytes or 80 lines, whichever is smaller. Defaults to File. Cannot be updated.
           *
           * Possible enum values:
           *  - `"FallbackToLogsOnError"` will read the most recent contents of the container logs for the container status message when the container exits with an error and the terminationMessagePath has no contents.
           *  - `"File"` is the default behavior and will set the container status message to the contents of the container's terminationMessagePath when the container exits.
           */
          terminationMessagePolicy?: 'FallbackToLogsOnError' | 'File';
          /**
           * List of sources to populate environment variables in the container. The keys defined within a source must be a C_IDENTIFIER. All invalid keys will be reported as an event when the container is starting. When a key exists in multiple sources, the value associated with the last source will take precedence. Values defined by an Env with a duplicate key will take precedence. Cannot be updated.
           */
          envFrom?: {
            /**
             * ConfigMapEnvSource selects a ConfigMap to populate the environment variables with.
             *
             * The contents of the target ConfigMap's Data field will represent the key-value pairs as environment variables.
             */
            configMapRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              /**
               * Specify whether the ConfigMap must be defined
               */
              optional?: boolean;
              [k: string]: unknown;
            };
            /**
             * Optional text to prepend to the name of each environment variable. Must be a C_IDENTIFIER.
             */
            prefix?: string;
            /**
             * SecretEnvSource selects a Secret to populate the environment variables with.
             *
             * The contents of the target Secret's Data field will represent the key-value pairs as environment variables.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              /**
               * Specify whether the Secret must be defined
               */
              optional?: boolean;
              [k: string]: unknown;
            };
            [k: string]: unknown;
          }[];
          /**
           * Container image name. More info: https://kubernetes.io/docs/concepts/containers/images This field is optional to allow higher level config management to default or override container images in workload controllers like Deployments and StatefulSets.
           */
          image?: string;
          /**
           * List of ports to expose from the container. Not specifying a port here DOES NOT prevent that port from being exposed. Any port which is listening on the default "0.0.0.0" address inside a container will be accessible from the network. Modifying this array with strategic merge patch may corrupt the data. For more information See https://github.com/kubernetes/kubernetes/issues/108255. Cannot be updated.
           */
          ports?: {
            /**
             * Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this.
             */
            hostPort?: number;
            /**
             * If specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services.
             */
            name?: string;
            /**
             * Protocol for port. Must be UDP, TCP, or SCTP. Defaults to "TCP".
             *
             * Possible enum values:
             *  - `"SCTP"` is the SCTP protocol.
             *  - `"TCP"` is the TCP protocol.
             *  - `"UDP"` is the UDP protocol.
             */
            protocol?: 'SCTP' | 'TCP' | 'UDP';
            /**
             * Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.
             */
            containerPort: number;
            /**
             * What host IP to bind the external port to.
             */
            hostIP?: string;
            [k: string]: unknown;
          }[];
          /**
           * Whether this container should allocate a buffer for stdin in the container runtime. If this is not set, reads from stdin in the container will always result in EOF. Default is false.
           */
          stdin?: boolean;
          /**
           * Optional: Path at which the file to which the container's termination message will be written is mounted into the container's filesystem. Message written is intended to be brief final status, such as an assertion failure message. Will be truncated by the node if greater than 4096 bytes. The total message length across all containers will be limited to 12kb. Defaults to /dev/termination-log. Cannot be updated.
           */
          terminationMessagePath?: string;
          /**
           * Arguments to the entrypoint. The container image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
           */
          args?: string[];
          /**
           * Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
           *
           * Possible enum values:
           *  - `"Always"` means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.
           *  - `"IfNotPresent"` means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
           *  - `"Never"` means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn't present
           */
          imagePullPolicy?: 'Always' | 'IfNotPresent' | 'Never';
          /**
           * Lifecycle describes actions that the management system should take in response to container lifecycle events. For the PostStart and PreStop lifecycle handlers, management of the container blocks until the action is complete, unless the container process fails, in which case the handler is aborted.
           */
          lifecycle?: {
            /**
             * LifecycleHandler defines a specific action that should be taken in a lifecycle hook. One and only one of the fields, except TCPSocket must be specified.
             */
            postStart?: {
              /**
               * ExecAction describes a "run in container" action.
               */
              exec?: {
                /**
                 * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
                 */
                command?: string[];
                [k: string]: unknown;
              };
              /**
               * HTTPGetAction describes an action based on HTTP Get requests.
               */
              httpGet?: {
                /**
                 * Custom headers to set in the request. HTTP allows repeated headers.
                 */
                httpHeaders?: {
                  /**
                   * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                   */
                  name: string;
                  /**
                   * The header field value
                   */
                  value: string;
                  [k: string]: unknown;
                }[];
                /**
                 * Path to access on the HTTP server.
                 */
                path?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                /**
                 * Scheme to use for connecting to the host. Defaults to HTTP.
                 *
                 * Possible enum values:
                 *  - `"HTTP"` means that the scheme used will be http://
                 *  - `"HTTPS"` means that the scheme used will be https://
                 */
                scheme?: 'HTTP' | 'HTTPS';
                /**
                 * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
                 */
                host?: string;
                [k: string]: unknown;
              };
              /**
               * SleepAction describes a "sleep" action.
               */
              sleep?: {
                /**
                 * Seconds is the number of seconds to sleep.
                 */
                seconds: number;
                [k: string]: unknown;
              };
              /**
               * TCPSocketAction describes an action based on opening a socket
               */
              tcpSocket?: {
                /**
                 * Optional: Host name to connect to, defaults to the pod IP.
                 */
                host?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            /**
             * LifecycleHandler defines a specific action that should be taken in a lifecycle hook. One and only one of the fields, except TCPSocket must be specified.
             */
            preStop?: {
              /**
               * SleepAction describes a "sleep" action.
               */
              sleep?: {
                /**
                 * Seconds is the number of seconds to sleep.
                 */
                seconds: number;
                [k: string]: unknown;
              };
              /**
               * TCPSocketAction describes an action based on opening a socket
               */
              tcpSocket?: {
                /**
                 * Optional: Host name to connect to, defaults to the pod IP.
                 */
                host?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                [k: string]: unknown;
              };
              /**
               * ExecAction describes a "run in container" action.
               */
              exec?: {
                /**
                 * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
                 */
                command?: string[];
                [k: string]: unknown;
              };
              /**
               * HTTPGetAction describes an action based on HTTP Get requests.
               */
              httpGet?: {
                /**
                 * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
                 */
                host?: string;
                /**
                 * Custom headers to set in the request. HTTP allows repeated headers.
                 */
                httpHeaders?: {
                  /**
                   * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                   */
                  name: string;
                  /**
                   * The header field value
                   */
                  value: string;
                  [k: string]: unknown;
                }[];
                /**
                 * Path to access on the HTTP server.
                 */
                path?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                /**
                 * Scheme to use for connecting to the host. Defaults to HTTP.
                 *
                 * Possible enum values:
                 *  - `"HTTP"` means that the scheme used will be http://
                 *  - `"HTTPS"` means that the scheme used will be https://
                 */
                scheme?: 'HTTP' | 'HTTPS';
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            /**
             * StopSignal defines which signal will be sent to a container when it is being stopped. If not specified, the default is defined by the container runtime in use. StopSignal can only be set for Pods with a non-empty .spec.os.name
             *
             * Possible enum values:
             *  - `"SIGABRT"`
             *  - `"SIGALRM"`
             *  - `"SIGBUS"`
             *  - `"SIGCHLD"`
             *  - `"SIGCLD"`
             *  - `"SIGCONT"`
             *  - `"SIGFPE"`
             *  - `"SIGHUP"`
             *  - `"SIGILL"`
             *  - `"SIGINT"`
             *  - `"SIGIO"`
             *  - `"SIGIOT"`
             *  - `"SIGKILL"`
             *  - `"SIGPIPE"`
             *  - `"SIGPOLL"`
             *  - `"SIGPROF"`
             *  - `"SIGPWR"`
             *  - `"SIGQUIT"`
             *  - `"SIGRTMAX"`
             *  - `"SIGRTMAX-1"`
             *  - `"SIGRTMAX-10"`
             *  - `"SIGRTMAX-11"`
             *  - `"SIGRTMAX-12"`
             *  - `"SIGRTMAX-13"`
             *  - `"SIGRTMAX-14"`
             *  - `"SIGRTMAX-2"`
             *  - `"SIGRTMAX-3"`
             *  - `"SIGRTMAX-4"`
             *  - `"SIGRTMAX-5"`
             *  - `"SIGRTMAX-6"`
             *  - `"SIGRTMAX-7"`
             *  - `"SIGRTMAX-8"`
             *  - `"SIGRTMAX-9"`
             *  - `"SIGRTMIN"`
             *  - `"SIGRTMIN+1"`
             *  - `"SIGRTMIN+10"`
             *  - `"SIGRTMIN+11"`
             *  - `"SIGRTMIN+12"`
             *  - `"SIGRTMIN+13"`
             *  - `"SIGRTMIN+14"`
             *  - `"SIGRTMIN+15"`
             *  - `"SIGRTMIN+2"`
             *  - `"SIGRTMIN+3"`
             *  - `"SIGRTMIN+4"`
             *  - `"SIGRTMIN+5"`
             *  - `"SIGRTMIN+6"`
             *  - `"SIGRTMIN+7"`
             *  - `"SIGRTMIN+8"`
             *  - `"SIGRTMIN+9"`
             *  - `"SIGSEGV"`
             *  - `"SIGSTKFLT"`
             *  - `"SIGSTOP"`
             *  - `"SIGSYS"`
             *  - `"SIGTERM"`
             *  - `"SIGTRAP"`
             *  - `"SIGTSTP"`
             *  - `"SIGTTIN"`
             *  - `"SIGTTOU"`
             *  - `"SIGURG"`
             *  - `"SIGUSR1"`
             *  - `"SIGUSR2"`
             *  - `"SIGVTALRM"`
             *  - `"SIGWINCH"`
             *  - `"SIGXCPU"`
             *  - `"SIGXFSZ"`
             */
            stopSignal?:
              | 'SIGABRT'
              | 'SIGALRM'
              | 'SIGBUS'
              | 'SIGCHLD'
              | 'SIGCLD'
              | 'SIGCONT'
              | 'SIGFPE'
              | 'SIGHUP'
              | 'SIGILL'
              | 'SIGINT'
              | 'SIGIO'
              | 'SIGIOT'
              | 'SIGKILL'
              | 'SIGPIPE'
              | 'SIGPOLL'
              | 'SIGPROF'
              | 'SIGPWR'
              | 'SIGQUIT'
              | 'SIGRTMAX'
              | 'SIGRTMAX-1'
              | 'SIGRTMAX-10'
              | 'SIGRTMAX-11'
              | 'SIGRTMAX-12'
              | 'SIGRTMAX-13'
              | 'SIGRTMAX-14'
              | 'SIGRTMAX-2'
              | 'SIGRTMAX-3'
              | 'SIGRTMAX-4'
              | 'SIGRTMAX-5'
              | 'SIGRTMAX-6'
              | 'SIGRTMAX-7'
              | 'SIGRTMAX-8'
              | 'SIGRTMAX-9'
              | 'SIGRTMIN'
              | 'SIGRTMIN+1'
              | 'SIGRTMIN+10'
              | 'SIGRTMIN+11'
              | 'SIGRTMIN+12'
              | 'SIGRTMIN+13'
              | 'SIGRTMIN+14'
              | 'SIGRTMIN+15'
              | 'SIGRTMIN+2'
              | 'SIGRTMIN+3'
              | 'SIGRTMIN+4'
              | 'SIGRTMIN+5'
              | 'SIGRTMIN+6'
              | 'SIGRTMIN+7'
              | 'SIGRTMIN+8'
              | 'SIGRTMIN+9'
              | 'SIGSEGV'
              | 'SIGSTKFLT'
              | 'SIGSTOP'
              | 'SIGSYS'
              | 'SIGTERM'
              | 'SIGTRAP'
              | 'SIGTSTP'
              | 'SIGTTIN'
              | 'SIGTTOU'
              | 'SIGURG'
              | 'SIGUSR1'
              | 'SIGUSR2'
              | 'SIGVTALRM'
              | 'SIGWINCH'
              | 'SIGXCPU'
              | 'SIGXFSZ';
            [k: string]: unknown;
          };
          /**
           * Name of the container specified as a DNS_LABEL. Each container in a pod must have a unique name (DNS_LABEL). Cannot be updated.
           */
          name: string;
          /**
           * ResourceRequirements describes the compute resource requirements.
           */
          resources?: {
            /**
             * Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
             */
            limits?: {
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
             * Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
             */
            requests?: {
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
             * Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container.
             *
             * This is an alpha field and requires enabling the DynamicResourceAllocation feature gate.
             *
             * This field is immutable. It can only be set for containers.
             */
            claims?: {
              /**
               * Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
               */
              name: string;
              /**
               * Request is the name chosen for a request in the referenced claim. If empty, everything from the claim is made available, otherwise only the result of this request.
               */
              request?: string;
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          [k: string]: unknown;
        }[];
        /**
         * NodeSelector is a selector which must be true for the pod to fit on a node. Selector which must match a node's labels for the pod to be scheduled on that node. More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
         */
        nodeSelector?: {
          [k: string]: string;
        };
        /**
         * PodOS defines the OS parameters of a pod.
         */
        os?: {
          /**
           * Name is the name of the operating system. The currently supported values are linux and windows. Additional value may be defined in future and can be one of: https://github.com/opencontainers/runtime-spec/blob/master/config.md#platform-specific-configuration Clients should expect to handle additional values and treat unrecognized values in this field as os: null
           */
          name: string;
          [k: string]: unknown;
        };
        /**
         * If specified, the pod will be dispatched by specified scheduler. If not specified, the pod will be dispatched by default scheduler.
         */
        schedulerName?: string;
        /**
         * PodSecurityContext holds pod-level security attributes and common container settings. Some fields are also present in container.securityContext.  Field values of container.securityContext take precedence over field values of PodSecurityContext.
         */
        securityContext?: {
          /**
           * Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in SecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
           */
          runAsNonRoot?: boolean;
          /**
           * SeccompProfile defines a pod/container's seccomp profile settings. Only one profile source may be set.
           */
          seccompProfile?: {
            /**
             * localhostProfile indicates a profile defined in a file on the node should be used. The profile must be preconfigured on the node to work. Must be a descending path, relative to the kubelet's configured seccomp profile location. Must be set if type is "Localhost". Must NOT be set for any other type.
             */
            localhostProfile?: string;
            /**
             * type indicates which kind of seccomp profile will be applied. Valid options are:
             *
             * Localhost - a profile defined in a file on the node should be used. RuntimeDefault - the container runtime default profile should be used. Unconfined - no profile should be applied.
             *
             * Possible enum values:
             *  - `"Localhost"` indicates a profile defined in a file on the node should be used. The file's location relative to <kubelet-root-dir>/seccomp.
             *  - `"RuntimeDefault"` represents the default container runtime seccomp profile.
             *  - `"Unconfined"` indicates no seccomp profile is applied (A.K.A. unconfined).
             */
            type: 'Localhost' | 'RuntimeDefault' | 'Unconfined';
            [k: string]: unknown;
          };
          /**
           * A list of groups applied to the first process run in each container, in addition to the container's primary GID and fsGroup (if specified).  If the SupplementalGroupsPolicy feature is enabled, the supplementalGroupsPolicy field determines whether these are in addition to or instead of any group memberships defined in the container image. If unspecified, no additional groups are added, though group memberships defined in the container image may still be used, depending on the supplementalGroupsPolicy field. Note that this field cannot be set when spec.os.name is windows.
           */
          supplementalGroups?: number[];
          /**
           * Sysctls hold a list of namespaced sysctls used for the pod. Pods with unsupported sysctls (by the container runtime) might fail to launch. Note that this field cannot be set when spec.os.name is windows.
           */
          sysctls?: {
            /**
             * Name of a property to set
             */
            name: string;
            /**
             * Value of a property to set
             */
            value: string;
            [k: string]: unknown;
          }[];
          /**
           * WindowsSecurityContextOptions contain Windows-specific options and credentials.
           */
          windowsOptions?: {
            /**
             * The UserName in Windows to run the entrypoint of the container process. Defaults to the user specified in image metadata if unspecified. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
             */
            runAsUserName?: string;
            /**
             * GMSACredentialSpec is where the GMSA admission webhook (https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the GMSA credential spec named by the GMSACredentialSpecName field.
             */
            gmsaCredentialSpec?: string;
            /**
             * GMSACredentialSpecName is the name of the GMSA credential spec to use.
             */
            gmsaCredentialSpecName?: string;
            /**
             * HostProcess determines if a container should be run as a 'Host Process' container. All of a Pod's containers must have the same effective HostProcess value (it is not allowed to have a mix of HostProcess containers and non-HostProcess containers). In addition, if HostProcess is true then HostNetwork must also be set to true.
             */
            hostProcess?: boolean;
            [k: string]: unknown;
          };
          /**
           * AppArmorProfile defines a pod or container's AppArmor settings.
           */
          appArmorProfile?: {
            /**
             * type indicates which kind of AppArmor profile will be applied. Valid options are:
             *   Localhost - a profile pre-loaded on the node.
             *   RuntimeDefault - the container runtime's default profile.
             *   Unconfined - no AppArmor enforcement.
             *
             * Possible enum values:
             *  - `"Localhost"` indicates that a profile pre-loaded on the node should be used.
             *  - `"RuntimeDefault"` indicates that the container runtime's default AppArmor profile should be used.
             *  - `"Unconfined"` indicates that no AppArmor profile should be enforced.
             */
            type: 'Localhost' | 'RuntimeDefault' | 'Unconfined';
            /**
             * localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is "Localhost".
             */
            localhostProfile?: string;
            [k: string]: unknown;
          };
          /**
           * fsGroupChangePolicy defines behavior of changing ownership and permission of the volume before being exposed inside Pod. This field will only apply to volume types which support fsGroup based ownership(and permissions). It will have no effect on ephemeral volume types such as: secret, configmaps and emptydir. Valid values are "OnRootMismatch" and "Always". If not specified, "Always" is used. Note that this field cannot be set when spec.os.name is windows.
           *
           * Possible enum values:
           *  - `"Always"` indicates that volume's ownership and permissions should always be changed whenever volume is mounted inside a Pod. This the default behavior.
           *  - `"OnRootMismatch"` indicates that volume's ownership and permissions will be changed only when permission and ownership of root directory does not match with expected permissions on the volume. This can help shorten the time it takes to change ownership and permissions of a volume.
           */
          fsGroupChangePolicy?: 'Always' | 'OnRootMismatch';
          /**
           * The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in SecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container. Note that this field cannot be set when spec.os.name is windows.
           */
          runAsGroup?: number;
          /**
           * The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in SecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container. Note that this field cannot be set when spec.os.name is windows.
           */
          runAsUser?: number;
          /**
           * seLinuxChangePolicy defines how the container's SELinux label is applied to all volumes used by the Pod. It has no effect on nodes that do not support SELinux or to volumes does not support SELinux. Valid values are "MountOption" and "Recursive".
           *
           * "Recursive" means relabeling of all files on all Pod volumes by the container runtime. This may be slow for large volumes, but allows mixing privileged and unprivileged Pods sharing the same volume on the same node.
           *
           * "MountOption" mounts all eligible Pod volumes with `-o context` mount option. This requires all Pods that share the same volume to use the same SELinux label. It is not possible to share the same volume among privileged and unprivileged Pods. Eligible volumes are in-tree FibreChannel and iSCSI volumes, and all CSI volumes whose CSI driver announces SELinux support by setting spec.seLinuxMount: true in their CSIDriver instance. Other volumes are always re-labelled recursively. "MountOption" value is allowed only when SELinuxMount feature gate is enabled.
           *
           * If not specified and SELinuxMount feature gate is enabled, "MountOption" is used. If not specified and SELinuxMount feature gate is disabled, "MountOption" is used for ReadWriteOncePod volumes and "Recursive" for all other volumes.
           *
           * This field affects only Pods that have SELinux label set, either in PodSecurityContext or in SecurityContext of all containers.
           *
           * All Pods that use the same volume should use the same seLinuxChangePolicy, otherwise some pods can get stuck in ContainerCreating state. Note that this field cannot be set when spec.os.name is windows.
           */
          seLinuxChangePolicy?: string;
          /**
           * SELinuxOptions are the labels to be applied to the container
           */
          seLinuxOptions?: {
            /**
             * User is a SELinux user label that applies to the container.
             */
            user?: string;
            /**
             * Level is SELinux level label that applies to the container.
             */
            level?: string;
            /**
             * Role is a SELinux role label that applies to the container.
             */
            role?: string;
            /**
             * Type is a SELinux type label that applies to the container.
             */
            type?: string;
            [k: string]: unknown;
          };
          /**
           * Defines how supplemental groups of the first container processes are calculated. Valid values are "Merge" and "Strict". If not specified, "Merge" is used. (Alpha) Using the field requires the SupplementalGroupsPolicy feature gate to be enabled and the container runtime must implement support for this feature. Note that this field cannot be set when spec.os.name is windows.
           *
           * Possible enum values:
           *  - `"Merge"` means that the container's provided SupplementalGroups and FsGroup (specified in SecurityContext) will be merged with the primary user's groups as defined in the container image (in /etc/group).
           *  - `"Strict"` means that the container's provided SupplementalGroups and FsGroup (specified in SecurityContext) will be used instead of any groups defined in the container image.
           */
          supplementalGroupsPolicy?: 'Merge' | 'Strict';
          /**
           * A special supplemental group that applies to all containers in a pod. Some volume types allow the Kubelet to change the ownership of that volume to be owned by the pod:
           *
           * 1. The owning GID will be the FSGroup 2. The setgid bit is set (new files created in the volume will be owned by FSGroup) 3. The permission bits are OR'd with rw-rw----
           *
           * If unset, the Kubelet will not modify the ownership and permissions of any volume. Note that this field cannot be set when spec.os.name is windows.
           */
          fsGroup?: number;
          [k: string]: unknown;
        };
        /**
         * If specified, the pod's tolerations.
         */
        tolerations?: {
          /**
           * Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute.
           *
           * Possible enum values:
           *  - `"NoExecute"` Evict any already-running pods that do not tolerate the taint. Currently enforced by NodeController.
           *  - `"NoSchedule"` Do not allow new pods to schedule onto the node unless they tolerate the taint, but allow all pods submitted to Kubelet without going through the scheduler to start, and allow all already-running pods to continue running. Enforced by the scheduler.
           *  - `"PreferNoSchedule"` Like TaintEffectNoSchedule, but the scheduler tries not to schedule new pods onto the node, rather than prohibiting new pods from scheduling onto the node entirely. Enforced by the scheduler.
           */
          effect?: 'NoExecute' | 'NoSchedule' | 'PreferNoSchedule';
          /**
           * Key is the taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys.
           */
          key?: string;
          /**
           * Operator represents a key's relationship to the value. Valid operators are Exists and Equal. Defaults to Equal. Exists is equivalent to wildcard for value, so that a pod can tolerate all taints of a particular category.
           *
           * Possible enum values:
           *  - `"Equal"`
           *  - `"Exists"`
           */
          operator?: 'Equal' | 'Exists';
          /**
           * TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system.
           */
          tolerationSeconds?: number;
          /**
           * Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string.
           */
          value?: string;
          [k: string]: unknown;
        }[];
        /**
         * AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.
         */
        automountServiceAccountToken?: boolean;
        /**
         * Use the host's ipc namespace. Optional: Default to false.
         */
        hostIPC?: boolean;
        /**
         * ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec. If specified, these secrets will be passed to individual puller implementations for them to use. More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
         */
        imagePullSecrets?: {
          /**
           * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
           */
          name?: string;
          [k: string]: unknown;
        }[];
        /**
         * PreemptionPolicy is the Policy for preempting pods with lower priority. One of Never, PreemptLowerPriority. Defaults to PreemptLowerPriority if unset.
         *
         * Possible enum values:
         *  - `"Never"` means that pod never preempts other pods with lower priority.
         *  - `"PreemptLowerPriority"` means that pod can preempt other pods with lower priority.
         */
        preemptionPolicy?: 'Never' | 'PreemptLowerPriority';
        /**
         * List of initialization containers belonging to the pod. Init containers are executed in order prior to containers being started. If any init container fails, the pod is considered to have failed and is handled according to its restartPolicy. The name for an init container or normal container must be unique among all containers. Init containers may not have Lifecycle actions, Readiness probes, Liveness probes, or Startup probes. The resourceRequirements of an init container are taken into account during scheduling by finding the highest request/limit for each resource type, and then using the max of that value or the sum of the normal containers. Limits are applied to init containers in a similar fashion. Init containers cannot currently be added or removed. Cannot be updated. More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
         */
        initContainers?: {
          /**
           * Entrypoint array. Not executed within a shell. The container image's ENTRYPOINT is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
           */
          command?: string[];
          /**
           * List of sources to populate environment variables in the container. The keys defined within a source must be a C_IDENTIFIER. All invalid keys will be reported as an event when the container is starting. When a key exists in multiple sources, the value associated with the last source will take precedence. Values defined by an Env with a duplicate key will take precedence. Cannot be updated.
           */
          envFrom?: {
            /**
             * ConfigMapEnvSource selects a ConfigMap to populate the environment variables with.
             *
             * The contents of the target ConfigMap's Data field will represent the key-value pairs as environment variables.
             */
            configMapRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              /**
               * Specify whether the ConfigMap must be defined
               */
              optional?: boolean;
              [k: string]: unknown;
            };
            /**
             * Optional text to prepend to the name of each environment variable. Must be a C_IDENTIFIER.
             */
            prefix?: string;
            /**
             * SecretEnvSource selects a Secret to populate the environment variables with.
             *
             * The contents of the target Secret's Data field will represent the key-value pairs as environment variables.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              /**
               * Specify whether the Secret must be defined
               */
              optional?: boolean;
              [k: string]: unknown;
            };
            [k: string]: unknown;
          }[];
          /**
           * Container image name. More info: https://kubernetes.io/docs/concepts/containers/images This field is optional to allow higher level config management to default or override container images in workload controllers like Deployments and StatefulSets.
           */
          image?: string;
          /**
           * Lifecycle describes actions that the management system should take in response to container lifecycle events. For the PostStart and PreStop lifecycle handlers, management of the container blocks until the action is complete, unless the container process fails, in which case the handler is aborted.
           */
          lifecycle?: {
            /**
             * LifecycleHandler defines a specific action that should be taken in a lifecycle hook. One and only one of the fields, except TCPSocket must be specified.
             */
            postStart?: {
              /**
               * ExecAction describes a "run in container" action.
               */
              exec?: {
                /**
                 * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
                 */
                command?: string[];
                [k: string]: unknown;
              };
              /**
               * HTTPGetAction describes an action based on HTTP Get requests.
               */
              httpGet?: {
                /**
                 * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
                 */
                host?: string;
                /**
                 * Custom headers to set in the request. HTTP allows repeated headers.
                 */
                httpHeaders?: {
                  /**
                   * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                   */
                  name: string;
                  /**
                   * The header field value
                   */
                  value: string;
                  [k: string]: unknown;
                }[];
                /**
                 * Path to access on the HTTP server.
                 */
                path?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                /**
                 * Scheme to use for connecting to the host. Defaults to HTTP.
                 *
                 * Possible enum values:
                 *  - `"HTTP"` means that the scheme used will be http://
                 *  - `"HTTPS"` means that the scheme used will be https://
                 */
                scheme?: 'HTTP' | 'HTTPS';
                [k: string]: unknown;
              };
              /**
               * SleepAction describes a "sleep" action.
               */
              sleep?: {
                /**
                 * Seconds is the number of seconds to sleep.
                 */
                seconds: number;
                [k: string]: unknown;
              };
              /**
               * TCPSocketAction describes an action based on opening a socket
               */
              tcpSocket?: {
                /**
                 * Optional: Host name to connect to, defaults to the pod IP.
                 */
                host?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            /**
             * LifecycleHandler defines a specific action that should be taken in a lifecycle hook. One and only one of the fields, except TCPSocket must be specified.
             */
            preStop?: {
              /**
               * ExecAction describes a "run in container" action.
               */
              exec?: {
                /**
                 * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
                 */
                command?: string[];
                [k: string]: unknown;
              };
              /**
               * HTTPGetAction describes an action based on HTTP Get requests.
               */
              httpGet?: {
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                /**
                 * Scheme to use for connecting to the host. Defaults to HTTP.
                 *
                 * Possible enum values:
                 *  - `"HTTP"` means that the scheme used will be http://
                 *  - `"HTTPS"` means that the scheme used will be https://
                 */
                scheme?: 'HTTP' | 'HTTPS';
                /**
                 * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
                 */
                host?: string;
                /**
                 * Custom headers to set in the request. HTTP allows repeated headers.
                 */
                httpHeaders?: {
                  /**
                   * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                   */
                  name: string;
                  /**
                   * The header field value
                   */
                  value: string;
                  [k: string]: unknown;
                }[];
                /**
                 * Path to access on the HTTP server.
                 */
                path?: string;
                [k: string]: unknown;
              };
              /**
               * SleepAction describes a "sleep" action.
               */
              sleep?: {
                /**
                 * Seconds is the number of seconds to sleep.
                 */
                seconds: number;
                [k: string]: unknown;
              };
              /**
               * TCPSocketAction describes an action based on opening a socket
               */
              tcpSocket?: {
                /**
                 * Optional: Host name to connect to, defaults to the pod IP.
                 */
                host?: string;
                /**
                 * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
                 */
                port: number | string;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            /**
             * StopSignal defines which signal will be sent to a container when it is being stopped. If not specified, the default is defined by the container runtime in use. StopSignal can only be set for Pods with a non-empty .spec.os.name
             *
             * Possible enum values:
             *  - `"SIGABRT"`
             *  - `"SIGALRM"`
             *  - `"SIGBUS"`
             *  - `"SIGCHLD"`
             *  - `"SIGCLD"`
             *  - `"SIGCONT"`
             *  - `"SIGFPE"`
             *  - `"SIGHUP"`
             *  - `"SIGILL"`
             *  - `"SIGINT"`
             *  - `"SIGIO"`
             *  - `"SIGIOT"`
             *  - `"SIGKILL"`
             *  - `"SIGPIPE"`
             *  - `"SIGPOLL"`
             *  - `"SIGPROF"`
             *  - `"SIGPWR"`
             *  - `"SIGQUIT"`
             *  - `"SIGRTMAX"`
             *  - `"SIGRTMAX-1"`
             *  - `"SIGRTMAX-10"`
             *  - `"SIGRTMAX-11"`
             *  - `"SIGRTMAX-12"`
             *  - `"SIGRTMAX-13"`
             *  - `"SIGRTMAX-14"`
             *  - `"SIGRTMAX-2"`
             *  - `"SIGRTMAX-3"`
             *  - `"SIGRTMAX-4"`
             *  - `"SIGRTMAX-5"`
             *  - `"SIGRTMAX-6"`
             *  - `"SIGRTMAX-7"`
             *  - `"SIGRTMAX-8"`
             *  - `"SIGRTMAX-9"`
             *  - `"SIGRTMIN"`
             *  - `"SIGRTMIN+1"`
             *  - `"SIGRTMIN+10"`
             *  - `"SIGRTMIN+11"`
             *  - `"SIGRTMIN+12"`
             *  - `"SIGRTMIN+13"`
             *  - `"SIGRTMIN+14"`
             *  - `"SIGRTMIN+15"`
             *  - `"SIGRTMIN+2"`
             *  - `"SIGRTMIN+3"`
             *  - `"SIGRTMIN+4"`
             *  - `"SIGRTMIN+5"`
             *  - `"SIGRTMIN+6"`
             *  - `"SIGRTMIN+7"`
             *  - `"SIGRTMIN+8"`
             *  - `"SIGRTMIN+9"`
             *  - `"SIGSEGV"`
             *  - `"SIGSTKFLT"`
             *  - `"SIGSTOP"`
             *  - `"SIGSYS"`
             *  - `"SIGTERM"`
             *  - `"SIGTRAP"`
             *  - `"SIGTSTP"`
             *  - `"SIGTTIN"`
             *  - `"SIGTTOU"`
             *  - `"SIGURG"`
             *  - `"SIGUSR1"`
             *  - `"SIGUSR2"`
             *  - `"SIGVTALRM"`
             *  - `"SIGWINCH"`
             *  - `"SIGXCPU"`
             *  - `"SIGXFSZ"`
             */
            stopSignal?:
              | 'SIGABRT'
              | 'SIGALRM'
              | 'SIGBUS'
              | 'SIGCHLD'
              | 'SIGCLD'
              | 'SIGCONT'
              | 'SIGFPE'
              | 'SIGHUP'
              | 'SIGILL'
              | 'SIGINT'
              | 'SIGIO'
              | 'SIGIOT'
              | 'SIGKILL'
              | 'SIGPIPE'
              | 'SIGPOLL'
              | 'SIGPROF'
              | 'SIGPWR'
              | 'SIGQUIT'
              | 'SIGRTMAX'
              | 'SIGRTMAX-1'
              | 'SIGRTMAX-10'
              | 'SIGRTMAX-11'
              | 'SIGRTMAX-12'
              | 'SIGRTMAX-13'
              | 'SIGRTMAX-14'
              | 'SIGRTMAX-2'
              | 'SIGRTMAX-3'
              | 'SIGRTMAX-4'
              | 'SIGRTMAX-5'
              | 'SIGRTMAX-6'
              | 'SIGRTMAX-7'
              | 'SIGRTMAX-8'
              | 'SIGRTMAX-9'
              | 'SIGRTMIN'
              | 'SIGRTMIN+1'
              | 'SIGRTMIN+10'
              | 'SIGRTMIN+11'
              | 'SIGRTMIN+12'
              | 'SIGRTMIN+13'
              | 'SIGRTMIN+14'
              | 'SIGRTMIN+15'
              | 'SIGRTMIN+2'
              | 'SIGRTMIN+3'
              | 'SIGRTMIN+4'
              | 'SIGRTMIN+5'
              | 'SIGRTMIN+6'
              | 'SIGRTMIN+7'
              | 'SIGRTMIN+8'
              | 'SIGRTMIN+9'
              | 'SIGSEGV'
              | 'SIGSTKFLT'
              | 'SIGSTOP'
              | 'SIGSYS'
              | 'SIGTERM'
              | 'SIGTRAP'
              | 'SIGTSTP'
              | 'SIGTTIN'
              | 'SIGTTOU'
              | 'SIGURG'
              | 'SIGUSR1'
              | 'SIGUSR2'
              | 'SIGVTALRM'
              | 'SIGWINCH'
              | 'SIGXCPU'
              | 'SIGXFSZ';
            [k: string]: unknown;
          };
          /**
           * Whether the container runtime should close the stdin channel after it has been opened by a single attach. When stdin is true the stdin stream will remain open across multiple attach sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the first client attaches to stdin, and then remains open and accepts data until the client disconnects, at which time stdin is closed and remains closed until the container is restarted. If this flag is false, a container processes that reads from stdin will never receive an EOF. Default is false
           */
          stdinOnce?: boolean;
          /**
           * Optional: Path at which the file to which the container's termination message will be written is mounted into the container's filesystem. Message written is intended to be brief final status, such as an assertion failure message. Will be truncated by the node if greater than 4096 bytes. The total message length across all containers will be limited to 12kb. Defaults to /dev/termination-log. Cannot be updated.
           */
          terminationMessagePath?: string;
          /**
           * volumeDevices is the list of block devices to be used by the container.
           */
          volumeDevices?: {
            /**
             * devicePath is the path inside of the container that the device will be mapped to.
             */
            devicePath: string;
            /**
             * name must match the name of a persistentVolumeClaim in the pod
             */
            name: string;
            [k: string]: unknown;
          }[];
          /**
           * Arguments to the entrypoint. The container image's CMD is used if this is not provided. Variable references $(VAR_NAME) are expanded using the container's environment. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Cannot be updated. More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
           */
          args?: string[];
          /**
           * Image pull policy. One of Always, Never, IfNotPresent. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise. Cannot be updated. More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
           *
           * Possible enum values:
           *  - `"Always"` means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.
           *  - `"IfNotPresent"` means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
           *  - `"Never"` means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn't present
           */
          imagePullPolicy?: 'Always' | 'IfNotPresent' | 'Never';
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          livenessProbe?: {
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                /**
                 * The header field value
                 */
                value: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              [k: string]: unknown;
            };
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            [k: string]: unknown;
          };
          /**
           * ResourceRequirements describes the compute resource requirements.
           */
          resources?: {
            /**
             * Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container.
             *
             * This is an alpha field and requires enabling the DynamicResourceAllocation feature gate.
             *
             * This field is immutable. It can only be set for containers.
             */
            claims?: {
              /**
               * Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
               */
              name: string;
              /**
               * Request is the name chosen for a request in the referenced claim. If empty, everything from the claim is made available, otherwise only the result of this request.
               */
              request?: string;
              [k: string]: unknown;
            }[];
            /**
             * Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
             */
            limits?: {
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
             * Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
             */
            requests?: {
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
            [k: string]: unknown;
          };
          /**
           * SecurityContext holds security configuration that will be applied to a container. Some fields are present in both SecurityContext and PodSecurityContext.  When both are set, the values in SecurityContext take precedence.
           */
          securityContext?: {
            /**
             * Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
             */
            runAsNonRoot?: boolean;
            /**
             * The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
             */
            runAsUser?: number;
            /**
             * SELinuxOptions are the labels to be applied to the container
             */
            seLinuxOptions?: {
              /**
               * Role is a SELinux role label that applies to the container.
               */
              role?: string;
              /**
               * Type is a SELinux type label that applies to the container.
               */
              type?: string;
              /**
               * User is a SELinux user label that applies to the container.
               */
              user?: string;
              /**
               * Level is SELinux level label that applies to the container.
               */
              level?: string;
              [k: string]: unknown;
            };
            /**
             * procMount denotes the type of proc mount to use for the containers. The default value is Default which uses the container runtime defaults for readonly paths and masked paths. This requires the ProcMountType feature flag to be enabled. Note that this field cannot be set when spec.os.name is windows.
             *
             * Possible enum values:
             *  - `"Default"` uses the container runtime defaults for readonly and masked paths for /proc. Most container runtimes mask certain paths in /proc to avoid accidental security exposure of special devices or information.
             *  - `"Unmasked"` bypasses the default masking behavior of the container runtime and ensures the newly created /proc the container stays in tact with no modifications.
             */
            procMount?: 'Default' | 'Unmasked';
            /**
             * The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
             */
            runAsGroup?: number;
            /**
             * SeccompProfile defines a pod/container's seccomp profile settings. Only one profile source may be set.
             */
            seccompProfile?: {
              /**
               * localhostProfile indicates a profile defined in a file on the node should be used. The profile must be preconfigured on the node to work. Must be a descending path, relative to the kubelet's configured seccomp profile location. Must be set if type is "Localhost". Must NOT be set for any other type.
               */
              localhostProfile?: string;
              /**
               * type indicates which kind of seccomp profile will be applied. Valid options are:
               *
               * Localhost - a profile defined in a file on the node should be used. RuntimeDefault - the container runtime default profile should be used. Unconfined - no profile should be applied.
               *
               * Possible enum values:
               *  - `"Localhost"` indicates a profile defined in a file on the node should be used. The file's location relative to <kubelet-root-dir>/seccomp.
               *  - `"RuntimeDefault"` represents the default container runtime seccomp profile.
               *  - `"Unconfined"` indicates no seccomp profile is applied (A.K.A. unconfined).
               */
              type: 'Localhost' | 'RuntimeDefault' | 'Unconfined';
              [k: string]: unknown;
            };
            /**
             * WindowsSecurityContextOptions contain Windows-specific options and credentials.
             */
            windowsOptions?: {
              /**
               * GMSACredentialSpec is where the GMSA admission webhook (https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the GMSA credential spec named by the GMSACredentialSpecName field.
               */
              gmsaCredentialSpec?: string;
              /**
               * GMSACredentialSpecName is the name of the GMSA credential spec to use.
               */
              gmsaCredentialSpecName?: string;
              /**
               * HostProcess determines if a container should be run as a 'Host Process' container. All of a Pod's containers must have the same effective HostProcess value (it is not allowed to have a mix of HostProcess containers and non-HostProcess containers). In addition, if HostProcess is true then HostNetwork must also be set to true.
               */
              hostProcess?: boolean;
              /**
               * The UserName in Windows to run the entrypoint of the container process. Defaults to the user specified in image metadata if unspecified. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
               */
              runAsUserName?: string;
              [k: string]: unknown;
            };
            /**
             * AllowPrivilegeEscalation controls whether a process can gain more privileges than its parent process. This bool directly controls if the no_new_privs flag will be set on the container process. AllowPrivilegeEscalation is true always when the container is: 1) run as Privileged 2) has CAP_SYS_ADMIN Note that this field cannot be set when spec.os.name is windows.
             */
            allowPrivilegeEscalation?: boolean;
            /**
             * AppArmorProfile defines a pod or container's AppArmor settings.
             */
            appArmorProfile?: {
              /**
               * localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is "Localhost".
               */
              localhostProfile?: string;
              /**
               * type indicates which kind of AppArmor profile will be applied. Valid options are:
               *   Localhost - a profile pre-loaded on the node.
               *   RuntimeDefault - the container runtime's default profile.
               *   Unconfined - no AppArmor enforcement.
               *
               * Possible enum values:
               *  - `"Localhost"` indicates that a profile pre-loaded on the node should be used.
               *  - `"RuntimeDefault"` indicates that the container runtime's default AppArmor profile should be used.
               *  - `"Unconfined"` indicates that no AppArmor profile should be enforced.
               */
              type: 'Localhost' | 'RuntimeDefault' | 'Unconfined';
              [k: string]: unknown;
            };
            /**
             * Adds and removes POSIX capabilities from running containers.
             */
            capabilities?: {
              /**
               * Added capabilities
               */
              add?: string[];
              /**
               * Removed capabilities
               */
              drop?: string[];
              [k: string]: unknown;
            };
            /**
             * Run container in privileged mode. Processes in privileged containers are essentially equivalent to root on the host. Defaults to false. Note that this field cannot be set when spec.os.name is windows.
             */
            privileged?: boolean;
            /**
             * Whether this container has a read-only root filesystem. Default is false. Note that this field cannot be set when spec.os.name is windows.
             */
            readOnlyRootFilesystem?: boolean;
            [k: string]: unknown;
          };
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          startupProbe?: {
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                /**
                 * The header field value
                 */
                value: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              [k: string]: unknown;
            };
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            [k: string]: unknown;
          };
          /**
           * Whether this container should allocate a buffer for stdin in the container runtime. If this is not set, reads from stdin in the container will always result in EOF. Default is false.
           */
          stdin?: boolean;
          /**
           * Whether this container should allocate a TTY for itself, also requires 'stdin' to be true. Default is false.
           */
          tty?: boolean;
          /**
           * List of environment variables to set in the container. Cannot be updated.
           */
          env?: {
            /**
             * Name of the environment variable. Must be a C_IDENTIFIER.
             */
            name: string;
            /**
             * Variable references $(VAR_NAME) are expanded using the previously defined environment variables in the container and any service environment variables. If a variable cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless of whether the variable exists or not. Defaults to "".
             */
            value?: string;
            /**
             * EnvVarSource represents a source for the value of an EnvVar.
             */
            valueFrom?: {
              /**
               * Selects a key from a ConfigMap.
               */
              configMapKeyRef?: {
                /**
                 * The key to select.
                 */
                key: string;
                /**
                 * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                 */
                name?: string;
                /**
                 * Specify whether the ConfigMap or its key must be defined
                 */
                optional?: boolean;
                [k: string]: unknown;
              };
              /**
               * ObjectFieldSelector selects an APIVersioned field of an object.
               */
              fieldRef?: {
                /**
                 * Version of the schema the FieldPath is written in terms of, defaults to "v1".
                 */
                apiVersion?: string;
                /**
                 * Path of the field to select in the specified API version.
                 */
                fieldPath: string;
                [k: string]: unknown;
              };
              /**
               * ResourceFieldSelector represents container resources (cpu, memory) and their output format
               */
              resourceFieldRef?: {
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
                divisor?: string | number;
                /**
                 * Required: resource to select
                 */
                resource: string;
                /**
                 * Container name: required for volumes, optional for env vars
                 */
                containerName?: string;
                [k: string]: unknown;
              };
              /**
               * SecretKeySelector selects a key of a Secret.
               */
              secretKeyRef?: {
                /**
                 * The key of the secret to select from.  Must be a valid secret key.
                 */
                key: string;
                /**
                 * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                 */
                name?: string;
                /**
                 * Specify whether the Secret or its key must be defined
                 */
                optional?: boolean;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            [k: string]: unknown;
          }[];
          /**
           * Name of the container specified as a DNS_LABEL. Each container in a pod must have a unique name (DNS_LABEL). Cannot be updated.
           */
          name: string;
          /**
           * List of ports to expose from the container. Not specifying a port here DOES NOT prevent that port from being exposed. Any port which is listening on the default "0.0.0.0" address inside a container will be accessible from the network. Modifying this array with strategic merge patch may corrupt the data. For more information See https://github.com/kubernetes/kubernetes/issues/108255. Cannot be updated.
           */
          ports?: {
            /**
             * Number of port to expose on the host. If specified, this must be a valid port number, 0 < x < 65536. If HostNetwork is specified, this must match ContainerPort. Most containers do not need this.
             */
            hostPort?: number;
            /**
             * If specified, this must be an IANA_SVC_NAME and unique within the pod. Each named port in a pod must have a unique name. Name for the port that can be referred to by services.
             */
            name?: string;
            /**
             * Protocol for port. Must be UDP, TCP, or SCTP. Defaults to "TCP".
             *
             * Possible enum values:
             *  - `"SCTP"` is the SCTP protocol.
             *  - `"TCP"` is the TCP protocol.
             *  - `"UDP"` is the UDP protocol.
             */
            protocol?: 'SCTP' | 'TCP' | 'UDP';
            /**
             * Number of port to expose on the pod's IP address. This must be a valid port number, 0 < x < 65536.
             */
            containerPort: number;
            /**
             * What host IP to bind the external port to.
             */
            hostIP?: string;
            [k: string]: unknown;
          }[];
          /**
           * Resources resize policy for the container.
           */
          resizePolicy?: {
            /**
             * Restart policy to apply when specified resource is resized. If not specified, it defaults to NotRequired.
             */
            restartPolicy: string;
            /**
             * Name of the resource to which this resource resize policy applies. Supported values: cpu, memory.
             */
            resourceName: string;
            [k: string]: unknown;
          }[];
          /**
           * RestartPolicy defines the restart behavior of individual containers in a pod. This field may only be set for init containers, and the only allowed value is "Always". For non-init containers or when this field is not specified, the restart behavior is defined by the Pod's restart policy and the container type. Setting the RestartPolicy as "Always" for the init container will have the following effect: this init container will be continually restarted on exit until all regular containers have terminated. Once all regular containers have completed, all init containers with restartPolicy "Always" will be shut down. This lifecycle differs from normal init containers and is often referred to as a "sidecar" container. Although this init container still starts in the init container sequence, it does not wait for the container to complete before proceeding to the next init container. Instead, the next init container starts immediately after this init container is started, or after any startupProbe has successfully completed.
           */
          restartPolicy?: string;
          /**
           * Indicate how the termination message should be populated. File will use the contents of terminationMessagePath to populate the container status message on both success and failure. FallbackToLogsOnError will use the last chunk of container log output if the termination message file is empty and the container exited with an error. The log output is limited to 2048 bytes or 80 lines, whichever is smaller. Defaults to File. Cannot be updated.
           *
           * Possible enum values:
           *  - `"FallbackToLogsOnError"` will read the most recent contents of the container logs for the container status message when the container exits with an error and the terminationMessagePath has no contents.
           *  - `"File"` is the default behavior and will set the container status message to the contents of the container's terminationMessagePath when the container exits.
           */
          terminationMessagePolicy?: 'FallbackToLogsOnError' | 'File';
          /**
           * Pod volumes to mount into the container's filesystem. Cannot be updated.
           */
          volumeMounts?: {
            /**
             * mountPropagation determines how mounts are propagated from the host to container and the other way around. When not set, MountPropagationNone is used. This field is beta in 1.10. When RecursiveReadOnly is set to IfPossible or to Enabled, MountPropagation must be None or unspecified (which defaults to None).
             *
             * Possible enum values:
             *  - `"Bidirectional"` means that the volume in a container will receive new mounts from the host or other containers, and its own mounts will be propagated from the container to the host or other containers. Note that this mode is recursively applied to all mounts in the volume ("rshared" in Linux terminology).
             *  - `"HostToContainer"` means that the volume in a container will receive new mounts from the host or other containers, but filesystems mounted inside the container won't be propagated to the host or other containers. Note that this mode is recursively applied to all mounts in the volume ("rslave" in Linux terminology).
             *  - `"None"` means that the volume in a container will not receive new mounts from the host or other containers, and filesystems mounted inside the container won't be propagated to the host or other containers. Note that this mode corresponds to "private" in Linux terminology.
             */
            mountPropagation?: 'Bidirectional' | 'HostToContainer' | 'None';
            /**
             * This must match the Name of a Volume.
             */
            name: string;
            /**
             * Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false.
             */
            readOnly?: boolean;
            /**
             * RecursiveReadOnly specifies whether read-only mounts should be handled recursively.
             *
             * If ReadOnly is false, this field has no meaning and must be unspecified.
             *
             * If ReadOnly is true, and this field is set to Disabled, the mount is not made recursively read-only.  If this field is set to IfPossible, the mount is made recursively read-only, if it is supported by the container runtime.  If this field is set to Enabled, the mount is made recursively read-only if it is supported by the container runtime, otherwise the pod will not be started and an error will be generated to indicate the reason.
             *
             * If this field is set to IfPossible or Enabled, MountPropagation must be set to None (or be unspecified, which defaults to None).
             *
             * If this field is not specified, it is treated as an equivalent of Disabled.
             */
            recursiveReadOnly?: string;
            /**
             * Path within the volume from which the container's volume should be mounted. Defaults to "" (volume's root).
             */
            subPath?: string;
            /**
             * Expanded path within the volume from which the container's volume should be mounted. Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment. Defaults to "" (volume's root). SubPathExpr and SubPath are mutually exclusive.
             */
            subPathExpr?: string;
            /**
             * Path within the container at which the volume should be mounted.  Must not contain ':'.
             */
            mountPath: string;
            [k: string]: unknown;
          }[];
          /**
           * Container's working directory. If not specified, the container runtime's default will be used, which might be configured in the container image. Cannot be updated.
           */
          workingDir?: string;
          /**
           * Probe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.
           */
          readinessProbe?: {
            /**
             * TCPSocketAction describes an action based on opening a socket
             */
            tcpSocket?: {
              /**
               * Optional: Host name to connect to, defaults to the pod IP.
               */
              host?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
             */
            terminationGracePeriodSeconds?: number;
            /**
             * Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            timeoutSeconds?: number;
            /**
             * ExecAction describes a "run in container" action.
             */
            exec?: {
              /**
               * Command is the command line to execute inside the container, the working directory for the command  is root ('/') in the container's filesystem. The command is simply exec'd, it is not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use a shell, you need to explicitly call out to that shell. Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
               */
              command?: string[];
              [k: string]: unknown;
            };
            /**
             * Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
             */
            failureThreshold?: number;
            /**
             * GRPCAction specifies an action involving a GRPC service.
             */
            grpc?: {
              /**
               * Port number of the gRPC service. Number must be in the range 1 to 65535.
               */
              port: number;
              /**
               * Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
               *
               * If this is not specified, the default behavior is defined by gRPC.
               */
              service?: string;
              [k: string]: unknown;
            };
            /**
             * HTTPGetAction describes an action based on HTTP Get requests.
             */
            httpGet?: {
              /**
               * Scheme to use for connecting to the host. Defaults to HTTP.
               *
               * Possible enum values:
               *  - `"HTTP"` means that the scheme used will be http://
               *  - `"HTTPS"` means that the scheme used will be https://
               */
              scheme?: 'HTTP' | 'HTTPS';
              /**
               * Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
               */
              host?: string;
              /**
               * Custom headers to set in the request. HTTP allows repeated headers.
               */
              httpHeaders?: {
                /**
                 * The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
                 */
                name: string;
                /**
                 * The header field value
                 */
                value: string;
                [k: string]: unknown;
              }[];
              /**
               * Path to access on the HTTP server.
               */
              path?: string;
              /**
               * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
               */
              port: number | string;
              [k: string]: unknown;
            };
            /**
             * Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
             */
            initialDelaySeconds?: number;
            /**
             * How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
             */
            periodSeconds?: number;
            /**
             * Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
             */
            successThreshold?: number;
            [k: string]: unknown;
          };
          [k: string]: unknown;
        }[];
        /**
         * If specified, all readiness gates will be evaluated for pod readiness. A pod is ready when all its containers are ready AND all conditions specified in the readiness gates have status equal to "True" More info: https://git.k8s.io/enhancements/keps/sig-network/580-pod-readiness-gates
         */
        readinessGates?: {
          /**
           * ConditionType refers to a condition in the pod's condition list with matching type.
           */
          conditionType: string;
          [k: string]: unknown;
        }[];
        /**
         * List of volumes that can be mounted by containers belonging to the pod. More info: https://kubernetes.io/docs/concepts/storage/volumes
         */
        volumes?: {
          /**
           * AzureFile represents an Azure File Service mount on the host and bind mount to the pod.
           */
          azureFile?: {
            /**
             * readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
             */
            readOnly?: boolean;
            /**
             * secretName is the  name of secret that contains Azure Storage Account Name and Key
             */
            secretName: string;
            /**
             * shareName is the azure share Name
             */
            shareName: string;
            [k: string]: unknown;
          };
          /**
           * Adapts a ConfigMap into a volume.
           *
           * The contents of the target ConfigMap's Data field will be presented in a volume as files using the keys in the Data field as the file names, unless the items element is populated with specific mappings of keys to paths. ConfigMap volumes support ownership management and SELinux relabeling.
           */
          configMap?: {
            /**
             * items if unspecified, each key-value pair in the Data field of the referenced ConfigMap will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the ConfigMap, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the '..' path or start with '..'.
             */
            items?: {
              /**
               * path is the relative path of the file to map the key to. May not be an absolute path. May not contain the path element '..'. May not start with the string '..'.
               */
              path: string;
              /**
               * key is the key to project.
               */
              key: string;
              /**
               * mode is Optional: mode bits used to set permissions on this file. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. If not specified, the volume defaultMode will be used. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
               */
              mode?: number;
              [k: string]: unknown;
            }[];
            /**
             * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
             */
            name?: string;
            /**
             * optional specify whether the ConfigMap or its keys must be defined
             */
            optional?: boolean;
            /**
             * defaultMode is optional: mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
             */
            defaultMode?: number;
            [k: string]: unknown;
          };
          /**
           * Represents a source location of a volume to mount, managed by an external CSI driver
           */
          csi?: {
            /**
             * fsType to mount. Ex. "ext4", "xfs", "ntfs". If not provided, the empty value is passed to the associated CSI driver which will determine the default filesystem to apply.
             */
            fsType?: string;
            /**
             * LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
             */
            nodePublishSecretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              [k: string]: unknown;
            };
            /**
             * readOnly specifies a read-only configuration for the volume. Defaults to false (read/write).
             */
            readOnly?: boolean;
            /**
             * volumeAttributes stores driver-specific properties that are passed to the CSI driver. Consult your driver's documentation for supported values.
             */
            volumeAttributes?: {
              [k: string]: string;
            };
            /**
             * driver is the name of the CSI driver that handles this volume. Consult with your admin for the correct name as registered in the cluster.
             */
            driver: string;
            [k: string]: unknown;
          };
          /**
           * Represents an empty directory for a pod. Empty directory volumes support ownership management and SELinux relabeling.
           */
          emptyDir?: {
            /**
             * medium represents what type of storage medium should back this directory. The default is "" which means to use the node's default medium. Must be an empty string (default) or Memory. More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
             */
            medium?: string;
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
            sizeLimit?: string | number;
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
            [k: string]: unknown;
          };
          /**
           * Represents a Persistent Disk resource in Google Compute Engine.
           *
           * A GCE PD must exist before mounting to a container. The disk must also be in the same GCE project and zone as the kubelet. A GCE PD can only be mounted as read/write once or read-only many times. GCE PDs support ownership management and SELinux relabeling.
           */
          gcePersistentDisk?: {
            /**
             * pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
             */
            pdName: string;
            /**
             * readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
             */
            readOnly?: boolean;
            /**
             * fsType is filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
             */
            fsType?: string;
            /**
             * partition is the partition in the volume that you want to mount. If omitted, the default is to mount by volume name. Examples: For volume /dev/sda1, you specify the partition as "1". Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty). More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
             */
            partition?: number;
            [k: string]: unknown;
          };
          /**
           * PersistentVolumeClaimVolumeSource references the user's PVC in the same namespace. This volume finds the bound PV and mounts that volume for the pod. A PersistentVolumeClaimVolumeSource is, essentially, a wrapper around another type of volume that is owned by someone else (the system).
           */
          persistentVolumeClaim?: {
            /**
             * claimName is the name of a PersistentVolumeClaim in the same namespace as the pod using this volume. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
             */
            claimName: string;
            /**
             * readOnly Will force the ReadOnly setting in VolumeMounts. Default false.
             */
            readOnly?: boolean;
            [k: string]: unknown;
          };
          /**
           * Represents an ephemeral volume that is handled by a normal storage driver.
           */
          ephemeral?: {
            /**
             * PersistentVolumeClaimTemplate is used to produce PersistentVolumeClaim objects as part of an EphemeralVolumeSource.
             */
            volumeClaimTemplate?: {
              /**
               * ObjectMeta is metadata that all persisted resources must have, which includes all objects users must create.
               */
              metadata?: {
                /**
                 * Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
                 */
                deletionGracePeriodSeconds?: number;
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
                 * ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field. A workflow can be the user's name, a controller's name, or the name of a specific apply path like "ci-cd". The set of fields is always in the version that the workflow used when modifying the object.
                 */
                managedFields?: {
                  /**
                   * Manager is an identifier of the workflow managing these fields.
                   */
                  manager?: string;
                  /**
                   * Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are 'Apply' and 'Update'.
                   */
                  operation?: string;
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
                  [k: string]: unknown;
                }[];
                /**
                 * List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller.
                 */
                ownerReferences?: {
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
                  /**
                   * API version of the referent.
                   */
                  apiVersion: string;
                  [k: string]: unknown;
                }[];
                /**
                 * UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.
                 *
                 * Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
                 */
                uid?: string;
                /**
                 * Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
                 */
                annotations?: {
                  [k: string]: string;
                };
                /**
                 * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
                 */
                creationTimestamp?: string;
                /**
                 * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
                 */
                deletionTimestamp?: string;
                /**
                 * Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
                 */
                labels?: {
                  [k: string]: string;
                };
                /**
                 * Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.
                 */
                selfLink?: string;
                /**
                 * A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.
                 */
                generation?: number;
                /**
                 * Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
                 */
                name?: string;
                /**
                 * An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.
                 *
                 * Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
                 */
                resourceVersion?: string;
                /**
                 * Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
                 *
                 * Must be a DNS_LABEL. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces
                 */
                namespace?: string;
                [k: string]: unknown;
              };
              /**
               * PersistentVolumeClaimSpec describes the common attributes of storage devices and allows a Source for provider-specific attributes
               */
              spec: {
                /**
                 * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
                 */
                selector?: {
                  /**
                   * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                   */
                  matchExpressions?: {
                    /**
                     * key is the label key that the selector applies to.
                     */
                    key: string;
                    /**
                     * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                     */
                    operator: string;
                    /**
                     * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                     */
                    values?: string[];
                    [k: string]: unknown;
                  }[];
                  /**
                   * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                   */
                  matchLabels?: {
                    [k: string]: string;
                  };
                  [k: string]: unknown;
                };
                /**
                 * storageClassName is the name of the StorageClass required by the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
                 */
                storageClassName?: string;
                /**
                 * volumeName is the binding reference to the PersistentVolume backing this claim.
                 */
                volumeName?: string;
                /**
                 * accessModes contains the desired access modes the volume should have. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
                 */
                accessModes?: ('ReadOnlyMany' | 'ReadWriteMany' | 'ReadWriteOnce' | 'ReadWriteOncePod')[];
                /**
                 * TypedLocalObjectReference contains enough information to let you locate the typed referenced object inside the same namespace.
                 */
                dataSource?: {
                  /**
                   * APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required.
                   */
                  apiGroup?: string;
                  /**
                   * Kind is the type of resource being referenced
                   */
                  kind: string;
                  /**
                   * Name is the name of resource being referenced
                   */
                  name: string;
                  [k: string]: unknown;
                };
                /**
                 * VolumeResourceRequirements describes the storage resource requirements for a volume.
                 */
                resources?: {
                  /**
                   * Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
                   */
                  limits?: {
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
                   * Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
                   */
                  requests?: {
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
                  [k: string]: unknown;
                };
                /**
                 * volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim. If specified, the CSI driver will create or update the volume with the attributes defined in the corresponding VolumeAttributesClass. This has a different purpose than storageClassName, it can be changed after the claim is created. An empty string value means that no VolumeAttributesClass will be applied to the claim but it's not allowed to reset this field to empty string once it is set. If unspecified and the PersistentVolumeClaim is unbound, the default VolumeAttributesClass will be set by the persistentvolume controller if it exists. If the resource referred to by volumeAttributesClass does not exist, this PersistentVolumeClaim will be set to a Pending state, as reflected by the modifyVolumeStatus field, until such as a resource exists. More info: https://kubernetes.io/docs/concepts/storage/volume-attributes-classes/ (Beta) Using this field requires the VolumeAttributesClass feature gate to be enabled (off by default).
                 */
                volumeAttributesClassName?: string;
                /**
                 * volumeMode defines what type of volume is required by the claim. Value of Filesystem is implied when not included in claim spec.
                 *
                 * Possible enum values:
                 *  - `"Block"` means the volume will not be formatted with a filesystem and will remain a raw block device.
                 *  - `"Filesystem"` means the volume will be or is formatted with a filesystem.
                 */
                volumeMode?: 'Block' | 'Filesystem';
                /**
                 * TypedObjectReference contains enough information to let you locate the typed referenced object
                 */
                dataSourceRef?: {
                  /**
                   * APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required.
                   */
                  apiGroup?: string;
                  /**
                   * Kind is the type of resource being referenced
                   */
                  kind: string;
                  /**
                   * Name is the name of resource being referenced
                   */
                  name: string;
                  /**
                   * Namespace is the namespace of resource being referenced Note that when a namespace is specified, a gateway.networking.k8s.io/ReferenceGrant object is required in the referent namespace to allow that namespace's owner to accept the reference. See the ReferenceGrant documentation for details. (Alpha) This field requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
                   */
                  namespace?: string;
                  [k: string]: unknown;
                };
                [k: string]: unknown;
              };
              [k: string]: unknown;
            };
            [k: string]: unknown;
          };
          /**
           * AzureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
           */
          azureDisk?: {
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
            [k: string]: unknown;
          };
          /**
           * Represents a volume that is populated with the contents of a git repository. Git repo volumes do not support ownership management. Git repo volumes support SELinux relabeling.
           *
           * DEPRECATED: GitRepo is deprecated. To provision a container with a git repo, mount an EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir into the Pod's container.
           */
          gitRepo?: {
            /**
             * directory is the target directory name. Must not contain or start with '..'.  If '.' is supplied, the volume directory will be the git repository.  Otherwise, if specified, the volume will contain the git repository in the subdirectory with the given name.
             */
            directory?: string;
            /**
             * repository is the URL
             */
            repository: string;
            /**
             * revision is the commit hash for the specified revision.
             */
            revision?: string;
            [k: string]: unknown;
          };
          /**
           * Represents an ISCSI disk. ISCSI volumes can only be mounted as read/write once. ISCSI volumes support ownership management and SELinux relabeling.
           */
          iscsi?: {
            /**
             * iscsiInterface is the interface Name that uses an iSCSI transport. Defaults to 'default' (tcp).
             */
            iscsiInterface?: string;
            /**
             * lun represents iSCSI Target Lun number.
             */
            lun: number;
            /**
             * portals is the iSCSI Target Portal List. The portal is either an IP or ip_addr:port if the port is other than default (typically TCP ports 860 and 3260).
             */
            portals?: string[];
            /**
             * readOnly here will force the ReadOnly setting in VolumeMounts. Defaults to false.
             */
            readOnly?: boolean;
            /**
             * LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              [k: string]: unknown;
            };
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
             * fsType is the filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#iscsi
             */
            fsType?: string;
            /**
             * initiatorName is the custom iSCSI Initiator Name. If initiatorName is specified with iscsiInterface simultaneously, new iSCSI interface <target portal>:<volume name> will be created for the connection.
             */
            initiatorName?: string;
            /**
             * iqn is the target iSCSI Qualified Name.
             */
            iqn: string;
            [k: string]: unknown;
          };
          /**
           * name of the volume. Must be a DNS_LABEL and unique within the pod. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
           */
          name: string;
          /**
           * Represents an NFS mount that lasts the lifetime of a pod. NFS volumes do not support ownership management or SELinux relabeling.
           */
          nfs?: {
            /**
             * path that is exported by the NFS server. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
             */
            path: string;
            /**
             * readOnly here will force the NFS export to be mounted with read-only permissions. Defaults to false. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
             */
            readOnly?: boolean;
            /**
             * server is the hostname or IP address of the NFS server. More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
             */
            server: string;
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
          /**
           * Represents a Fibre Channel volume. Fibre Channel volumes can only be mounted as read/write once. Fibre Channel volumes support ownership management and SELinux relabeling.
           */
          fc?: {
            /**
             * targetWWNs is Optional: FC target worldwide names (WWNs)
             */
            targetWWNs?: string[];
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
            [k: string]: unknown;
          };
          /**
           * DownwardAPIVolumeSource represents a volume containing downward API info. Downward API volumes support ownership management and SELinux relabeling.
           */
          downwardAPI?: {
            /**
             * Optional: mode bits to use on created files by default. Must be a Optional: mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
             */
            defaultMode?: number;
            /**
             * Items is a list of downward API volume file
             */
            items?: {
              /**
               * ObjectFieldSelector selects an APIVersioned field of an object.
               */
              fieldRef?: {
                /**
                 * Path of the field to select in the specified API version.
                 */
                fieldPath: string;
                /**
                 * Version of the schema the FieldPath is written in terms of, defaults to "v1".
                 */
                apiVersion?: string;
                [k: string]: unknown;
              };
              /**
               * Optional: mode bits used to set permissions on this file, must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. If not specified, the volume defaultMode will be used. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
               */
              mode?: number;
              /**
               * Required: Path is  the relative path name of the file to be created. Must not be absolute or contain the '..' path. Must be utf-8 encoded. The first item of the relative path must not start with '..'
               */
              path: string;
              /**
               * ResourceFieldSelector represents container resources (cpu, memory) and their output format
               */
              resourceFieldRef?: {
                /**
                 * Container name: required for volumes, optional for env vars
                 */
                containerName?: string;
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
                divisor?: string | number;
                /**
                 * Required: resource to select
                 */
                resource: string;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          /**
           * ImageVolumeSource represents a image volume resource.
           */
          image?: {
            /**
             * Policy for pulling OCI objects. Possible values are: Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails. Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present. IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails. Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
             *
             * Possible enum values:
             *  - `"Always"` means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.
             *  - `"IfNotPresent"` means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
             *  - `"Never"` means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn't present
             */
            pullPolicy?: 'Always' | 'IfNotPresent' | 'Never';
            /**
             * Required: Image or artifact reference to be used. Behaves in the same way as pod.spec.containers[*].image. Pull secrets will be assembled in the same way as for the container image by looking up node credentials, SA image pull secrets, and pod spec image pull secrets. More info: https://kubernetes.io/docs/concepts/containers/images This field is optional to allow higher level config management to default or override container images in workload controllers like Deployments and StatefulSets.
             */
            reference?: string;
            [k: string]: unknown;
          };
          /**
           * ScaleIOVolumeSource represents a persistent ScaleIO volume
           */
          scaleIO?: {
            /**
             * readOnly Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
             */
            readOnly?: boolean;
            /**
             * LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
             */
            secretRef: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              [k: string]: unknown;
            };
            /**
             * storageMode indicates whether the storage for a volume should be ThickProvisioned or ThinProvisioned. Default is ThinProvisioned.
             */
            storageMode?: string;
            /**
             * system is the name of the storage system as configured in ScaleIO.
             */
            system: string;
            /**
             * volumeName is the name of a volume already created in the ScaleIO system that is associated with this volume source.
             */
            volumeName?: string;
            /**
             * fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". Default is "xfs".
             */
            fsType?: string;
            /**
             * sslEnabled Flag enable/disable SSL communication with Gateway, default false
             */
            sslEnabled?: boolean;
            /**
             * storagePool is the ScaleIO Storage Pool associated with the protection domain.
             */
            storagePool?: string;
            /**
             * gateway is the host address of the ScaleIO API Gateway.
             */
            gateway: string;
            /**
             * protectionDomain is the name of the ScaleIO Protection Domain for the configured storage.
             */
            protectionDomain?: string;
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
             * LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
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
           * FlexVolume represents a generic volume resource that is provisioned/attached using an exec based plugin.
           */
          flexVolume?: {
            /**
             * driver is the name of the driver to use for this volume.
             */
            driver: string;
            /**
             * fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs", "ntfs". The default filesystem depends on FlexVolume script.
             */
            fsType?: string;
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
             * LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              [k: string]: unknown;
            };
            [k: string]: unknown;
          };
          /**
           * Represents a Ceph Filesystem mount that lasts the lifetime of a pod Cephfs volumes do not support ownership management or SELinux relabeling.
           */
          cephfs?: {
            /**
             * readOnly is Optional: Defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts. More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
             */
            readOnly?: boolean;
            /**
             * secretFile is Optional: SecretFile is the path to key ring for User, default is /etc/ceph/user.secret More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
             */
            secretFile?: string;
            /**
             * LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              [k: string]: unknown;
            };
            /**
             * user is optional: User is the rados user name, default is admin More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
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
            [k: string]: unknown;
          };
          /**
           * Represents a cinder volume resource in Openstack. A Cinder volume must exist before mounting to a container. The volume must also be in the same region as the kubelet. Cinder volumes support ownership management and SELinux relabeling.
           */
          cinder?: {
            /**
             * fsType is the filesystem type to mount. Must be a filesystem type supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
             */
            fsType?: string;
            /**
             * readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
             */
            readOnly?: boolean;
            /**
             * LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
              [k: string]: unknown;
            };
            /**
             * volumeID used to identify the volume in cinder. More info: https://examples.k8s.io/mysql-cinder-pd/README.md
             */
            volumeID: string;
            [k: string]: unknown;
          };
          /**
           * Adapts a Secret into a volume.
           *
           * The contents of the target Secret's Data field will be presented in a volume as files using the keys in the Data field as the file names. Secret volumes support ownership management and SELinux relabeling.
           */
          secret?: {
            /**
             * defaultMode is Optional: mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
             */
            defaultMode?: number;
            /**
             * items If unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the Secret, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the '..' path or start with '..'.
             */
            items?: {
              /**
               * key is the key to project.
               */
              key: string;
              /**
               * mode is Optional: mode bits used to set permissions on this file. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. If not specified, the volume defaultMode will be used. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
               */
              mode?: number;
              /**
               * path is the relative path of the file to map the key to. May not be an absolute path. May not contain the path element '..'. May not start with the string '..'.
               */
              path: string;
              [k: string]: unknown;
            }[];
            /**
             * optional field specify whether the Secret or its keys must be defined
             */
            optional?: boolean;
            /**
             * secretName is the name of the secret in the pod's namespace to use. More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
             */
            secretName?: string;
            [k: string]: unknown;
          };
          /**
           * Represents a vSphere volume resource.
           */
          vsphereVolume?: {
            /**
             * storagePolicyID is the storage Policy Based Management (SPBM) profile ID associated with the StoragePolicyName.
             */
            storagePolicyID?: string;
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
            [k: string]: unknown;
          };
          /**
           * Represents a Persistent Disk resource in AWS.
           *
           * An AWS EBS disk must exist before mounting to a container. The disk must also be in the same AWS zone as the kubelet. An AWS EBS disk can only be mounted as read/write once. AWS EBS volumes support ownership management and SELinux relabeling.
           */
          awsElasticBlockStore?: {
            /**
             * fsType is the filesystem type of the volume that you want to mount. Tip: Ensure that the filesystem type is supported by the host operating system. Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified. More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
             */
            fsType?: string;
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
           * PortworxVolumeSource represents a Portworx volume resource.
           */
          portworxVolume?: {
            /**
             * fSType represents the filesystem type to mount Must be a filesystem type supported by the host operating system. Ex. "ext4", "xfs". Implicitly inferred to be "ext4" if unspecified.
             */
            fsType?: string;
            /**
             * readOnly defaults to false (read/write). ReadOnly here will force the ReadOnly setting in VolumeMounts.
             */
            readOnly?: boolean;
            /**
             * volumeID uniquely identifies a Portworx volume
             */
            volumeID: string;
            [k: string]: unknown;
          };
          /**
           * Represents a Rados Block Device mount that lasts the lifetime of a pod. RBD volumes support ownership management and SELinux relabeling.
           */
          rbd?: {
            /**
             * image is the rados image name. More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
             */
            image: string;
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
             * LocalObjectReference contains enough information to let you locate the referenced object inside the same namespace.
             */
            secretRef?: {
              /**
               * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
               */
              name?: string;
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
            [k: string]: unknown;
          };
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
           * Represents a projected volume source
           */
          projected?: {
            /**
             * defaultMode are the mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
             */
            defaultMode?: number;
            /**
             * sources is the list of volume projections. Each entry in this list handles one source.
             */
            sources?: {
              /**
               * Adapts a ConfigMap into a projected volume.
               *
               * The contents of the target ConfigMap's Data field will be presented in a projected volume as files using the keys in the Data field as the file names, unless the items element is populated with specific mappings of keys to paths. Note that this is identical to a configmap volume source without the default mode.
               */
              configMap?: {
                /**
                 * items if unspecified, each key-value pair in the Data field of the referenced ConfigMap will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the ConfigMap, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the '..' path or start with '..'.
                 */
                items?: {
                  /**
                   * key is the key to project.
                   */
                  key: string;
                  /**
                   * mode is Optional: mode bits used to set permissions on this file. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. If not specified, the volume defaultMode will be used. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
                   */
                  mode?: number;
                  /**
                   * path is the relative path of the file to map the key to. May not be an absolute path. May not contain the path element '..'. May not start with the string '..'.
                   */
                  path: string;
                  [k: string]: unknown;
                }[];
                /**
                 * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                 */
                name?: string;
                /**
                 * optional specify whether the ConfigMap or its keys must be defined
                 */
                optional?: boolean;
                [k: string]: unknown;
              };
              /**
               * Represents downward API info for projecting into a projected volume. Note that this is identical to a downwardAPI volume source without the default mode.
               */
              downwardAPI?: {
                /**
                 * Items is a list of DownwardAPIVolume file
                 */
                items?: {
                  /**
                   * ObjectFieldSelector selects an APIVersioned field of an object.
                   */
                  fieldRef?: {
                    /**
                     * Version of the schema the FieldPath is written in terms of, defaults to "v1".
                     */
                    apiVersion?: string;
                    /**
                     * Path of the field to select in the specified API version.
                     */
                    fieldPath: string;
                    [k: string]: unknown;
                  };
                  /**
                   * Optional: mode bits used to set permissions on this file, must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. If not specified, the volume defaultMode will be used. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
                   */
                  mode?: number;
                  /**
                   * Required: Path is  the relative path name of the file to be created. Must not be absolute or contain the '..' path. Must be utf-8 encoded. The first item of the relative path must not start with '..'
                   */
                  path: string;
                  /**
                   * ResourceFieldSelector represents container resources (cpu, memory) and their output format
                   */
                  resourceFieldRef?: {
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
                    divisor?: string | number;
                    /**
                     * Required: resource to select
                     */
                    resource: string;
                    /**
                     * Container name: required for volumes, optional for env vars
                     */
                    containerName?: string;
                    [k: string]: unknown;
                  };
                  [k: string]: unknown;
                }[];
                [k: string]: unknown;
              };
              /**
               * Adapts a secret into a projected volume.
               *
               * The contents of the target Secret's Data field will be presented in a projected volume as files using the keys in the Data field as the file names. Note that this is identical to a secret volume source without the default mode.
               */
              secret?: {
                /**
                 * items if unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the Secret, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the '..' path or start with '..'.
                 */
                items?: {
                  /**
                   * key is the key to project.
                   */
                  key: string;
                  /**
                   * mode is Optional: mode bits used to set permissions on this file. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. If not specified, the volume defaultMode will be used. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
                   */
                  mode?: number;
                  /**
                   * path is the relative path of the file to map the key to. May not be an absolute path. May not contain the path element '..'. May not start with the string '..'.
                   */
                  path: string;
                  [k: string]: unknown;
                }[];
                /**
                 * Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
                 */
                name?: string;
                /**
                 * optional field specify whether the Secret or its key must be defined
                 */
                optional?: boolean;
                [k: string]: unknown;
              };
              /**
               * ServiceAccountTokenProjection represents a projected service account token volume. This projection can be used to insert a service account token into the pods runtime filesystem for use against APIs (Kubernetes API Server or otherwise).
               */
              serviceAccountToken?: {
                /**
                 * audience is the intended audience of the token. A recipient of a token must identify itself with an identifier specified in the audience of the token, and otherwise should reject the token. The audience defaults to the identifier of the apiserver.
                 */
                audience?: string;
                /**
                 * expirationSeconds is the requested duration of validity of the service account token. As the token approaches expiration, the kubelet volume plugin will proactively rotate the service account token. The kubelet will start trying to rotate the token if the token is older than 80 percent of its time to live or if the token is older than 24 hours.Defaults to 1 hour and must be at least 10 minutes.
                 */
                expirationSeconds?: number;
                /**
                 * path is the path relative to the mount point of the file to project the token into.
                 */
                path: string;
                [k: string]: unknown;
              };
              /**
               * ClusterTrustBundleProjection describes how to select a set of ClusterTrustBundle objects and project their contents into the pod filesystem.
               */
              clusterTrustBundle?: {
                /**
                 * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
                 */
                labelSelector?: {
                  /**
                   * matchExpressions is a list of label selector requirements. The requirements are ANDed.
                   */
                  matchExpressions?: {
                    /**
                     * key is the label key that the selector applies to.
                     */
                    key: string;
                    /**
                     * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
                     */
                    operator: string;
                    /**
                     * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
                     */
                    values?: string[];
                    [k: string]: unknown;
                  }[];
                  /**
                   * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
                   */
                  matchLabels?: {
                    [k: string]: string;
                  };
                  [k: string]: unknown;
                };
                /**
                 * Select a single ClusterTrustBundle by object name.  Mutually-exclusive with signerName and labelSelector.
                 */
                name?: string;
                /**
                 * If true, don't block pod startup if the referenced ClusterTrustBundle(s) aren't available.  If using name, then the named ClusterTrustBundle is allowed not to exist.  If using signerName, then the combination of signerName and labelSelector is allowed to match zero ClusterTrustBundles.
                 */
                optional?: boolean;
                /**
                 * Relative path from the volume root to write the bundle.
                 */
                path: string;
                /**
                 * Select all ClusterTrustBundles that match this signer name. Mutually-exclusive with name.  The contents of all selected ClusterTrustBundles will be unified and deduplicated.
                 */
                signerName?: string;
                [k: string]: unknown;
              };
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          /**
           * Represents a Quobyte mount that lasts the lifetime of a pod. Quobyte volumes do not support ownership management or SELinux relabeling.
           */
          quobyte?: {
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
            [k: string]: unknown;
          };
          [k: string]: unknown;
        }[];
        /**
         * Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers. Value must be a positive integer.
         */
        activeDeadlineSeconds?: number;
        /**
         * Use the host's pid namespace. Optional: Default to false.
         */
        hostPID?: boolean;
        /**
         * Use the host's user namespace. Optional: Default to true. If set to true or not present, the pod will be run in the host user namespace, useful for when the pod needs a feature only available to the host user namespace, such as loading a kernel module with CAP_SYS_MODULE. When set to false, a new userns is created for the pod. Setting false is useful for mitigating container breakout vulnerabilities even allowing users to run their containers as root without actually having root privileges on the host. This field is alpha-level and is only honored by servers that enable the UserNamespacesSupport feature.
         */
        hostUsers?: boolean;
        /**
         * If specified, indicates the pod's priority. "system-node-critical" and "system-cluster-critical" are two special keywords which indicate the highest priorities with the former being the highest priority. Any other name must be defined by creating a PriorityClass object with that name. If not specified, the pod priority will be default or zero if there is no default.
         */
        priorityClassName?: string;
        /**
         * RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run. If unset or empty, the "legacy" RuntimeClass will be used, which is an implicit class with an empty definition that uses the default runtime handler. More info: https://git.k8s.io/enhancements/keps/sig-node/585-runtime-class
         */
        runtimeClassName?: string;
        /**
         * If specified, the fully qualified Pod hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>". If not specified, the pod will not have a domainname at all.
         */
        subdomain?: string;
        /**
         * SchedulingGates is an opaque list of values that if specified will block scheduling the pod. If schedulingGates is not empty, the pod will stay in the SchedulingGated state and the scheduler will not attempt to schedule the pod.
         *
         * SchedulingGates can only be set at pod creation time, and be removed only afterwards.
         */
        schedulingGates?: {
          /**
           * Name of the scheduling gate. Each scheduling gate must have a unique name field.
           */
          name: string;
          [k: string]: unknown;
        }[];
        /**
         * If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default). In Linux containers, this means setting the FQDN in the hostname field of the kernel (the nodename field of struct utsname). In Windows containers, this means setting the registry value of hostname for the registry key HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters to FQDN. If a pod does not have FQDN, this has no effect. Default to false.
         */
        setHostnameAsFQDN?: boolean;
        /**
         * TopologySpreadConstraints describes how a group of pods ought to spread across topology domains. Scheduler will schedule pods in a way which abides by the constraints. All topologySpreadConstraints are ANDed.
         */
        topologySpreadConstraints?: {
          /**
           * NodeTaintsPolicy indicates how we will treat node taints when calculating pod topology spread skew. Options are: - Honor: nodes without taints, along with tainted nodes for which the incoming pod has a toleration, are included. - Ignore: node taints are ignored. All nodes are included.
           *
           * If this value is nil, the behavior is equivalent to the Ignore policy.
           *
           * Possible enum values:
           *  - `"Honor"` means use this scheduling directive when calculating pod topology spread skew.
           *  - `"Ignore"` means ignore this scheduling directive when calculating pod topology spread skew.
           */
          nodeTaintsPolicy?: 'Honor' | 'Ignore';
          /**
           * TopologyKey is the key of node labels. Nodes that have a label with this key and identical values are considered to be in the same topology. We consider each <key, value> as a "bucket", and try to put balanced number of pods into each bucket. We define a domain as a particular instance of a topology. Also, we define an eligible domain as a domain whose nodes meet the requirements of nodeAffinityPolicy and nodeTaintsPolicy. e.g. If TopologyKey is "kubernetes.io/hostname", each Node is a domain of that topology. And, if TopologyKey is "topology.kubernetes.io/zone", each zone is a domain of that topology. It's a required field.
           */
          topologyKey: string;
          /**
           * WhenUnsatisfiable indicates how to deal with a pod if it doesn't satisfy the spread constraint. - DoNotSchedule (default) tells the scheduler not to schedule it. - ScheduleAnyway tells the scheduler to schedule the pod in any location,
           *   but giving higher precedence to topologies that would help reduce the
           *   skew.
           * A constraint is considered "Unsatisfiable" for an incoming pod if and only if every possible node assignment for that pod would violate "MaxSkew" on some topology. For example, in a 3-zone cluster, MaxSkew is set to 1, and pods with the same labelSelector spread as 3/1/1: | zone1 | zone2 | zone3 | | P P P |   P   |   P   | If WhenUnsatisfiable is set to DoNotSchedule, incoming pod can only be scheduled to zone2(zone3) to become 3/2/1(3/1/2) as ActualSkew(2-1) on zone2(zone3) satisfies MaxSkew(1). In other words, the cluster can still be imbalanced, but scheduler won't make it *more* imbalanced. It's a required field.
           *
           * Possible enum values:
           *  - `"DoNotSchedule"` instructs the scheduler not to schedule the pod when constraints are not satisfied.
           *  - `"ScheduleAnyway"` instructs the scheduler to schedule the pod even if constraints are not satisfied.
           */
          whenUnsatisfiable: 'DoNotSchedule' | 'ScheduleAnyway';
          /**
           * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
           */
          labelSelector?: {
            /**
             * matchExpressions is a list of label selector requirements. The requirements are ANDed.
             */
            matchExpressions?: {
              /**
               * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
               */
              operator: string;
              /**
               * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
               */
              values?: string[];
              /**
               * key is the label key that the selector applies to.
               */
              key: string;
              [k: string]: unknown;
            }[];
            /**
             * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
             */
            matchLabels?: {
              [k: string]: string;
            };
            [k: string]: unknown;
          };
          /**
           * MatchLabelKeys is a set of pod label keys to select the pods over which spreading will be calculated. The keys are used to lookup values from the incoming pod labels, those key-value labels are ANDed with labelSelector to select the group of existing pods over which spreading will be calculated for the incoming pod. The same key is forbidden to exist in both MatchLabelKeys and LabelSelector. MatchLabelKeys cannot be set when LabelSelector isn't set. Keys that don't exist in the incoming pod labels will be ignored. A null or empty list means only match against labelSelector.
           *
           * This is a beta field and requires the MatchLabelKeysInPodTopologySpread feature gate to be enabled (enabled by default).
           */
          matchLabelKeys?: string[];
          /**
           * MaxSkew describes the degree to which pods may be unevenly distributed. When `whenUnsatisfiable=DoNotSchedule`, it is the maximum permitted difference between the number of matching pods in the target topology and the global minimum. The global minimum is the minimum number of matching pods in an eligible domain or zero if the number of eligible domains is less than MinDomains. For example, in a 3-zone cluster, MaxSkew is set to 1, and pods with the same labelSelector spread as 2/2/1: In this case, the global minimum is 1. | zone1 | zone2 | zone3 | |  P P  |  P P  |   P   | - if MaxSkew is 1, incoming pod can only be scheduled to zone3 to become 2/2/2; scheduling it onto zone1(zone2) would make the ActualSkew(3-1) on zone1(zone2) violate MaxSkew(1). - if MaxSkew is 2, incoming pod can be scheduled onto any zone. When `whenUnsatisfiable=ScheduleAnyway`, it is used to give higher precedence to topologies that satisfy it. It's a required field. Default value is 1 and 0 is not allowed.
           */
          maxSkew: number;
          /**
           * MinDomains indicates a minimum number of eligible domains. When the number of eligible domains with matching topology keys is less than minDomains, Pod Topology Spread treats "global minimum" as 0, and then the calculation of Skew is performed. And when the number of eligible domains with matching topology keys equals or greater than minDomains, this value has no effect on scheduling. As a result, when the number of eligible domains is less than minDomains, scheduler won't schedule more than maxSkew Pods to those domains. If value is nil, the constraint behaves as if MinDomains is equal to 1. Valid values are integers greater than 0. When value is not nil, WhenUnsatisfiable must be DoNotSchedule.
           *
           * For example, in a 3-zone cluster, MaxSkew is set to 2, MinDomains is set to 5 and pods with the same labelSelector spread as 2/2/2: | zone1 | zone2 | zone3 | |  P P  |  P P  |  P P  | The number of domains is less than 5(MinDomains), so "global minimum" is treated as 0. In this situation, new pod with the same labelSelector cannot be scheduled, because computed skew will be 3(3 - 0) if new Pod is scheduled to any of the three zones, it will violate MaxSkew.
           */
          minDomains?: number;
          /**
           * NodeAffinityPolicy indicates how we will treat Pod's nodeAffinity/nodeSelector when calculating pod topology spread skew. Options are: - Honor: only nodes matching nodeAffinity/nodeSelector are included in the calculations. - Ignore: nodeAffinity/nodeSelector are ignored. All nodes are included in the calculations.
           *
           * If this value is nil, the behavior is equivalent to the Honor policy.
           *
           * Possible enum values:
           *  - `"Honor"` means use this scheduling directive when calculating pod topology spread skew.
           *  - `"Ignore"` means ignore this scheduling directive when calculating pod topology spread skew.
           */
          nodeAffinityPolicy?: 'Honor' | 'Ignore';
          [k: string]: unknown;
        }[];
        [k: string]: unknown;
      };
      [k: string]: unknown;
    };
    /**
     * StatefulSetUpdateStrategy indicates the strategy that the StatefulSet controller will use to perform updates. It includes any additional parameters necessary to perform the update for the indicated strategy.
     */
    updateStrategy?: {
      /**
       * RollingUpdateStatefulSetStrategy is used to communicate parameter for RollingUpdateStatefulSetStrategyType.
       */
      rollingUpdate?: {
        /**
         * IntOrString is a type that can hold an int32 or a string.  When used in JSON or YAML marshalling and unmarshalling, it produces or consumes the inner type.  This allows you to have, for example, a JSON field that can accept a name or number.
         */
        maxUnavailable?: number | string;
        /**
         * Partition indicates the ordinal at which the StatefulSet should be partitioned for updates. During a rolling update, all pods from ordinal Replicas-1 to Partition are updated. All pods from ordinal Partition-1 to 0 remain untouched. This is helpful in being able to do a canary based deployment. The default value is 0.
         */
        partition?: number;
        [k: string]: unknown;
      };
      /**
       * Type indicates the type of the StatefulSetUpdateStrategy. Default is RollingUpdate.
       *
       * Possible enum values:
       *  - `"OnDelete"` triggers the legacy behavior. Version tracking and ordered rolling restarts are disabled. Pods are recreated from the StatefulSetSpec when they are manually deleted. When a scale operation is performed with this strategy,specification version indicated by the StatefulSet's currentRevision.
       *  - `"RollingUpdate"` indicates that update will be applied to all Pods in the StatefulSet with respect to the StatefulSet ordering constraints. When a scale operation is performed with this strategy, new Pods will be created from the specification version indicated by the StatefulSet's updateRevision.
       */
      type?: 'OnDelete' | 'RollingUpdate';
      [k: string]: unknown;
    };
    /**
     * StatefulSetOrdinals describes the policy used for replica ordinal assignment in this StatefulSet.
     */
    ordinals?: {
      /**
       * start is the number representing the first replica's index. It may be used to number replicas from an alternate index (eg: 1-indexed) over the default 0-indexed names, or to orchestrate progressive movement of replicas from one StatefulSet to another. If set, replica indices will be in the range:
       *   [.spec.ordinals.start, .spec.ordinals.start + .spec.replicas).
       * If unset, defaults to 0. Replica indices will be in the range:
       *   [0, .spec.replicas).
       */
      start?: number;
      [k: string]: unknown;
    };
    /**
     * serviceName is the name of the service that governs this StatefulSet. This service must exist before the StatefulSet, and is responsible for the network identity of the set. Pods get DNS/hostnames that follow the pattern: pod-specific-string.serviceName.default.svc.cluster.local where "pod-specific-string" is managed by the StatefulSet controller.
     */
    serviceName?: string;
    /**
     * volumeClaimTemplates is a list of claims that pods are allowed to reference. The StatefulSet controller is responsible for mapping network identities to claims in a way that maintains the identity of a pod. Every claim in this list must have at least one matching (by name) volumeMount in one container in the template. A claim in this list takes precedence over any volumes in the template, with the same name.
     */
    volumeClaimTemplates?: {
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
         * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
         */
        deletionTimestamp?: string;
        /**
         * Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.
         */
        selfLink?: string;
        /**
         * Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
         */
        annotations?: {
          [k: string]: string;
        };
        /**
         * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
         */
        creationTimestamp?: string;
        /**
         * Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
         */
        deletionGracePeriodSeconds?: number;
        /**
         * Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order.  Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list.
         */
        finalizers?: string[];
        /**
         * A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.
         */
        generation?: number;
        /**
         * ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field. A workflow can be the user's name, a controller's name, or the name of a specific apply path like "ci-cd". The set of fields is always in the version that the workflow used when modifying the object.
         */
        managedFields?: {
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
          /**
           * Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
           */
          subresource?: string;
          [k: string]: unknown;
        }[];
        /**
         * Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
         */
        name?: string;
        /**
         * GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
         *
         * If this field is specified and the generated name exists, the server will return a 409.
         *
         * Applied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency
         */
        generateName?: string;
        /**
         * An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.
         *
         * Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
         */
        resourceVersion?: string;
        /**
         * Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
         */
        labels?: {
          [k: string]: string;
        };
        /**
         * Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
         *
         * Must be a DNS_LABEL. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces
         */
        namespace?: string;
        /**
         * List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller.
         */
        ownerReferences?: {
          /**
           * Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
           */
          name: string;
          /**
           * UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
           */
          uid: string;
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
          [k: string]: unknown;
        }[];
        /**
         * UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.
         *
         * Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
         */
        uid?: string;
        [k: string]: unknown;
      };
      /**
       * PersistentVolumeClaimSpec describes the common attributes of storage devices and allows a Source for provider-specific attributes
       */
      spec?: {
        /**
         * accessModes contains the desired access modes the volume should have. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
         */
        accessModes?: ('ReadOnlyMany' | 'ReadWriteMany' | 'ReadWriteOnce' | 'ReadWriteOncePod')[];
        /**
         * TypedObjectReference contains enough information to let you locate the typed referenced object
         */
        dataSourceRef?: {
          /**
           * Name is the name of resource being referenced
           */
          name: string;
          /**
           * Namespace is the namespace of resource being referenced Note that when a namespace is specified, a gateway.networking.k8s.io/ReferenceGrant object is required in the referent namespace to allow that namespace's owner to accept the reference. See the ReferenceGrant documentation for details. (Alpha) This field requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
           */
          namespace?: string;
          /**
           * APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required.
           */
          apiGroup?: string;
          /**
           * Kind is the type of resource being referenced
           */
          kind: string;
          [k: string]: unknown;
        };
        /**
         * VolumeResourceRequirements describes the storage resource requirements for a volume.
         */
        resources?: {
          /**
           * Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
           */
          limits?: {
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
           * Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
           */
          requests?: {
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
          [k: string]: unknown;
        };
        /**
         * A label selector is a label query over a set of resources. The result of matchLabels and matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
         */
        selector?: {
          /**
           * matchExpressions is a list of label selector requirements. The requirements are ANDed.
           */
          matchExpressions?: {
            /**
             * key is the label key that the selector applies to.
             */
            key: string;
            /**
             * operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
             */
            operator: string;
            /**
             * values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
             */
            values?: string[];
            [k: string]: unknown;
          }[];
          /**
           * matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
           */
          matchLabels?: {
            [k: string]: string;
          };
          [k: string]: unknown;
        };
        /**
         * storageClassName is the name of the StorageClass required by the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
         */
        storageClassName?: string;
        /**
         * volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim. If specified, the CSI driver will create or update the volume with the attributes defined in the corresponding VolumeAttributesClass. This has a different purpose than storageClassName, it can be changed after the claim is created. An empty string value means that no VolumeAttributesClass will be applied to the claim but it's not allowed to reset this field to empty string once it is set. If unspecified and the PersistentVolumeClaim is unbound, the default VolumeAttributesClass will be set by the persistentvolume controller if it exists. If the resource referred to by volumeAttributesClass does not exist, this PersistentVolumeClaim will be set to a Pending state, as reflected by the modifyVolumeStatus field, until such as a resource exists. More info: https://kubernetes.io/docs/concepts/storage/volume-attributes-classes/ (Beta) Using this field requires the VolumeAttributesClass feature gate to be enabled (off by default).
         */
        volumeAttributesClassName?: string;
        /**
         * TypedLocalObjectReference contains enough information to let you locate the typed referenced object inside the same namespace.
         */
        dataSource?: {
          /**
           * APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required.
           */
          apiGroup?: string;
          /**
           * Kind is the type of resource being referenced
           */
          kind: string;
          /**
           * Name is the name of resource being referenced
           */
          name: string;
          [k: string]: unknown;
        };
        /**
         * volumeMode defines what type of volume is required by the claim. Value of Filesystem is implied when not included in claim spec.
         *
         * Possible enum values:
         *  - `"Block"` means the volume will not be formatted with a filesystem and will remain a raw block device.
         *  - `"Filesystem"` means the volume will be or is formatted with a filesystem.
         */
        volumeMode?: 'Block' | 'Filesystem';
        /**
         * volumeName is the binding reference to the PersistentVolume backing this claim.
         */
        volumeName?: string;
        [k: string]: unknown;
      };
      /**
       * PersistentVolumeClaimStatus is the current status of a persistent volume claim.
       */
      status?: {
        /**
         * allocatedResourceStatuses stores status of resource being resized for the given PVC. Key names follow standard Kubernetes label syntax. Valid values are either:
         * 	* Un-prefixed keys:
         * 		- storage - the capacity of the volume.
         * 	* Custom resources must use implementation-defined prefixed names such as "example.com/my-custom-resource"
         * Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used.
         *
         * ClaimResourceStatus can be in any of following states:
         * 	- ControllerResizeInProgress:
         * 		State set when resize controller starts resizing the volume in control-plane.
         * 	- ControllerResizeFailed:
         * 		State set when resize has failed in resize controller with a terminal error.
         * 	- NodeResizePending:
         * 		State set when resize controller has finished resizing the volume but further resizing of
         * 		volume is needed on the node.
         * 	- NodeResizeInProgress:
         * 		State set when kubelet starts resizing the volume.
         * 	- NodeResizeFailed:
         * 		State set when resizing has failed in kubelet with a terminal error. Transient errors don't set
         * 		NodeResizeFailed.
         * For example: if expanding a PVC for more capacity - this field can be one of the following states:
         * 	- pvc.status.allocatedResourceStatus['storage'] = "ControllerResizeInProgress"
         *      - pvc.status.allocatedResourceStatus['storage'] = "ControllerResizeFailed"
         *      - pvc.status.allocatedResourceStatus['storage'] = "NodeResizePending"
         *      - pvc.status.allocatedResourceStatus['storage'] = "NodeResizeInProgress"
         *      - pvc.status.allocatedResourceStatus['storage'] = "NodeResizeFailed"
         * When this field is not set, it means that no resize operation is in progress for the given PVC.
         *
         * A controller that receives PVC update with previously unknown resourceName or ClaimResourceStatus should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC.
         *
         * This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
         */
        allocatedResourceStatuses?: {
          [k: string]:
            | 'ControllerResizeInProgress'
            | 'ControllerResizeInfeasible'
            | 'NodeResizeInProgress'
            | 'NodeResizeInfeasible'
            | 'NodeResizePending';
        };
        /**
         * allocatedResources tracks the resources allocated to a PVC including its capacity. Key names follow standard Kubernetes label syntax. Valid values are either:
         * 	* Un-prefixed keys:
         * 		- storage - the capacity of the volume.
         * 	* Custom resources must use implementation-defined prefixed names such as "example.com/my-custom-resource"
         * Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used.
         *
         * Capacity reported here may be larger than the actual capacity when a volume expansion operation is requested. For storage quota, the larger value from allocatedResources and PVC.spec.resources is used. If allocatedResources is not set, PVC.spec.resources alone is used for quota calculation. If a volume expansion capacity request is lowered, allocatedResources is only lowered if there are no expansion operations in progress and if the actual volume capacity is equal or lower than the requested capacity.
         *
         * A controller that receives PVC update with previously unknown resourceName should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC.
         *
         * This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
         */
        allocatedResources?: {
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
         * capacity represents the actual resources of the underlying volume.
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
         * conditions is the current Condition of persistent volume claim. If underlying persistent volume is being resized then the Condition will be set to 'Resizing'.
         */
        conditions?: {
          /**
           * Status is the status of the condition. Can be True, False, Unknown. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=state%20of%20pvc-,conditions.status,-(string)%2C%20required
           */
          status: string;
          /**
           * Type is the type of the condition. More info: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#:~:text=set%20to%20%27ResizeStarted%27.-,PersistentVolumeClaimCondition,-contains%20details%20about
           */
          type: string;
          /**
           * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
           */
          lastProbeTime?: string;
          /**
           * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
           */
          lastTransitionTime?: string;
          /**
           * message is the human-readable message indicating details about last transition.
           */
          message?: string;
          /**
           * reason is a unique, this should be a short, machine understandable string that gives the reason for condition's last transition. If it reports "Resizing" that means the underlying persistent volume is being resized.
           */
          reason?: string;
          [k: string]: unknown;
        }[];
        /**
         * currentVolumeAttributesClassName is the current name of the VolumeAttributesClass the PVC is using. When unset, there is no VolumeAttributeClass applied to this PersistentVolumeClaim This is a beta field and requires enabling VolumeAttributesClass feature (off by default).
         */
        currentVolumeAttributesClassName?: string;
        /**
         * ModifyVolumeStatus represents the status object of ControllerModifyVolume operation
         */
        modifyVolumeStatus?: {
          /**
           * status is the status of the ControllerModifyVolume operation. It can be in any of following states:
           *  - Pending
           *    Pending indicates that the PersistentVolumeClaim cannot be modified due to unmet requirements, such as
           *    the specified VolumeAttributesClass not existing.
           *  - InProgress
           *    InProgress indicates that the volume is being modified.
           *  - Infeasible
           *   Infeasible indicates that the request has been rejected as invalid by the CSI driver. To
           * 	  resolve the error, a valid VolumeAttributesClass needs to be specified.
           * Note: New statuses can be added in the future. Consumers should check for unknown statuses and fail appropriately.
           *
           * Possible enum values:
           *  - `"InProgress"` InProgress indicates that the volume is being modified
           *  - `"Infeasible"` Infeasible indicates that the request has been rejected as invalid by the CSI driver. To resolve the error, a valid VolumeAttributesClass needs to be specified
           *  - `"Pending"` Pending indicates that the PersistentVolumeClaim cannot be modified due to unmet requirements, such as the specified VolumeAttributesClass not existing
           */
          status: 'InProgress' | 'Infeasible' | 'Pending';
          /**
           * targetVolumeAttributesClassName is the name of the VolumeAttributesClass the PVC currently being reconciled
           */
          targetVolumeAttributesClassName?: string;
          [k: string]: unknown;
        };
        /**
         * phase represents the current phase of PersistentVolumeClaim.
         *
         * Possible enum values:
         *  - `"Bound"` used for PersistentVolumeClaims that are bound
         *  - `"Lost"` used for PersistentVolumeClaims that lost their underlying PersistentVolume. The claim was bound to a PersistentVolume and this volume does not exist any longer and all data on it was lost.
         *  - `"Pending"` used for PersistentVolumeClaims that are not yet bound
         */
        phase?: 'Bound' | 'Lost' | 'Pending';
        /**
         * accessModes contains the actual access modes the volume backing the PVC has. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
         */
        accessModes?: ('ReadOnlyMany' | 'ReadWriteMany' | 'ReadWriteOnce' | 'ReadWriteOncePod')[];
        [k: string]: unknown;
      };
      [k: string]: unknown;
    }[];
    [k: string]: unknown;
  };
  /**
   * StatefulSetStatus represents the current state of a StatefulSet.
   */
  status?: {
    /**
     * Represents the latest available observations of a statefulset's current state.
     */
    conditions?: {
      /**
       * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
       */
      lastTransitionTime?: string;
      /**
       * A human readable message indicating details about the transition.
       */
      message?: string;
      /**
       * The reason for the condition's last transition.
       */
      reason?: string;
      /**
       * Status of the condition, one of True, False, Unknown.
       */
      status: string;
      /**
       * Type of statefulset condition.
       */
      type: string;
      [k: string]: unknown;
    }[];
    /**
     * currentReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version indicated by currentRevision.
     */
    currentReplicas?: number;
    /**
     * currentRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence [0,currentReplicas).
     */
    currentRevision?: string;
    /**
     * observedGeneration is the most recent generation observed for this StatefulSet. It corresponds to the StatefulSet's generation, which is updated on mutation by the API Server.
     */
    observedGeneration?: number;
    /**
     * readyReplicas is the number of pods created for this StatefulSet with a Ready Condition.
     */
    readyReplicas?: number;
    /**
     * updatedReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version indicated by updateRevision.
     */
    updatedReplicas?: number;
    /**
     * Total number of available pods (ready for at least minReadySeconds) targeted by this statefulset.
     */
    availableReplicas?: number;
    /**
     * collisionCount is the count of hash collisions for the StatefulSet. The StatefulSet controller uses this field as a collision avoidance mechanism when it needs to create the name for the newest ControllerRevision.
     */
    collisionCount?: number;
    /**
     * replicas is the number of Pods created by the StatefulSet controller.
     */
    replicas: number;
    /**
     * updateRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence [replicas-updatedReplicas,replicas)
     */
    updateRevision?: string;
    [k: string]: unknown;
  };
  [k: string]: unknown;
}

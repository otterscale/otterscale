/** Generated from Remote JSON Schema for inference.networking.k8s.io.v1.InferencePool */

/**
 * InferencePool is the Schema for the InferencePools API.
 *
 */
export interface InferenceNetworkingK8SIoV1InferencePool {
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
     * Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.
     */
    selfLink?: string;
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
      /**
       * APIVersion defines the version of this resource that this field set applies to. The format is "group/version" just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
       */
      apiVersion?: string;
      /**
       * FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: "FieldsV1"
       */
      fieldsType?: string;
      [k: string]: unknown;
    }[];
    /**
     * UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.
     *
     * Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
     */
    uid?: string;
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
    /**
     * Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
     */
    annotations?: {
      [k: string]: string;
    };
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
    [k: string]: unknown;
  };
  /**
   * Spec defines the desired state of the InferencePool.
   */
  spec: {
    /**
     * Selector determines which Pods are members of this inference pool.
     * It matches Pods by their labels only within the same namespace; cross-namespace
     * selection is not supported.
     *
     * The structure of this LabelSelector is intentionally simple to be compatible
     * with Kubernetes Service selectors, as some implementations may translate
     * this configuration into a Service resource.
     */
    selector: {
      /**
       * MatchLabels contains a set of required {key,value} pairs.
       * An object must match every label in this map to be selected.
       * The matching logic is an AND operation on all entries.
       */
      matchLabels: {
        /**
         * LabelValue is the value of a label. This is used for validation
         * of maps. This matches the Kubernetes label validation rules:
         * * must be 63 characters or less (can be empty),
         * * unless empty, must begin and end with an alphanumeric character ([a-z0-9A-Z]),
         * * could contain dashes (-), underscores (_), dots (.), and alphanumerics between.
         *
         * Valid values include:
         *
         * * MyValue
         * * my.name
         * * 123-my-value
         */
        [k: string]: string;
      };
      [k: string]: unknown;
    };
    /**
     * TargetPorts defines a list of ports that are exposed by this InferencePool.
     * Every port will be treated as a distinctive endpoint by EPP,
     * addressable as a 'podIP:portNumber' combination.
     *
     * @minItems 1
     * @maxItems 8
     */
    targetPorts: {
      /**
       * Number defines the port number to access the selected model server Pods.
       * The number must be in the range 1 to 65535.
       */
      number: number;
      [k: string]: unknown;
    }[];
    /**
     * EndpointPickerRef is a reference to the Endpoint Picker extension and its
     * associated configuration.
     */
    endpointPickerRef: {
      /**
       * FailureMode configures how the parent handles the case when the Endpoint Picker extension
       * is non-responsive. When unspecified, defaults to "FailClose".
       */
      failureMode?: 'FailOpen' | 'FailClose';
      /**
       * Group is the group of the referent API object. When unspecified, the default value
       * is "", representing the Core API group.
       */
      group?: string;
      /**
       * Kind is the Kubernetes resource kind of the referent.
       *
       * Required if the referent is ambiguous, e.g. service with multiple ports.
       *
       * Defaults to "Service" when not specified.
       *
       * ExternalName services can refer to CNAME DNS records that may live
       * outside of the cluster and as such are difficult to reason about in
       * terms of conformance. They also may not be safe to forward to (see
       * CVE-2021-25740 for more information). Implementations MUST NOT
       * support ExternalName Services.
       */
      kind?: string;
      /**
       * Name is the name of the referent API object.
       */
      name: string;
      /**
       * Port is the port of the Endpoint Picker extension service.
       *
       * Port is required when the referent is a Kubernetes Service. In this
       * case, the port number is the service port number, not the target port.
       * For other resources, destination port might be derived from the referent
       * resource or this field.
       */
      port?: {
        /**
         * Number defines the port number to access the selected model server Pods.
         * The number must be in the range 1 to 65535.
         */
        number: number;
        [k: string]: unknown;
      };
      [k: string]: unknown;
    };
    [k: string]: unknown;
  };
  /**
   * Status defines the observed state of the InferencePool.
   */
  status?: {
    /**
     * Parents is a list of parent resources, typically Gateways, that are associated with
     * the InferencePool, and the status of the InferencePool with respect to each parent.
     *
     * A controller that manages the InferencePool, must add an entry for each parent it manages
     * and remove the parent entry when the controller no longer considers the InferencePool to
     * be associated with that parent.
     *
     * A maximum of 32 parents will be represented in this list. When the list is empty,
     * it indicates that the InferencePool is not associated with any parents.
     *
     * @maxItems 32
     */
    parents?: {
      /**
       * ParentRef is used to identify the parent resource that this status
       * is associated with. It is used to match the InferencePool with the parent
       * resource, such as a Gateway.
       */
      parentRef: {
        /**
         * Name is the name of the referent API object.
         */
        name: string;
        /**
         * Namespace is the namespace of the referenced object. When unspecified, the local
         * namespace is inferred.
         *
         * Note that when a namespace different than the local namespace is specified,
         * a ReferenceGrant object is required in the referent namespace to allow that
         * namespace's owner to accept the reference. See the ReferenceGrant
         * documentation for details: https://gateway-api.sigs.k8s.io/api-types/referencegrant/
         */
        namespace?: string;
        /**
         * Group is the group of the referent API object. When unspecified, the referent is assumed
         * to be in the "gateway.networking.k8s.io" API group.
         */
        group?: string;
        /**
         * Kind is the kind of the referent API object. When unspecified, the referent is assumed
         * to be a "Gateway" kind.
         */
        kind?: string;
        [k: string]: unknown;
      };
      /**
       * Conditions is a list of status conditions that provide information about the observed
       * state of the InferencePool. This field is required to be set by the controller that
       * manages the InferencePool.
       *
       * Supported condition types are:
       *
       * * "Accepted"
       * * "ResolvedRefs"
       *
       * @maxItems 8
       */
      conditions?: {
        /**
         * lastTransitionTime is the last time the condition transitioned from one status to another.
         * This should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.
         */
        lastTransitionTime: string;
        /**
         * message is a human readable message indicating details about the transition.
         * This may be an empty string.
         */
        message: string;
        /**
         * observedGeneration represents the .metadata.generation that the condition was set based upon.
         * For instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date
         * with respect to the current state of the instance.
         */
        observedGeneration?: number;
        /**
         * reason contains a programmatic identifier indicating the reason for the condition's last transition.
         * Producers of specific condition types may define expected values and meanings for this field,
         * and whether the values are considered a guaranteed API.
         * The value should be a CamelCase string.
         * This field may not be empty.
         */
        reason: string;
        /**
         * status of the condition, one of True, False, Unknown.
         */
        status: 'True' | 'False' | 'Unknown';
        /**
         * type of condition in CamelCase or in foo.example.com/CamelCase.
         */
        type: string;
        [k: string]: unknown;
      }[];
      /**
       * ControllerName is a domain/path string that indicates the name of the controller that
       * wrote this status. This corresponds with the GatewayClass controllerName field when the
       * parentRef references a Gateway kind.
       *
       * Example: "example.net/gateway-controller".
       *
       * The format of this field is DOMAIN "/" PATH, where DOMAIN and PATH are valid Kubernetes names:
       *
       *  https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
       *
       * Controllers MAY populate this field when writing status. When populating this field, controllers
       * should ensure that entries to status populated with their ControllerName are cleaned up when they
       * are no longer necessary.
       */
      controllerName?: string;
      [k: string]: unknown;
    }[];
    [k: string]: unknown;
  };
  [k: string]: unknown;
}

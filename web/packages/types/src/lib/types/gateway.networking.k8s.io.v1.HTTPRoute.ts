/** Generated from Remote JSON Schema for gateway.networking.k8s.io.v1.HTTPRoute */

/**
 * HTTPRoute provides a way to route HTTP requests. This includes the capability
 * to match requests by hostname, path, header, or query param. Filters can be
 * used to specify additional processing steps. Backends specify where matching
 * requests should be routed.
 */
export interface GatewayNetworkingK8SIoV1HTTPRoute {
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
     * UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.
     *
     * Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
     */
    uid?: string;
    /**
     * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
     */
    creationTimestamp?: string;
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
     * Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
     */
    name?: string;
    /**
     * Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
     */
    deletionGracePeriodSeconds?: number;
    /**
     * Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.
     */
    deletionTimestamp?: string;
    /**
     * Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order.  Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list.
     */
    finalizers?: string[];
    /**
     * Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
     *
     * Must be a DNS_LABEL. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces
     */
    namespace?: string;
    /**
     * GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
     *
     * If this field is specified and the generated name exists, the server will return a 409.
     *
     * Applied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency
     */
    generateName?: string;
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
    [k: string]: unknown;
  };
  /**
   * Spec defines the desired state of HTTPRoute.
   */
  spec: {
    /**
     * Hostnames defines a set of hostnames that should match against the HTTP Host
     * header to select a HTTPRoute used to process the request. Implementations
     * MUST ignore any port value specified in the HTTP Host header while
     * performing a match and (absent of any applicable header modification
     * configuration) MUST forward this header unmodified to the backend.
     *
     * Valid values for Hostnames are determined by RFC 1123 definition of a
     * hostname with 2 notable exceptions:
     *
     * 1. IPs are not allowed.
     * 2. A hostname may be prefixed with a wildcard label (`*.`). The wildcard
     *    label must appear by itself as the first label.
     *
     * If a hostname is specified by both the Listener and HTTPRoute, there
     * must be at least one intersecting hostname for the HTTPRoute to be
     * attached to the Listener. For example:
     *
     * * A Listener with `test.example.com` as the hostname matches HTTPRoutes
     *   that have either not specified any hostnames, or have specified at
     *   least one of `test.example.com` or `*.example.com`.
     * * A Listener with `*.example.com` as the hostname matches HTTPRoutes
     *   that have either not specified any hostnames or have specified at least
     *   one hostname that matches the Listener hostname. For example,
     *   `*.example.com`, `test.example.com`, and `foo.test.example.com` would
     *   all match. On the other hand, `example.com` and `test.example.net` would
     *   not match.
     *
     * Hostnames that are prefixed with a wildcard label (`*.`) are interpreted
     * as a suffix match. That means that a match for `*.example.com` would match
     * both `test.example.com`, and `foo.test.example.com`, but not `example.com`.
     *
     * If both the Listener and HTTPRoute have specified hostnames, any
     * HTTPRoute hostnames that do not match the Listener hostname MUST be
     * ignored. For example, if a Listener specified `*.example.com`, and the
     * HTTPRoute specified `test.example.com` and `test.example.net`,
     * `test.example.net` must not be considered for a match.
     *
     * If both the Listener and HTTPRoute have specified hostnames, and none
     * match with the criteria above, then the HTTPRoute is not accepted. The
     * implementation must raise an 'Accepted' Condition with a status of
     * `False` in the corresponding RouteParentStatus.
     *
     * In the event that multiple HTTPRoutes specify intersecting hostnames (e.g.
     * overlapping wildcard matching and exact matching hostnames), precedence must
     * be given to rules from the HTTPRoute with the largest number of:
     *
     * * Characters in a matching non-wildcard hostname.
     * * Characters in a matching hostname.
     *
     * If ties exist across multiple Routes, the matching precedence rules for
     * HTTPRouteMatches takes over.
     *
     * Support: Core
     *
     * @maxItems 16
     */
    hostnames?: string[];
    /**
     * ParentRefs references the resources (usually Gateways) that a Route wants
     * to be attached to. Note that the referenced parent resource needs to
     * allow this for the attachment to be complete. For Gateways, that means
     * the Gateway needs to allow attachment from Routes of this kind and
     * namespace. For Services, that means the Service must either be in the same
     * namespace for a "producer" route, or the mesh implementation must support
     * and allow "consumer" routes for the referenced Service. ReferenceGrant is
     * not applicable for governing ParentRefs to Services - it is not possible to
     * create a "producer" route for a Service in a different namespace from the
     * Route.
     *
     * There are two kinds of parent resources with "Core" support:
     *
     * * Gateway (Gateway conformance profile)
     * * Service (Mesh conformance profile, ClusterIP Services only)
     *
     * This API may be extended in the future to support additional kinds of parent
     * resources.
     *
     * ParentRefs must be _distinct_. This means either that:
     *
     * * They select different objects.  If this is the case, then parentRef
     *   entries are distinct. In terms of fields, this means that the
     *   multi-part key defined by `group`, `kind`, `namespace`, and `name` must
     *   be unique across all parentRef entries in the Route.
     * * They do not select different objects, but for each optional field used,
     *   each ParentRef that selects the same object must set the same set of
     *   optional fields to different values. If one ParentRef sets a
     *   combination of optional fields, all must set the same combination.
     *
     * Some examples:
     *
     * * If one ParentRef sets `sectionName`, all ParentRefs referencing the
     *   same object must also set `sectionName`.
     * * If one ParentRef sets `port`, all ParentRefs referencing the same
     *   object must also set `port`.
     * * If one ParentRef sets `sectionName` and `port`, all ParentRefs
     *   referencing the same object must also set `sectionName` and `port`.
     *
     * It is possible to separately reference multiple distinct objects that may
     * be collapsed by an implementation. For example, some implementations may
     * choose to merge compatible Gateway Listeners together. If that is the
     * case, the list of routes attached to those resources should also be
     * merged.
     *
     * Note that for ParentRefs that cross namespace boundaries, there are specific
     * rules. Cross-namespace references are only valid if they are explicitly
     * allowed by something in the namespace they are referring to. For example,
     * Gateway has the AllowedRoutes field, and ReferenceGrant provides a
     * generic way to enable other kinds of cross-namespace reference.
     *
     * @maxItems 32
     */
    parentRefs?: {
      /**
       * Group is the group of the referent.
       * When unspecified, "gateway.networking.k8s.io" is inferred.
       * To set the core API group (such as for a "Service" kind referent),
       * Group must be explicitly set to "" (empty string).
       *
       * Support: Core
       */
      group?: string;
      /**
       * Kind is kind of the referent.
       *
       * There are two kinds of parent resources with "Core" support:
       *
       * * Gateway (Gateway conformance profile)
       * * Service (Mesh conformance profile, ClusterIP Services only)
       *
       * Support for other resources is Implementation-Specific.
       */
      kind?: string;
      /**
       * Name is the name of the referent.
       *
       * Support: Core
       */
      name: string;
      /**
       * Namespace is the namespace of the referent. When unspecified, this refers
       * to the local namespace of the Route.
       *
       * Note that there are specific rules for ParentRefs which cross namespace
       * boundaries. Cross-namespace references are only valid if they are explicitly
       * allowed by something in the namespace they are referring to. For example:
       * Gateway has the AllowedRoutes field, and ReferenceGrant provides a
       * generic way to enable any other kind of cross-namespace reference.
       *
       * Support: Core
       */
      namespace?: string;
      /**
       * Port is the network port this Route targets. It can be interpreted
       * differently based on the type of parent resource.
       *
       * When the parent resource is a Gateway, this targets all listeners
       * listening on the specified port that also support this kind of Route(and
       * select this Route). It's not recommended to set `Port` unless the
       * networking behaviors specified in a Route must apply to a specific port
       * as opposed to a listener(s) whose port(s) may be changed. When both Port
       * and SectionName are specified, the name and port of the selected listener
       * must match both specified values.
       *
       * Implementations MAY choose to support other parent resources.
       * Implementations supporting other types of parent resources MUST clearly
       * document how/if Port is interpreted.
       *
       * For the purpose of status, an attachment is considered successful as
       * long as the parent resource accepts it partially. For example, Gateway
       * listeners can restrict which Routes can attach to them by Route kind,
       * namespace, or hostname. If 1 of 2 Gateway listeners accept attachment
       * from the referencing Route, the Route MUST be considered successfully
       * attached. If no Gateway listeners accept attachment from this Route,
       * the Route MUST be considered detached from the Gateway.
       *
       * Support: Extended
       */
      port?: number;
      /**
       * SectionName is the name of a section within the target resource. In the
       * following resources, SectionName is interpreted as the following:
       *
       * * Gateway: Listener name. When both Port (experimental) and SectionName
       * are specified, the name and port of the selected listener must match
       * both specified values.
       * * Service: Port name. When both Port (experimental) and SectionName
       * are specified, the name and port of the selected listener must match
       * both specified values.
       *
       * Implementations MAY choose to support attaching Routes to other resources.
       * If that is the case, they MUST clearly document how SectionName is
       * interpreted.
       *
       * When unspecified (empty string), this will reference the entire resource.
       * For the purpose of status, an attachment is considered successful if at
       * least one section in the parent resource accepts it. For example, Gateway
       * listeners can restrict which Routes can attach to them by Route kind,
       * namespace, or hostname. If 1 of 2 Gateway listeners accept attachment from
       * the referencing Route, the Route MUST be considered successfully
       * attached. If no Gateway listeners accept attachment from this Route, the
       * Route MUST be considered detached from the Gateway.
       *
       * Support: Core
       */
      sectionName?: string;
      [k: string]: unknown;
    }[];
    /**
     * Rules are a list of HTTP matchers, filters and actions.
     *
     * @maxItems 16
     */
    rules?: {
      /**
       * Name is the name of the route rule. This name MUST be unique within a Route if it is set.
       *
       * Support: Extended
       */
      name?: string;
      /**
       * Timeouts defines the timeouts that can be configured for an HTTP request.
       *
       * Support: Extended
       */
      timeouts?: {
        /**
         * BackendRequest specifies a timeout for an individual request from the gateway
         * to a backend. This covers the time from when the request first starts being
         * sent from the gateway to when the full response has been received from the backend.
         *
         * Setting a timeout to the zero duration (e.g. "0s") SHOULD disable the timeout
         * completely. Implementations that cannot completely disable the timeout MUST
         * instead interpret the zero duration as the longest possible value to which
         * the timeout can be set.
         *
         * An entire client HTTP transaction with a gateway, covered by the Request timeout,
         * may result in more than one call from the gateway to the destination backend,
         * for example, if automatic retries are supported.
         *
         * The value of BackendRequest must be a Gateway API Duration string as defined by
         * GEP-2257.  When this field is unspecified, its behavior is implementation-specific;
         * when specified, the value of BackendRequest must be no more than the value of the
         * Request timeout (since the Request timeout encompasses the BackendRequest timeout).
         *
         * Support: Extended
         */
        backendRequest?: string;
        /**
         * Request specifies the maximum duration for a gateway to respond to an HTTP request.
         * If the gateway has not been able to respond before this deadline is met, the gateway
         * MUST return a timeout error.
         *
         * For example, setting the `rules.timeouts.request` field to the value `10s` in an
         * `HTTPRoute` will cause a timeout if a client request is taking longer than 10 seconds
         * to complete.
         *
         * Setting a timeout to the zero duration (e.g. "0s") SHOULD disable the timeout
         * completely. Implementations that cannot completely disable the timeout MUST
         * instead interpret the zero duration as the longest possible value to which
         * the timeout can be set.
         *
         * This timeout is intended to cover as close to the whole request-response transaction
         * as possible although an implementation MAY choose to start the timeout after the entire
         * request stream has been received instead of immediately after the transaction is
         * initiated by the client.
         *
         * The value of Request is a Gateway API Duration string as defined by GEP-2257. When this
         * field is unspecified, request timeout behavior is implementation-specific.
         *
         * Support: Extended
         */
        request?: string;
        [k: string]: unknown;
      };
      /**
       * BackendRefs defines the backend(s) where matching requests should be
       * sent.
       *
       * Failure behavior here depends on how many BackendRefs are specified and
       * how many are invalid.
       *
       * If *all* entries in BackendRefs are invalid, and there are also no filters
       * specified in this route rule, *all* traffic which matches this rule MUST
       * receive a 500 status code.
       *
       * See the HTTPBackendRef definition for the rules about what makes a single
       * HTTPBackendRef invalid.
       *
       * When a HTTPBackendRef is invalid, 500 status codes MUST be returned for
       * requests that would have otherwise been routed to an invalid backend. If
       * multiple backends are specified, and some are invalid, the proportion of
       * requests that would otherwise have been routed to an invalid backend
       * MUST receive a 500 status code.
       *
       * For example, if two backends are specified with equal weights, and one is
       * invalid, 50 percent of traffic must receive a 500. Implementations may
       * choose how that 50 percent is determined.
       *
       * When a HTTPBackendRef refers to a Service that has no ready endpoints,
       * implementations SHOULD return a 503 for requests to that backend instead.
       * If an implementation chooses to do this, all of the above rules for 500 responses
       * MUST also apply for responses that return a 503.
       *
       * Support: Core for Kubernetes Service
       *
       * Support: Extended for Kubernetes ServiceImport
       *
       * Support: Implementation-specific for any other resource
       *
       * Support for weight: Core
       *
       * @maxItems 16
       */
      backendRefs?: {
        /**
         * Filters defined at this level should be executed if and only if the
         * request is being forwarded to the backend defined here.
         *
         * Support: Implementation-specific (For broader support of filters, use the
         * Filters field in HTTPRouteRule.)
         *
         * @maxItems 16
         */
        filters?: {
          /**
           * ExtensionRef is an optional, implementation-specific extension to the
           * "filter" behavior.  For example, resource "myroutefilter" in group
           * "networking.example.net"). ExtensionRef MUST NOT be used for core and
           * extended filters.
           *
           * This filter can be used multiple times within the same rule.
           *
           * Support: Implementation-specific
           */
          extensionRef?: {
            /**
             * Group is the group of the referent. For example, "gateway.networking.k8s.io".
             * When unspecified or empty string, core API group is inferred.
             */
            group: string;
            /**
             * Kind is kind of the referent. For example "HTTPRoute" or "Service".
             */
            kind: string;
            /**
             * Name is the name of the referent.
             */
            name: string;
            [k: string]: unknown;
          };
          /**
           * RequestHeaderModifier defines a schema for a filter that modifies request
           * headers.
           *
           * Support: Core
           */
          requestHeaderModifier?: {
            /**
             * Add adds the given header(s) (name, value) to the request
             * before the action. It appends to any existing values associated
             * with the header name.
             *
             * Input:
             *   GET /foo HTTP/1.1
             *   my-header: foo
             *
             * Config:
             *   add:
             *   - name: "my-header"
             *     value: "bar,baz"
             *
             * Output:
             *   GET /foo HTTP/1.1
             *   my-header: foo,bar,baz
             *
             * @maxItems 16
             */
            add?: {
              /**
               * Name is the name of the HTTP Header to be matched. Name matching MUST be
               * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
               *
               * If multiple entries specify equivalent header names, the first entry with
               * an equivalent name MUST be considered for a match. Subsequent entries
               * with an equivalent header name MUST be ignored. Due to the
               * case-insensitivity of header names, "foo" and "Foo" are considered
               * equivalent.
               */
              name: string;
              /**
               * Value is the value of HTTP Header to be matched.
               */
              value: string;
              [k: string]: unknown;
            }[];
            /**
             * Remove the given header(s) from the HTTP request before the action. The
             * value of Remove is a list of HTTP header names. Note that the header
             * names are case-insensitive (see
             * https://datatracker.ietf.org/doc/html/rfc2616#section-4.2).
             *
             * Input:
             *   GET /foo HTTP/1.1
             *   my-header1: foo
             *   my-header2: bar
             *   my-header3: baz
             *
             * Config:
             *   remove: ["my-header1", "my-header3"]
             *
             * Output:
             *   GET /foo HTTP/1.1
             *   my-header2: bar
             *
             * @maxItems 16
             */
            remove?: string[];
            /**
             * Set overwrites the request with the given header (name, value)
             * before the action.
             *
             * Input:
             *   GET /foo HTTP/1.1
             *   my-header: foo
             *
             * Config:
             *   set:
             *   - name: "my-header"
             *     value: "bar"
             *
             * Output:
             *   GET /foo HTTP/1.1
             *   my-header: bar
             *
             * @maxItems 16
             */
            set?: {
              /**
               * Name is the name of the HTTP Header to be matched. Name matching MUST be
               * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
               *
               * If multiple entries specify equivalent header names, the first entry with
               * an equivalent name MUST be considered for a match. Subsequent entries
               * with an equivalent header name MUST be ignored. Due to the
               * case-insensitivity of header names, "foo" and "Foo" are considered
               * equivalent.
               */
              name: string;
              /**
               * Value is the value of HTTP Header to be matched.
               */
              value: string;
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          /**
           * RequestMirror defines a schema for a filter that mirrors requests.
           * Requests are sent to the specified destination, but responses from
           * that destination are ignored.
           *
           * This filter can be used multiple times within the same rule. Note that
           * not all implementations will be able to support mirroring to multiple
           * backends.
           *
           * Support: Extended
           */
          requestMirror?: {
            /**
             * BackendRef references a resource where mirrored requests are sent.
             *
             * Mirrored requests must be sent only to a single destination endpoint
             * within this BackendRef, irrespective of how many endpoints are present
             * within this BackendRef.
             *
             * If the referent cannot be found, this BackendRef is invalid and must be
             * dropped from the Gateway. The controller must ensure the "ResolvedRefs"
             * condition on the Route status is set to `status: False` and not configure
             * this backend in the underlying implementation.
             *
             * If there is a cross-namespace reference to an *existing* object
             * that is not allowed by a ReferenceGrant, the controller must ensure the
             * "ResolvedRefs"  condition on the Route is set to `status: False`,
             * with the "RefNotPermitted" reason and not configure this backend in the
             * underlying implementation.
             *
             * In either error case, the Message of the `ResolvedRefs` Condition
             * should be used to provide more detail about the problem.
             *
             * Support: Extended for Kubernetes Service
             *
             * Support: Implementation-specific for any other resource
             */
            backendRef: {
              /**
               * Name is the name of the referent.
               */
              name: string;
              /**
               * Namespace is the namespace of the backend. When unspecified, the local
               * namespace is inferred.
               *
               * Note that when a namespace different than the local namespace is specified,
               * a ReferenceGrant object is required in the referent namespace to allow that
               * namespace's owner to accept the reference. See the ReferenceGrant
               * documentation for details.
               *
               * Support: Core
               */
              namespace?: string;
              /**
               * Port specifies the destination port number to use for this resource.
               * Port is required when the referent is a Kubernetes Service. In this
               * case, the port number is the service port number, not the target port.
               * For other resources, destination port might be derived from the referent
               * resource or this field.
               */
              port?: number;
              /**
               * Group is the group of the referent. For example, "gateway.networking.k8s.io".
               * When unspecified or empty string, core API group is inferred.
               */
              group?: string;
              /**
               * Kind is the Kubernetes resource kind of the referent. For example
               * "Service".
               *
               * Defaults to "Service" when not specified.
               *
               * ExternalName services can refer to CNAME DNS records that may live
               * outside of the cluster and as such are difficult to reason about in
               * terms of conformance. They also may not be safe to forward to (see
               * CVE-2021-25740 for more information). Implementations SHOULD NOT
               * support ExternalName Services.
               *
               * Support: Core (Services with a type other than ExternalName)
               *
               * Support: Implementation-specific (Services with type ExternalName)
               */
              kind?: string;
              [k: string]: unknown;
            };
            /**
             * Fraction represents the fraction of requests that should be
             * mirrored to BackendRef.
             *
             * Only one of Fraction or Percent may be specified. If neither field
             * is specified, 100% of requests will be mirrored.
             */
            fraction?: {
              denominator?: number;
              numerator: number;
              [k: string]: unknown;
            };
            /**
             * Percent represents the percentage of requests that should be
             * mirrored to BackendRef. Its minimum value is 0 (indicating 0% of
             * requests) and its maximum value is 100 (indicating 100% of requests).
             *
             * Only one of Fraction or Percent may be specified. If neither field
             * is specified, 100% of requests will be mirrored.
             */
            percent?: number;
            [k: string]: unknown;
          };
          /**
           * RequestRedirect defines a schema for a filter that responds to the
           * request with an HTTP redirection.
           *
           * Support: Core
           */
          requestRedirect?: {
            /**
             * Hostname is the hostname to be used in the value of the `Location`
             * header in the response.
             * When empty, the hostname in the `Host` header of the request is used.
             *
             * Support: Core
             */
            hostname?: string;
            /**
             * Path defines parameters used to modify the path of the incoming request.
             * The modified path is then used to construct the `Location` header. When
             * empty, the request path is used as-is.
             *
             * Support: Extended
             */
            path?: {
              /**
               * ReplaceFullPath specifies the value with which to replace the full path
               * of a request during a rewrite or redirect.
               */
              replaceFullPath?: string;
              /**
               * ReplacePrefixMatch specifies the value with which to replace the prefix
               * match of a request during a rewrite or redirect. For example, a request
               * to "/foo/bar" with a prefix match of "/foo" and a ReplacePrefixMatch
               * of "/xyz" would be modified to "/xyz/bar".
               *
               * Note that this matches the behavior of the PathPrefix match type. This
               * matches full path elements. A path element refers to the list of labels
               * in the path split by the `/` separator. When specified, a trailing `/` is
               * ignored. For example, the paths `/abc`, `/abc/`, and `/abc/def` would all
               * match the prefix `/abc`, but the path `/abcd` would not.
               *
               * ReplacePrefixMatch is only compatible with a `PathPrefix` HTTPRouteMatch.
               * Using any other HTTPRouteMatch type on the same HTTPRouteRule will result in
               * the implementation setting the Accepted Condition for the Route to `status: False`.
               *
               * Request Path | Prefix Match | Replace Prefix | Modified Path
               */
              replacePrefixMatch?: string;
              /**
               * Type defines the type of path modifier. Additional types may be
               * added in a future release of the API.
               *
               * Note that values may be added to this enum, implementations
               * must ensure that unknown values will not cause a crash.
               *
               * Unknown values here must result in the implementation setting the
               * Accepted Condition for the Route to `status: False`, with a
               * Reason of `UnsupportedValue`.
               */
              type: 'ReplaceFullPath' | 'ReplacePrefixMatch';
              [k: string]: unknown;
            };
            /**
             * Port is the port to be used in the value of the `Location`
             * header in the response.
             *
             * If no port is specified, the redirect port MUST be derived using the
             * following rules:
             *
             * * If redirect scheme is not-empty, the redirect port MUST be the well-known
             *   port associated with the redirect scheme. Specifically "http" to port 80
             *   and "https" to port 443. If the redirect scheme does not have a
             *   well-known port, the listener port of the Gateway SHOULD be used.
             * * If redirect scheme is empty, the redirect port MUST be the Gateway
             *   Listener port.
             *
             * Implementations SHOULD NOT add the port number in the 'Location'
             * header in the following cases:
             *
             * * A Location header that will use HTTP (whether that is determined via
             *   the Listener protocol or the Scheme field) _and_ use port 80.
             * * A Location header that will use HTTPS (whether that is determined via
             *   the Listener protocol or the Scheme field) _and_ use port 443.
             *
             * Support: Extended
             */
            port?: number;
            /**
             * Scheme is the scheme to be used in the value of the `Location` header in
             * the response. When empty, the scheme of the request is used.
             *
             * Scheme redirects can affect the port of the redirect, for more information,
             * refer to the documentation for the port field of this filter.
             *
             * Note that values may be added to this enum, implementations
             * must ensure that unknown values will not cause a crash.
             *
             * Unknown values here must result in the implementation setting the
             * Accepted Condition for the Route to `status: False`, with a
             * Reason of `UnsupportedValue`.
             *
             * Support: Extended
             */
            scheme?: 'http' | 'https';
            /**
             * StatusCode is the HTTP status code to be used in response.
             *
             * Note that values may be added to this enum, implementations
             * must ensure that unknown values will not cause a crash.
             *
             * Unknown values here must result in the implementation setting the
             * Accepted Condition for the Route to `status: False`, with a
             * Reason of `UnsupportedValue`.
             *
             * Support: Core
             */
            statusCode?: 301 | 302;
            [k: string]: unknown;
          };
          /**
           * ResponseHeaderModifier defines a schema for a filter that modifies response
           * headers.
           *
           * Support: Extended
           */
          responseHeaderModifier?: {
            /**
             * Add adds the given header(s) (name, value) to the request
             * before the action. It appends to any existing values associated
             * with the header name.
             *
             * Input:
             *   GET /foo HTTP/1.1
             *   my-header: foo
             *
             * Config:
             *   add:
             *   - name: "my-header"
             *     value: "bar,baz"
             *
             * Output:
             *   GET /foo HTTP/1.1
             *   my-header: foo,bar,baz
             *
             * @maxItems 16
             */
            add?: {
              /**
               * Name is the name of the HTTP Header to be matched. Name matching MUST be
               * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
               *
               * If multiple entries specify equivalent header names, the first entry with
               * an equivalent name MUST be considered for a match. Subsequent entries
               * with an equivalent header name MUST be ignored. Due to the
               * case-insensitivity of header names, "foo" and "Foo" are considered
               * equivalent.
               */
              name: string;
              /**
               * Value is the value of HTTP Header to be matched.
               */
              value: string;
              [k: string]: unknown;
            }[];
            /**
             * Remove the given header(s) from the HTTP request before the action. The
             * value of Remove is a list of HTTP header names. Note that the header
             * names are case-insensitive (see
             * https://datatracker.ietf.org/doc/html/rfc2616#section-4.2).
             *
             * Input:
             *   GET /foo HTTP/1.1
             *   my-header1: foo
             *   my-header2: bar
             *   my-header3: baz
             *
             * Config:
             *   remove: ["my-header1", "my-header3"]
             *
             * Output:
             *   GET /foo HTTP/1.1
             *   my-header2: bar
             *
             * @maxItems 16
             */
            remove?: string[];
            /**
             * Set overwrites the request with the given header (name, value)
             * before the action.
             *
             * Input:
             *   GET /foo HTTP/1.1
             *   my-header: foo
             *
             * Config:
             *   set:
             *   - name: "my-header"
             *     value: "bar"
             *
             * Output:
             *   GET /foo HTTP/1.1
             *   my-header: bar
             *
             * @maxItems 16
             */
            set?: {
              /**
               * Name is the name of the HTTP Header to be matched. Name matching MUST be
               * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
               *
               * If multiple entries specify equivalent header names, the first entry with
               * an equivalent name MUST be considered for a match. Subsequent entries
               * with an equivalent header name MUST be ignored. Due to the
               * case-insensitivity of header names, "foo" and "Foo" are considered
               * equivalent.
               */
              name: string;
              /**
               * Value is the value of HTTP Header to be matched.
               */
              value: string;
              [k: string]: unknown;
            }[];
            [k: string]: unknown;
          };
          /**
           * Type identifies the type of filter to apply. As with other API fields,
           * types are classified into three conformance levels:
           *
           * - Core: Filter types and their corresponding configuration defined by
           *   "Support: Core" in this package, e.g. "RequestHeaderModifier". All
           *   implementations must support core filters.
           *
           * - Extended: Filter types and their corresponding configuration defined by
           *   "Support: Extended" in this package, e.g. "RequestMirror". Implementers
           *   are encouraged to support extended filters.
           *
           * - Implementation-specific: Filters that are defined and supported by
           *   specific vendors.
           *   In the future, filters showing convergence in behavior across multiple
           *   implementations will be considered for inclusion in extended or core
           *   conformance levels. Filter-specific configuration for such filters
           *   is specified using the ExtensionRef field. `Type` should be set to
           *   "ExtensionRef" for custom filters.
           *
           * Implementers are encouraged to define custom implementation types to
           * extend the core API with implementation-specific behavior.
           *
           * If a reference to a custom filter type cannot be resolved, the filter
           * MUST NOT be skipped. Instead, requests that would have been processed by
           * that filter MUST receive a HTTP error response.
           *
           * Note that values may be added to this enum, implementations
           * must ensure that unknown values will not cause a crash.
           *
           * Unknown values here must result in the implementation setting the
           * Accepted Condition for the Route to `status: False`, with a
           * Reason of `UnsupportedValue`.
           */
          type:
            | 'RequestHeaderModifier'
            | 'ResponseHeaderModifier'
            | 'RequestMirror'
            | 'RequestRedirect'
            | 'URLRewrite'
            | 'ExtensionRef';
          /**
           * URLRewrite defines a schema for a filter that modifies a request during forwarding.
           *
           * Support: Extended
           */
          urlRewrite?: {
            /**
             * Hostname is the value to be used to replace the Host header value during
             * forwarding.
             *
             * Support: Extended
             */
            hostname?: string;
            /**
             * Path defines a path rewrite.
             *
             * Support: Extended
             */
            path?: {
              /**
               * ReplacePrefixMatch specifies the value with which to replace the prefix
               * match of a request during a rewrite or redirect. For example, a request
               * to "/foo/bar" with a prefix match of "/foo" and a ReplacePrefixMatch
               * of "/xyz" would be modified to "/xyz/bar".
               *
               * Note that this matches the behavior of the PathPrefix match type. This
               * matches full path elements. A path element refers to the list of labels
               * in the path split by the `/` separator. When specified, a trailing `/` is
               * ignored. For example, the paths `/abc`, `/abc/`, and `/abc/def` would all
               * match the prefix `/abc`, but the path `/abcd` would not.
               *
               * ReplacePrefixMatch is only compatible with a `PathPrefix` HTTPRouteMatch.
               * Using any other HTTPRouteMatch type on the same HTTPRouteRule will result in
               * the implementation setting the Accepted Condition for the Route to `status: False`.
               *
               * Request Path | Prefix Match | Replace Prefix | Modified Path
               */
              replacePrefixMatch?: string;
              /**
               * Type defines the type of path modifier. Additional types may be
               * added in a future release of the API.
               *
               * Note that values may be added to this enum, implementations
               * must ensure that unknown values will not cause a crash.
               *
               * Unknown values here must result in the implementation setting the
               * Accepted Condition for the Route to `status: False`, with a
               * Reason of `UnsupportedValue`.
               */
              type: 'ReplaceFullPath' | 'ReplacePrefixMatch';
              /**
               * ReplaceFullPath specifies the value with which to replace the full path
               * of a request during a rewrite or redirect.
               */
              replaceFullPath?: string;
              [k: string]: unknown;
            };
            [k: string]: unknown;
          };
          [k: string]: unknown;
        }[];
        /**
         * Group is the group of the referent. For example, "gateway.networking.k8s.io".
         * When unspecified or empty string, core API group is inferred.
         */
        group?: string;
        /**
         * Kind is the Kubernetes resource kind of the referent. For example
         * "Service".
         *
         * Defaults to "Service" when not specified.
         *
         * ExternalName services can refer to CNAME DNS records that may live
         * outside of the cluster and as such are difficult to reason about in
         * terms of conformance. They also may not be safe to forward to (see
         * CVE-2021-25740 for more information). Implementations SHOULD NOT
         * support ExternalName Services.
         *
         * Support: Core (Services with a type other than ExternalName)
         *
         * Support: Implementation-specific (Services with type ExternalName)
         */
        kind?: string;
        /**
         * Name is the name of the referent.
         */
        name: string;
        /**
         * Namespace is the namespace of the backend. When unspecified, the local
         * namespace is inferred.
         *
         * Note that when a namespace different than the local namespace is specified,
         * a ReferenceGrant object is required in the referent namespace to allow that
         * namespace's owner to accept the reference. See the ReferenceGrant
         * documentation for details.
         *
         * Support: Core
         */
        namespace?: string;
        /**
         * Port specifies the destination port number to use for this resource.
         * Port is required when the referent is a Kubernetes Service. In this
         * case, the port number is the service port number, not the target port.
         * For other resources, destination port might be derived from the referent
         * resource or this field.
         */
        port?: number;
        /**
         * Weight specifies the proportion of requests forwarded to the referenced
         * backend. This is computed as weight/(sum of all weights in this
         * BackendRefs list). For non-zero values, there may be some epsilon from
         * the exact proportion defined here depending on the precision an
         * implementation supports. Weight is not a percentage and the sum of
         * weights does not need to equal 100.
         *
         * If only one backend is specified and it has a weight greater than 0, 100%
         * of the traffic is forwarded to that backend. If weight is set to 0, no
         * traffic should be forwarded for this entry. If unspecified, weight
         * defaults to 1.
         *
         * Support for this field varies based on the context where used.
         */
        weight?: number;
        [k: string]: unknown;
      }[];
      /**
       * Filters define the filters that are applied to requests that match
       * this rule.
       *
       * Wherever possible, implementations SHOULD implement filters in the order
       * they are specified.
       *
       * Implementations MAY choose to implement this ordering strictly, rejecting
       * any combination or order of filters that cannot be supported. If implementations
       * choose a strict interpretation of filter ordering, they MUST clearly document
       * that behavior.
       *
       * To reject an invalid combination or order of filters, implementations SHOULD
       * consider the Route Rules with this configuration invalid. If all Route Rules
       * in a Route are invalid, the entire Route would be considered invalid. If only
       * a portion of Route Rules are invalid, implementations MUST set the
       * "PartiallyInvalid" condition for the Route.
       *
       * Conformance-levels at this level are defined based on the type of filter:
       *
       * - ALL core filters MUST be supported by all implementations.
       * - Implementers are encouraged to support extended filters.
       * - Implementation-specific custom filters have no API guarantees across
       *   implementations.
       *
       * Specifying the same filter multiple times is not supported unless explicitly
       * indicated in the filter.
       *
       * All filters are expected to be compatible with each other except for the
       * URLRewrite and RequestRedirect filters, which may not be combined. If an
       * implementation cannot support other combinations of filters, they must clearly
       * document that limitation. In cases where incompatible or unsupported
       * filters are specified and cause the `Accepted` condition to be set to status
       * `False`, implementations may use the `IncompatibleFilters` reason to specify
       * this configuration error.
       *
       * Support: Core
       *
       * @maxItems 16
       */
      filters?: {
        /**
         * URLRewrite defines a schema for a filter that modifies a request during forwarding.
         *
         * Support: Extended
         */
        urlRewrite?: {
          /**
           * Hostname is the value to be used to replace the Host header value during
           * forwarding.
           *
           * Support: Extended
           */
          hostname?: string;
          /**
           * Path defines a path rewrite.
           *
           * Support: Extended
           */
          path?: {
            /**
             * ReplacePrefixMatch specifies the value with which to replace the prefix
             * match of a request during a rewrite or redirect. For example, a request
             * to "/foo/bar" with a prefix match of "/foo" and a ReplacePrefixMatch
             * of "/xyz" would be modified to "/xyz/bar".
             *
             * Note that this matches the behavior of the PathPrefix match type. This
             * matches full path elements. A path element refers to the list of labels
             * in the path split by the `/` separator. When specified, a trailing `/` is
             * ignored. For example, the paths `/abc`, `/abc/`, and `/abc/def` would all
             * match the prefix `/abc`, but the path `/abcd` would not.
             *
             * ReplacePrefixMatch is only compatible with a `PathPrefix` HTTPRouteMatch.
             * Using any other HTTPRouteMatch type on the same HTTPRouteRule will result in
             * the implementation setting the Accepted Condition for the Route to `status: False`.
             *
             * Request Path | Prefix Match | Replace Prefix | Modified Path
             */
            replacePrefixMatch?: string;
            /**
             * Type defines the type of path modifier. Additional types may be
             * added in a future release of the API.
             *
             * Note that values may be added to this enum, implementations
             * must ensure that unknown values will not cause a crash.
             *
             * Unknown values here must result in the implementation setting the
             * Accepted Condition for the Route to `status: False`, with a
             * Reason of `UnsupportedValue`.
             */
            type: 'ReplaceFullPath' | 'ReplacePrefixMatch';
            /**
             * ReplaceFullPath specifies the value with which to replace the full path
             * of a request during a rewrite or redirect.
             */
            replaceFullPath?: string;
            [k: string]: unknown;
          };
          [k: string]: unknown;
        };
        /**
         * ExtensionRef is an optional, implementation-specific extension to the
         * "filter" behavior.  For example, resource "myroutefilter" in group
         * "networking.example.net"). ExtensionRef MUST NOT be used for core and
         * extended filters.
         *
         * This filter can be used multiple times within the same rule.
         *
         * Support: Implementation-specific
         */
        extensionRef?: {
          /**
           * Group is the group of the referent. For example, "gateway.networking.k8s.io".
           * When unspecified or empty string, core API group is inferred.
           */
          group: string;
          /**
           * Kind is kind of the referent. For example "HTTPRoute" or "Service".
           */
          kind: string;
          /**
           * Name is the name of the referent.
           */
          name: string;
          [k: string]: unknown;
        };
        /**
         * RequestHeaderModifier defines a schema for a filter that modifies request
         * headers.
         *
         * Support: Core
         */
        requestHeaderModifier?: {
          /**
           * Add adds the given header(s) (name, value) to the request
           * before the action. It appends to any existing values associated
           * with the header name.
           *
           * Input:
           *   GET /foo HTTP/1.1
           *   my-header: foo
           *
           * Config:
           *   add:
           *   - name: "my-header"
           *     value: "bar,baz"
           *
           * Output:
           *   GET /foo HTTP/1.1
           *   my-header: foo,bar,baz
           *
           * @maxItems 16
           */
          add?: {
            /**
             * Name is the name of the HTTP Header to be matched. Name matching MUST be
             * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
             *
             * If multiple entries specify equivalent header names, the first entry with
             * an equivalent name MUST be considered for a match. Subsequent entries
             * with an equivalent header name MUST be ignored. Due to the
             * case-insensitivity of header names, "foo" and "Foo" are considered
             * equivalent.
             */
            name: string;
            /**
             * Value is the value of HTTP Header to be matched.
             */
            value: string;
            [k: string]: unknown;
          }[];
          /**
           * Remove the given header(s) from the HTTP request before the action. The
           * value of Remove is a list of HTTP header names. Note that the header
           * names are case-insensitive (see
           * https://datatracker.ietf.org/doc/html/rfc2616#section-4.2).
           *
           * Input:
           *   GET /foo HTTP/1.1
           *   my-header1: foo
           *   my-header2: bar
           *   my-header3: baz
           *
           * Config:
           *   remove: ["my-header1", "my-header3"]
           *
           * Output:
           *   GET /foo HTTP/1.1
           *   my-header2: bar
           *
           * @maxItems 16
           */
          remove?: string[];
          /**
           * Set overwrites the request with the given header (name, value)
           * before the action.
           *
           * Input:
           *   GET /foo HTTP/1.1
           *   my-header: foo
           *
           * Config:
           *   set:
           *   - name: "my-header"
           *     value: "bar"
           *
           * Output:
           *   GET /foo HTTP/1.1
           *   my-header: bar
           *
           * @maxItems 16
           */
          set?: {
            /**
             * Value is the value of HTTP Header to be matched.
             */
            value: string;
            /**
             * Name is the name of the HTTP Header to be matched. Name matching MUST be
             * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
             *
             * If multiple entries specify equivalent header names, the first entry with
             * an equivalent name MUST be considered for a match. Subsequent entries
             * with an equivalent header name MUST be ignored. Due to the
             * case-insensitivity of header names, "foo" and "Foo" are considered
             * equivalent.
             */
            name: string;
            [k: string]: unknown;
          }[];
          [k: string]: unknown;
        };
        /**
         * RequestMirror defines a schema for a filter that mirrors requests.
         * Requests are sent to the specified destination, but responses from
         * that destination are ignored.
         *
         * This filter can be used multiple times within the same rule. Note that
         * not all implementations will be able to support mirroring to multiple
         * backends.
         *
         * Support: Extended
         */
        requestMirror?: {
          /**
           * Percent represents the percentage of requests that should be
           * mirrored to BackendRef. Its minimum value is 0 (indicating 0% of
           * requests) and its maximum value is 100 (indicating 100% of requests).
           *
           * Only one of Fraction or Percent may be specified. If neither field
           * is specified, 100% of requests will be mirrored.
           */
          percent?: number;
          /**
           * BackendRef references a resource where mirrored requests are sent.
           *
           * Mirrored requests must be sent only to a single destination endpoint
           * within this BackendRef, irrespective of how many endpoints are present
           * within this BackendRef.
           *
           * If the referent cannot be found, this BackendRef is invalid and must be
           * dropped from the Gateway. The controller must ensure the "ResolvedRefs"
           * condition on the Route status is set to `status: False` and not configure
           * this backend in the underlying implementation.
           *
           * If there is a cross-namespace reference to an *existing* object
           * that is not allowed by a ReferenceGrant, the controller must ensure the
           * "ResolvedRefs"  condition on the Route is set to `status: False`,
           * with the "RefNotPermitted" reason and not configure this backend in the
           * underlying implementation.
           *
           * In either error case, the Message of the `ResolvedRefs` Condition
           * should be used to provide more detail about the problem.
           *
           * Support: Extended for Kubernetes Service
           *
           * Support: Implementation-specific for any other resource
           */
          backendRef: {
            /**
             * Namespace is the namespace of the backend. When unspecified, the local
             * namespace is inferred.
             *
             * Note that when a namespace different than the local namespace is specified,
             * a ReferenceGrant object is required in the referent namespace to allow that
             * namespace's owner to accept the reference. See the ReferenceGrant
             * documentation for details.
             *
             * Support: Core
             */
            namespace?: string;
            /**
             * Port specifies the destination port number to use for this resource.
             * Port is required when the referent is a Kubernetes Service. In this
             * case, the port number is the service port number, not the target port.
             * For other resources, destination port might be derived from the referent
             * resource or this field.
             */
            port?: number;
            /**
             * Group is the group of the referent. For example, "gateway.networking.k8s.io".
             * When unspecified or empty string, core API group is inferred.
             */
            group?: string;
            /**
             * Kind is the Kubernetes resource kind of the referent. For example
             * "Service".
             *
             * Defaults to "Service" when not specified.
             *
             * ExternalName services can refer to CNAME DNS records that may live
             * outside of the cluster and as such are difficult to reason about in
             * terms of conformance. They also may not be safe to forward to (see
             * CVE-2021-25740 for more information). Implementations SHOULD NOT
             * support ExternalName Services.
             *
             * Support: Core (Services with a type other than ExternalName)
             *
             * Support: Implementation-specific (Services with type ExternalName)
             */
            kind?: string;
            /**
             * Name is the name of the referent.
             */
            name: string;
            [k: string]: unknown;
          };
          /**
           * Fraction represents the fraction of requests that should be
           * mirrored to BackendRef.
           *
           * Only one of Fraction or Percent may be specified. If neither field
           * is specified, 100% of requests will be mirrored.
           */
          fraction?: {
            denominator?: number;
            numerator: number;
            [k: string]: unknown;
          };
          [k: string]: unknown;
        };
        /**
         * RequestRedirect defines a schema for a filter that responds to the
         * request with an HTTP redirection.
         *
         * Support: Core
         */
        requestRedirect?: {
          /**
           * Port is the port to be used in the value of the `Location`
           * header in the response.
           *
           * If no port is specified, the redirect port MUST be derived using the
           * following rules:
           *
           * * If redirect scheme is not-empty, the redirect port MUST be the well-known
           *   port associated with the redirect scheme. Specifically "http" to port 80
           *   and "https" to port 443. If the redirect scheme does not have a
           *   well-known port, the listener port of the Gateway SHOULD be used.
           * * If redirect scheme is empty, the redirect port MUST be the Gateway
           *   Listener port.
           *
           * Implementations SHOULD NOT add the port number in the 'Location'
           * header in the following cases:
           *
           * * A Location header that will use HTTP (whether that is determined via
           *   the Listener protocol or the Scheme field) _and_ use port 80.
           * * A Location header that will use HTTPS (whether that is determined via
           *   the Listener protocol or the Scheme field) _and_ use port 443.
           *
           * Support: Extended
           */
          port?: number;
          /**
           * Scheme is the scheme to be used in the value of the `Location` header in
           * the response. When empty, the scheme of the request is used.
           *
           * Scheme redirects can affect the port of the redirect, for more information,
           * refer to the documentation for the port field of this filter.
           *
           * Note that values may be added to this enum, implementations
           * must ensure that unknown values will not cause a crash.
           *
           * Unknown values here must result in the implementation setting the
           * Accepted Condition for the Route to `status: False`, with a
           * Reason of `UnsupportedValue`.
           *
           * Support: Extended
           */
          scheme?: 'http' | 'https';
          /**
           * StatusCode is the HTTP status code to be used in response.
           *
           * Note that values may be added to this enum, implementations
           * must ensure that unknown values will not cause a crash.
           *
           * Unknown values here must result in the implementation setting the
           * Accepted Condition for the Route to `status: False`, with a
           * Reason of `UnsupportedValue`.
           *
           * Support: Core
           */
          statusCode?: 301 | 302;
          /**
           * Hostname is the hostname to be used in the value of the `Location`
           * header in the response.
           * When empty, the hostname in the `Host` header of the request is used.
           *
           * Support: Core
           */
          hostname?: string;
          /**
           * Path defines parameters used to modify the path of the incoming request.
           * The modified path is then used to construct the `Location` header. When
           * empty, the request path is used as-is.
           *
           * Support: Extended
           */
          path?: {
            /**
             * ReplaceFullPath specifies the value with which to replace the full path
             * of a request during a rewrite or redirect.
             */
            replaceFullPath?: string;
            /**
             * ReplacePrefixMatch specifies the value with which to replace the prefix
             * match of a request during a rewrite or redirect. For example, a request
             * to "/foo/bar" with a prefix match of "/foo" and a ReplacePrefixMatch
             * of "/xyz" would be modified to "/xyz/bar".
             *
             * Note that this matches the behavior of the PathPrefix match type. This
             * matches full path elements. A path element refers to the list of labels
             * in the path split by the `/` separator. When specified, a trailing `/` is
             * ignored. For example, the paths `/abc`, `/abc/`, and `/abc/def` would all
             * match the prefix `/abc`, but the path `/abcd` would not.
             *
             * ReplacePrefixMatch is only compatible with a `PathPrefix` HTTPRouteMatch.
             * Using any other HTTPRouteMatch type on the same HTTPRouteRule will result in
             * the implementation setting the Accepted Condition for the Route to `status: False`.
             *
             * Request Path | Prefix Match | Replace Prefix | Modified Path
             */
            replacePrefixMatch?: string;
            /**
             * Type defines the type of path modifier. Additional types may be
             * added in a future release of the API.
             *
             * Note that values may be added to this enum, implementations
             * must ensure that unknown values will not cause a crash.
             *
             * Unknown values here must result in the implementation setting the
             * Accepted Condition for the Route to `status: False`, with a
             * Reason of `UnsupportedValue`.
             */
            type: 'ReplaceFullPath' | 'ReplacePrefixMatch';
            [k: string]: unknown;
          };
          [k: string]: unknown;
        };
        /**
         * ResponseHeaderModifier defines a schema for a filter that modifies response
         * headers.
         *
         * Support: Extended
         */
        responseHeaderModifier?: {
          /**
           * Remove the given header(s) from the HTTP request before the action. The
           * value of Remove is a list of HTTP header names. Note that the header
           * names are case-insensitive (see
           * https://datatracker.ietf.org/doc/html/rfc2616#section-4.2).
           *
           * Input:
           *   GET /foo HTTP/1.1
           *   my-header1: foo
           *   my-header2: bar
           *   my-header3: baz
           *
           * Config:
           *   remove: ["my-header1", "my-header3"]
           *
           * Output:
           *   GET /foo HTTP/1.1
           *   my-header2: bar
           *
           * @maxItems 16
           */
          remove?: string[];
          /**
           * Set overwrites the request with the given header (name, value)
           * before the action.
           *
           * Input:
           *   GET /foo HTTP/1.1
           *   my-header: foo
           *
           * Config:
           *   set:
           *   - name: "my-header"
           *     value: "bar"
           *
           * Output:
           *   GET /foo HTTP/1.1
           *   my-header: bar
           *
           * @maxItems 16
           */
          set?: {
            /**
             * Name is the name of the HTTP Header to be matched. Name matching MUST be
             * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
             *
             * If multiple entries specify equivalent header names, the first entry with
             * an equivalent name MUST be considered for a match. Subsequent entries
             * with an equivalent header name MUST be ignored. Due to the
             * case-insensitivity of header names, "foo" and "Foo" are considered
             * equivalent.
             */
            name: string;
            /**
             * Value is the value of HTTP Header to be matched.
             */
            value: string;
            [k: string]: unknown;
          }[];
          /**
           * Add adds the given header(s) (name, value) to the request
           * before the action. It appends to any existing values associated
           * with the header name.
           *
           * Input:
           *   GET /foo HTTP/1.1
           *   my-header: foo
           *
           * Config:
           *   add:
           *   - name: "my-header"
           *     value: "bar,baz"
           *
           * Output:
           *   GET /foo HTTP/1.1
           *   my-header: foo,bar,baz
           *
           * @maxItems 16
           */
          add?: {
            /**
             * Name is the name of the HTTP Header to be matched. Name matching MUST be
             * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
             *
             * If multiple entries specify equivalent header names, the first entry with
             * an equivalent name MUST be considered for a match. Subsequent entries
             * with an equivalent header name MUST be ignored. Due to the
             * case-insensitivity of header names, "foo" and "Foo" are considered
             * equivalent.
             */
            name: string;
            /**
             * Value is the value of HTTP Header to be matched.
             */
            value: string;
            [k: string]: unknown;
          }[];
          [k: string]: unknown;
        };
        /**
         * Type identifies the type of filter to apply. As with other API fields,
         * types are classified into three conformance levels:
         *
         * - Core: Filter types and their corresponding configuration defined by
         *   "Support: Core" in this package, e.g. "RequestHeaderModifier". All
         *   implementations must support core filters.
         *
         * - Extended: Filter types and their corresponding configuration defined by
         *   "Support: Extended" in this package, e.g. "RequestMirror". Implementers
         *   are encouraged to support extended filters.
         *
         * - Implementation-specific: Filters that are defined and supported by
         *   specific vendors.
         *   In the future, filters showing convergence in behavior across multiple
         *   implementations will be considered for inclusion in extended or core
         *   conformance levels. Filter-specific configuration for such filters
         *   is specified using the ExtensionRef field. `Type` should be set to
         *   "ExtensionRef" for custom filters.
         *
         * Implementers are encouraged to define custom implementation types to
         * extend the core API with implementation-specific behavior.
         *
         * If a reference to a custom filter type cannot be resolved, the filter
         * MUST NOT be skipped. Instead, requests that would have been processed by
         * that filter MUST receive a HTTP error response.
         *
         * Note that values may be added to this enum, implementations
         * must ensure that unknown values will not cause a crash.
         *
         * Unknown values here must result in the implementation setting the
         * Accepted Condition for the Route to `status: False`, with a
         * Reason of `UnsupportedValue`.
         */
        type:
          | 'RequestHeaderModifier'
          | 'ResponseHeaderModifier'
          | 'RequestMirror'
          | 'RequestRedirect'
          | 'URLRewrite'
          | 'ExtensionRef';
        [k: string]: unknown;
      }[];
      /**
       * Matches define conditions used for matching the rule against incoming
       * HTTP requests. Each match is independent, i.e. this rule will be matched
       * if **any** one of the matches is satisfied.
       *
       * For example, take the following matches configuration:
       *
       * ```
       * matches:
       * - path:
       *     value: "/foo"
       *   headers:
       *   - name: "version"
       *     value: "v2"
       * - path:
       *     value: "/v2/foo"
       * ```
       *
       * For a request to match against this rule, a request must satisfy
       * EITHER of the two conditions:
       *
       * - path prefixed with `/foo` AND contains the header `version: v2`
       * - path prefix of `/v2/foo`
       *
       * See the documentation for HTTPRouteMatch on how to specify multiple
       * match conditions that should be ANDed together.
       *
       * If no matches are specified, the default is a prefix
       * path match on "/", which has the effect of matching every
       * HTTP request.
       *
       * Proxy or Load Balancer routing configuration generated from HTTPRoutes
       * MUST prioritize matches based on the following criteria, continuing on
       * ties. Across all rules specified on applicable Routes, precedence must be
       * given to the match having:
       *
       * * "Exact" path match.
       * * "Prefix" path match with largest number of characters.
       * * Method match.
       * * Largest number of header matches.
       * * Largest number of query param matches.
       *
       * Note: The precedence of RegularExpression path matches are implementation-specific.
       *
       * If ties still exist across multiple Routes, matching precedence MUST be
       * determined in order of the following criteria, continuing on ties:
       *
       * * The oldest Route based on creation timestamp.
       * * The Route appearing first in alphabetical order by
       *   "{namespace}/{name}".
       *
       * If ties still exist within an HTTPRoute, matching precedence MUST be granted
       * to the FIRST matching rule (in list order) with a match meeting the above
       * criteria.
       *
       * When no rules matching a request have been successfully attached to the
       * parent a request is coming from, a HTTP 404 status code MUST be returned.
       *
       * @maxItems 64
       */
      matches?: {
        /**
         * Headers specifies HTTP request header matchers. Multiple match values are
         * ANDed together, meaning, a request must match all the specified headers
         * to select the route.
         *
         * @maxItems 16
         */
        headers?: {
          /**
           * Value is the value of HTTP Header to be matched.
           */
          value: string;
          /**
           * Name is the name of the HTTP Header to be matched. Name matching MUST be
           * case-insensitive. (See https://tools.ietf.org/html/rfc7230#section-3.2).
           *
           * If multiple entries specify equivalent header names, only the first
           * entry with an equivalent name MUST be considered for a match. Subsequent
           * entries with an equivalent header name MUST be ignored. Due to the
           * case-insensitivity of header names, "foo" and "Foo" are considered
           * equivalent.
           *
           * When a header is repeated in an HTTP request, it is
           * implementation-specific behavior as to how this is represented.
           * Generally, proxies should follow the guidance from the RFC:
           * https://www.rfc-editor.org/rfc/rfc7230.html#section-3.2.2 regarding
           * processing a repeated header, with special handling for "Set-Cookie".
           */
          name: string;
          /**
           * Type specifies how to match against the value of the header.
           *
           * Support: Core (Exact)
           *
           * Support: Implementation-specific (RegularExpression)
           *
           * Since RegularExpression HeaderMatchType has implementation-specific
           * conformance, implementations can support POSIX, PCRE or any other dialects
           * of regular expressions. Please read the implementation's documentation to
           * determine the supported dialect.
           */
          type?: 'Exact' | 'RegularExpression';
          [k: string]: unknown;
        }[];
        /**
         * Method specifies HTTP method matcher.
         * When specified, this route will be matched only if the request has the
         * specified method.
         *
         * Support: Extended
         */
        method?: 'GET' | 'HEAD' | 'POST' | 'PUT' | 'DELETE' | 'CONNECT' | 'OPTIONS' | 'TRACE' | 'PATCH';
        /**
         * Path specifies a HTTP request path matcher. If this field is not
         * specified, a default prefix match on the "/" path is provided.
         */
        path?: {
          /**
           * Type specifies how to match against the path Value.
           *
           * Support: Core (Exact, PathPrefix)
           *
           * Support: Implementation-specific (RegularExpression)
           */
          type?: 'Exact' | 'PathPrefix' | 'RegularExpression';
          /**
           * Value of the HTTP path to match against.
           */
          value?: string;
          [k: string]: unknown;
        };
        /**
         * QueryParams specifies HTTP query parameter matchers. Multiple match
         * values are ANDed together, meaning, a request must match all the
         * specified query parameters to select the route.
         *
         * Support: Extended
         *
         * @maxItems 16
         */
        queryParams?: {
          /**
           * Value is the value of HTTP query param to be matched.
           */
          value: string;
          /**
           * Name is the name of the HTTP query param to be matched. This must be an
           * exact string match. (See
           * https://tools.ietf.org/html/rfc7230#section-2.7.3).
           *
           * If multiple entries specify equivalent query param names, only the first
           * entry with an equivalent name MUST be considered for a match. Subsequent
           * entries with an equivalent query param name MUST be ignored.
           *
           * If a query param is repeated in an HTTP request, the behavior is
           * purposely left undefined, since different data planes have different
           * capabilities. However, it is *recommended* that implementations should
           * match against the first value of the param if the data plane supports it,
           * as this behavior is expected in other load balancing contexts outside of
           * the Gateway API.
           *
           * Users SHOULD NOT route traffic based on repeated query params to guard
           * themselves against potential differences in the implementations.
           */
          name: string;
          /**
           * Type specifies how to match against the value of the query parameter.
           *
           * Support: Extended (Exact)
           *
           * Support: Implementation-specific (RegularExpression)
           *
           * Since RegularExpression QueryParamMatchType has Implementation-specific
           * conformance, implementations can support POSIX, PCRE or any other
           * dialects of regular expressions. Please read the implementation's
           * documentation to determine the supported dialect.
           */
          type?: 'Exact' | 'RegularExpression';
          [k: string]: unknown;
        }[];
        [k: string]: unknown;
      }[];
      [k: string]: unknown;
    }[];
    [k: string]: unknown;
  };
  /**
   * Status defines the current state of HTTPRoute.
   */
  status?: {
    /**
     * Parents is a list of parent resources (usually Gateways) that are
     * associated with the route, and the status of the route with respect to
     * each parent. When this route attaches to a parent, the controller that
     * manages the parent must add an entry to this list when the controller
     * first sees the route and should update the entry as appropriate when the
     * route or gateway is modified.
     *
     * Note that parent references that cannot be resolved by an implementation
     * of this API will not be added to this list. Implementations of this API
     * can only populate Route status for the Gateways/parent resources they are
     * responsible for.
     *
     * A maximum of 32 Gateways will be represented in this list. An empty list
     * means the route has not been attached to any Gateway.
     *
     * @maxItems 32
     */
    parents: {
      /**
       * Conditions describes the status of the route with respect to the Gateway.
       * Note that the route's availability is also subject to the Gateway's own
       * status conditions and listener status.
       *
       * If the Route's ParentRef specifies an existing Gateway that supports
       * Routes of this kind AND that Gateway's controller has sufficient access,
       * then that Gateway's controller MUST set the "Accepted" condition on the
       * Route, to indicate whether the route has been accepted or rejected by the
       * Gateway, and why.
       *
       * A Route MUST be considered "Accepted" if at least one of the Route's
       * rules is implemented by the Gateway.
       *
       * There are a number of cases where the "Accepted" condition may not be set
       * due to lack of controller visibility, that includes when:
       *
       * * The Route refers to a nonexistent parent.
       * * The Route is of a type that the controller does not support.
       * * The Route is in a namespace the controller does not have access to.
       *
       * @minItems 1
       * @maxItems 8
       */
      conditions: {
        /**
         * status of the condition, one of True, False, Unknown.
         */
        status: 'True' | 'False' | 'Unknown';
        /**
         * type of condition in CamelCase or in foo.example.com/CamelCase.
         */
        type: string;
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
        [k: string]: unknown;
      }[];
      /**
       * ControllerName is a domain/path string that indicates the name of the
       * controller that wrote this status. This corresponds with the
       * controllerName field on GatewayClass.
       *
       * Example: "example.net/gateway-controller".
       *
       * The format of this field is DOMAIN "/" PATH, where DOMAIN and PATH are
       * valid Kubernetes names
       * (https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names).
       *
       * Controllers MUST populate this field when writing status. Controllers should ensure that
       * entries to status populated with their ControllerName are cleaned up when they are no
       * longer necessary.
       */
      controllerName: string;
      /**
       * ParentRef corresponds with a ParentRef in the spec that this
       * RouteParentStatus struct describes the status of.
       */
      parentRef: {
        /**
         * Group is the group of the referent.
         * When unspecified, "gateway.networking.k8s.io" is inferred.
         * To set the core API group (such as for a "Service" kind referent),
         * Group must be explicitly set to "" (empty string).
         *
         * Support: Core
         */
        group?: string;
        /**
         * Kind is kind of the referent.
         *
         * There are two kinds of parent resources with "Core" support:
         *
         * * Gateway (Gateway conformance profile)
         * * Service (Mesh conformance profile, ClusterIP Services only)
         *
         * Support for other resources is Implementation-Specific.
         */
        kind?: string;
        /**
         * Name is the name of the referent.
         *
         * Support: Core
         */
        name: string;
        /**
         * Namespace is the namespace of the referent. When unspecified, this refers
         * to the local namespace of the Route.
         *
         * Note that there are specific rules for ParentRefs which cross namespace
         * boundaries. Cross-namespace references are only valid if they are explicitly
         * allowed by something in the namespace they are referring to. For example:
         * Gateway has the AllowedRoutes field, and ReferenceGrant provides a
         * generic way to enable any other kind of cross-namespace reference.
         *
         * Support: Core
         */
        namespace?: string;
        /**
         * Port is the network port this Route targets. It can be interpreted
         * differently based on the type of parent resource.
         *
         * When the parent resource is a Gateway, this targets all listeners
         * listening on the specified port that also support this kind of Route(and
         * select this Route). It's not recommended to set `Port` unless the
         * networking behaviors specified in a Route must apply to a specific port
         * as opposed to a listener(s) whose port(s) may be changed. When both Port
         * and SectionName are specified, the name and port of the selected listener
         * must match both specified values.
         *
         * Implementations MAY choose to support other parent resources.
         * Implementations supporting other types of parent resources MUST clearly
         * document how/if Port is interpreted.
         *
         * For the purpose of status, an attachment is considered successful as
         * long as the parent resource accepts it partially. For example, Gateway
         * listeners can restrict which Routes can attach to them by Route kind,
         * namespace, or hostname. If 1 of 2 Gateway listeners accept attachment
         * from the referencing Route, the Route MUST be considered successfully
         * attached. If no Gateway listeners accept attachment from this Route,
         * the Route MUST be considered detached from the Gateway.
         *
         * Support: Extended
         */
        port?: number;
        /**
         * SectionName is the name of a section within the target resource. In the
         * following resources, SectionName is interpreted as the following:
         *
         * * Gateway: Listener name. When both Port (experimental) and SectionName
         * are specified, the name and port of the selected listener must match
         * both specified values.
         * * Service: Port name. When both Port (experimental) and SectionName
         * are specified, the name and port of the selected listener must match
         * both specified values.
         *
         * Implementations MAY choose to support attaching Routes to other resources.
         * If that is the case, they MUST clearly document how SectionName is
         * interpreted.
         *
         * When unspecified (empty string), this will reference the entire resource.
         * For the purpose of status, an attachment is considered successful if at
         * least one section in the parent resource accepts it. For example, Gateway
         * listeners can restrict which Routes can attach to them by Route kind,
         * namespace, or hostname. If 1 of 2 Gateway listeners accept attachment from
         * the referencing Route, the Route MUST be considered successfully
         * attached. If no Gateway listeners accept attachment from this Route, the
         * Route MUST be considered detached from the Gateway.
         *
         * Support: Core
         */
        sectionName?: string;
        [k: string]: unknown;
      };
      [k: string]: unknown;
    }[];
    [k: string]: unknown;
  };
  [k: string]: unknown;
}

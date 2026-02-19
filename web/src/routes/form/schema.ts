const schema = {
    description:
        'Workspace is the Schema for the workspaces API.\nA Workspace represents a logical isolation unit (Namespace) with associated policies, quotas, and user access.',
    type: 'object',
    required: ['spec'],
    properties: {
        apiVersion: {
            description:
                'APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources',
            type: 'string'
        },
        kind: {
            description:
                'Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds',
            type: 'string'
        },
        metadata: {
            description:
                'ObjectMeta is metadata that all persisted resources must have, which includes all objects users must create.',
            type: 'object',
            properties: {
                deletionGracePeriodSeconds: {
                    description:
                        'Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.',
                    type: 'integer',
                    format: 'int64'
                },
                deletionTimestamp: {
                    description:
                        'Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.',
                    type: 'string',
                    format: 'date-time'
                },
                generateName: {
                    type: 'string',
                    description:
                        'GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.\n\nIf this field is specified and the generated name exists, the server will return a 409.\n\nApplied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency'
                },
                managedFields: {
                    description:
                        "ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field. A workflow can be the user's name, a controller's name, or the name of a specific apply path like \"ci-cd\". The set of fields is always in the version that the workflow used when modifying the object.",
                    type: 'array',
                    items: {
                        properties: {
                            apiVersion: {
                                type: 'string',
                                description:
                                    'APIVersion defines the version of this resource that this field set applies to. The format is "group/version" just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.'
                            },
                            fieldsType: {
                                description:
                                    'FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: "FieldsV1"',
                                type: 'string'
                            },
                            fieldsV1: {
                                description:
                                    "FieldsV1 stores a set of fields in a data structure like a Trie, in JSON format.\n\nEach key is either a '.' representing the field itself, and will always map to an empty set, or a string representing a sub-field or item. The string will follow one of these four formats: 'f:<name>', where <name> is the name of a field in a struct, or key in a map 'v:<value>', where <value> is the exact json formatted value of a list item 'i:<index>', where <index> is position of a item in a list 'k:<keys>', where <keys> is a map of  a list item's key fields to their unique values If a key maps to an empty Fields value, the field that key represents is part of the set.\n\nThe exact format is defined in sigs.k8s.io/structured-merge-diff",
                                type: 'object'
                            },
                            manager: {
                                description: 'Manager is an identifier of the workflow managing these fields.',
                                type: 'string'
                            },
                            operation: {
                                description:
                                    "Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are 'Apply' and 'Update'.",
                                type: 'string'
                            },
                            subresource: {
                                type: 'string',
                                description:
                                    'Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.'
                            },
                            time: {
                                description:
                                    'Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.',
                                type: 'string',
                                format: 'date-time'
                            }
                        },
                        description:
                            'ManagedFieldsEntry is a workflow-id, a FieldSet and the group version of the resource that the fieldset applies to.',
                        type: 'object'
                    },
                    'x-kubernetes-list-type': 'atomic'
                },
                creationTimestamp: {
                    description:
                        'Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.',
                    type: 'string',
                    format: 'date-time'
                },
                labels: {
                    description:
                        'Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels',
                    type: 'object',
                    additionalProperties: {
                        type: 'string',
                        default: ''
                    }
                },
                ownerReferences: {
                    'x-kubernetes-list-map-keys': ['uid'],
                    'x-kubernetes-list-type': 'map',
                    'x-kubernetes-patch-merge-key': 'uid',
                    'x-kubernetes-patch-strategy': 'merge',
                    description:
                        'List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller.',
                    type: 'array',
                    items: {
                        description:
                            'OwnerReference contains enough information to let you identify an owning object. An owning object must be in the same namespace as the dependent, or be cluster-scoped, so there is no namespace field.',
                        type: 'object',
                        required: ['apiVersion', 'kind', 'name', 'uid'],
                        properties: {
                            apiVersion: {
                                description: 'API version of the referent.',
                                type: 'string',
                                default: ''
                            },
                            blockOwnerDeletion: {
                                description:
                                    'If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. See https://kubernetes.io/docs/concepts/architecture/garbage-collection/#foreground-deletion for how the garbage collector interacts with this field and enforces the foreground deletion. Defaults to false. To set this field, a user needs "delete" permission of the owner, otherwise 422 (Unprocessable Entity) will be returned.',
                                type: 'boolean'
                            },
                            controller: {
                                description: 'If true, this reference points to the managing controller.',
                                type: 'boolean'
                            },
                            kind: {
                                description:
                                    'Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds',
                                type: 'string',
                                default: ''
                            },
                            name: {
                                description:
                                    'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names',
                                type: 'string',
                                default: ''
                            },
                            uid: {
                                description:
                                    'UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids',
                                type: 'string',
                                default: ''
                            }
                        },
                        'x-kubernetes-map-type': 'atomic'
                    }
                },
                resourceVersion: {
                    description:
                        'An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.\n\nPopulated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency',
                    type: 'string'
                },
                selfLink: {
                    description:
                        'Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.',
                    type: 'string'
                },
                annotations: {
                    description:
                        'Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations',
                    type: 'object',
                    additionalProperties: {
                        type: 'string',
                        default: ''
                    }
                },
                finalizers: {
                    items: {
                        type: 'string',
                        default: ''
                    },
                    'x-kubernetes-list-type': 'set',
                    'x-kubernetes-patch-strategy': 'merge',
                    description:
                        'Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order.  Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list.',
                    type: 'array'
                },
                generation: {
                    description:
                        'A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.',
                    type: 'integer',
                    format: 'int64'
                },
                name: {
                    description:
                        'Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names',
                    type: 'string'
                },
                namespace: {
                    description:
                        'Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.\n\nMust be a DNS_LABEL. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces',
                    type: 'string'
                },
                uid: {
                    description:
                        'UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.\n\nPopulated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids',
                    type: 'string'
                }
            }
        },
        spec: {
            description: 'Spec defines the desired behavior of the Workspace.',
            type: 'object',
            required: ['namespace', 'users'],
            properties: {
                limitRange: {
                    description:
                        'LimitRange describes the default resource limits and requests for pods in the workspace.',
                    type: 'object',
                    required: ['limits'],
                    properties: {
                        limits: {
                            'x-kubernetes-list-type': 'atomic',
                            description: 'Limits is the list of LimitRangeItem objects that are enforced.',
                            type: 'array',
                            items: {
                                description:
                                    'LimitRangeItem defines a min/max usage limit for any resource that matches on kind.',
                                type: 'object',
                                required: ['type'],
                                properties: {
                                    default: {
                                        description:
                                            'Default resource requirement limit value by resource name if resource limit is omitted.',
                                        type: 'object',
                                        additionalProperties: {
                                            pattern:
                                                '^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$',
                                            anyOf: [
                                                {
                                                    type: 'integer'
                                                },
                                                {
                                                    type: 'string'
                                                }
                                            ],
                                            'x-kubernetes-int-or-string': true
                                        }
                                    },
                                    defaultRequest: {
                                        description:
                                            'DefaultRequest is the default resource requirement request value by resource name if resource request is omitted.',
                                        type: 'object',
                                        additionalProperties: {
                                            anyOf: [
                                                {
                                                    type: 'integer'
                                                },
                                                {
                                                    type: 'string'
                                                }
                                            ],
                                            'x-kubernetes-int-or-string': true,
                                            pattern:
                                                '^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$'
                                        }
                                    },
                                    max: {
                                        description: 'Max usage constraints on this kind by resource name.',
                                        type: 'object',
                                        additionalProperties: {
                                            'x-kubernetes-int-or-string': true,
                                            pattern:
                                                '^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$',
                                            anyOf: [
                                                {
                                                    type: 'integer'
                                                },
                                                {
                                                    type: 'string'
                                                }
                                            ]
                                        }
                                    },
                                    maxLimitRequestRatio: {
                                        type: 'object',
                                        additionalProperties: {
                                            pattern:
                                                '^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$',
                                            anyOf: [
                                                {
                                                    type: 'integer'
                                                },
                                                {
                                                    type: 'string'
                                                }
                                            ],
                                            'x-kubernetes-int-or-string': true
                                        },
                                        description:
                                            'MaxLimitRequestRatio if specified, the named resource must have a request and limit that are both non-zero where limit divided by request is less than or equal to the enumerated value; this represents the max burst for the named resource.'
                                    },
                                    min: {
                                        type: 'object',
                                        additionalProperties: {
                                            pattern:
                                                '^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$',
                                            anyOf: [
                                                {
                                                    type: 'integer'
                                                },
                                                {
                                                    type: 'string'
                                                }
                                            ],
                                            'x-kubernetes-int-or-string': true
                                        },
                                        description: 'Min usage constraints on this kind by resource name.'
                                    },
                                    type: {
                                        description: 'Type of resource that this limit applies to.',
                                        type: 'string'
                                    }
                                }
                            }
                        }
                    }
                },
                namespace: {
                    minLength: 1,
                    pattern: '^([a-z0-9]([-a-z0-9]*[a-z0-9])?)$',
                    'x-kubernetes-validations': [
                        {
                            message: 'namespace is immutable',
                            rule: 'self == oldSelf'
                        }
                    ],
                    description:
                        'Namespace is the name of the Kubernetes Namespace to be created for this workspace.\nIt must be unique across all Workspaces.',
                    type: 'string',
                    maxLength: 63
                },
                networkIsolation: {
                    description: 'NetworkIsolation defines the ingress traffic rules for the workspace.',
                    type: 'object',
                    properties: {
                        allowedNamespaces: {
                            description:
                                "AllowedNamespaces specifies a list of external namespaces permitted to access this workspace\nwhen isolation is enabled. Essential system namespaces (e.g., 'istio-system', 'monitoring')\nshould be included here if required.",
                            type: 'array',
                            items: {
                                type: 'string'
                            },
                            'x-kubernetes-list-type': 'set'
                        },
                        enabled: {
                            description:
                                'Enabled toggles the enforcement of network isolation.\nIf true, default deny-all ingress rules are applied except for allowed namespaces.',
                            type: 'boolean'
                        }
                    }
                },
                resourceQuota: {
                    description:
                        'ResourceQuota describes the compute resource constraints (CPU, Memory, etc.) applied to the underlying namespace.',
                    type: 'object',
                    properties: {
                        scopeSelector: {
                            description:
                                'scopeSelector is also a collection of filters like scopes that must match each object tracked by a quota\nbut expressed using ScopeSelectorOperator in combination with possible values.\nFor a resource to match, both scopes AND scopeSelector (if specified in spec), must be matched.',
                            type: 'object',
                            properties: {
                                matchExpressions: {
                                    type: 'array',
                                    items: {
                                        description:
                                            'A scoped-resource selector requirement is a selector that contains values, a scope name, and an operator\nthat relates the scope name and values.',
                                        type: 'object',
                                        required: ['operator', 'scopeName'],
                                        properties: {
                                            operator: {
                                                description:
                                                    "Represents a scope's relationship to a set of values.\nValid operators are In, NotIn, Exists, DoesNotExist.",
                                                type: 'string'
                                            },
                                            scopeName: {
                                                description: 'The name of the scope that the selector applies to.',
                                                type: 'string'
                                            },
                                            values: {
                                                'x-kubernetes-list-type': 'atomic',
                                                description:
                                                    'An array of string values. If the operator is In or NotIn,\nthe values array must be non-empty. If the operator is Exists or DoesNotExist,\nthe values array must be empty.\nThis array is replaced during a strategic merge patch.',
                                                type: 'array',
                                                items: {
                                                    type: 'string'
                                                }
                                            }
                                        }
                                    },
                                    'x-kubernetes-list-type': 'atomic',
                                    description: 'A list of scope selector requirements by scope of the resources.'
                                }
                            },
                            'x-kubernetes-map-type': 'atomic'
                        },
                        scopes: {
                            type: 'array',
                            items: {
                                description:
                                    'A ResourceQuotaScope defines a filter that must match each object tracked by a quota',
                                type: 'string'
                            },
                            'x-kubernetes-list-type': 'atomic',
                            description:
                                'A collection of filters that must match each object tracked by a quota.\nIf not specified, the quota matches all objects.'
                        },
                        hard: {
                            description:
                                'hard is the set of desired hard limits for each named resource.\nMore info: https://kubernetes.io/docs/concepts/policy/resource-quotas/',
                            type: 'object',
                            additionalProperties: {
                                pattern:
                                    '^(\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))(([KMGTPE]i)|[numkMGTPE]|([eE](\\+|-)?(([0-9]+(\\.[0-9]*)?)|(\\.[0-9]+))))?$',
                                anyOf: [
                                    {
                                        type: 'integer'
                                    },
                                    {
                                        type: 'string'
                                    }
                                ],
                                'x-kubernetes-int-or-string': true
                            }
                        }
                    }
                },
                users: {
                    description: 'Users is the list of users granted access to this workspace.',
                    type: 'array',
                    minItems: 1,
                    items: {
                        type: 'object',
                        required: ['role', 'subject'],
                        properties: {
                            name: {
                                description: 'Name is the human-readable display name of the user.',
                                type: 'string'
                            },
                            role: {
                                description: 'Role defines the authorization level (Admin, Edit, View).',
                                type: 'string',
                                enum: ['admin', 'edit', 'view']
                            },
                            subject: {
                                description:
                                    'Subject is the unique identifier of the user (e.g., OIDC subject or username).\nThis identifier maps directly to the Kubernetes RBAC Subject.',
                                type: 'string',
                                maxLength: 63,
                                minLength: 1,
                                pattern: '^([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$'
                            }
                        },
                        description: 'WorkspaceUser defines a single user entity associated with a workspace.'
                    },
                    'x-kubernetes-list-map-keys': ['subject'],
                    'x-kubernetes-list-type': 'map'
                }
            }
        },
        status: {
            description: 'Status represents the current information about the Workspace.',
            type: 'object',
            properties: {
                networkPolicyRef: {
                    description:
                        'NetworkPolicyRef is a reference to the corev1.NetworkPolicy enforcing network isolation.',
                    type: 'object',
                    properties: {
                        apiVersion: {
                            description: 'API version of the referent.',
                            type: 'string'
                        },
                        fieldPath: {
                            description:
                                'If referring to a piece of an object instead of an entire object, this string\nshould contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].\nFor example, if the object reference is to a container within a pod, this would take on a value like:\n"spec.containers{name}" (where "name" refers to the name of the container that triggered\nthe event) or if no container name is specified "spec.containers[2]" (container with\nindex 2 in this pod). This syntax is chosen only to have some well-defined way of\nreferencing a part of an object.',
                            type: 'string'
                        },
                        kind: {
                            description:
                                'Kind of the referent.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds',
                            type: 'string'
                        },
                        name: {
                            description:
                                'Name of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names',
                            type: 'string'
                        },
                        namespace: {
                            description:
                                'Namespace of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/',
                            type: 'string'
                        },
                        resourceVersion: {
                            description:
                                'Specific resourceVersion to which this reference is made, if any.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency',
                            type: 'string'
                        },
                        uid: {
                            description:
                                'UID of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids',
                            type: 'string'
                        }
                    },
                    'x-kubernetes-map-type': 'atomic'
                },
                peerAuthenticationRef: {
                    description:
                        'PeerAuthenticationRef is a reference to the Istio PeerAuthentication resource for mTLS settings.',
                    type: 'object',
                    properties: {
                        resourceVersion: {
                            description:
                                'Specific resourceVersion to which this reference is made, if any.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency',
                            type: 'string'
                        },
                        uid: {
                            description:
                                'UID of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids',
                            type: 'string'
                        },
                        apiVersion: {
                            description: 'API version of the referent.',
                            type: 'string'
                        },
                        fieldPath: {
                            description:
                                'If referring to a piece of an object instead of an entire object, this string\nshould contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].\nFor example, if the object reference is to a container within a pod, this would take on a value like:\n"spec.containers{name}" (where "name" refers to the name of the container that triggered\nthe event) or if no container name is specified "spec.containers[2]" (container with\nindex 2 in this pod). This syntax is chosen only to have some well-defined way of\nreferencing a part of an object.',
                            type: 'string'
                        },
                        kind: {
                            description:
                                'Kind of the referent.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds',
                            type: 'string'
                        },
                        name: {
                            description:
                                'Name of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names',
                            type: 'string'
                        },
                        namespace: {
                            type: 'string',
                            description:
                                'Namespace of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/'
                        }
                    },
                    'x-kubernetes-map-type': 'atomic'
                },
                resourceQuotaRef: {
                    type: 'object',
                    properties: {
                        namespace: {
                            description:
                                'Namespace of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/',
                            type: 'string'
                        },
                        resourceVersion: {
                            description:
                                'Specific resourceVersion to which this reference is made, if any.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency',
                            type: 'string'
                        },
                        uid: {
                            description:
                                'UID of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids',
                            type: 'string'
                        },
                        apiVersion: {
                            description: 'API version of the referent.',
                            type: 'string'
                        },
                        fieldPath: {
                            description:
                                'If referring to a piece of an object instead of an entire object, this string\nshould contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].\nFor example, if the object reference is to a container within a pod, this would take on a value like:\n"spec.containers{name}" (where "name" refers to the name of the container that triggered\nthe event) or if no container name is specified "spec.containers[2]" (container with\nindex 2 in this pod). This syntax is chosen only to have some well-defined way of\nreferencing a part of an object.',
                            type: 'string'
                        },
                        kind: {
                            description:
                                'Kind of the referent.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds',
                            type: 'string'
                        },
                        name: {
                            type: 'string',
                            description:
                                'Name of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names'
                        }
                    },
                    'x-kubernetes-map-type': 'atomic',
                    description:
                        'ResourceQuotaRef is a reference to the corev1.ResourceQuota managed by this Workspace.'
                },
                roleBindingRefs: {
                    description:
                        'RoleBindingRefs contains references to all RBAC RoleBindings created for the workspace users.',
                    type: 'array',
                    items: {
                        description:
                            'ObjectReference contains enough information to let you inspect or modify the referred object.',
                        type: 'object',
                        properties: {
                            apiVersion: {
                                description: 'API version of the referent.',
                                type: 'string'
                            },
                            fieldPath: {
                                description:
                                    'If referring to a piece of an object instead of an entire object, this string\nshould contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].\nFor example, if the object reference is to a container within a pod, this would take on a value like:\n"spec.containers{name}" (where "name" refers to the name of the container that triggered\nthe event) or if no container name is specified "spec.containers[2]" (container with\nindex 2 in this pod). This syntax is chosen only to have some well-defined way of\nreferencing a part of an object.',
                                type: 'string'
                            },
                            kind: {
                                type: 'string',
                                description:
                                    'Kind of the referent.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                            },
                            name: {
                                description:
                                    'Name of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names',
                                type: 'string'
                            },
                            namespace: {
                                description:
                                    'Namespace of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/',
                                type: 'string'
                            },
                            resourceVersion: {
                                description:
                                    'Specific resourceVersion to which this reference is made, if any.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency',
                                type: 'string'
                            },
                            uid: {
                                description:
                                    'UID of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids',
                                type: 'string'
                            }
                        },
                        'x-kubernetes-map-type': 'atomic'
                    },
                    'x-kubernetes-list-type': 'atomic'
                },
                authorizationPolicyRef: {
                    description:
                        'AuthorizationPolicyRef is a reference to the Istio AuthorizationPolicy enforcing network isolation.',
                    type: 'object',
                    properties: {
                        apiVersion: {
                            description: 'API version of the referent.',
                            type: 'string'
                        },
                        fieldPath: {
                            description:
                                'If referring to a piece of an object instead of an entire object, this string\nshould contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].\nFor example, if the object reference is to a container within a pod, this would take on a value like:\n"spec.containers{name}" (where "name" refers to the name of the container that triggered\nthe event) or if no container name is specified "spec.containers[2]" (container with\nindex 2 in this pod). This syntax is chosen only to have some well-defined way of\nreferencing a part of an object.',
                            type: 'string'
                        },
                        kind: {
                            description:
                                'Kind of the referent.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds',
                            type: 'string'
                        },
                        name: {
                            description:
                                'Name of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names',
                            type: 'string'
                        },
                        namespace: {
                            description:
                                'Namespace of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/',
                            type: 'string'
                        },
                        resourceVersion: {
                            description:
                                'Specific resourceVersion to which this reference is made, if any.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency',
                            type: 'string'
                        },
                        uid: {
                            description:
                                'UID of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids',
                            type: 'string'
                        }
                    },
                    'x-kubernetes-map-type': 'atomic'
                },
                conditions: {
                    'x-kubernetes-list-map-keys': ['type'],
                    'x-kubernetes-list-type': 'map',
                    description:
                        'Conditions store the status conditions of the Workspace (e.g., Ready, Failed).',
                    type: 'array',
                    items: {
                        description:
                            'Condition contains details for one aspect of the current state of this API Resource.',
                        type: 'object',
                        required: ['lastTransitionTime', 'message', 'reason', 'status', 'type'],
                        properties: {
                            lastTransitionTime: {
                                format: 'date-time',
                                description:
                                    'lastTransitionTime is the last time the condition transitioned from one status to another.\nThis should be when the underlying condition changed.  If that is not known, then using the time when the API field changed is acceptable.',
                                type: 'string'
                            },
                            message: {
                                description:
                                    'message is a human readable message indicating details about the transition.\nThis may be an empty string.',
                                type: 'string',
                                maxLength: 32768
                            },
                            observedGeneration: {
                                description:
                                    'observedGeneration represents the .metadata.generation that the condition was set based upon.\nFor instance, if .metadata.generation is currently 12, but the .status.conditions[x].observedGeneration is 9, the condition is out of date\nwith respect to the current state of the instance.',
                                type: 'integer',
                                format: 'int64',
                                minimum: 0
                            },
                            reason: {
                                description:
                                    "reason contains a programmatic identifier indicating the reason for the condition's last transition.\nProducers of specific condition types may define expected values and meanings for this field,\nand whether the values are considered a guaranteed API.\nThe value should be a CamelCase string.\nThis field may not be empty.",
                                type: 'string',
                                maxLength: 1024,
                                minLength: 1,
                                pattern: '^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$'
                            },
                            status: {
                                description: 'status of the condition, one of True, False, Unknown.',
                                type: 'string',
                                enum: ['True', 'False', 'Unknown']
                            },
                            type: {
                                type: 'string',
                                maxLength: 316,
                                pattern:
                                    '^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$',
                                description: 'type of condition in CamelCase or in foo.example.com/CamelCase.'
                            }
                        }
                    }
                },
                limitRangeRef: {
                    description:
                        'LimitRangeRef is a reference to the corev1.LimitRange managed by this Workspace.',
                    type: 'object',
                    properties: {
                        namespace: {
                            description:
                                'Namespace of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/',
                            type: 'string'
                        },
                        resourceVersion: {
                            description:
                                'Specific resourceVersion to which this reference is made, if any.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency',
                            type: 'string'
                        },
                        uid: {
                            description:
                                'UID of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids',
                            type: 'string'
                        },
                        apiVersion: {
                            description: 'API version of the referent.',
                            type: 'string'
                        },
                        fieldPath: {
                            description:
                                'If referring to a piece of an object instead of an entire object, this string\nshould contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].\nFor example, if the object reference is to a container within a pod, this would take on a value like:\n"spec.containers{name}" (where "name" refers to the name of the container that triggered\nthe event) or if no container name is specified "spec.containers[2]" (container with\nindex 2 in this pod). This syntax is chosen only to have some well-defined way of\nreferencing a part of an object.',
                            type: 'string'
                        },
                        kind: {
                            description:
                                'Kind of the referent.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds',
                            type: 'string'
                        },
                        name: {
                            description:
                                'Name of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names',
                            type: 'string'
                        }
                    },
                    'x-kubernetes-map-type': 'atomic'
                },
                namespaceRef: {
                    description:
                        'NamespaceRef is a reference to the corev1.Namespace managed by this Workspace.',
                    type: 'object',
                    properties: {
                        namespace: {
                            description:
                                'Namespace of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/',
                            type: 'string'
                        },
                        resourceVersion: {
                            type: 'string',
                            description:
                                'Specific resourceVersion to which this reference is made, if any.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency'
                        },
                        uid: {
                            description:
                                'UID of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids',
                            type: 'string'
                        },
                        apiVersion: {
                            description: 'API version of the referent.',
                            type: 'string'
                        },
                        fieldPath: {
                            description:
                                'If referring to a piece of an object instead of an entire object, this string\nshould contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].\nFor example, if the object reference is to a container within a pod, this would take on a value like:\n"spec.containers{name}" (where "name" refers to the name of the container that triggered\nthe event) or if no container name is specified "spec.containers[2]" (container with\nindex 2 in this pod). This syntax is chosen only to have some well-defined way of\nreferencing a part of an object.',
                            type: 'string'
                        },
                        kind: {
                            type: 'string',
                            description:
                                'Kind of the referent.\nMore info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                        },
                        name: {
                            description:
                                'Name of the referent.\nMore info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names',
                            type: 'string'
                        }
                    },
                    'x-kubernetes-map-type': 'atomic'
                }
            }
        }
    },
    'x-kubernetes-group-version-kind': [
        {
            version: 'v1alpha1',
            group: 'tenant.otterscale.io',
            kind: 'Workspace'
        }
    ]
};

export { schema }
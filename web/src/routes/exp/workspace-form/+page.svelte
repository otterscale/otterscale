<script lang="ts">
	import {
		type GroupedFields,
		type K8sOpenAPISchema,
		MultiStepSchemaForm
	} from '$lib/components/custom/schema-form';

	import workspaceSchema from './workspace_api.json';

	// Grouped fields for multi-step form
	const groupedFields: GroupedFields = {
		// Step 1: Basic Metadata
		'Basic Metadata': {
			'metadata.name': { title: 'Workspace Name' },
			'metadata.labels.app': { title: 'App Label' }
		},
		// Step 2: Namespace & Users
		'Namespace & Users': {
			'spec.namespace': { title: 'Namespace', showDescription: true },
			'spec.users': { title: 'Users' },
			'spec.users.subject': { title: 'User Subject' },
			'spec.users.name': { title: 'Display Name' },
			'spec.users.role': { title: 'Role' }
		},
		// Step 3: Resource Quota
		'Resource Quota': {
			'spec.resourceQuota.hard.requests.cpu': { title: 'Requests CPU' },
			'spec.resourceQuota.hard.requests.memory': { title: 'Requests Memory' },
			'spec.resourceQuota.hard.limits.cpu': { title: 'Limits CPU' },
			'spec.resourceQuota.hard.limits.memory': { title: 'Limits Memory' }
		},
		// Step 4: Limit Range
		'Limit Range': {
			'spec.limitRange.limits.type': { title: 'Type' },
			'spec.limitRange.limits.default.cpu': { title: 'Default CPU Limit' },
			'spec.limitRange.limits.default.memory': { title: 'Default Memory Limit' },
			'spec.limitRange.limits.defaultRequest.cpu': { title: 'Default CPU Request' },
			'spec.limitRange.limits.defaultRequest.memory': { title: 'Default Memory Request' }
		},
		// Step 5: Network Isolation
		'Network Isolation': {
			'spec.networkIsolation': { title: 'Network Isolation' },
			'spec.networkIsolation.enabled': {
				title: 'Enable Network Isolation',
				uiSchema: {
					'ui:components': {
						checkboxWidget: 'switchWidget'
					}
				}
			},
			'spec.networkIsolation.allowedNamespaces': { title: 'Allowed Namespaces' }
		}
	};

	function handleMultiStepSubmit(data: Record<string, unknown>) {
		console.log('Workspace form submitted:', data);
		alert('Workspace created! Check console for data.');
	}
</script>

<div class="container mx-auto py-10">
	<!-- Multi-Step Form -->
	<div class="mb-12">
		<div class="rounded border bg-card p-4 text-card-foreground">
			<MultiStepSchemaForm
				apiSchema={workspaceSchema as K8sOpenAPISchema}
				fields={groupedFields}
				title="Create Workspace"
				onSubmit={handleMultiStepSubmit}
			/>
		</div>
	</div>
</div>

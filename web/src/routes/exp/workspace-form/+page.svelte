<script lang="ts">
	import {
		type GroupedFields,
		type K8sOpenAPISchema,
		MultiStepSchemaForm
	} from '$lib/components/custom/schema-form';

	import workspaceSchema from './workspace_api.json';

	// Default values for Resource Quota and Limit Range
	const initialData = {
		spec: {
			resourceQuota: {
				hard: {
					'requests.cpu': '2',
					'requests.memory': '4Gi',
					'limits.cpu': '4',
					'limits.memory': '8Gi'
				}
			},
			limitRange: {
				limits: [
					{
						type: 'Container',
						default: {
							cpu: '500m',
							memory: '512Mi'
						},
						defaultRequest: {
							cpu: '100m',
							memory: '128Mi'
						}
					}
				]
			}
		}
	};

	// Grouped fields for multi-step form (3 pages)
	const groupedFields: GroupedFields = {
		// Step 1: Workspace & Users
		'Workspace & Users': {
			'metadata.name': { title: 'Workspace Name' },
			'spec.namespace': { title: 'Namespace', showDescription: true },
			'spec.users': { title: 'Users' },
			'spec.users.subject': { title: 'User Subject' },
			'spec.users.name': { title: 'Display Name' },
			'spec.users.role': { title: 'Role' }
		},
		// Step 2: Network Isolation
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
		},
		// Step 3: Default Resource Settings (read-only with preset values)
		'Default Resource Settings': {
			'spec.resourceQuota.hard.requests.cpu': { title: 'Requests CPU', disabled: true },
			'spec.resourceQuota.hard.requests.memory': { title: 'Requests Memory', disabled: true },
			'spec.resourceQuota.hard.limits.cpu': { title: 'Limits CPU', disabled: true },
			'spec.resourceQuota.hard.limits.memory': { title: 'Limits Memory', disabled: true },
			'spec.limitRange.limits': {
				title: 'Limit Range',
				uiSchema: {
					'ui:options': {
						addable: false,
						removable: false,
						orderable: false
					}
				}
			},
			'spec.limitRange.limits.type': { title: 'Type', disabled: true },
			'spec.limitRange.limits.default.cpu': { title: 'Default CPU Limit', disabled: true },
			'spec.limitRange.limits.default.memory': { title: 'Default Memory Limit', disabled: true },
			'spec.limitRange.limits.defaultRequest.cpu': { title: 'Default CPU Request', disabled: true },
			'spec.limitRange.limits.defaultRequest.memory': {
				title: 'Default Memory Request',
				disabled: true
			}
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
				{initialData}
				title="Create Workspace"
				onSubmit={handleMultiStepSubmit}
			/>
		</div>
	</div>
</div>

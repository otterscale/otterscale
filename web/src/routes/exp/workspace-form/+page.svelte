<script lang="ts">
	import { toJson } from '@bufbuild/protobuf';
	import { StructSchema } from '@bufbuild/protobuf/wkt';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import {
		type GroupedFields,
		type K8sOpenAPISchema,
		MultiStepSchemaForm,
		UserSelectWidget
	} from '$lib/components/custom/schema-form';

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	let apiSchema: K8sOpenAPISchema | undefined = $state();

	// Default values for Resource Quota and Limit Range
	const initialData = {
		spec: {
			resourceQuota: {
				hard: {
					'requests.cpu': '2',
					'requests.memory': '4Gi',
					'requests.nvidia.com/gpu': '0',
					'limits.cpu': '4',
					'limits.memory': '8Gi',
					'limits.nvidia.com/gpu': '0'
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
			'spec.users': {
				title: 'Users',
				uiSchema: {
					items: {
						'ui:components': {
							objectField: UserSelectWidget
						}
					}
				}
			}
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
			'spec.resourceQuota.hard.requests.nvidia.com/gpu': { title: 'Requests GPU', disabled: true },
			'spec.resourceQuota.hard.limits.cpu': { title: 'Limits CPU', disabled: true },
			'spec.resourceQuota.hard.limits.memory': { title: 'Limits Memory', disabled: true },
			'spec.resourceQuota.hard.limits.nvidia.com/gpu': { title: 'Limits GPU', disabled: true },
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

	async function handleMultiStepSubmit(data: Record<string, unknown>) {
		// Construct the full resource object
		const resourceObject = {
			apiVersion: 'tenant.otterscale.io/v1alpha1',
			kind: 'Workspace',
			...data
		};

		const name = (data.metadata as { name: string })?.name;

		toast.promise(
			async () => {
				const manifest = new TextEncoder().encode(JSON.stringify(resourceObject));

				await resourceClient.create({
					cluster: 'aaa',
					group: 'tenant.otterscale.io',
					version: 'v1alpha1',
					resource: 'workspaces',
					manifest
				});
			},
			{
				loading: `Creating workspace ${name}...`,
				success: `Successfully created workspace ${name}`,
				error: (err) => {
					console.error('Failed to create workspace:', err);
					return `Failed to create workspace: ${(err as ConnectError).message}`;
				}
			}
		);
	}

	onMount(async () => {
		try {
			const res = await resourceClient.schema({
				cluster: 'aaa',
				group: 'tenant.otterscale.io',
				version: 'v1alpha1',
				kind: 'Workspace'
			});
			// Convert Protobuf Struct to plain JSON object
			apiSchema = toJson(StructSchema, res) as K8sOpenAPISchema;
		} catch (err) {
			console.error('Failed to fetch workspace schema:', err);
			toast.error(`Failed to fetch workspace schema: ${(err as ConnectError).message}`);
		}
	});

</script>

<div class="container mx-auto py-10">
	<!-- Multi-Step Form -->
	<div class="mb-12">
		<div class="rounded border bg-card p-4 text-card-foreground">
			{#if apiSchema}
				<MultiStepSchemaForm
					{apiSchema}
					fields={groupedFields}
					{initialData}
					title="Create Workspace"
					onSubmit={handleMultiStepSubmit}
				/>
			{:else}
				<div class="flex h-32 items-center justify-center">
					<p class="text-muted-foreground">Loading schema...</p>
				</div>
			{/if}
		</div>
	</div>
</div>


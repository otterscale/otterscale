<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import {
		type GroupedFields,
		type K8sOpenAPISchema,
		MultiStepSchemaForm,
		UserSelectWidget
	} from '$lib/components/custom/schema-form';

	let {
		name,
		schema,
		object,
		onsuccess
	}: {
		name: string;
		schema: K8sOpenAPISchema;
		object: Record<string, unknown>;
		onsuccess?: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	const cluster = $derived(page.params.scope ?? '');

	let isSubmitting = $state(false);

	// Grouped fields for multi-step form (3 pages)
	const groupedFields: GroupedFields = {
		// Step 1: Workspace & Users
		'Workspace & Users': {
			'metadata.name': { title: 'Workspace Name', disabled: true }, // Name is immutable
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
			'spec.resourceQuota.hard.requests.cpu': { title: 'Requests CPU' },
			'spec.resourceQuota.hard.requests.memory': { title: 'Requests Memory' },
			'spec.resourceQuota.hard.requests.otterscale.com/vgpu': {
				title: 'Requests GPU',
				disabled: true
			},
			'spec.resourceQuota.hard.requests.otterscale.com/vgpumem': {
				title: 'Requests GPU Memory',
				disabled: true
			},
			'spec.resourceQuota.hard.requests.otterscale.com/vgpumem-percentage': {
				title: 'Requests GPU Memory Percentage',
				disabled: true
			}
		}
	};

	function transformFormData(data: Record<string, unknown>) {
		const spec = data.spec as Record<string, any>;

		// Handle Resource Quota Logic: limits align with requests, strict defaults
		if (spec?.resourceQuota?.hard) {
			const hard = spec.resourceQuota.hard;
			// Sync limits with requests
			if (hard['requests.cpu']) hard['limits.cpu'] = hard['requests.cpu'];
			if (hard['requests.memory']) hard['limits.memory'] = hard['requests.memory'];
		}

		return data;
	}

	async function handleMultiStepSubmit(data: Record<string, unknown>) {
		if (isSubmitting) return;
		isSubmitting = true;

		// Construct the full resource object - metadata.name should already be in data or we ensure it
		const resourceObject: Record<string, any> = {
			apiVersion: 'tenant.otterscale.io/v1alpha1',
			kind: 'Workspace',
			...data
		};

		// Ensure name is correct (though it's disabled in form)
		if (!resourceObject.metadata) resourceObject.metadata = {};
		(resourceObject.metadata as any).name = name;

		toast.promise(
			async () => {
				const manifest = new TextEncoder().encode(JSON.stringify(resourceObject));

				await resourceClient.apply({
					cluster,
					group: 'tenant.otterscale.io',
					version: 'v1alpha1',
					resource: 'workspaces',
					manifest,
					fieldManager: 'otterscale-web-ui',
					force: true
				});
			},
			{
				loading: `Updating workspace ${name}...`,
				success: () => {
					isSubmitting = false;
					onsuccess?.();
					return `Successfully updated workspace ${name}`;
				},
				error: (err) => {
					isSubmitting = false;
					console.error('Failed to update workspace:', err);
					return `Failed to update workspace: ${(err as ConnectError).message}`;
				}
			}
		);
	}
</script>

<div class="h-full w-full">
	<MultiStepSchemaForm
		apiSchema={schema}
		fields={groupedFields}
		initialData={object}
		title={`Edit Workspace: ${name}`}
		onSubmit={handleMultiStepSubmit}
		transformData={transformFormData}
	/>
</div>

<script lang="ts">
	import { toJson } from '@bufbuild/protobuf';
	import { StructSchema } from '@bufbuild/protobuf/wkt';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext, onMount } from 'svelte';
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
		onsuccess
	}: {
		name: string;
		onsuccess?: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	const cluster = $derived(page.params.scope ?? '');

	let apiSchema: K8sOpenAPISchema | undefined = $state();
	let isSubmitting = $state(false);
	let initialData: Record<string, unknown> | undefined = $state();

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
		if (isSubmitting) return;
		isSubmitting = true;

		// Construct the full resource object - metadata.name should already be in data or we ensure it
		const resourceObject = {
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

	onMount(async () => {
		try {
			// Fetch Schema
			const schemaRes = await resourceClient.schema({
				cluster: 'aaa',
				group: 'tenant.otterscale.io',
				version: 'v1alpha1',
				kind: 'Workspace'
			});
			apiSchema = toJson(StructSchema, schemaRes) as K8sOpenAPISchema;

			// Fetch Existing Resource
			const resourceRes = await resourceClient.get({
				cluster: 'aaa',
				group: 'tenant.otterscale.io',
				version: 'v1alpha1',
				resource: 'workspaces',
				name: name
			});
			console.log('schemaRes', schemaRes);
			console.log('resourceRes', resourceRes);
			if (resourceRes.object) {
				// resourceRes.object is already a plain JS object, no conversion needed
				initialData = resourceRes.object as Record<string, unknown>;
			}
		} catch (err) {
			console.error('Failed to load workspace data:', err);
			toast.error(`Failed to load workspace data: ${(err as ConnectError).message}`);
		}
	});
</script>

<div class="h-full w-full">
	{#if apiSchema && initialData}
		<MultiStepSchemaForm
			{apiSchema}
			fields={groupedFields}
			{initialData}
			title={`Edit Workspace: ${name}`}
			onSubmit={handleMultiStepSubmit}
		/>
	{:else}
		<div class="flex h-32 items-center justify-center">
			<p class="text-muted-foreground">Loading workspace data...</p>
		</div>
	{/if}
</div>

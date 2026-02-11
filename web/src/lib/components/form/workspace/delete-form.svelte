<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import {
		type GroupedFields,
		type K8sOpenAPISchema,
		MultiStepSchemaForm
	} from '$lib/components/custom/schema-form';

	let {
		schema: apiSchema,
		object
	}: {
		schema: K8sOpenAPISchema;
		object: Record<string, unknown>;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	const cluster = $derived(page.params.cluster ?? page.params.scope ?? '');

	function getCleanedObject() {
		const copy = structuredClone($state.snapshot(object) as Record<string, unknown>);
		if (copy.metadata && typeof copy.metadata === 'object') {
			delete (copy.metadata as Record<string, unknown>).managedFields;
		}
		return copy;
	}

	const fields: GroupedFields = {
		Confirm: {
			'metadata.name': { title: 'Workspace Name' }
		}
	};

	let isSubmitting = $state(false);
	async function handleSubmit(data: Record<string, unknown>) {
		if (isSubmitting) return;
		isSubmitting = true;

		toast.promise(
			async () => {
				await resourceClient.delete({
					cluster,
					group: 'tenant.otterscale.io',
					version: 'v1alpha1',
					resource: 'workspaces',
					name: (data as TenantOtterscaleIoV1Alpha1Workspace)?.metadata?.name
				});
			},
			{
				loading: `Deleting workspace ${name}...`,
				success: () => {
					return `Successfully deleted workspace ${name}`;
				},
				error: (err) => {
					console.error('Failed to delete workspace:', err);
					return `Failed to delete workspace: ${(err as ConnectError).message}`;
				},
				finally: () => {
					isSubmitting = false;
				}
			}
		);
	}
</script>

<div class="h-full w-full">
	<MultiStepSchemaForm
		{apiSchema}
		{fields}
		initialData={getCleanedObject()}
		title="Delete Workspace"
		onSubmit={handleSubmit}
	/>
</div>

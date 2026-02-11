<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { Trash2 } from '@lucide/svelte';
	import type { TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import { type GroupedFields, MultiStepSchemaForm } from '$lib/components/custom/schema-form';
	import * as Dialog from '$lib/components/ui/dialog';
	import * as Item from '$lib/components/ui/item';

	let {
		schema: apiSchema,
		onOpenChangeComplete
	}: {
		schema: any;
		onOpenChangeComplete: () => void;
	} = $props();

	const cluster = $derived(page.params.cluster ?? page.params.scope ?? '');

	let open = $state(false);

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);

	function getCleanedObject() {
		const copy = structuredClone({} as Record<string, unknown>);
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
					open = false;
				}
			}
		);
	}
</script>

<Dialog.Root bind:open {onOpenChangeComplete}>
	<Dialog.Trigger class="w-full">
		<Item.Root class="p-0 text-xs **:text-destructive" size="sm">
			<Item.Media>
				<Trash2 />
			</Item.Media>
			<Item.Content>
				<Item.Title>Delete</Item.Title>
			</Item.Content>
		</Item.Root>
	</Dialog.Trigger>
	<Dialog.Content class="min-h-[38vh] min-w-[38vw]">
		<div class="h-full w-full">
			<MultiStepSchemaForm
				{apiSchema}
				{fields}
				initialData={getCleanedObject()}
				title="Delete Workspace"
				onSubmit={handleSubmit}
			/>
		</div>
	</Dialog.Content>
</Dialog.Root>

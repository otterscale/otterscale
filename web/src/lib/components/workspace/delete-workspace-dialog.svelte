<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { buttonVariants } from '$lib/components/ui/button';

	let {
		open = $bindable(false),
		name,
		onsuccess
	}: {
		open: boolean;
		name: string;
		onsuccess?: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	const cluster = $derived(page.params.scope ?? '');

	let isDeleting = $state(false);

	async function handleDelete() {
		if (isDeleting) return;
		isDeleting = true;

		toast.promise(
			async () => {
				await resourceClient.delete({
					cluster,
					group: 'tenant.otterscale.io',
					version: 'v1alpha1',
					resource: 'workspaces',
					name: name
				});
			},
			{
				loading: `Deleting workspace ${name}...`,
				success: () => {
					isDeleting = false;
					open = false;
					onsuccess?.();
					return `Successfully deleted workspace ${name}`;
				},
				error: (err) => {
					isDeleting = false;
					console.error('Failed to delete workspace:', err);
					return `Failed to delete workspace: ${(err as ConnectError).message}`;
				}
			}
		);
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Are you absolutely sure?</AlertDialog.Title>
			<AlertDialog.Description>
				This action cannot be undone. This will permanently delete the workspace
				<span class="font-semibold text-foreground">{name}</span> and remove all associated data.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				class={buttonVariants({ variant: 'destructive' })}
				onclick={handleDelete}
				disabled={isDeleting}
			>
				{isDeleting ? 'Deleting...' : 'Delete'}
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

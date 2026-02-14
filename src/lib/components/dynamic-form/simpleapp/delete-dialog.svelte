<script lang="ts">
	import type { ConnectError } from '@connectrpc/connect';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { Trash2 } from '@lucide/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/stores';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';

	let {
		name,
		onOpenChangeComplete,
		onsuccess
	}: {
		name: string;
		onOpenChangeComplete: () => void;
		onsuccess?: () => void;
	} = $props();

	const transport: Transport = getContext('transport');
	const resourceClient = createClient(ResourceService, transport);
	const cluster = $derived($page.params.cluster ?? $page.params.scope ?? '');

	let open = $state(false);
	let confirmName = $state('');
	let invalid = $state(false);
	let isDeleting = $state(false);

	function init() {
		confirmName = '';
		invalid = false;
		isDeleting = false;
	}

	async function handleDelete() {
		if (isDeleting) return;
		isDeleting = true;

		toast.promise(
			async () => {
				await resourceClient.delete({
					cluster,
					namespace: $page.url.searchParams.get('namespace') ?? '',
					group: 'apps.otterscale.io',
					version: 'v1alpha1',
					resource: 'simpleapps',
					name: name
				});
			},
			{
				loading: `Deleting simpleapp ${name}...`,
				success: () => {
					isDeleting = false;
					open = false;
					onsuccess?.();
					// Redirect after delete
					goto(
						resolve(
							`/(auth)/${cluster}/SimpleApp?group=apps.otterscale.io&version=v1alpha1&namespace=${$page.url.searchParams.get('namespace') ?? ''}&resource=simpleapps`
						)
					);
					return `Successfully deleted simpleapp ${name}`;
				},
				error: (err) => {
					isDeleting = false;
					console.error('Failed to delete simpleapp:', err);
					return `Failed to delete simpleapp: ${(err as ConnectError).message}`;
				}
			}
		);
	}
</script>

<AlertDialog.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
	{onOpenChangeComplete}
>
	<AlertDialog.Trigger class="flex w-full items-center gap-2">
		<Trash2 size={16} />
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Delete SimpleApp</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>SimpleApp Name</Form.Label>
					<Form.Help>
						This action cannot be undone. Please type <strong>{name}</strong> to confirm deletion.
					</Form.Help>
					<SingleInput.Confirm required target={name} bind:value={confirmName} bind:invalid />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action disabled={invalid || isDeleting} onclick={handleDelete}>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

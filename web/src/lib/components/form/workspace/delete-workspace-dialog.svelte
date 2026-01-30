<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { Trash2 } from '@lucide/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { Button } from '$lib/components/ui/button';
	import { m } from '$lib/paraglide/messages';

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

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger>
		<Button variant="destructive" size="icon">
			<Trash2 />
		</Button>
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete Workspace</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Workspace Name</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: 'Workspace' })}
					</Form.Help>
					<SingleInput.Confirm required target={name} bind:value={confirmName} bind:invalid />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action disabled={invalid || isDeleting} variant="destructive" onclick={handleDelete}>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

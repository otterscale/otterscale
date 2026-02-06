<script lang="ts">
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { Trash2 } from '@lucide/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { ResourceService } from '$lib/api/resource/v1/resource_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import * as AlertDialog from '$lib/components/ui/alert-dialog';
	import { m } from '$lib/paraglide/messages';

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
	const cluster = $derived(page.params.cluster ?? page.params.scope ?? '');

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
					namespace: page.url.searchParams.get('namespace') ?? '',
					group: 'batch', // Changed from tenant.otterscale.io
					version: 'v1', // Changed from v1alpha1
					resource: 'cronjobs', // Changed from workspaces
					name: name
				});
			},
			{
				loading: `Deleting cronjob ${name}...`,
				success: () => {
					isDeleting = false;
					open = false;
					onsuccess?.();
					// Redirect after delete - using scope root for now
					goto(
						resolve(
							`/(auth)/${cluster}/CronJob?group=batch&version=v1&namespace=${page.url.searchParams.get('namespace') ?? ''}&resource=cronjobs`
						)
					);
					return `Successfully deleted cronjob ${name}`;
				},
				error: (err) => {
					isDeleting = false;
					console.error('Failed to delete cronjob:', err);
					return `Failed to delete cronjob: ${(err as ConnectError).message}`;
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
		<AlertDialog.Header>Delete CronJob</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>CronJob Name</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: 'CronJob' })}
					</Form.Help>
					<SingleInput.Confirm required target={name} bind:value={confirmName} bind:invalid />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>
				{m.cancel()}
			</AlertDialog.Cancel>
			<AlertDialog.Action disabled={invalid || isDeleting} onclick={handleDelete}>
				{m.confirm()}
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

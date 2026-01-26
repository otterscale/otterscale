<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { DeletePoolRequest, Pool } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		pool,
		scope,
		reloadManager,
		closeActions
	}: {
		pool: Pool;
		scope: string;
		reloadManager: ReloadManager;
		closeActions: () => void;
	} = $props();

	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);
	let invalid: boolean | undefined = $state();

	let request = $state({} as DeletePoolRequest);
	function init() {
		request = {
			scope: scope
		} as DeletePoolRequest;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) {
			closeActions();
		}
	}}
>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		{m.delete()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete_pool()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.Confirm
						id="deletion"
						required
						target={pool.name}
						bind:value={request.poolName}
						bind:invalid
					/>
				</Form.Field>
				<Form.Help>
					{m.deletion_warning({
						identifier: m.pool_name()
					})}
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.deletePool(request), {
							loading: `Deleting ${request.poolName}...`,
							success: () => {
								reloadManager.force();
								return `Delete ${request.poolName}`;
							},
							error: (error) => {
								let message = `Fail to delete ${request.poolName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY,
									closeButton: true
								});
								return message;
							}
						});
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

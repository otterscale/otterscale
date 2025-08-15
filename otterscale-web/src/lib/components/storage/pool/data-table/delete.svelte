<script lang="ts" module>
	import type { DeletePoolRequest, Pool } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { currentCeph } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
</script>

<script lang="ts">
	let {
		pool
	}: {
		pool: Pool;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name
	} as DeletePoolRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new StateController(false);

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('ReloadManager');
	const storageClient = createClient(StorageService, transport);

	let invalid: boolean | undefined = $state();
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		Delete
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete Pool</Modal.Header>
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
					Please type the pool name exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.info(`Deleting ${request.poolName}...`);
						storageClient
							.deletePool(request)
							.then((r) => {
								reloadManager.force();
								toast.success(`Delete ${request.poolName}`);
							})
							.catch((e) => {
								toast.error(`Fail to delete pool: ${e.toString()}`);
							});
						stateController.close();
					}}
				>
					Delete
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

<script lang="ts" module>
	import type { DeletePoolRequest, Pool } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		pool,
		data = $bindable()
	}: {
		selectedScopeUuid: string;
		pool: Pool;
		data: Writable<Pool[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScopeUuid,
		facilityName: 'ceph-mon'
	} as DeletePoolRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new StateController(false);

	const transport: Transport = getContext('transport');
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
					<SingleInput.DeletionConfirm
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
								toast.success(`Delete ${request.poolName}`);
								storageClient
									.listPools({ scopeUuid: selectedScopeUuid, facilityName: 'ceph-mon' })
									.then((r) => {
										data.set(r.pools);
									});
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

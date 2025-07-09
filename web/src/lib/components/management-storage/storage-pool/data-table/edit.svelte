<script lang="ts" module>
	import type { Pool, UpdatePoolRequest } from '$gen/api/storage/v1/storage_pb';
	import { StorageService } from '$gen/api/storage/v1/storage_pb';
	import * as AlertDialog from '$lib/components/custom/alert-dialog';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { DialogStateController } from '$lib/components/custom/utils.svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
	import { QUOTAS_BYTES_HELP_TEXT, QUOTAS_OBJECTS_HELP_TEXT } from './helper';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		pool,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		pool: Pool;
		data: Writable<Pool[]>;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: selectedScope,
		facilityName: selectedFacility,
		poolName: pool.name,
		quotaBytes: pool.quotaBytes,
		quotaObjects: pool.quotaObjects
	} as UpdatePoolRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new DialogStateController(false);

	const transport: Transport = getContext('transport');
	const storageClient = createClient(StorageService, transport);
</script>

<AlertDialog.Root bind:open={stateController.state}>
	<AlertDialog.Trigger class="flex h-full w-full items-center gap-2">
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>Update Pool</AlertDialog.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.poolName} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Quotas</Form.Legend>

				<Form.Field>
					<Form.Label>Bytes</Form.Label>
					<SingleInput.General type="number" bind:value={request.quotaBytes} />
				</Form.Field>
				<Form.Help>
					{QUOTAS_BYTES_HELP_TEXT}
				</Form.Help>

				<Form.Field>
					<Form.Label>Objects</Form.Label>
					<SingleInput.General type="number" bind:value={request.quotaObjects} />
				</Form.Field>
				<Form.Help>
					{QUOTAS_OBJECTS_HELP_TEXT}
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset}>Cancel</AlertDialog.Cancel>
			<AlertDialog.ActionsGroup>
				<AlertDialog.Action
					onclick={() => {
						console.log(request);
						stateController.close();
						storageClient
							.updatePool(request)
							.then((r) => {
								toast.success(`Update ${request.poolName}`);
								storageClient
									.listPools({ scopeUuid: selectedScope, facilityName: selectedFacility })
									.then((r) => {
										data.set(r.pools);
									});
							})
							.catch((e) => {
								toast.error(`Fail to update pool: ${e.toString()}`);
								console.log(e);
							})
							.finally(() => {
								reset();
							});
					}}
				>
					Edit
				</AlertDialog.Action>
			</AlertDialog.ActionsGroup>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

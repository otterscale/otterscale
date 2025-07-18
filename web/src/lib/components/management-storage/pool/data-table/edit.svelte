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
	console.log('DEFAULT_REQUEST', DEFAULT_REQUEST);
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
					<SingleInput.Measurement
						bind:value={request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 0), label: 'B' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 1), label: 'KB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 2), label: 'MB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 5), label: 'PB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>
				<Form.Help>
					{QUOTAS_BYTES_HELP_TEXT}
				</Form.Help>

				<Form.Field>
					<Form.Label>Objects</Form.Label>
					<SingleInput.BigInteger bind:value={request.quotaObjects} />
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
						console.log('request', request);
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
								stateController.close();
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

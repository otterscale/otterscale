<script lang="ts" module>
	import type { Pool, UpdatePoolRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as SingleStepModal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { currentCeph } from '$lib/stores';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { QUOTAS_BYTES_HELP_TEXT, QUOTAS_OBJECTS_HELP_TEXT } from './helper';
</script>

<script lang="ts">
	let {
		pool
	}: {
		pool: Pool;
	} = $props();

	const DEFAULT_REQUEST = {
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		poolName: pool.name,
		quotaBytes: pool.quotaBytes,
		quotaObjects: pool.quotaObjects
	} as UpdatePoolRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	const stateController = new StateController(false);

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('ReloadManager');
	const storageClient = createClient(StorageService, transport);

	let invalid = $state(false);
</script>

<SingleStepModal.Root bind:open={stateController.state}>
	<SingleStepModal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		Edit
	</SingleStepModal.Trigger>
	<SingleStepModal.Content>
		<SingleStepModal.Header>Update Pool</SingleStepModal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General required type="text" bind:value={request.poolName} bind:invalid />
				</Form.Field>

				<Form.Field>
					<Form.Label>Quota Size</Form.Label>
					<Form.Help>
						{QUOTAS_BYTES_HELP_TEXT}
					</Form.Help>
					<SingleInput.Measurement
						bind:value={request.quotaBytes}
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Quota Objects</Form.Label>
					<Form.Help>
						{QUOTAS_OBJECTS_HELP_TEXT}
					</Form.Help>
					<SingleInput.General bind:value={request.quotaObjects} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<SingleStepModal.Footer>
			<SingleStepModal.Cancel onclick={reset}>Cancel</SingleStepModal.Cancel>
			<SingleStepModal.ActionsGroup>
				<SingleStepModal.Action
					disabled={invalid}
					onclick={() => {
						toast.info(`Updating ${request.poolName}...`);
						storageClient
							.updatePool(request)
							.then((r) => {
								reloadManager.force();
								toast.success(`Update ${request.poolName}`);
							})
							.catch((e) => {
								toast.error(`Fail to update pool: ${e.toString()}`);
								console.log(e);
							})
							.finally(() => {
								reset();
							});
						stateController.close();
					}}
				>
					Edit
				</SingleStepModal.Action>
			</SingleStepModal.ActionsGroup>
		</SingleStepModal.Footer>
	</SingleStepModal.Content>
</SingleStepModal.Root>

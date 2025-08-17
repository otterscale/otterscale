<script lang="ts" module>
	import type { Pool, UpdatePoolRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as SingleStepModal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { currentCeph } from '$lib/stores';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { QUOTAS_BYTES_HELP_TEXT, QUOTAS_OBJECTS_HELP_TEXT } from '../helper';
</script>

<script lang="ts">
	let {
		pool
	}: {
		pool: Pool;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let invalid = $state(false);

	const requestManager = new RequestManager<UpdatePoolRequest>({
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		poolName: pool.name,
		quotaBytes: pool.quotaBytes,
		quotaObjects: pool.quotaObjects
	} as UpdatePoolRequest);
	const stateController = new StateController(false);
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
					<SingleInput.General
						required
						type="text"
						bind:value={requestManager.request.poolName}
						bind:invalid
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>Quota Size</Form.Label>
					<Form.Help>
						{QUOTAS_BYTES_HELP_TEXT}
					</Form.Help>
					<SingleInput.Measurement
						bind:value={requestManager.request.quotaBytes}
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
					<SingleInput.General bind:value={requestManager.request.quotaObjects} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<SingleStepModal.Footer>
			<SingleStepModal.Cancel
				onclick={() => {
					requestManager.reset();
				}}
			>
				Cancel
			</SingleStepModal.Cancel>
			<SingleStepModal.ActionsGroup>
				<SingleStepModal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => storageClient.updatePool(requestManager.request), {
							loading: `Updating ${requestManager.request.poolName}...`,
							success: (response) => {
								reloadManager.force();
								return `Update ${requestManager.request.poolName}`;
							},
							error: (error) => {
								let message = `Fail to update ${requestManager.request.poolName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						requestManager.reset();
						stateController.close();
					}}
				>
					Edit
				</SingleStepModal.Action>
			</SingleStepModal.ActionsGroup>
		</SingleStepModal.Footer>
	</SingleStepModal.Content>
</SingleStepModal.Root>

<script lang="ts" module>
	import type { Image, UpdateImageRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { currentCeph } from '$lib/stores';
	import { cn } from '$lib/utils';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
</script>

<script lang="ts">
	let {
		image
	}: {
		image: Image;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalid = $state(false);
	const storageClient = createClient(StorageService, transport);
	const requestManager = new RequestManager<UpdateImageRequest>({
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name,
		poolName: image.poolName,
		imageName: image.name,
		quotaBytes: image.quotaBytes
	} as UpdateImageRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class={cn('flex h-full w-full items-center gap-2')}>
		<Icon icon="ph:pencil" />
		Edit
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Updatge RADOS Block Device</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Label>Quota Size</Form.Label>
				<Form.Field>
					<SingleInput.Measurement
						required
						transformer={(value) => String(value)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
						bind:value={requestManager.request.quotaBytes}
						bind:invalid
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					requestManager.reset();
				}}
			>
				Cancel
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(storageClient.updateImage(requestManager.request), {
							loading: `Updating ${requestManager.request.imageName}...`,
							success: (response) => {
								reloadManager.force();
								return `Updated ${requestManager.request.imageName}`;
							},
							error: (error) => {
								let message = `Fail to updated ${requestManager.request.imageName}`;
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
					Update
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

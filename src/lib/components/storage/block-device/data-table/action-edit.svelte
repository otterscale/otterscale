<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { Image, UpdateImageRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		image,
		scope,
		reloadManager,
		closeActions
	}: {
		image: Image;
		scope: string;
		reloadManager: ReloadManager;
		closeActions: () => void;
	} = $props();
	const transport: Transport = getContext('transport');

	const storageClient = createClient(StorageService, transport);

	let request = $state({} as UpdateImageRequest);
	function init() {
		request = {
			scope: scope,
			poolName: image.poolName,
			imageName: image.name,
			quotaBytes: image.quotaBytes
		} as UpdateImageRequest;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) {
			closeActions();
		}
	}}
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		{m.edit()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit_rbd()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Label>{m.quota_size()}</Form.Label>
				<Form.Field>
					<SingleInput.Measurement
						transformer={(value) => (typeof value === 'number' ? BigInt(value) : undefined)}
						units={[
							{ value: Math.pow(2, 10 * 3), label: 'GB' } as SingleInput.UnitType,
							{ value: Math.pow(2, 10 * 4), label: 'TB' } as SingleInput.UnitType
						]}
						bind:value={request.quotaBytes}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					onclick={() => {
						toast.promise(storageClient.updateImage(request), {
							loading: `Updating ${request.imageName}...`,
							success: () => {
								reloadManager.force();
								return `Updated ${request.imageName}`;
							},
							error: (error) => {
								let message = `Fail to updated ${request.imageName}`;
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

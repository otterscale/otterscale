<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import type { DeleteImageRequest, Image } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
	import { currentCeph } from '$lib/stores';
</script>

<script lang="ts">
	let {
		image
	}: {
		image: Image;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let isPoolNameInvalid = $state(false);
	let isImageNameInvalid = $state(false);
	const storageClient = createClient(StorageService, transport);
	const defaults = {
		scope: $currentCeph?.scope,
		facility: $currentCeph?.name
	} as DeleteImageRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		{m.delete()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete_rbd()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.image_name()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: m.image_name() })}
					</Form.Help>
					<SingleInput.Confirm
						required
						target={image.name}
						bind:value={request.imageName}
						bind:invalid={isImageNameInvalid}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>{m.pool_name()}</Form.Label>
					<Form.Help>
						{m.deletion_warning({ identifier: m.pool_name() })}
					</Form.Help>
					<SingleInput.Confirm
						required
						target={image.poolName}
						bind:value={request.poolName}
						bind:invalid={isPoolNameInvalid}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={isPoolNameInvalid || isImageNameInvalid}
					onclick={() => {
						toast.promise(() => storageClient.deleteImage(request), {
							loading: `Deleting ${request.imageName}...`,
							success: () => {
								reloadManager.force();
								return `Delete ${request.imageName}`;
							},
							error: (error) => {
								let message = `Fail to delete ${request.imageName}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});
						reset();
						close();
					}}
				>
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

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
	import type { Booleanified } from '$lib/components/custom/modal/single-step/type';
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

	let request = $state({} as DeleteImageRequest);
	function init() {
		request = {
			scope: scope
		} as DeleteImageRequest;
	}

	let invalidity = $state({} as Booleanified<DeleteImageRequest>);
	const invalid = $derived(invalidity.poolName || invalidity.imageName);

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
						bind:invalid={invalidity.imageName}
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
						bind:invalid={invalidity.poolName}
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
					disabled={invalid}
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

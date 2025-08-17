<script lang="ts" module>
	import type { Bucket, DeleteBucketRequest } from '$lib/api/storage/v1/storage_pb';
	import { StorageService } from '$lib/api/storage/v1/storage_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { currentCeph } from '$lib/stores';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
</script>

<script lang="ts">
	let {
		bucket
	}: {
		bucket: Bucket;
	} = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const storageClient = createClient(StorageService, transport);
	let invalid = $state(false);

	const requestManager = new RequestManager<DeleteBucketRequest>({
		scopeUuid: $currentCeph?.scopeUuid,
		facilityName: $currentCeph?.name
	} as DeleteBucketRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="text-destructive flex h-full w-full items-center gap-2">
		<Icon icon="ph:trash" />
		Delete
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete Bucket</Modal.Header>
		<Form.Root bind:invalid>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.Confirm
						required
						id="deletion"
						target={bucket.name}
						bind:value={requestManager.request.bucketName}
					/>
				</Form.Field>
				<Form.Help>
					Please type the bucket name exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
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
						toast.promise(() => storageClient.deleteBucket(requestManager.request), {
							loading: `Deleting ${requestManager.request.bucketName}...`,
							success: (response) => {
								reloadManager.force();
								return `Delete ${requestManager.request.bucketName}`;
							},
							error: (error) => {
								let message = `Fail to delete ${requestManager.request.bucketName}`;
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
					Delete
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

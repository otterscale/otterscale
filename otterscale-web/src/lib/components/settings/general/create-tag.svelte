<script lang="ts" module>
	import { TagService, type CreateTagRequest, type Tag } from '$lib/api/tag/v1/tag_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let { tags = $bindable() }: { tags: Writable<Tag[]> } = $props();

	const transport: Transport = getContext('transport');

	const client = createClient(TagService, transport);
	const requestManager = new RequestManager<CreateTagRequest>({} as CreateTagRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create Tag</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General type="text" bind:value={requestManager.request.name} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Comment</Form.Label>
					<SingleInput.General type="text" bind:value={requestManager.request.comment} />
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
					onclick={() => {
						toast.promise(() => client.createTag(requestManager.request), {
							loading: 'Loading...',
							success: () => {
								client.listTags({}).then((response) => {
									tags.set(response.tags);
								});
								return `Create ${requestManager.request.name} success`;
							},
							error: (error) => {
								let message = `Fail to create ${requestManager.request.name}`;
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
					Create
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

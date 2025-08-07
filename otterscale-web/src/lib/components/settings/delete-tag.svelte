<script lang="ts" module>
	import { TagService, type DeleteTagRequest, type Tag } from '$lib/api/tag/v1/tag_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let { tag, tags = $bindable() }: { tag: Tag; tags: Writable<Tag[]> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(TagService, transport);

	const DEFAULT_REQUEST = {} as DeleteTagRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let invalid: boolean | undefined = $state();

	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		Delete
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete Fabric</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.DeletionConfirm
						required
						target={tag.name}
						value={request.name}
						bind:invalid
					/>
				</Form.Field>
				<Form.Help>
					Please type the tag name exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.deleteTag(request), {
							loading: 'Loading...',
							success: () => {
								client.listTags({}).then((response) => {
									tags.set(response.tags);
								});
								return `Delete ${tag.name} success`;
							},
							error: (error) => {
								let message = `Fail to delete ${tag.name}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
								});
								return message;
							}
						});

						reset();
						stateController.close();
					}}
				>
					Delete
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

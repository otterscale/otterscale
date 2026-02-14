<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { type DeleteTagRequest, MachineService, type Tag } from '$lib/api/machine/v1/machine_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { tag, tags }: { tag: Tag; tags: Writable<Tag[]> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(MachineService, transport);
	const isDisabled = $derived(tag.comment.toLocaleLowerCase().includes('built-in'));

	let request = $state({} as DeleteTagRequest);
	let invalid: boolean | undefined = $state();
	let open = $state(false);

	function init() {
		request = {} as DeleteTagRequest;
	}

	function close() {
		open = false;
	}
</script>

<span>
	<Modal.Root
		bind:open
		onOpenChange={(isOpen) => {
			if (isOpen) {
				init();
			}
		}}
	>
		<Modal.Trigger variant="destructive" disabled={isDisabled}>
			<Icon icon="ph:trash" />
			{m.delete()}
		</Modal.Trigger>
		<Modal.Content>
			<Modal.Header>{m.delete_fabric()}</Modal.Header>
			<Form.Root>
				<Form.Fieldset>
					<Form.Field>
						<SingleInput.Confirm
							required
							target={tag.name}
							bind:value={request.name}
							bind:invalid
						/>
					</Form.Field>
					<Form.Help>
						{m.deletion_warning({ identifier: m.name() })}
					</Form.Help>
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
</span>

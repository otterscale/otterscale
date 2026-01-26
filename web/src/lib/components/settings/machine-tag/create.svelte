<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import { toast } from 'svelte-sonner';

	import { type CreateTagRequest, MachineService, type Tag } from '$lib/api/machine/v1/machine_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { tags }: { tags: Writable<Tag[]> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(MachineService, transport);

	let request = $state({} as CreateTagRequest);
	let open = $state(false);

	function init() {
		request = {} as CreateTagRequest;
	}

	function close() {
		open = false;
	}
</script>

<Modal.Root
	bind:open
	onOpenChange={(isOpen) => {
		if (isOpen) {
			init();
		}
	}}
>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_machine_tag()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.GeneralRule type="text" bind:value={request.name} validateRule="rfc1123" />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.comment()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.comment} />
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
						toast.promise(() => client.createTag(request), {
							loading: 'Loading...',
							success: () => {
								client.listTags({}).then((response) => {
									tags.set(response.tags);
								});
								return `Create ${request.name} success`;
							},
							error: (error) => {
								let message = `Fail to create ${request.name}`;
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

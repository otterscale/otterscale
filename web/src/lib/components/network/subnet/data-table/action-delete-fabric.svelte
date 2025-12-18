<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type DeleteNetworkRequest,
		type Network_Fabric,
		NetworkService
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { fabric, reloadManager }: { fabric: Network_Fabric; reloadManager: ReloadManager } =
		$props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	let request = $state({} as DeleteNetworkRequest);
	let invalid: boolean | undefined = $state();
	let open = $state(false);

	function init() {
		request = {
			id: fabric.id
		} as DeleteNetworkRequest;
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
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		{m.delete_fabric()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.delete_fabric()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<SingleInput.Confirm id="deletion" required target={fabric.name} bind:invalid />
				</Form.Field>
				<Form.Help>
					{m.deletion_warning({ identifier: m.name() })}
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel>{m.cancel()}</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.deleteNetwork(request), {
							loading: 'Loading...',
							success: () => {
								reloadManager.force();
								return `Delete ${fabric.name} success`;
							},
							error: (error) => {
								let message = `Fail to delete ${fabric.name}`;
								toast.error(message, {
									description: (error as ConnectError).message.toString(),
									duration: Number.POSITIVE_INFINITY
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

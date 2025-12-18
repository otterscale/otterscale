<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type Network_Fabric,
		NetworkService,
		type UpdateFabricRequest
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let {
		fabric,
		reloadManager,
		closeActions
	}: { fabric: Network_Fabric; reloadManager: ReloadManager; closeActions: () => void } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	let request = $state({} as UpdateFabricRequest);
	let invalid: boolean | undefined = $state();
	let open = $state(false);

	function init() {
		request = {
			id: fabric.id,
			name: fabric.name
		} as UpdateFabricRequest;
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
	onOpenChangeComplete={(isOpen) => {
		if (!isOpen) {
			closeActions();
		}
	}}
>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		{m.edit_fabric()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit_fabric()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General type="text" required bind:value={request.name} bind:invalid />
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
						toast.promise(() => client.updateFabric(request), {
							loading: 'Loading...',
							success: () => {
								reloadManager.force();
								return `Update ${fabric.name} success`;
							},
							error: (error) => {
								let message = `Fail to update ${fabric.name}`;
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

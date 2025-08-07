<script lang="ts" module>
	import {
		NetworkService,
		type DeleteNetworkRequest,
		type Network,
		type Network_Fabric
	} from '$lib/api/network/v1/network_pb';
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
	let {
		fabric,
		networks = $bindable()
	}: { fabric: Network_Fabric; networks: Writable<Network[]> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	const DEFAULT_REQUEST = { id: fabric.id } as DeleteNetworkRequest;

	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let name = $state('');
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
						id="deletion"
						required
						target={fabric.name}
						value={name}
						bind:invalid
					/>
				</Form.Field>
				<Form.Help>
					Please type the fabric id exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.deleteNetwork(request), {
							loading: 'Loading...',
							success: () => {
								client.listNetworks({}).then((response) => {
									networks.set(response.networks);
								});
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

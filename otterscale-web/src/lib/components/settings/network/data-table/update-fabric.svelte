<script lang="ts" module>
	import {
		NetworkService,
		type Network,
		type Network_Fabric,
		type UpdateFabricRequest
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

	const DEFAULT_REQUEST = { id: fabric.id, name: fabric.name } as UpdateFabricRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let invalid: boolean | undefined = $state();

	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		Edit Fabric
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit Fabric</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General type="text" required value={request.name} bind:invalid />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.updateFabric(request), {
							loading: 'Loading...',
							success: () => {
								client.listNetworks({}).then((response) => {
									networks.set(response.networks);
								});
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

						reset();
						stateController.close();
					}}
				>
					Edit
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

<script lang="ts" module>
	import {
		NetworkService,
		type Network,
		type Network_IPRange,
		type UpdateIPRangeRequest
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
		ipRange,
		networks = $bindable()
	}: {
		ipRange: Network_IPRange;
		networks: Writable<Network[]>;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	const DEFAULT_REQUEST = {
		id: ipRange.id,
		startIp: ipRange.startIp,
		endIp: ipRange.endIp,
		comment: ipRange.comment
	} as UpdateIPRangeRequest;
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
		Edit
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit IP Range</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Start</Form.Label>
					<SingleInput.General type="text" bind:value={request.startIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>End</Form.Label>
					<SingleInput.General type="text" bind:value={request.endIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Coment</Form.Label>
					<SingleInput.General type="text" bind:value={request.comment} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset} class="mr-auto">Cancel</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => client.updateIPRange(request), {
						loading: 'Loading...',
						success: () => {
							client.listNetworks({}).then((response) => {
								networks.set(response.networks);
							});
							return `Update ${request.startIp} - ${request.endIp} success`;
						},
						error: (error) => {
							let message = `Fail to update ${request.startIp} - ${request.endIp}`;
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
				Create
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

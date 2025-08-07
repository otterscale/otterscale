<script lang="ts" module>
	import {
		NetworkService,
		type DeleteIPRangeRequest,
		type Network,
		type Network_IPRange
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
	}: { ipRange: Network_IPRange; networks: Writable<Network[]> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	const DEFAULT_REQUEST = { id: ipRange.id } as DeleteIPRangeRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let startIP = $state('');
	let endIP = $state('');
	let invalidStartIP: boolean | undefined = $state();
	let invalidEndIP: boolean | undefined = $state();

	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="destructive">
		<Icon icon="ph:trash" />
		Delete
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Delete IP Range</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Start</Form.Label>
					<SingleInput.DeletionConfirm
						required
						target={ipRange.startIp}
						value={startIP}
						bind:invalid={invalidStartIP}
					/>
				</Form.Field>
				<Form.Field>
					<Form.Label>End</Form.Label>
					<SingleInput.DeletionConfirm
						required
						target={ipRange.endIp}
						value={endIP}
						bind:invalid={invalidEndIP}
					/>
				</Form.Field>
				<Form.Help>
					Please type the ip range exactly to confirm deletion. This action cannot be undone.
				</Form.Help>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalidStartIP || invalidEndIP}
					onclick={() => {
						toast.promise(() => client.deleteIPRange(request), {
							loading: 'Loading...',
							success: () => {
								client.listNetworks({}).then((response) => {
									networks.set(response.networks);
								});
								return `Delete ${ipRange.id} success`;
							},
							error: (error) => {
								let message = `Fail to delete ${ipRange.id}`;
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

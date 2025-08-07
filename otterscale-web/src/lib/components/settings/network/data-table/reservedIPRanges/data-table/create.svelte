<script lang="ts" module>
	import {
		NetworkService,
		type CreateIPRangeRequest,
		type Network,
		type Network_Subnet
	} from '$lib/api/network/v1/network_pb';
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
		subnet,
		networks = $bindable()
	}: {
		subnet: Network_Subnet;
		networks: Writable<Network[]>;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	const DEFAULT_REQUEST = { subnetId: subnet.id } as CreateIPRangeRequest;
	let request = $state(DEFAULT_REQUEST);
	function reset() {
		request = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create IP Range</Modal.Header>
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
					toast.promise(() => client.createIPRange(request), {
						loading: 'Loading...',
						success: () => {
							client.listNetworks({}).then((response) => {
								networks.set(response.networks);
							});
							return `Create ${request.startIp} - ${request.endIp} success`;
						},
						error: (error) => {
							let message = `Fail to create ${request.startIp} - ${request.endIp}`;
							toast.error(message, {
								description: (error as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return message;
						}
					});

					reset();
					close();
				}}
			>
				Create
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

<script lang="ts" module>
	import {
		NetworkService,
		type CreateNetworkRequest,
		type Network
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let { networks = $bindable() }: { networks: Writable<Network[]> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	const DEFAULT_REQUEST = { dhcpOn: true, dnsServers: [] as string[] } as CreateNetworkRequest;
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
		<Modal.Header>Create Network</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Legend>VLAN</Form.Legend>
				<Form.Field>
					<SingleInput.Boolean descriptor={() => 'DHCP ON'} bind:value={request.dhcpOn} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Subnet</Form.Legend>
				<Form.Field>
					<Form.Label>CIDR</Form.Label>
					<SingleInput.General type="text" bind:value={request.cidr} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Gateway</Form.Label>
					<SingleInput.General type="text" bind:value={request.gatewayIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>DNS</Form.Label>
					<MultipleInput.Root type="text" bind:values={request.dnsServers}>
						<MultipleInput.Viewer />
						<MultipleInput.Controller>
							<MultipleInput.Input />
							<MultipleInput.Add />
							<MultipleInput.Clear />
						</MultipleInput.Controller>
					</MultipleInput.Root>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset} class="mr-auto">Cancel</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					toast.promise(() => client.createNetwork(request), {
						loading: 'Loading...',
						success: () => {
							client.listNetworks({}).then((response) => {
								networks.set(response.networks);
							});
							return `Create ${request.cidr} success`;
						},
						error: (error) => {
							let message = `Fail to create ${request.cidr}`;
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

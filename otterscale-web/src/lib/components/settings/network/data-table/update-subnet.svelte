<script lang="ts" module>
	import {
		NetworkService,
		type Network,
		type Network_Subnet,
		type UpdateSubnetRequest
	} from '$lib/api/network/v1/network_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
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
	let {
		subnet,
		networks = $bindable()
	}: { subnet: Network_Subnet; networks: Writable<Network[]> } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	const DEFAULT_REQUEST = {
		id: subnet.id,
		name: subnet.name,
		cidr: subnet.cidr,
		gatewayIp: subnet.gatewayIp,
		dnsServers: subnet.dnsServers,
		description: subnet.description,
		allowDnsResolution: subnet.allowDnsResolution
	} as UpdateSubnetRequest;
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
		Edit Subnet
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit Subnet</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General type="text" required value={request.name} bind:invalid />
				</Form.Field>

				<Form.Field>
					<Form.Label>Description</Form.Label>
					<SingleInput.General type="text" value={request.description} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Network</Form.Legend>

				<Form.Field>
					<Form.Label>CIDR</Form.Label>
					<SingleInput.General type="text" value={request.cidr} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Gateway</Form.Label>
					<SingleInput.General type="text" value={request.gatewayIp} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>DNS</Form.Legend>

				<Form.Field>
					<Form.Label>Server</Form.Label>
					<MultipleInput.Root type="text" bind:values={request.dnsServers}>
						<MultipleInput.Viewer />
						<MultipleInput.Controller>
							<MultipleInput.Input />
							<MultipleInput.Add />
							<MultipleInput.Clear />
						</MultipleInput.Controller>
					</MultipleInput.Root>
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean
						descriptor={() => 'Allow DNS Resolution'}
						bind:value={request.allowDnsResolution}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel onclick={reset}>Cancel</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.updateSubnet(request), {
							loading: 'Loading...',
							success: () => {
								client.listNetworks({}).then((response) => {
									networks.set(response.networks);
								});
								return `Update ${subnet.name} success`;
							},
							error: (error) => {
								let message = `Fail to update ${subnet.name}`;
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

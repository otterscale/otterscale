<script lang="ts">
	import {
		NetworkService,
		type Network_Subnet,
		type UpdateSubnetRequest
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	let { subnet }: { subnet: Network_Subnet } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalid: boolean | undefined = $state();

	const client = createClient(NetworkService, transport);
	const defaults = {
		id: subnet.id,
		name: subnet.name,
		cidr: subnet.cidr,
		gatewayIp: subnet.gatewayIp,
		dnsServers: subnet.dnsServers,
		description: subnet.description,
		allowDnsResolution: subnet.allowDnsResolution
	} as UpdateSubnetRequest;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
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
			<Modal.Cancel
				onclick={() => {
					reset();
				}}>Cancel</Modal.Cancel
			>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.updateSubnet(request), {
							loading: 'Loading...',
							success: () => {
								reloadManager.force();
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
						close();
					}}
				>
					Edit
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

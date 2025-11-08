<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		type Network_Subnet,
		NetworkService,
		type UpdateSubnetRequest
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
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
		{m.edit_subnet()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit_subnet()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General type="text" required bind:value={request.name} bind:invalid />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.description()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.description} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>{m.network()}</Form.Legend>

				<Form.Field>
					<Form.Label>{m.cidr()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.cidr} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.gateway()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.gatewayIp} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>{m.dns()}</Form.Legend>

				<Form.Field>
					<Form.Label>{m.dns_server()}</Form.Label>
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
						descriptor={m.allow_dns_resolution}
						bind:value={request.allowDnsResolution}
					/>
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}>{m.cancel()}</Modal.Cancel
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
					{m.confirm()}
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

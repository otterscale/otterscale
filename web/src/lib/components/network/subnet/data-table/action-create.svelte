<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import { type CreateNetworkRequest, NetworkService } from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { reloadManager }: { reloadManager: ReloadManager } = $props();
	
	const transport: Transport = getContext('transport');

	const client = createClient(NetworkService, transport);
	const defaults = {
		dhcpOn: true,
		dnsServers: [] as string[]
	} as CreateNetworkRequest;
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
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create_network()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Legend>{m.vlan()}</Form.Legend>
				<Form.Field>
					<SingleInput.Boolean descriptor={() => m.dhcp_on()} bind:value={request.dhcpOn} />
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>{m.subnet()}</Form.Legend>
				<Form.Field>
					<Form.Label>{m.cidr()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.cidr} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.gateway()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.gatewayIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.dns()}</Form.Label>
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
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
				class="mr-auto">{m.cancel()}</Modal.Cancel
			>
			<Modal.Action
				onclick={() => {
					toast.promise(() => client.createNetwork(request), {
						loading: 'Loading...',
						success: () => {
							reloadManager.force();
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
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

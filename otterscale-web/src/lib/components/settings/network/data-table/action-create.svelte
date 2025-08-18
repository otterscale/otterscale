<script lang="ts" module>
	import { NetworkService, type CreateNetworkRequest } from '$lib/api/network/v1/network_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import * as Form from '$lib/components/custom/form';
	import { RequestManager } from '$lib/components/custom/form';
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
</script>

<script lang="ts">
	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	const client = createClient(NetworkService, transport);
	const requestManager = new RequestManager<CreateNetworkRequest>({
		dhcpOn: true,
		dnsServers: [] as string[]
	} as CreateNetworkRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		Create
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Create Network</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Legend>VLAN</Form.Legend>
				<Form.Field>
					<SingleInput.Boolean
						descriptor={() => 'DHCP ON'}
						bind:value={requestManager.request.dhcpOn}
					/>
				</Form.Field>
			</Form.Fieldset>

			<Form.Fieldset>
				<Form.Legend>Subnet</Form.Legend>
				<Form.Field>
					<Form.Label>CIDR</Form.Label>
					<SingleInput.General type="text" bind:value={requestManager.request.cidr} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Gateway</Form.Label>
					<SingleInput.General type="text" bind:value={requestManager.request.gatewayIp} />
				</Form.Field>

				<Form.Field>
					<Form.Label>DNS</Form.Label>
					<MultipleInput.Root type="text" bind:values={requestManager.request.dnsServers}>
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
					requestManager.reset();
				}}
				class="mr-auto">Cancel</Modal.Cancel
			>
			<Modal.Action
				onclick={() => {
					toast.promise(() => client.createNetwork(requestManager.request), {
						loading: 'Loading...',
						success: () => {
							reloadManager.force();
							return `Create ${requestManager.request.cidr} success`;
						},
						error: (error) => {
							let message = `Fail to create ${requestManager.request.cidr}`;
							toast.error(message, {
								description: (error as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return message;
						}
					});

					requestManager.reset();
					stateController.close();
				}}
			>
				Create
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

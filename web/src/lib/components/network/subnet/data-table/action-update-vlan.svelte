<script lang="ts" module>
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		NetworkService,
		type Network_Fabric,
		type Network_VLAN,
		type UpdateVLANRequest
	} from '$lib/api/network/v1/network_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { fabric, vlan }: { fabric: Network_Fabric; vlan: Network_VLAN } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('reloadManager');

	let invalid: boolean | undefined = $state();

	const client = createClient(NetworkService, transport);
	const defaults = {
		fabricId: fabric.id,
		vid: vlan.vid,
		name: vlan.name,
		mtu: vlan.mtu,
		description: vlan.description,
		dhcpOn: vlan.dhcpOn
	} as UpdateVLANRequest;
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
		{m.edit_vlan()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.edit_vlan()}</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>{m.name()}</Form.Label>
					<SingleInput.General type="text" required bind:value={request.name} bind:invalid />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.mtu()}</Form.Label>
					<SingleInput.General type="number" bind:value={request.mtu} />
				</Form.Field>

				<Form.Field>
					<Form.Label>{m.description()}</Form.Label>
					<SingleInput.General type="text" bind:value={request.description} />
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean descriptor={() => 'DHCP ON'} bind:value={request.dhcpOn} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					reset();
				}}
			>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.updateVLAN(request), {
							loading: 'Loading...',
							success: () => {
								reloadManager.force();
								return `Update ${vlan.name} success`;
							},
							error: (error) => {
								let message = `Fail to update ${vlan.name}`;
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

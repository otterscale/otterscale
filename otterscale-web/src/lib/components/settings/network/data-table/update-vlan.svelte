<script lang="ts" module>
	import {
		NetworkService,
		type Network,
		type Network_Fabric,
		type Network_VLAN,
		type UpdateVLANRequest
	} from '$lib/api/network/v1/network_pb';
	import { StateController } from '$lib/components/custom/alert-dialog/utils.svelte';
	import { RequestManager } from '$lib/components/custom/form';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let { fabric, vlan }: { fabric: Network_Fabric; vlan: Network_VLAN } = $props();

	const transport: Transport = getContext('transport');
	const reloadManager: ReloadManager = getContext('ReloadManager');

	let invalid: boolean | undefined = $state();

	const client = createClient(NetworkService, transport);
	const requestManager = new RequestManager<UpdateVLANRequest>({
		fabricId: fabric.id,
		vid: vlan.id,
		name: vlan.name,
		mtu: vlan.mtu,
		description: vlan.description,
		dhcpOn: vlan.dhcpOn
	} as UpdateVLANRequest);
	const stateController = new StateController(false);
</script>

<Modal.Root bind:open={stateController.state}>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:pencil" />
		Edit VLAN
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>Edit VLAN</Modal.Header>
		<Form.Root>
			<Form.Fieldset>
				<Form.Field>
					<Form.Label>Name</Form.Label>
					<SingleInput.General
						type="text"
						required
						value={requestManager.request.name}
						bind:invalid
					/>
				</Form.Field>

				<Form.Field>
					<Form.Label>NTU</Form.Label>
					<SingleInput.General type="number" value={requestManager.request.mtu} />
				</Form.Field>

				<Form.Field>
					<Form.Label>Description</Form.Label>
					<SingleInput.General type="text" value={requestManager.request.description} />
				</Form.Field>

				<Form.Field>
					<SingleInput.Boolean descriptor={() => 'DHCP ON'} value={requestManager.request.dhcpOn} />
				</Form.Field>
			</Form.Fieldset>
		</Form.Root>
		<Modal.Footer>
			<Modal.Cancel
				onclick={() => {
					requestManager.reset();
				}}>Cancel</Modal.Cancel
			>
			<Modal.ActionsGroup>
				<Modal.Action
					disabled={invalid}
					onclick={() => {
						toast.promise(() => client.updateVLAN(requestManager.request), {
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

						requestManager.reset();
						stateController.close();
					}}
				>
					Edit
				</Modal.Action>
			</Modal.ActionsGroup>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

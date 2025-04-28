<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { toast } from 'svelte-sonner';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import {
		Nexus,
		type Network_Fabric,
		type UpdateVLANRequest,
		type Network_VLAN
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';

	let { fabric, vlan }: { fabric: Network_Fabric; vlan: Network_VLAN } = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		fabricId: fabric.id,
		vid: vlan.id,
		name: vlan.name,
		mtu: vlan.mtu,
		description: vlan.description,
		dhcpOn: vlan.dhcpOn
	} as UpdateVLANRequest;

	let updateVLANRequest = $state(DEFAULT_REQUEST);

	function reset() {
		updateVLANRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:pencil" class="hover:scale-105" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Update VLAN {vlan.name}</AlertDialog.Title>
			<AlertDialog.Description>
				<div class="grid gap-3">
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<legend class="text-sm font-semibold">VLAN Settings</legend>
						<div>
							<Label>Name</Label>
							<Input bind:value={updateVLANRequest.name} />
						</div>
						<div>
							<Label>MTU</Label>
							<Input bind:value={updateVLANRequest.mtu} type="number" />
						</div>
						<div>
							<Label>Description</Label>
							<Input bind:value={updateVLANRequest.description} />
						</div>
						<div class="flex justify-between">
							<Label>DHCP</Label>
							<Switch bind:checked={updateVLANRequest.dhcpOn} />
						</div>
					</fieldset>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.updateVLAN(updateVLANRequest).then((r) => {
						toast.info(`Update ${vlan.name}`);
					});
					// toast.info(`Update ${vlan.name}`);
					console.log(updateVLANRequest);
					reset();
					close();
				}}
			>
				Update
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

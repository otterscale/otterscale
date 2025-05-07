<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import {
		Nexus,
		type UpdateFabricRequest,
		type Network_Fabric,
		type Network
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';

	let { networks = $bindable(), fabric }: { networks: Network[]; fabric: Network_Fabric } =
		$props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = { id: fabric.id, name: fabric.name } as UpdateFabricRequest;

	let updateFabricRequest = $state(DEFAULT_REQUEST);

	function reset() {
		updateFabricRequest = DEFAULT_REQUEST;
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
			<AlertDialog.Title>Update Fabric {fabric.name}</AlertDialog.Title>
			<AlertDialog.Description>
				<fieldset class="grid items-center gap-2 rounded-lg bg-muted/50 p-4">
					<Label>Name</Label>
					<Input bind:value={updateFabricRequest.name} />
				</fieldset>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.updateFabric(updateFabricRequest)
						.then((r) => {
							toast.success(`Update ${r.name} success`);
							client.listNetworks({}).then((r) => {
								networks = r.networks;
							});
						})
						.catch((e) => {
							toast.error(`Fail to update ${fabric.name}: ${e.toString()}`);
						});
					reset();
					close();
				}}
			>
				Update
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

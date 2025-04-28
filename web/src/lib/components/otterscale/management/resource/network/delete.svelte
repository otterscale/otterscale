<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { toast } from 'svelte-sonner';

	import {
		Nexus,
		type DeleteNetworkRequest,
		type Network_Fabric
	} from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';

	let { fabric }: { fabric: Network_Fabric } = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = { id: fabric.id } as DeleteNetworkRequest;

	let deleteNetworkRequest = $state(DEFAULT_REQUEST);

	function reset() {
		deleteNetworkRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<span class="flex items-center gap-1"
			><Icon icon="ph:trash" class="hover:scale-105" />
			<p>Delete</p></span
		>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Delete {fabric.name}</AlertDialog.Title>
			<AlertDialog.Description class="grid gap-2 rounded-lg bg-muted/50 p-4">
				<p>Are you sure you want to delete {fabric.name}?</p>
				<p>This cannot be undone.</p>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.deleteNetwork(deleteNetworkRequest).then((r) => {
						toast.info(`Delete ${fabric.name}`);
					});
					// toast.info(`Delete ${fabric.name}`);
					console.log(deleteNetworkRequest);
					reset();
					close();
				}}
			>
				Delete
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

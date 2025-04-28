<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		Nexus,
		type DeleteIPRangeRequest,
		type Network_IPRange
	} from '$gen/api/nexus/v1/nexus_pb';

	let {
		ipRange
	}: {
		ipRange: Network_IPRange;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = { id: ipRange.id } as DeleteIPRangeRequest;

	let deleteIPRangeRequest = $state(DEFAULT_REQUEST);

	function reset() {
		deleteIPRangeRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:trash" />
		Delete
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Delete IP Range {ipRange.startIp} - {ipRange.endIp}</AlertDialog.Title>
			<AlertDialog.Description class="rounded-lg bg-muted/50 p-4">
				Are you sure you want to delete IP range {ipRange.startIp} - {ipRange.endIp}? This action
				cannot be undone.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.deleteIPRange(deleteIPRangeRequest).then((r) => {
						toast.info(`Delete ip range ${ipRange.startIp} - ${ipRange.endIp}`);
					});
					// toast.info(`Delete ip range ${ipRange.startIp} - ${ipRange.endIp}`);
					console.log(deleteIPRangeRequest);
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

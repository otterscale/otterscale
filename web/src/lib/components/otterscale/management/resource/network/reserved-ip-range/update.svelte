<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		NetworkService,
		type Network_IPRange,
		type UpdateIPRangeRequest
	} from '$gen/api/network/v1/network_pb';

	let {
		ipRange
	}: {
		ipRange: Network_IPRange;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(NetworkService, transport);

	const DEFAULT_REQUEST = {
		id: ipRange.id,
		startIp: ipRange.startIp,
		endIp: ipRange.endIp,
		comment: ipRange.comment
	} as UpdateIPRangeRequest;

	let updateIPRangeRequest = $state(DEFAULT_REQUEST);

	function reset() {
		updateIPRangeRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:pencil" />
		<p>Edit</p>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>
				Update {ipRange.startIp} - {ipRange.endIp}
			</AlertDialog.Title>
			<AlertDialog.Description>
				<div class="grid gap-3"></div>
				<fieldset class="grid items-center gap-3 rounded-lg border p-3">
					<legend class="text-sm font-semibold">Basic Settings</legend>
					<div>
						<Label>Start IP</Label>
						<Input bind:value={updateIPRangeRequest.startIp} />
					</div>
					<div>
						<Label>End IP</Label>
						<Input bind:value={updateIPRangeRequest.endIp} />
					</div>
				</fieldset>

				<fieldset class="grid items-center gap-3 rounded-lg border p-3">
					<legend class="text-sm font-semibold">Advanced Settings</legend>
					<div>
						<Label>Coment</Label>
						<Input bind:value={updateIPRangeRequest.comment} />
					</div>
				</fieldset>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.updateIPRange(updateIPRangeRequest)
						.then((r) => {
							toast.success(`Update ip range ${ipRange.startIp} - ${ipRange.endIp}`);
						})
						.catch((e) => {
							toast.error(`Update ip range ${ipRange.startIp} - ${ipRange.endIp} fail`);
						});
					console.log(updateIPRangeRequest);
					reset();
					close();
				}}
			>
				Update
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

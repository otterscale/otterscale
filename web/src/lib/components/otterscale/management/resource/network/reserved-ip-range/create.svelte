<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	import {
		Nexus,
		type CreateIPRangeRequest,
		type Network_Subnet
	} from '$gen/api/nexus/v1/nexus_pb';

	let { subnet }: { subnet: Network_Subnet } = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = { subnetId: subnet.id } as CreateIPRangeRequest;

	let createIPRangeRequest = $state(DEFAULT_REQUEST);

	function reset() {
		createIPRangeRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:plus" />
		<p>IP Range</p>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title><div class="text-center">Create IP Range</div></AlertDialog.Title>
			<AlertDialog.Description>
				<div class="grid gap-3"></div>
				<fieldset class="grid items-center gap-3 rounded-lg border p-3">
					<legend class="text-sm font-semibold">Basic Settings</legend>
					<div>
						<Label>Subnet ID</Label>
						<Input bind:value={createIPRangeRequest.subnetId} />
					</div>
					<div>
						<Label>Start IP</Label>
						<Input bind:value={createIPRangeRequest.startIp} />
					</div>
					<div>
						<Label>End IP</Label>
						<Input bind:value={createIPRangeRequest.endIp} />
					</div>
				</fieldset>

				<fieldset class="grid items-center gap-3 rounded-lg border p-3">
					<legend class="text-sm font-semibold">Advanced Settings</legend>
					<div>
						<Label>Coment</Label>
						<Input bind:value={createIPRangeRequest.comment} />
					</div>
				</fieldset>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.createIPRange(createIPRangeRequest).then((r) => {
						toast.info(`Create reserved ip range to ${subnet.cidr}`);
					});
					// toast.info(`Create reserved ip range to ${subnet.cidr}`);
					console.log(createIPRangeRequest);
					reset();
					close();
				}}
			>
				Add
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

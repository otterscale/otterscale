<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Select from '$lib/components/ui/select/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { toast } from 'svelte-sonner';
	import { Nexus, type Machine, type DeleteMachineRequest } from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';

	let {
		machine
	}: {
		machine: Machine;
	} = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		id: machine.id,
		force: false
	} as DeleteMachineRequest;

	let deleteMachineRequest = $state(DEFAULT_REQUEST);

	function reset() {
		deleteMachineRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:trash" />
		Delete Machine
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Remove {machine.fqdn}</AlertDialog.Title>
			<AlertDialog.Description class="grid gap-4 rounded-lg bg-muted/50 p-2">
				<p class="px-2">
					This action will permanently delete the machine {machine.fqdn}. This action cannot be
					undone.
				</p>
				<div class="ml-auto">
					<span class="flex justify-between space-x-4">
						<legent>Force</legent>
						<Switch bind:checked={deleteMachineRequest.force} />
					</span>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client.deleteMachine(deleteMachineRequest).then((r) => {
						toast.info(`Delete ${machine.fqdn}`);
					});
					// toast.info(`Delete ${machine.fqdn}`);
					console.log(deleteMachineRequest);
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { toast } from 'svelte-sonner';
	import { Nexus, type Machine, type PowerOffMachineRequest } from '$gen/api/nexus/v1/nexus_pb';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';

	let {
		machine = $bindable(),
		disabled
	}: {
		machine: Machine;
		disabled: boolean;
	} = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = { id: machine.id } as PowerOffMachineRequest;

	let powerOffMachineRequest = $state(DEFAULT_REQUEST);

	function reset() {
		powerOffMachineRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger
		class="flex items-center gap-1 disabled:pointer-events-auto disabled:cursor-not-allowed disabled:text-muted-foreground"
		{disabled}
	>
		<Icon icon="ph:power" />
		Power Off
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>
				Power Off {machine.fqdn}
			</AlertDialog.Title>
			<AlertDialog.Description class="grid gap-2 rounded-lg  bg-muted/50 p-4">
				<p>Are you sure you want to turn on this machine?</p>
				<div class="my-2 grid gap-2">
					<Label>Comment</Label>
					<Input bind:value={powerOffMachineRequest.comment} />
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					toast.promise(() => client.powerOffMachine(powerOffMachineRequest), {
						loading: 'Loading...',
						success: (r) => {
							client.getMachine({ id: machine.id }).then((r) => {
								machine = r;
							});
							return `Turn off ${machine.fqdn}`;
						},
						error: (e) => {
							let msg = `Fail to turn off ${machine.fqdn}`;
							toast.error(msg, {
								description: (e as ConnectError).message.toString(),
								duration: Number.POSITIVE_INFINITY
							});
							return msg;
						}
					});

					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

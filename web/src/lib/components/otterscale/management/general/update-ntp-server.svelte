<script lang="ts">
	import { getContext } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { toast } from 'svelte-sonner';

	import {
		Nexus,
		type Configuration,
		type UpdateNTPServerRequest
	} from '$gen/api/nexus/v1/nexus_pb';

	let { configuration = $bindable() }: { configuration: Configuration } = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = { addresses: [] as string[] } as UpdateNTPServerRequest;

	let updateNTPServerRequest = $state(DEFAULT_REQUEST);

	function initiate(configuration: Configuration) {
		if (configuration.ntpServer) {
			updateNTPServerRequest.addresses = configuration.ntpServer.addresses;
		} else {
			updateNTPServerRequest.addresses = DEFAULT_REQUEST.addresses;
		}
	}

	function reset() {
		updateNTPServerRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger
		onclick={() => {
			initiate(configuration);
		}}
	>
		<Button variant="ghost">
			<Icon icon="ph:pencil" /> Edit
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Update NTP Servers</AlertDialog.Title>
			<AlertDialog.Description class="flex flex-col gap-2 rounded-lg border p-4">
				<span class="flex flex-wrap items-center gap-2">
					{#each updateNTPServerRequest.addresses as address}
						<Badge
							variant="secondary"
							class="flex gap-1 text-sm hover:cursor-pointer"
							onclick={() => {
								updateNTPServerRequest.addresses = updateNTPServerRequest.addresses.filter(
									(_, i) => i !== updateNTPServerRequest.addresses.indexOf(address)
								);
							}}
						>
							{address}
							<Icon icon="ph:x" class="h-3 w-3" />
						</Badge>
					{/each}
				</span>
				<div class="flex items-center justify-between gap-3">
					<Input
						onkeydown={(e) => {
							if (e.key === 'Enter') {
								updateNTPServerRequest.addresses = [
									...updateNTPServerRequest.addresses,
									e.currentTarget.value
								];
								e.currentTarget.value = '';
							}
						}}
					/>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.updateNTPServer(updateNTPServerRequest)
						.then((r) => {
							toast.success(`Update NTP Servers`);
							client.getConfiguration({}).then((r) => {
								configuration = r;
							});
						})
						.catch((e) => {
							toast.error(`Update NTP Servers fail`);
						});
					console.log(updateNTPServerRequest);
					reset();
					close();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

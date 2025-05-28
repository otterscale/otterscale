<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { toast } from 'svelte-sonner';
	import { Nexus, type CreateNetworkRequest, type Network } from '$gen/api/nexus/v1/nexus_pb';
	import { ConnectError, createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';

	let { networks = $bindable() }: { networks: Network[] } = $props();

	const transport: Transport = getContext('transport');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = { dhcpOn: true, dnsServers: [] as string[] } as CreateNetworkRequest;

	let createNetworkRequest = $state(DEFAULT_REQUEST);

	function reset() {
		createNetworkRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger>
		<Button class="flex items-center gap-1">
			<Icon icon="ph:plus" /> Network
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Description>
				<div class="grid gap-3">
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<legend class="text-sm font-semibold">VLAN Settings</legend>
						<div class="flex justify-between">
							<Label>DHCP</Label>
							<Switch bind:checked={createNetworkRequest.dhcpOn} />
						</div>
					</fieldset>
					<fieldset class="grid items-center gap-3 rounded-lg border p-3">
						<legend class="text-sm font-semibold">Subnet Settings</legend>
						<Label>CIDR</Label>
						<Input bind:value={createNetworkRequest.cidr} />
						<Label>Gateway IP</Label>
						<Input bind:value={createNetworkRequest.gatewayIp} />
						<Label>DNS Servers</Label>
						<div class="flex flex-col justify-between gap-1">
							<span class="flex items-center gap-1">
								{#each createNetworkRequest.dnsServers as dnsServer}
									<Badge
										class="flex w-fit items-center gap-1 text-sm hover:cursor-pointer"
										variant="outline"
										onclick={() => {
											createNetworkRequest.dnsServers = createNetworkRequest.dnsServers.filter(
												(s) => s !== dnsServer
											);
										}}
									>
										{dnsServer}
										<Icon icon="ph:x" class="p-0 " />
									</Badge>
								{/each}
							</span>
							<Input
								placeholder="Add DNS Server"
								onkeydown={(e) => {
									if (e.key === 'Enter' && e.currentTarget.value.trim()) {
										createNetworkRequest.dnsServers = [
											...(createNetworkRequest.dnsServers || []),
											e.currentTarget.value.trim()
										];
										e.currentTarget.value = '';
									}
								}}
							/>
						</div>
					</fieldset>
				</div>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					toast.promise(() => client.createNetwork(createNetworkRequest), {
						loading: 'Loading...',
						success: (r) => {
							client.listNetworks({}).then((r) => {
								networks = r.networks;
							});
							return `Create ${createNetworkRequest.cidr} success`;
						},
						error: (e) => {
							let msg = `Fail to create ${createNetworkRequest.cidr}`;
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
				Create
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

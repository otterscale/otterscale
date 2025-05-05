<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Nexus, type Network_Subnet, type UpdateSubnetRequest } from '$gen/api/nexus/v1/nexus_pb';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	let { subnet }: { subnet: Network_Subnet } = $props();

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const DEFAULT_REQUEST = {
		id: subnet.id,
		name: subnet.name,
		cidr: subnet.cidr,
		gatewayIp: subnet.gatewayIp,
		dnsServers: subnet.dnsServers,
		description: subnet.description,
		allowDnsResolution: subnet.allowDnsResolution
	} as UpdateSubnetRequest;

	let updateSubnetRequest = $state(DEFAULT_REQUEST);

	function reset() {
		updateSubnetRequest = DEFAULT_REQUEST;
	}

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class="flex items-center gap-1">
		<Icon icon="ph:pencil" />
		Edit
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Update Subnet {subnet.cidr}</AlertDialog.Title>
			<AlertDialog.Description>
				<div class="grid gap-3"></div>
				<fieldset class="grid items-center gap-3 rounded-lg border p-3">
					<legend class="text-sm font-semibold">Basic Settings</legend>
					<div>
						<Label>Name</Label>
						<Input bind:value={updateSubnetRequest.name} />
					</div>
					<div>
						<Label>Description</Label>
						<Input bind:value={updateSubnetRequest.description} />
					</div>
				</fieldset>

				<fieldset class="grid items-center gap-3 rounded-lg border p-3">
					<legend class="text-sm font-semibold">Network Settings</legend>
					<div>
						<Label>CIDR</Label>
						<Input bind:value={updateSubnetRequest.cidr} />
					</div>
					<div>
						<Label>Gateway IP</Label>
						<Input bind:value={updateSubnetRequest.gatewayIp} />
					</div>
				</fieldset>

				<fieldset class="grid items-center gap-3 rounded-lg border p-3">
					<legend class="text-sm font-semibold">DNS Settings</legend>
					<div class="flex flex-col justify-between gap-1">
						<span class="flex items-center gap-1">
							{#each updateSubnetRequest.dnsServers as dnsServer}
								<Badge
									class="flex w-fit items-center gap-1 text-sm hover:cursor-pointer"
									variant="outline"
									onclick={() => {
										updateSubnetRequest.dnsServers = updateSubnetRequest.dnsServers.filter(
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
									updateSubnetRequest.dnsServers = [
										...(updateSubnetRequest.dnsServers || []),
										e.currentTarget.value.trim()
									];
									e.currentTarget.value = '';
								}
							}}
						/>
					</div>
					<div class="flex justify-between">
						<Label>Allow DNS Resolution</Label>
						<Switch bind:checked={updateSubnetRequest.allowDnsResolution} />
					</div>
				</fieldset>
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel onclick={reset} class="mr-auto">Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					client
						.updateSubnet(updateSubnetRequest)
						.then((r) => {
							toast.info(`Update ${r.cidr} success`);
						})
						.catch((e) => {
							toast.error(`Fail to update ${subnet.cidr}: ${e.toString()}`);
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

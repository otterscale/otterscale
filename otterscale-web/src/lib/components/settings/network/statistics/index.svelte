<script lang="ts">
	import { type Network } from '$lib/api/network/v1/network_pb';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatBigNumber as formatNumber, formatHealthColor } from '$lib/formatter';

	let {
		networks
	}: {
		networks: Network[];
	} = $props();

	// Calculate statistics
	const fabricCount = $derived(networks.filter((n) => n).map((n) => n.fabric).length);
	const vlanList = $derived(networks.filter((n) => n).map((n) => n.vlan));
	const dhcpEnabledCount = $derived(
		networks.reduce((a, network) => a + (network.vlan?.dhcpOn ? 1 : 0), 0)
	);
	const availableSubnetCount = $derived(
		networks.reduce(
			(a, network) =>
				a + Number(network.subnet?.statistics ? network.subnet.statistics.available : 0),
			0
		)
	);
	const totalSubnetCount = $derived(
		networks.reduce(
			(a, network) => a + Number(network.subnet?.statistics ? network.subnet.statistics.total : 0),
			0
		)
	);
	const availabilityPercentage = $derived((availableSubnetCount * 100) / totalSubnetCount || 0);
	const dhcpEnabledPercentage = $derived((dhcpEnabledCount * 100) / fabricCount || 0);
</script>

<div class="grid grid-cols-4 gap-3">
	<Card.Root>
		<Card.Header class="h-10">
			<Card.Title>NETWORK</Card.Title>
		</Card.Header>
		<Card.Content class="h-30">
			<span class="flex items-end gap-1">
				<p class="text-7xl">{fabricCount}</p>
				<p>fabrics</p>
			</span>
		</Card.Content>
		<Card.Footer>
			<div class="flex flex-wrap gap-1">
				{#each [...new Set(vlanList.map((v) => v?.name))] as vlanName}
					<Badge variant="outline">
						{vlanName}: {networks.reduce(
							(a, network) => a + (network.vlan?.name === vlanName ? 1 : 0),
							0
						)}
					</Badge>
				{/each}
			</div>
		</Card.Footer>
	</Card.Root>
	<Card.Root>
		<Card.Header class="h-10">
			<Card.Title>DHCP</Card.Title>
		</Card.Header>
		<Card.Content class="h-30">
			<p class="text-3xl">
				{Math.round(dhcpEnabledPercentage)}%
			</p>
			<p class="text-muted-foreground text-xs">
				{dhcpEnabledCount} Enabled over {fabricCount} fabrics
			</p>
		</Card.Content>
		<Card.Footer>
			<Progress value={dhcpEnabledPercentage} max={100} />
		</Card.Footer>
	</Card.Root>
	<Card.Root>
		<Card.Header class="h-10">
			<Card.Title>AVAILABILITY</Card.Title>
		</Card.Header>
		<Card.Content class="h-30">
			<p class="text-3xl">
				{Math.round(availabilityPercentage)}%
			</p>
			<p class="text-muted-foreground text-xs">
				{formatNumber(availableSubnetCount)} available over {formatNumber(totalSubnetCount)} units
			</p>
		</Card.Content>
		<Card.Footer>
			<Progress
				value={availabilityPercentage}
				max={100}
				class={formatHealthColor(availabilityPercentage)}
			/>
		</Card.Footer>
	</Card.Root>
</div>

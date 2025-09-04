<script lang="ts">
	import { type Network } from '$lib/api/network/v1/network_pb';
	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import ContentSubtitle from '$lib/components/custom/chart/content/text/text-with-subtitle.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatBigNumber as formatNumber, formatProgressColor } from '$lib/formatter';

	let {
		networks,
	}: {
		networks: Network[];
	} = $props();

	// Calculate statistics
	const fabricCount = $derived(networks.map((n) => n.fabric).length);
	const isDhcpEnabled = $derived(networks.some((n) => n.vlan?.dhcpOn === true));
	const dhcpOnFabricName = $derived(
		networks
			.filter((n) => n.vlan?.dhcpOn === true)
			.map((n) => n.fabric?.name)
			.filter(Boolean)
			.join(', '),
	);
	const availableSubnetCount = $derived(
		networks.reduce(
			(a, network) => a + Number(network.subnet?.statistics ? network.subnet.statistics.available : 0),
			0,
		),
	);
	const totalSubnetCount = $derived(
		networks.reduce(
			(a, network) => a + Number(network.subnet?.statistics ? network.subnet.statistics.total : 0),
			0,
		),
	);
	const availabilityPercentage = $derived((availableSubnetCount * 100) / totalSubnetCount || 0);
</script>

<div class="grid w-full gap-3 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-4 xl:grid-cols-5">
	<Layout>
		{#snippet title()}
			<Title title="NETWORK" />
		{/snippet}

		{#snippet content()}
			<Content value={fabricCount} />
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="DHCP" />
		{/snippet}

		{#snippet content()}
			<Content value={isDhcpEnabled ? 'ON' : 'OFF'} />
		{/snippet}

		{#snippet footer()}
			<p class="text-muted-foreground text-xs">
				{dhcpOnFabricName}
			</p>
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="AVAILABILITY" />
		{/snippet}

		{#snippet content()}
			<ContentSubtitle
				value={Math.round(availabilityPercentage)}
				unit="%"
				subtitle={`${formatNumber(availableSubnetCount)} available over ${formatNumber(totalSubnetCount)} units`}
			/>
		{/snippet}

		{#snippet footer()}
			<Progress value={availabilityPercentage} max={100} class={formatProgressColor(availabilityPercentage)} />
		{/snippet}
	</Layout>
</div>

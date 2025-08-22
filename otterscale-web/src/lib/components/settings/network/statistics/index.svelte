<script lang="ts">
	import { type Network } from '$lib/api/network/v1/network_pb';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import { formatBigNumber as formatNumber, formatHealthColor } from '$lib/formatter';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Content from '$lib/components/custom/chart/content/text/text-large.svelte';
	import ContentSubtitle from '$lib/components/custom/chart/content/text/text-with-subtitle.svelte';
	import Layout from '$lib/components/custom/chart/layout/small-flexible-height.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';

	let {
		networks
	}: {
		networks: Network[];
	} = $props();

	// Calculate statistics
	const fabricCount = $derived(networks.map((n) => n.fabric).length);
	const vlanList = $derived(networks.map((n) => n.vlan));
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
			<ContentSubtitle
				value={Math.round(dhcpEnabledPercentage)}
				unit={'%'}
				subtitle={`${dhcpEnabledCount} Enabled over ${fabricCount} fabrics`}
			/>
		{/snippet}

		{#snippet footer()}
			<Progress
				value={dhcpEnabledPercentage}
				max={100}
				class={formatHealthColor(dhcpEnabledPercentage)}
			/>
		{/snippet}
	</Layout>

	<Layout>
		{#snippet title()}
			<Title title="AVAILABILITY" />
		{/snippet}

		{#snippet content()}
			<ContentSubtitle
				value={Math.round(availabilityPercentage)}
				unit={'%'}
				subtitle={`${formatNumber(availableSubnetCount)} available over ${formatNumber(totalSubnetCount)} units`}
			/>
		{/snippet}

		{#snippet footer()}
			<Progress
				value={availabilityPercentage}
				max={100}
				class={formatHealthColor(availabilityPercentage)}
			/>
		{/snippet}
	</Layout>
</div>

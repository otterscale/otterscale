<script lang="ts">
	import Icon from '@iconify/svelte';
	import { PrometheusDriver } from 'prometheus-query';

	import * as Statistics from '$lib/components/custom/data-table/statistics/index';
	import { formatDuration } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	let { client, fqdn }: { client: PrometheusDriver; fqdn: string } = $props();

	const query = $derived(
		`
		node_time_seconds{instance=~"${fqdn}"}
		-
		node_boot_time_seconds{instance=~"${fqdn}"}
		`
	);
</script>

<Statistics.Root type="count">
	<Statistics.Header>
		<Statistics.Title>{m.uptime()}</Statistics.Title>
	</Statistics.Header>
	<Statistics.Content class="min-h-20">
		{#await client.instantQuery(query)}
			<div class="flex h-[200px] w-full items-center justify-center">
				<Icon icon="svg-spinners:3-dots-bounce" class="m-8 size-8" />
			</div>
		{:then response}
			{@const result = response.result}
			{#if result.length === 0}
				<div class="flex h-[200px] w-full flex-col items-center justify-center">
					<Icon icon="ph:chart-bar-fill" class="size-24 animate-pulse text-muted-foreground" />
					<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
				</div>
			{:else}
				{@const uptime = result[0].value.value}
				{@const duration = formatDuration(uptime)}
				<p class="flex h-[200px] items-center justify-center text-5xl font-semibold">
					{duration.value.toFixed(1)}
					{duration.unit}
				</p>
			{/if}
		{:catch}
			<div class="flex h-[200px] w-full flex-col items-center justify-center">
				<Icon icon="ph:chart-bar-fill" class="size-24 animate-pulse text-muted-foreground" />
				<p class="text-base text-muted-foreground">{m.no_data_display()}</p>
			</div>
		{/await}
	</Statistics.Content>
</Statistics.Root>

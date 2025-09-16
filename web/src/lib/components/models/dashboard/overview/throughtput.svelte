<script lang="ts">
	import { scaleUtc } from 'd3-scale';
	import { curveNatural } from 'd3-shape';
	import { Area, AreaChart, LinearGradient } from 'layerchart';
	import { PrometheusDriver, SampleValue } from 'prometheus-query';
	import { onMount } from 'svelte';

	import type { Scope } from '$lib/api/scope/v1/scope_pb';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Card from '$lib/components/ui/card';
	import * as Chart from '$lib/components/ui/chart';
	import { m } from '$lib/paraglide/messages';

	let {
		prometheusDriver,
		scope,
		isReloading = $bindable(),
	}: { prometheusDriver: PrometheusDriver; scope: Scope; isReloading: boolean } = $props();

	const thyroughputsConfiguration = {
		prompt: { label: 'Prompt', color: 'var(--chart-1)' },
		generation: { label: 'Generation', color: 'var(--chart-2)' },
	} satisfies Chart.ChartConfig;

	let prompts = $state([] as SampleValue[]);
	let generations = $state([] as SampleValue[]);
	const requests = $derived(
		prompts.map((sample, index) => ({
			time: sample.time,
			running: sample.value,
			waiting: generations[index]?.value ?? 0,
		})),
	);

	async function fetch() {
		prometheusDriver
			.rangeQuery(
				`sum(rate(vllm:prompt_tokens_total{scope_uuid="${scope.uuid}"}[2m]))`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				prompts = response.result[0]?.values;
			});
		prometheusDriver
			.rangeQuery(
				`sum(rate(vllm:generation_tokens_total{scope_uuid="${scope.uuid}"}[2m]))`,
				Date.now() - 24 * 60 * 60 * 1000,
				Date.now(),
				2 * 60,
			)
			.then((response) => {
				generations = response.result[0]?.values;
			});
	}

	const reloadManager = new ReloadManager(fetch);

	let isLoading = $state(true);
	onMount(async () => {
		await fetch();
		isLoading = false;
	});

	$effect(() => {
		if (isReloading) {
			reloadManager.restart();
		} else {
			reloadManager.stop();
		}
	});
</script>

{#if isLoading}
	Loading
{:else}
	<Card.Root class="h-full">
		<Card.Header>
			<Card.Title>{m.throughput()}</Card.Title>
			<Card.Description>{m.llm_dashboard_throighputs_tooltip()}</Card.Description>
		</Card.Header>
		<Card.Content>
			<Chart.Container config={thyroughputsConfiguration} class="h-[200px] w-full">
				<AreaChart
					data={requests}
					x="time"
					xScale={scaleUtc()}
					yPadding={[0, 25]}
					series={[
						{
							key: 'prompt',
							label: thyroughputsConfiguration.prompt.label,
							color: thyroughputsConfiguration.prompt.color,
						},
						{
							key: 'generation',
							label: thyroughputsConfiguration.generation.label,
							color: thyroughputsConfiguration.generation.color,
						},
					]}
					seriesLayout="stack"
					props={{
						area: {
							curve: curveNatural,
							'fill-opacity': 0.4,
							line: { class: 'stroke-1' },
							motion: 'tween',
						},
						xAxis: {
							format: (v: Date) => v.toLocaleDateString('en-US', { month: 'short' }),
						},
						yAxis: { format: () => '' },
					}}
				>
					{#snippet tooltip()}
						<Chart.Tooltip
							indicator="dot"
							labelFormatter={(v: Date) => {
								return v.toLocaleDateString('en-US', {
									year: 'numeric',
									month: 'short',
									day: 'numeric',
									hour: 'numeric',
									minute: 'numeric',
								});
							}}
						/>
					{/snippet}
					{#snippet marks({ series, getAreaProps })}
						{#each series as s, i (s.key)}
							<LinearGradient
								stops={[s.color ?? '', 'color-mix(in lch, ' + s.color + ' 10%, transparent)']}
								vertical
							>
								{#snippet children({ gradient })}
									<Area {...getAreaProps(s, i)} fill={gradient} />
								{/snippet}
							</LinearGradient>
						{/each}
					{/snippet}
				</AreaChart>
			</Chart.Container>
		</Card.Content>
	</Card.Root>
{/if}

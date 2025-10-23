<script lang="ts">
	import { type Table } from '@tanstack/table-core';

	import type { TestResult } from '$lib/api/configuration/v1/configuration_pb';
	import Content from '$lib/components/custom/chart/content/scatter/scatter.svelte';
	import Description from '$lib/components/custom/chart/description.svelte';
	import Title from '$lib/components/custom/chart/title.svelte';
	import { BistDashboardManager } from '$lib/components/settings/built-in-self-test/utils/bistManager';
	import * as Card from '$lib/components/ui/card/index.js';

	let { mode, table }: { mode: string; table: Table<TestResult> } = $props();
	const dashboardManager = new BistDashboardManager(table);
	const { get: getTmp, put: putTmp, delete: deleteTmp } = dashboardManager.getWarpOutputs();
</script>

<div class="mt-4 grid w-full gap-4 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-3">
	<Card.Root>
		<Card.Header class="h-[10px]">
			<Card.Title>
				<Title title="Throughput" />
			</Card.Title>
			<Card.Description>
				<Description description="Fastest" />
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Content
				xKey="bytesFastest"
				yKey="completedAt"
				data={mode === 'get' ? getTmp : mode === 'put' ? putTmp : deleteTmp}
			/>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header class="h-[10px]">
			<Card.Title>
				<Title title="Throughput" />
			</Card.Title>
			<Card.Description>
				<Description description="Slowest" />
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Content
				xKey="bytesSlowest"
				yKey="completedAt"
				data={mode === 'get' ? getTmp : mode === 'put' ? putTmp : deleteTmp}
			/>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header class="h-[10px]">
			<Card.Title>
				<Title title="Throughput" />
			</Card.Title>
			<Card.Description>
				<Description description="Median" />
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Content
				xKey="bytesMedian"
				yKey="completedAt"
				data={mode === 'get' ? getTmp : mode === 'put' ? putTmp : deleteTmp}
			/>
		</Card.Content>
	</Card.Root>
</div>

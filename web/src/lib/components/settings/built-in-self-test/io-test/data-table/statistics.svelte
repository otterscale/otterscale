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
	const { read: readTmp, write: writeTmp, trim: trimTmp } = dashboardManager.getFioOutputs();
</script>

<div class="mt-4 grid w-full gap-4 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-3">
	<Card.Root>
		<Card.Header class="h-[10px]">
			<Card.Title>
				<Title title="Bandwidth" />
			</Card.Title>
			<Card.Description>
				<Description description="Bandwidth Bytes" />
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Content
				xKey="bandwidthBytes"
				yKey="completedAt"
				data={mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp}
			/>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header class="h-[10px]">
			<Card.Title>
				<Title title="IOPS" />
			</Card.Title>
			<Card.Description>
				<Description description="IO Per Second" />
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Content
				xKey="ioPerSecond"
				yKey="completedAt"
				data={mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp}
			/>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header class="h-[10px]">
			<Card.Title>
				<Title title="Latency" />
			</Card.Title>
			<Card.Description>
				<Description description="Mean Latency" />
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Content
				xKey="latency"
				yKey="completedAt"
				data={mode === 'read' ? readTmp : mode === 'write' ? writeTmp : trimTmp}
			/>
		</Card.Content>
	</Card.Root>
</div>

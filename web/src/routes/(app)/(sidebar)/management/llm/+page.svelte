<script lang="ts">
	import * as Alert from '$lib/components/ui/alert/index.js';
	import { llmData, createDateSeries, type LLMModel } from './dataset';
	import { formatBigNumber as formatNumber } from '$lib/formatter';
	import * as Pagination from '$lib/components/ui/pagination/index.js';
	import * as Table from '$lib/components/ui/table';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import { PageLoading } from '$lib/components/otterscale/ui/index';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ManagementMachines } from '$lib/components/otterscale';
	import { Nexus, type Machine, type Network } from '$gen/api/nexus/v1/nexus_pb';
	import { getContext, onMount } from 'svelte';
	import { writable } from 'svelte/store';
	import { AreaChart, BarChart, LineChart, PieChart } from 'layerchart';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button } from '$lib/components/ui/button';
	import { Progress } from '$lib/components/ui/progress/index.js';
	import Icon from '@iconify/svelte';

	const transport: Transport = getContext('transportNEW');
	const client = createClient(Nexus, transport);

	const networksStore = writable<Network[]>([]);
	const networksLoading = writable(true);
	async function fetchNetworks() {
		try {
			const response = await client.listNetworks({});
			networksStore.set(response.networks);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			networksLoading.set(false);
		}
	}

	const machinesStore = writable<Machine[]>([]);
	const machinesLoading = writable(true);
	async function fetchMachines() {
		try {
			const response = await client.listMachines({});
			machinesStore.set(response.machines);
		} catch (error) {
			console.error('Error fetching:', error);
		} finally {
			machinesLoading.set(false);
		}
	}

	console.log(
		createDateSeries({
			count: 100,
			min: 0,
			max: 100,
			value: 'number',
			keys: ['value']
		})
	);

	let mounted = false;
	onMount(async () => {
		try {
			await fetchNetworks();
			await fetchMachines();
		} catch (error) {
			console.error('Error during initial data load:', error);
		}

		mounted = true;
	});
</script>

{#snippet Statistic()}
	{@const numberOfModels = llmData.length}
	{@const averageAccuracy =
		llmData.reduce((acc, model) => acc + model.metrics.accuracy, 0) / llmData.length}
	{@const maxRequests = Math.max(...llmData.map((model) => model.usageStats.requests))}
	<Card.Root>
		<Card.Header>
			<Card.Title>Models</Card.Title>
		</Card.Header>
		<Card.Content class="text-3xl">
			{numberOfModels}
		</Card.Content>
		<Card.Footer class="flex flex-col items-start"></Card.Footer>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title>Average Accuracy</Card.Title>
		</Card.Header>
		<Card.Content class="text-3xl">
			{Math.round(averageAccuracy * 100)} %
		</Card.Content>
		<Card.Footer class="flex flex-col items-start">
			<Progress value={averageAccuracy * 100} max={100} />
		</Card.Footer>
	</Card.Root>
	<Card.Root>
		<Card.Header>
			<Card.Title>Max Requests</Card.Title>
		</Card.Header>
		<Card.Content class="text-3xl">
			{formatNumber(maxRequests)}
		</Card.Content>
		<Card.Footer class="flex flex-col items-start"></Card.Footer>
	</Card.Root>
{/snippet}
<div class="flex-col space-y-4">
	<Alert.Root variant="destructive">
		<Icon icon="ph:warning-diamond" class="size-5" />
		<Alert.Title>Work In Progress</Alert.Title>
		<Alert.Description>This LLM management page is still under development.</Alert.Description>
	</Alert.Root>
	<div class="grid grid-cols-4 gap-3">
		{@render Statistic()}
	</div>
</div>

<div class="flex items-center justify-between">
	<span>
		<Input type="text" placeholder="Search models..." />
	</span>

	<span class="my-4 flex items-center gap-2">
		<Button>Filter</Button>
		<Select.Root type="multiple">
			<Select.Trigger class="w-[180px]">Select</Select.Trigger>
			<Select.Content>
				<Select.Item value="all">All Architectures</Select.Item>
				<Select.Item value="transformer">Transformer</Select.Item>
				<Select.Item value="lstm">LSTM</Select.Item>
				<Select.Item value="gpt">GPT</Select.Item>
			</Select.Content>
		</Select.Root>
	</span>
</div>
<Table.Root>
	<Table.Header>
		<Table.Row class="*:text-xs *:font-light">
			<Table.Head>NAME</Table.Head>
			<Table.Head>VERSION</Table.Head>
			<Table.Head>ARCHITECTURE</Table.Head>
			<Table.Head class="text-right">PARAMETERS</Table.Head>
			<Table.Head class="text-right">ACCURACY</Table.Head>
			<Table.Head class="text-right">SPEED</Table.Head>
			<Table.Head class="text-right">REQUESTS</Table.Head>
			<Table.Head class="text-right">LATENCY</Table.Head>
			<Table.Head class="text-right">THROUGHPUT</Table.Head>
		</Table.Row>
	</Table.Header>
	<Table.Body class="*:text-sm">
		{#each llmData as model}
			{@const latency = createDateSeries({
				count: 23,
				min: 0,
				max: 10000,
				value: 'number',
				keys: ['value', 'average']
			})}
			{@const throughput = createDateSeries({
				count: 13,
				min: 0,
				max: 1000000,
				value: 'number',
				keys: ['value', 'average']
			})}
			<Table.Row>
				<Table.Cell>{model.name}</Table.Cell>
				<Table.Cell><Badge variant="outline">{model.version}</Badge></Table.Cell>
				<Table.Cell><Badge variant="outline">{model.architecture}</Badge></Table.Cell>
				<Table.Cell class="text-right">{model.parameters}</Table.Cell>
				<Table.Cell class="text-right">{Math.round(model.metrics.accuracy * 100)}%</Table.Cell>
				<Table.Cell class="text-right">{model.metrics.speed}</Table.Cell>
				<Table.Cell class="text-right">{formatNumber(model.usageStats.requests)}</Table.Cell>
				<Table.Cell class="gap-1 text-right align-top">
					<p class="text-xs font-light">{latency[latency.length - 1].average.toFixed(0)} ms</p>
					<span class="inline-block h-[23px] w-[130px]">
						<LineChart data={latency} x="date" y="value" axis={false} />
					</span>
				</Table.Cell>
				<Table.Cell class="gap-1 text-right align-top">
					<p class="text-xs font-light">
						{throughput[throughput.length - 1].average.toFixed(0)} tpm
					</p>
					<span class="inline-block h-[23px] w-[130px]">
						<LineChart data={throughput} x="date" y="value" axis={false} />
					</span>
				</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>
<Pagination.Root count={10} perPage={10}>
	{#snippet children({ pages })}
		<Pagination.Content class="rounded-lg bg-muted p-1">
			<Pagination.Item>
				<Pagination.PrevButton />
			</Pagination.Item>
			{#each pages as page (page.key)}
				{#if page.type === 'ellipsis'}
					<Pagination.Item>
						<Pagination.Ellipsis />
					</Pagination.Item>
				{:else}
					<Pagination.Item>
						<Pagination.Link {page} isActive={page.value === page.value}>
							{page.value}
						</Pagination.Link>
					</Pagination.Item>
				{/if}
			{/each}
			<Pagination.Item>
				<Pagination.NextButton />
			</Pagination.Item>
		</Pagination.Content>
	{/snippet}
</Pagination.Root>

<script lang="ts">
	import * as Alert from '$lib/components/ui/alert';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import * as Pagination from '$lib/components/ui/pagination';
	import { Progress } from '$lib/components/ui/progress';
	import * as Select from '$lib/components/ui/select';
	import * as Table from '$lib/components/ui/table';
	import { formatBigNumber as formatNumber } from '$lib/formatter';
	import Icon from '@iconify/svelte';
	// import { startOfToday, subDays } from 'date-fns';
	import { LineChart } from 'layerchart';

	function getRandomInteger(min: number, max: number, includeMax = true) {
		min = Math.ceil(min);
		max = Math.floor(max);
		return Math.floor(Math.random() * (max - min + (includeMax ? 1 : 0)) + min);
	}

	function getRandomNumber(min: number, max: number) {
		return Math.random() * (max - min) + min;
	}

	// function createDateSeries<TKey extends string>(options: {
	// 	count?: number;
	// 	min: number;
	// 	max: number;
	// 	keys?: TKey[];
	// 	value?: 'number' | 'integer';
	// }) {
	// 	const now = startOfToday();

	// 	const count = options.count ?? 10;
	// 	const min = options.min;
	// 	const max = options.max;
	// 	const keys = options.keys ?? ['value'];

	// 	return Array.from({ length: count }).map((_, i) => {
	// 		return {
	// 			date: subDays(now, count - i - 1),
	// 			...Object.fromEntries(
	// 				keys.map((key) => {
	// 					return [
	// 						key,
	// 						options.value === 'integer' ? getRandomInteger(min, max) : getRandomNumber(min, max)
	// 					];
	// 				})
	// 			)
	// 		} as { date: Date } & { [K in TKey]: number };
	// 	});
	// }

	interface LLMModel {
		name: string;
		version: string;
		parameters: string;
		metrics: {
			accuracy: number;
			speed: number;
		};
		architecture: string;
		usageStats: {
			requests: number;
			uptime: number;
		};
	}

	export const llmData: LLMModel[] = [
		{
			name: 'GPT-3',
			version: '1.0',
			parameters: '175B',
			metrics: { accuracy: 0.92, speed: 1.2 },
			architecture: 'Transformer',
			usageStats: { requests: 1000000, uptime: 99.9 }
		},
		{
			name: 'BERT',
			version: '2.0',
			parameters: '340M',
			metrics: { accuracy: 0.89, speed: 1.5 },
			architecture: 'Bidirectional Transformer',
			usageStats: { requests: 500000, uptime: 99.5 }
		},
		{
			name: 'LLaMA',
			version: '2.0',
			parameters: '65B',
			metrics: { accuracy: 0.91, speed: 1.3 },
			architecture: 'Transformer',
			usageStats: { requests: 800000, uptime: 99.7 }
		},
		{
			name: 'RoBERTa',
			version: '1.5',
			parameters: '355M',
			metrics: { accuracy: 0.9, speed: 1.4 },
			architecture: 'Bidirectional Transformer',
			usageStats: { requests: 600000, uptime: 99.6 }
		},
		{
			name: 'T5',
			version: '1.1',
			parameters: '11B',
			metrics: { accuracy: 0.88, speed: 1.6 },
			architecture: 'Encoder-Decoder',
			usageStats: { requests: 400000, uptime: 99.3 }
		},
		{
			name: 'BLOOM',
			version: '1.0',
			parameters: '176B',
			metrics: { accuracy: 0.91, speed: 1.1 },
			architecture: 'Transformer',
			usageStats: { requests: 300000, uptime: 99.4 }
		},
		{
			name: 'PaLM',
			version: '2.0',
			parameters: '540B',
			metrics: { accuracy: 0.93, speed: 1.0 },
			architecture: 'Transformer',
			usageStats: { requests: 900000, uptime: 99.8 }
		},
		{
			name: 'Claude',
			version: '2.0',
			parameters: '100B',
			metrics: { accuracy: 0.92, speed: 1.2 },
			architecture: 'Constitutional AI',
			usageStats: { requests: 700000, uptime: 99.6 }
		},
		{
			name: 'Falcon',
			version: '1.0',
			parameters: '40B',
			metrics: { accuracy: 0.89, speed: 1.4 },
			architecture: 'Transformer',
			usageStats: { requests: 200000, uptime: 99.2 }
		},
		{
			name: 'OPT',
			version: '1.3',
			parameters: '175B',
			metrics: { accuracy: 0.9, speed: 1.3 },
			architecture: 'Transformer',
			usageStats: { requests: 450000, uptime: 99.5 }
		}
	];

	const numberOfModels = llmData.length;
	const averageAccuracy =
		llmData.reduce((acc, model) => acc + model.metrics.accuracy, 0) / llmData.length;
	const maxRequests = Math.max(...llmData.map((model) => model.usageStats.requests));
</script>

<div class="flex-col space-y-4">
	<Alert.Root variant="destructive">
		<Icon icon="ph:warning-diamond" class="size-5" />
		<Alert.Title>Work In Progress</Alert.Title>
		<Alert.Description>This LLM management page is still under development.</Alert.Description>
	</Alert.Root>
	<div class="grid grid-cols-4 gap-3">
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
	</div>
</div>
<div class="flex items-center justify-between">
	<span>
		<Input type="text" placeholder="Model" />
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
		</Table.Row>
	</Table.Header>
	<Table.Body class="*:text-sm">
		{#each llmData as model}
			<Table.Row>
				<Table.Cell>{model.name}</Table.Cell>
				<Table.Cell><Badge variant="outline">{model.version}</Badge></Table.Cell>
				<Table.Cell><Badge variant="outline">{model.architecture}</Badge></Table.Cell>
				<Table.Cell class="text-right">{model.parameters}</Table.Cell>
				<Table.Cell class="text-right">{Math.round(model.metrics.accuracy * 100)}%</Table.Cell>
				<Table.Cell class="text-right">{model.metrics.speed}</Table.Cell>
				<Table.Cell class="text-right">{formatNumber(model.usageStats.requests)}</Table.Cell>
			</Table.Row>
		{/each}
	</Table.Body>
</Table.Root>
<Pagination.Root count={10} perPage={10}>
	{#snippet children({ pages })}
		<Pagination.Content class="bg-muted rounded-lg p-1">
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

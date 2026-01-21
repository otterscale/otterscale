<script lang="ts">
	import Braces from '@lucide/svelte/icons/braces';
	import Circle from '@lucide/svelte/icons/circle';
	import File from '@lucide/svelte/icons/file';
	import FileCheck from '@lucide/svelte/icons/file-check';
	import FileClock from '@lucide/svelte/icons/file-clock';
	import FileCode from '@lucide/svelte/icons/file-code';
	import FileDigit from '@lucide/svelte/icons/file-digit';
	import FileText from '@lucide/svelte/icons/file-text';
	import Grid from '@lucide/svelte/icons/grid';
	import List from '@lucide/svelte/icons/list';
	import X from '@lucide/svelte/icons/x';
	import type { JsonObject, JsonValue } from '@openfeature/server-sdk';
	import type { Column, Row } from '@tanstack/table-core';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
	import Monaco from 'svelte-monaco';
	import { stringify } from 'yaml';

	import * as Code from '$lib/components/custom/code/index.js';
	import { Button } from '$lib/components/ui/button';
	import * as Collapsible from '$lib/components/ui/collapsible';
	import * as Empty from '$lib/components/ui/empty/index.js';
	import * as Item from '$lib/components/ui/item';
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { now } from '$lib/stores/now';

	let {
		ref = $bindable(null),
		column,
		row,
		fields,
		children,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		column: Column<Record<string, JsonValue>>;
		row: Row<Record<string, JsonValue>>;
		fields: Record<string, { description: string; type: string; format: string }>;
	} = $props();

	function format(value: string) {
		try {
			return JSON.stringify(JSON.parse(value), null, 4);
		} catch {
			return value;
		}
	}

	function getRelativeTime(now: number, timestamp: number) {
		const milliseconds = timestamp;

		const seconds = Math.floor((now - milliseconds) / 1000);
		if (seconds < 60) return { value: seconds, unit: 'second' };

		const minutes = Math.floor(seconds / 60);
		if (minutes < 60) return { value: minutes, unit: 'minute' };

		const hours = Math.floor(minutes / 60);
		if (hours < 24) return { value: hours, unit: 'hour' };

		const days = Math.floor(hours / 24);
		if (days < 7) return { value: days, unit: 'day' };

		const weeks = Math.floor(days / 7);
		if (weeks < 5) return { value: weeks, unit: 'week' };

		const months = Math.floor(days / 30);
		if (months < 12) return { value: months, unit: 'month' };

		const years = Math.floor(days / 365);
		return { value: years, unit: 'year' };
	}
</script>

<div class={className}>
	{#if children}
		{@render children()}
	{:else if fields[column.id].type === 'object'}
		{@render ObjectCell({ data: row.original[column.id] as JsonObject })}
	{:else if fields[column.id].type === 'array'}
		{@render ArrayCell({ data: row.original[column.id] as JsonValue[] })}
	{:else if fields[column.id].type === 'string' && fields[column.id].format === 'date'}
		{@render DateCell({ data: new Date(String(row.original[column.id])) })}
	{:else if fields[column.id].type === 'string' && fields[column.id].format === 'date-time'}
		{@render DatetimeCell({ data: new Date(String(row.original[column.id])) })}
	{:else if fields[column.id].type === 'number' || fields[column.id].type === 'integer'}
		{@render NumberCell({ data: Number(row.original[column.id]) })}
	{:else if fields[column.id].type === 'boolean'}
		{@render BooleanCell({ data: Boolean(row.original[column.id]) })}
	{:else if row.original[column.id]}
		{@render TextCell({ data: String(row.original[column.id]) })}
	{:else}
		{@render EmptyCell()}
	{/if}
</div>

{#snippet ObjectCell({ data }: { data: JsonObject })}
	<Sheet.Root>
		<Sheet.Trigger>
			{#if data && !Object.values(data).some((value) => value && typeof value === 'object')}
				<Button variant="ghost" class="hover:underline">
					{Object.keys(data).length}
				</Button>
			{:else if data}
				<Button variant="ghost">
					<FileCode />
				</Button>
			{/if}
		</Sheet.Trigger>
		<Sheet.Content
			side="right"
			class="flex h-full max-w-[62vw] min-w-[50vw] flex-col gap-0 overflow-y-auto p-4"
		>
			<Sheet.Header class="shrink-0 space-y-4">
				<Sheet.Title>{column.id}</Sheet.Title>
				<Sheet.Description>
					{fields[column.id].description}
				</Sheet.Description>
			</Sheet.Header>
			{#if data}
				{#if Object.values(data).some((value) => value && typeof value === 'object')}
					<div class="h-full p-4 pt-0">
						<Monaco
							value={stringify(data)}
							options={{
								language: 'yaml',
								padding: { top: 24 },
								automaticLayout: true,
								domReadOnly: true,
								readOnly: true
							}}
							theme="vs-dark"
						/>
					</div>
				{:else}
					<Tabs.Root value="grid" class="p-4">
						<Tabs.List class="ml-auto">
							<Tabs.Trigger value="grid">
								<Grid />
							</Tabs.Trigger>
							<Tabs.Trigger value="table">
								<List />
							</Tabs.Trigger>
						</Tabs.List>
						<Tabs.Content value="grid">
							<div class="space-y-0">
								{#each Object.entries(data) as [key, value], index (index)}
									{#if typeof value === 'string'}
										{@const data = format(value)}
										{@const isExpandable = data.split('\n').length > 2}
										<Collapsible.Root
											class="rounded-lg transition-colors duration-200 hover:bg-muted/50"
										>
											<Collapsible.Trigger
												disabled={!isExpandable}
												class="w-full transition-colors duration-200 hover:underline"
											>
												<Item.Root size="sm">
													<Item.Media variant="icon">
														{#if !value}
															<File />
														{:else if ['true', 'false'].includes(value)}
															<FileCheck />
														{:else if !isNaN(Number(value))}
															<FileDigit />
														{:else if !isNaN(Date.parse(value))}
															<FileClock />
														{:else}
															<FileText />
														{/if}
													</Item.Media>
													<Item.Content class="min-w-0 flex-1 text-left">
														<Item.Title class="w-full">
															{key}
														</Item.Title>
														<Item.Description class="wrap-break-words breaks-all">
															{value}
														</Item.Description>
													</Item.Content>
												</Item.Root>
											</Collapsible.Trigger>
											<Collapsible.Content
												class="overflow-hidden transition-all duration-300 ease-in-out"
											>
												<Code.Root
													lang="json"
													hideLines
													code={data}
													class="border-none bg-transparent px-8"
												/>
											</Collapsible.Content>
										</Collapsible.Root>
									{/if}
								{/each}
							</div>
						</Tabs.Content>
						<Tabs.Content value="table">
							<div>
								<Table.Root class="[&_td:first-child]:rounded-l-lg [&_td:last-child]:rounded-r-lg">
									<Table.Header>
										<Table.Row class="hover:[&>th,td]:bg-transparent!">
											<Table.Head>Key</Table.Head>
											<Table.Head>Value</Table.Head>
										</Table.Row>
									</Table.Header>
									<Table.Body>
										{#each Object.entries(data) as [key, value], index (index)}
											<Table.Row class="border-none align-top">
												<Table.Cell class="align-top">{key}</Table.Cell>
												<Table.Cell class="align-top">
													<p
														class="wrap-break-words max-w-3xl text-sm leading-normal font-normal text-balance break-all text-muted-foreground"
													>
														{value}
													</p>
												</Table.Cell>
											</Table.Row>
										{/each}
									</Table.Body>
								</Table.Root>
							</div>
						</Tabs.Content>
					</Tabs.Root>
				{/if}
			{:else}
				<Empty.Root class="m-4 bg-muted/50">
					<Empty.Header>
						<Empty.Media variant="icon">
							<Braces size={36} />
						</Empty.Media>
						<Empty.Title>No Data</Empty.Title>
						<Empty.Description>
							No data is currently available for this resource.
							<br />
							To populate this resource, please add properties or values through the resource editor.
						</Empty.Description>
					</Empty.Header>
					<Empty.Content></Empty.Content>
				</Empty.Root>
			{/if}
		</Sheet.Content>
	</Sheet.Root>
{/snippet}

{#snippet ArrayCell({ data }: { data: JsonValue[] })}
	{data.length}
{/snippet}

{#snippet DateCell({ data }: { data: Date })}
	{#if data && !isNaN(data.getTime())}
		{new Intl.DateTimeFormat('en-CA', {
			year: 'numeric',
			month: '2-digit',
			day: '2-digit'
		}).format(data)}
	{/if}
{/snippet}

{#snippet DatetimeCell({ data }: { data: Date })}
	{#if data && !isNaN(data.getTime())}
		{@const { value, unit } = getRelativeTime($now, data.getTime())}
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					{value}
					{unit}
				</Tooltip.Trigger>
				<Tooltip.Content>
					{new Intl.DateTimeFormat('en-CA', {
						year: 'numeric',
						month: '2-digit',
						day: '2-digit',
						hour: '2-digit',
						minute: '2-digit',
						second: '2-digit',
						hour12: false,
						timeZoneName: 'longOffset'
					}).format(data)}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{/if}
{/snippet}

{#snippet NumberCell({ data }: { data: number })}
	{data}
{/snippet}

{#snippet BooleanCell({ data }: { data: boolean })}
	{#if data === true}
		<Circle class="inline-block size-4 text-primary" />
	{:else if data === false}
		<X class="inline-block size-6 text-destructive" />
	{/if}
{/snippet}

{#snippet TextCell({ data }: { data: string })}
	<p class="truncate">
		{data}
	</p>
{/snippet}

{#snippet EmptyCell()}
	<p class="text-muted-foreground">no data</p>
{/snippet}

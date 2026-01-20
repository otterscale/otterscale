<script lang="ts">
	import Circle from '@lucide/svelte/icons/circle';
	import FileCode from '@lucide/svelte/icons/file-code';
	import X from '@lucide/svelte/icons/x';
	import type { JsonObject, JsonValue } from '@openfeature/server-sdk';
	import type { Column, Row } from '@tanstack/table-core';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
	import Monaco from 'svelte-monaco';
	import { stringify } from 'yaml';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet/index.js';

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
		<Sheet.Trigger class={buttonVariants({ variant: 'outline' })}>
			<FileCode />
		</Sheet.Trigger>
		<Sheet.Content side="right" class="flex h-full max-w-[62vw] min-w-[50vw] flex-col p-4">
			<Sheet.Header class="shrink-0 space-y-4">
				<Sheet.Title>YAML</Sheet.Title>
				<Sheet.Description>
					{fields[column.id].description}
				</Sheet.Description>
			</Sheet.Header>
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
		{new Intl.DateTimeFormat('en-CA', {
			year: 'numeric',
			month: '2-digit',
			day: '2-digit',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit',
			hour12: false
		}).format(data)}
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

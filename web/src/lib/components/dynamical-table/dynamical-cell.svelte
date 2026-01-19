<script lang="ts">
	import Circle from '@lucide/svelte/icons/circle';
	import FileCode from '@lucide/svelte/icons/file-code';
	import X from '@lucide/svelte/icons/x';
	import type { JsonObject, JsonValue } from '@openfeature/server-sdk';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
	import Monaco from 'svelte-monaco';
	import { stringify } from 'yaml';

	import { buttonVariants } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet/index.js';

	let {
		ref = $bindable(null),
		object,
		field,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		object: any;
		field: any;
	} = $props();
</script>

<div class={className}>
	{#if field?.snippet}
		{@render field.snippet()}
	{:else if field?.type === 'object'}
		{@const data = object as JsonObject}
		<Sheet.Root>
			<Sheet.Trigger class={buttonVariants({ variant: 'outline' })}>
				<FileCode />
			</Sheet.Trigger>
			<Sheet.Content side="right" class="flex h-full max-w-[62vw] min-w-[50vw] flex-col p-4">
				<Sheet.Header class="shrink-0 space-y-4">
					<Sheet.Title>YAML</Sheet.Title>
					<Sheet.Description>
						{field.description}
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
	{:else if field?.type === 'array'}
		{@render ArrayCell({ data: object })}
	{:else if field?.type === 'string' && field?.format === 'date'}
		{@render DateCell({ data: new Date(object) })}
	{:else if field?.type === 'string' && field?.format === 'date-time'}
		{@render DatetimeCell({ data: new Date(object) })}
	{:else if field?.type === 'number' || field?.type === 'integer'}
		{@render NumberCell({ data: Number(object) })}
	{:else if field?.type === 'boolean'}
		{@render BooleanCell({ data: Boolean(object) })}
	{:else if object}
		{@render TextCell({ data: object })}
	{:else}
		{@render EmptyCell()}
	{/if}
</div>

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

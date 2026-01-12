<script lang="ts">
	import Circle from '@lucide/svelte/icons/circle';
	import FileCode from '@lucide/svelte/icons/file-code';
	import X from '@lucide/svelte/icons/x';
	import type { JsonObject, JsonValue } from '@openfeature/server-sdk';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Code from '$lib/components/custom/code/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as Sheet from '$lib/components/ui/sheet/index.js';

	let {
		ref = $bindable(null),
		object,
		field,
		class: className
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		object: JsonValue;
		field: any;
	} = $props();
</script>

<div class={className}>
	{#if field?.type === 'object'}
		{@const data = object as JsonObject}
		<Sheet.Root>
			<Sheet.Trigger class={buttonVariants({ variant: 'outline' })}>
				<FileCode />
			</Sheet.Trigger>
			<Sheet.Content side="right" class="flex h-full max-w-[62vw] min-w-[50vw] flex-col p-6">
				<Sheet.Header class="shrink-0">
					<Sheet.Title>YAML</Sheet.Title>
					<Sheet.Description>
						{field.description}
					</Sheet.Description>
				</Sheet.Header>
				<Code.Root
					class="border-none bg-transparent wrap-break-word whitespace-pre-wrap"
					lang="json"
					code={JSON.stringify(data, null, 2)}
					hideLines
				/>
			</Sheet.Content>
		</Sheet.Root>
	{:else if field?.type === 'array'}
		{@const data = object as unknown as JsonValue[]}
		{#each data as datum, index (index)}
			<Badge variant="outline">
				{datum}
			</Badge>
		{/each}
	{:else if field?.type === 'string' && field?.format === 'date-time'}
		{@const data = object as unknown as string}
		{@const time = new Date(data)}
		{#if time && !isNaN(time.getTime())}
			{@const year=String(time.getFullYear())}
			{@const month=String(time.getMonth() + 1).padStart(2, '0')}
			{@const day=String(time.getDate()).padStart(2, '0')}
			{@const hour=String(time.getHours()).padStart(2, '0')}
			{@const minute=String(time.getMinutes()).padStart(2, '0')}
			{@const second=String(time.getSeconds()).padStart(2, '0')}
			<p>{year}-{month}-{day} {hour}:{minute}:{second}</p>
		{/if}
	{:else if field?.type === 'number' || field?.type === 'integer'}
		{@const number = object as unknown as number}
		<p>
			{number}
		</p>
	{:else if field?.type === 'boolean'}
		{@const data = object as unknown as string}
		{#if Boolean(data)}
			<Circle class="inline-block text-primary" />
		{:else}
			<X class="inline-block text-destructive" />
		{/if}
	{:else}
		{object}
	{/if}
</div>

<script lang="ts">
	import { cn } from '$lib/utils';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import Self from './display.svelte';

	let { data = $bindable() }: { data: any } = $props();
</script>

{#if typeof data !== 'object'}
	{data}
{:else if Array.isArray(data)}
	{@const anyItem = data[0]}
	{#if typeof anyItem === 'object'}
		<span class="flex flex-col">
			{#each data as item, i}
				<div class="m-2 rounded-lg bg-muted/50">
					<Self bind:data={data[i]} />
				</div>
			{/each}
		</span>
	{:else}
		{#each data as item, index}
			<Badge class="w-fit">{item}</Badge>
		{/each}
	{/if}
{:else}
	<div class="flex flex-col pl-4">
		{#each Object.entries(data) as [key, value]}
			{@const className =
				typeof value === 'object'
					? 'rounded-lg rounded-background border'
					: 'flex items-center justify-between gap-2'}
			<div class={cn(className, 'p-2')}>
				<p>{key[0].toUpperCase() + key.slice(1)}</p>
				<div>
					<Self bind:data={data[key]} />
				</div>
			</div>
		{/each}
	</div>
{/if}

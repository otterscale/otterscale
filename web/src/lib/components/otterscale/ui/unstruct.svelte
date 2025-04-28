<script lang="ts">
	import { cn } from '$lib/utils';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import DiscreteArrayInput from '$lib/components/otterscale/ui/discrete-array-input.svelte';
	import Self from '$lib/components/otterscale/ui/unstruct.svelte';
	let { data = $bindable() }: { data: any } = $props();
</script>

{#if typeof data !== 'object'}
	{#if typeof data === 'boolean'}
		<Switch bind:checked={data} />
	{:else if typeof data === 'string'}
		<Input type="text" bind:value={data} />
	{:else if typeof data === 'number'}
		<Input type="number" bind:value={data} />
	{:else}
		<Input type="text" bind:value={data} />
	{/if}
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
		<DiscreteArrayInput bind:items={data} />
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

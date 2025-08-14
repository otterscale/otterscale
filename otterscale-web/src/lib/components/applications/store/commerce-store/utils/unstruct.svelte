<script lang="ts">
	import { cn } from '$lib/utils';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Label } from '$lib/components/ui/label';
	import DiscreteArrayInput from './discrete-array-input.svelte';
	import Self from './unstruct.svelte';
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
				<div class="rounded-lg bg-muted/50">
					<Self bind:data={data[i]} />
				</div>
			{/each}
		</span>
	{:else}
		<DiscreteArrayInput bind:items={data} />
	{/if}
{:else}
	<div class="flex flex-col">
		{#each Object.entries(data) as [key, value]}
			{@const className =
				typeof value === 'object' ? 'rounded-lg rounded-background border' : 'grid gap-2'}
			<div class={cn(className, 'py-2')}>
				<Label
					>{key
						.split('.')
						.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
						.join(' ')}</Label
				>
				<div>
					<Self bind:data={data[key]} />
				</div>
			</div>
		{/each}
	</div>
{/if}

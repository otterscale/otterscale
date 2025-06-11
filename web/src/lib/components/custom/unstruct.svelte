<script lang="ts" module>
	import Self from '$lib/components/custom/unstruct.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { cn } from '$lib/utils';

	function camelToPascal(value: string): string {
		return value
			.replace(/([a-z])([A-Z])/g, '$1 $2')
			.split(' ')
			.map((word) => word.charAt(0).toUpperCase() + word.slice(1))
			.join(' ');
	}

	export { camelToPascal };
</script>

<script lang="ts">
	let {
		data = $bindable(),
		orientation = 'horizontal'
	}: { data: any; orientation?: 'vertical' | 'horizontal' } = $props();
</script>

<main class="text-xs">
	{#if typeof data !== 'object'}
		{#if typeof data === 'boolean'}
			<Badge variant="outline">{data}</Badge>
		{:else if data instanceof Date}
			<Badge variant="outline">{data.toDateString()}</Badge>
		{:else}
			{data}
		{/if}
	{:else if Array.isArray(data)}
		{@const [anyItem] = data}
		{#if typeof anyItem === 'object'}
			<span class="flex flex-col">
				{#each data as item, index}
					<div class="p-2">
						<Self bind:data={data[index]} {orientation} />
					</div>
				{/each}
			</span>
		{:else}
			<span class="flex items-center gap-2">
				{#each data as item, index}
					<Badge variant="outline">
						{data[index]}
					</Badge>
				{/each}
			</span>
		{/if}
	{:else}
		<div class={cn(orientation === 'vertical' ? 'flex flex-col gap-2' : 'flex gap-2')}>
			{#each Object.entries(data) as [key, value]}
				<div class="rounded-background rounded-lg border">
					<p class="bg-muted whitespace-nowrap p-2 font-bold">{camelToPascal(key)}</p>
					<Separator />
					<div class="p-2">
						<Self bind:data={data[key]} {orientation} />
					</div>
				</div>
			{/each}
		</div>
	{/if}
</main>

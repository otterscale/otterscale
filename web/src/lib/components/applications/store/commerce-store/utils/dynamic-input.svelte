<script lang="ts" module>
	import { Multiple as MultipleInput, Single as SingleInput } from '$lib/components/custom/input';
	import { Label } from '$lib/components/ui/label';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import Self from './dynamic-input.svelte';

	function getType(object: any) {
		if (typeof object === 'string') {
			return 'text';
		} else if (typeof object === 'number') {
			return 'number';
		} else {
			return 'text';
		}
	}
</script>

<script lang="ts">
	let { data = $bindable() }: { data: any } = $props();
</script>

{#if typeof data !== 'object'}
	{#if typeof data === 'boolean'}
		<Switch bind:checked={data} />
	{:else}
		<SingleInput.General type={getType(data)} bind:value={data} />
	{/if}
{:else if Array.isArray(data)}
	{@const [anyItem] = data}
	{#if typeof anyItem === 'object'}
		<span class="flex flex-col">
			{#each data as _, i}
				<div class="bg-muted/50 rounded-lg">
					<Self bind:data={data[i]} />
				</div>
			{/each}
		</span>
	{:else}
		<MultipleInput.Root type={getType(data)} bind:values={data}>
			<MultipleInput.Viewer />
			<MultipleInput.Controller>
				<MultipleInput.Input />
				<MultipleInput.Add />
				<MultipleInput.Clear />
			</MultipleInput.Controller>
		</MultipleInput.Root>
	{/if}
{:else}
	<div class="flex flex-col">
		{#each Object.entries(data) as [key]}
			<div class="space-y-2 p-2">
				<Label>
					{key
						.split('.')
						.map((part) => part.charAt(0).toUpperCase() + part.slice(1))
						.join(' ')}
				</Label>
				<Self bind:data={data[key]} />
			</div>
		{/each}
	</div>
{/if}

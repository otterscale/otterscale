<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { onMount } from 'svelte';

	import type { FilterManager } from './utils';

	import { page } from '$app/stores';
	import Input from '$lib/components/ui/input/input.svelte';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils';

	export const COLUMN_ID = 'name';
</script>

<script lang="ts">
	let { filterManager }: { filterManager: FilterManager } = $props();

	onMount(() => {
		filterManager.searchedName = $page.url.searchParams.get(COLUMN_ID) ?? '';
	});
</script>

<span class="relative h-8">
	<Input placeholder={m.name()} type="text" bind:value={filterManager.searchedName} class="h-8 w-40 pr-9 pl-9" />
	<Icon icon="ph:magnifying-glass" class="absolute top-1/2 left-3 -translate-y-1/2" />
	<button
		onclick={() => {
			filterManager.resetName();
		}}
	>
		<Icon
			icon="ph:x"
			class={cn(filterManager.searchedName ? 'visible' : 'hidden', 'absolute top-1/2 right-3 -translate-y-1/2')}
		/>
	</button>
</span>

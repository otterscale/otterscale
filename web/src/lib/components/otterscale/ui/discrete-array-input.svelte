<script lang="ts">
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Input } from '$lib/components/ui/input';

	let { items: items = $bindable() }: { items: string[] } = $props();

	let input = $state('');

	function trigger(event: KeyboardEvent | FocusEvent) {
		if (event.type !== 'blur' && (event as KeyboardEvent).key !== 'Enter') return;

		if (input === '') return;

		items = [...items, input];
		input = '';
		console.log(`inner: ${items}`);
	}

	function remove(index: number) {
		items.splice(index, 1);
		items = items;
	}
</script>

<div class="flex flex-wrap gap-1">
	<div class="flex flex-wrap gap-1">
		{#each items as value, index}
			<Badge variant="outline" onclick={() => remove(index)} class="text-sm">{value}</Badge>
		{/each}
	</div>
	<Input onkeyup={trigger} onblur={trigger} placeholder={items.join(' ')} bind:value={input} />
</div>

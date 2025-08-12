<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	let {
		items = $bindable()
	}: {
		items: string[];
	} = $props();
	function add() {
		if (input.trim()) {
			items = [...items, input.trim()];
			input = '';
		}
	}
	function remove(index: number) {
		items = items.filter((_, i) => i !== index);
	}
	let input = $state('');
</script>

<div class="flex items-start justify-between gap-2">
	<div class="flex flex-wrap gap-2 p-2">
		{#each items as item, index}
			<Badge class="w-fit" onclick={() => remove(index)}>{item}</Badge>
		{/each}
	</div>
	<div>
		<Input
			bind:value={input}
			placeholder="Enter"
			onkeydown={(e) => {
				if (e.key === 'Enter') add();
			}}
		/>
	</div>
</div>

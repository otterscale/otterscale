<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children?: Snippet } = $props();
</script>

<div class="flex">
	<span class="ml-auto space-x-2">
		{@render children?.()}
	</span>
</div>

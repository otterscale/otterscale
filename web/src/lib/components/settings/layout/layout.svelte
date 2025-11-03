<script lang="ts" module>
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<div class="flex flex-col items-center justify-center">
	{@render children?.()}
</div>

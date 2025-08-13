<script lang="ts" module>
	import * as Accordion from '$lib/components/ui/accordion';
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<Accordion.Root type="single" class="space-y-4">
	{@render children?.()}
</Accordion.Root>

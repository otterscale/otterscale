<script lang="ts">
	import * as Accordion from '$lib/components/ui/accordion';
	import type { Snippet } from 'svelte';

	let { children }: { children?: Snippet } = $props();
</script>

<Accordion.Root type="single" class="space-y-4">
	{@render children?.()}
</Accordion.Root>

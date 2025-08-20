<script lang="ts">
	import * as Accordion from '$lib/components/ui/accordion';
	import type { Snippet } from 'svelte';

	let { value, children }: { value?: string; children?: Snippet } = $props();
</script>

<Accordion.Root type="single" class="space-y-4" {value}>
	{@render children?.()}
</Accordion.Root>

<script lang="ts" module>
	import * as Accordion from '$lib/components/ui/accordion';
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { value, children }: { value?: string; children?: Snippet } = $props();
</script>

<Accordion.Root type="single" class="space-y-4" {value}>
	{@render children?.()}
</Accordion.Root>

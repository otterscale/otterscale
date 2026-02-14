<script lang="ts">
	import { stringify } from 'yaml';

	import * as Code from '$lib/components/custom/code';

	let { object }: { object: any } = $props();
</script>

<Code.Root
	variant="secondary"
	lang="yaml"
	class="no-shiki-limit w-full border-none"
	code={stringify(object)}
/>

<style>
	@reference '../../../../app.css';

	:global(.no-shiki-limit pre.shiki:not([data-code-overflow] *):not([data-code-overflow])) {
		overflow-y: visible !important;
		max-height: none !important;
	}
</style>

<script lang="ts">
	import type { ComponentProps } from '@sjsf/form';
	import Template from '@sjsf/form/templates/object.svelte';

	let props: ComponentProps['objectTemplate'] = $props();

	// Use #await to defer rendering to the next microtask, breaking the recursion stack
</script>

{#await Promise.resolve() then}
	<Template {...props} />
{/await}

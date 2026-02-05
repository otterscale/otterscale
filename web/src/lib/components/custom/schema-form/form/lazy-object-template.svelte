<script lang="ts">
	import type { ComponentProps } from '@sjsf/form';
	import Template from '@sjsf/form/templates/object.svelte';

	let props: ComponentProps['objectTemplate'] = $props();

	// Use #await with a resolved promise to defer rendering asynchronously, breaking the synchronous recursion stack
</script>

{#await Promise.resolve() then}
	<Template {...props} />
{/await}

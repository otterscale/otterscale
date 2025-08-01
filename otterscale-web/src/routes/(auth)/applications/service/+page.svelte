<script lang="ts">
	import { applicationsPath, applicationsServicePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [applicationsPath], current: applicationsServicePath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}

<script lang="ts">
	import { applicationsPath, applicationsStorePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: applicationsPath, current: applicationsStorePath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}

<script lang="ts">
	import { applicationsPath, homePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: homePath, current: applicationsPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}

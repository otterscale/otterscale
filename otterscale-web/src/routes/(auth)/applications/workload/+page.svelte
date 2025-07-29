<script lang="ts">
	import { applicationsPath, applicationsWorkloadPath } from '$lib/path';
	import { activeScope, breadcrumb, loadingScopes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: applicationsPath, current: applicationsWorkloadPath });
</script>

{#if $loadingScopes}
	loading scope
{:else}
	current scope: {$activeScope.uuid}
{/if}

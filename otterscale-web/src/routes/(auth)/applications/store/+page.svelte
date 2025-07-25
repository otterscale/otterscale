<script lang="ts">
	import { applicationsPath, applicationsStorePath } from '$lib/path';
	import { activeScope, breadcrumb, scopeLoading } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: applicationsPath, current: applicationsStorePath });
</script>

{#if $scopeLoading}
	loading scope
{:else}
	current scope: {$activeScope.uuid}
{/if}

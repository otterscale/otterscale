<script lang="ts">
	import { setupPath, setupScopePath } from '$lib/path';
	import { activeScope, breadcrumb, currentCeph, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: setupPath, current: setupScopePath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}

Ceph: {$currentCeph?.name}
Kubernetes: {$currentKubernetes?.name}

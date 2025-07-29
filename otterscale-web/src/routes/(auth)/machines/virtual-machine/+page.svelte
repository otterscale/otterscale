<script lang="ts">
	import { machinesPath, machinesVirtualMachinePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: machinesPath, current: machinesVirtualMachinePath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}

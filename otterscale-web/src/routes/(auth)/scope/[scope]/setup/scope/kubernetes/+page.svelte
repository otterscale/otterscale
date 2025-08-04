<script lang="ts">
	import { setupPath, setupScopeKubernetesPath, setupScopePath } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [setupPath, setupScopePath], current: setupScopeKubernetesPath });
</script>

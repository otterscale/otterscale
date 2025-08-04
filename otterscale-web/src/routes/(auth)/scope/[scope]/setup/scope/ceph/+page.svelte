<script lang="ts">
	import { setupPath, setupScopeCephPath, setupScopePath } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [setupPath, setupScopePath], current: setupScopeCephPath });
</script>

<script lang="ts">
	import { homePath, setupPath } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: homePath, current: setupPath });
</script>

Setup Enviroment

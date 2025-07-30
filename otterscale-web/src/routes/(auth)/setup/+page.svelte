<script lang="ts">
	import { homePath, setupPath } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [homePath], current: setupPath });
</script>

Setup Enviroment

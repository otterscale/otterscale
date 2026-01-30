<script lang="ts">
	import SheetCreateWorkspace from './sheet-create-workspace.svelte';

	let { open = $bindable(false) }: { open: boolean } = $props();
</script>

<!-- Directly render the sheet as we don't have a plan selection step yet -->
<SheetCreateWorkspace bind:open />

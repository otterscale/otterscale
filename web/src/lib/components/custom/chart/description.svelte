<script lang="ts">
	let { description }: { description?: string } = $props();
</script>

{description || 'No Data'}

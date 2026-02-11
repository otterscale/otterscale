<script lang="ts">
	import Create from '$lib/components/form/cronjob/create-sheet.svelte';

	let { schema }: { schema: any } = $props();
</script>

<Create {schema} />

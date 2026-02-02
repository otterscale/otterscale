<script lang="ts">
	import CreateSheet from '$lib/components/form/workspace/create-sheet.svelte';
	import { Button } from '$lib/components/ui/button';

	let open = $state(false);
</script>

<div class="container mx-auto py-10">
	<div class="mb-12">
		<h1 class="mb-4 text-2xl font-bold">Workspace Creation Dialog Test</h1>
		<Button onclick={() => (open = true)}>Create Workspace</Button>
		<CreateSheet bind:open />
	</div>
</div>

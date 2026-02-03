<script lang="ts">
	import CreateSheet from '$lib/components/form/cronjob/create-sheet.svelte';
	import { Button } from '$lib/components/ui/button';

	import schema from '../cron_api.json';

	let open = $state(false);
</script>

<div class="container mx-auto py-10">
	<div class="mb-12">
		<h1 class="mb-4 text-2xl font-bold">CronJob Creation Dialog Test</h1>
		<Button onclick={() => (open = true)}>Create CronJob</Button>
		<!-- Pass the imported schema to the CreateSheet -->
		<CreateSheet bind:open cluster="test-cluster" {schema} />
	</div>
</div>

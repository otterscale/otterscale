<script lang="ts">
	import DeleteDialog from '$lib/components/form/cronjob/delete-dialog.svelte';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	let cronJobName = $state('default-cronjob'); // Default value for testing
</script>

<div class="container mx-auto py-10">
	<div class="mb-12 space-y-4">
		<h1 class="mb-4 text-2xl font-bold">CronJob Delete Dialog Test</h1>

		<div class="grid w-full max-w-sm items-center gap-1.5">
			<Label for="cronjob-name">CronJob Name</Label>
			<Input
				type="text"
				id="cronjob-name"
				placeholder="Enter cronjob name"
				bind:value={cronJobName}
			/>
		</div>

		<DeleteDialog name={cronJobName} />
	</div>
</div>

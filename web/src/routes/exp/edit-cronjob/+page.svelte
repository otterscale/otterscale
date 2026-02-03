<script lang="ts">
	import EditSheet from '$lib/components/form/cronjob/edit-sheet.svelte';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	import schema from '../cron_api.json';

	let cronJobName = $state('my-cronjob'); // Default value for testing

	// Mock object
	const mockObject = {
		apiVersion: 'batch/v1',
		kind: 'CronJob',
		metadata: {
			name: 'my-cronjob',
			namespace: 'default',
			creationTimestamp: '2023-10-27T10:00:00Z',
			uid: '1234-5678-90ab-cdef'
		},
		spec: {
			schedule: '*/5 * * * *',
			jobTemplate: {
				spec: {
					template: {
						spec: {
							containers: [
								{
									name: 'hello',
									image: 'busybox',
									args: ['/bin/sh', '-c', 'date; echo Hello from the Kubernetes cluster']
								}
							],
							restartPolicy: 'OnFailure'
						}
					}
				}
			}
		}
	};
</script>

<div class="container mx-auto py-10">
	<div class="mb-12 space-y-4">
		<h1 class="mb-4 text-2xl font-bold">CronJob Edit Sheet Test</h1>

		<div class="grid w-full max-w-sm items-center gap-1.5">
			<Label for="cronjob-name">CronJob Name</Label>
			<Input
				type="text"
				id="cronjob-name"
				placeholder="Enter cronjob name"
				bind:value={cronJobName}
			/>
		</div>

		<!-- Pass the mock object and schema -->
		<EditSheet name={cronJobName} {schema} object={mockObject} />
	</div>
</div>

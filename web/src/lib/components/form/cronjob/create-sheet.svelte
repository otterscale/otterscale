<script lang="ts">
	import { Plus } from '@lucide/svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Sheet from '$lib/components/ui/sheet';

	import CreateCronJobForm from './create-form.svelte';

	let {
		schema = undefined,
		onsuccess
	}: {
		schema?: any;
		onsuccess?: (cronjob?: any) => void;
	} = $props();

	function handleCronJobSuccess(cronjob?: any) {
		open = false;
		if (cronjob?.metadata?.name) {
			onsuccess?.(cronjob);
			goto(
				resolve(
					`/(auth)/${page.params.cluster}/CronJob/cronjobs?group=batch&version=v1&name=${cronjob.metadata.name}&namespace=${page.url.searchParams.get('namespace') ?? ''}`
				)
			);
		}
	}

	let open = $state(false);
	function handleOpenChange(isOpen: boolean) {
		open = isOpen;
	}
</script>

<Sheet.Root bind:open onOpenChange={handleOpenChange}>
	<Sheet.Trigger>
		<Button variant="outline" size="icon">
			<Plus />
		</Button>
	</Sheet.Trigger>
	<Sheet.Content
		class="fixed top-1/2 left-1/2 h-[90vh] w-[90vw] max-w-4xl min-w-[800px] -translate-x-1/2 -translate-y-1/2 rounded-lg border bg-background shadow-lg"
	>
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full flex-col">
				<div class="flex-1 overflow-y-auto p-6">
					<CreateCronJobForm {schema} onsuccess={handleCronJobSuccess} />
				</div>
			</div>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>

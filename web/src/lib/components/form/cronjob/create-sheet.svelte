<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import BasicTierImage from '$lib/assets/basic-tier.jpg'; // Reusing image or replace if needed
	import * as Sheet from '$lib/components/ui/sheet';

	import CreateCronJobForm from './create-form.svelte';

	let {
		open = $bindable(false),
		cluster,
		schema = undefined,
		onsuccess
	}: {
		open: boolean;
		cluster: string;
		schema?: any;
		onsuccess?: (cronjob?: any) => void;
	} = $props();

	function handleCronJobSuccess(cronjob?: any) {
		open = false;
		if (cronjob?.metadata?.name) {
			onsuccess?.(cronjob);
			goto(
				resolve(
					`/(auth)/${cluster}/CronJob/cronjobs?group=batch&version=v1&name=${cronjob.metadata.name}`
				)
			);
		}
	}

	function handleOpenChange(isOpen: boolean) {
		open = isOpen;
	}
</script>

<Sheet.Root bind:open onOpenChange={handleOpenChange}>
	<Sheet.Content class="inset-y-auto bottom-0 h-9/10 rounded-tl-lg sm:max-w-4/5">
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full">
				<div class="flex-1 overflow-y-auto p-6">
					<CreateCronJobForm {schema} onsuccess={handleCronJobSuccess} />
				</div>
				<!-- Image -->
				<div class="relative hidden w-2/5 lg:block">
					<img src={BasicTierImage} alt="CronJob" class="absolute inset-0 size-full object-cover" />
				</div>
			</div>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>

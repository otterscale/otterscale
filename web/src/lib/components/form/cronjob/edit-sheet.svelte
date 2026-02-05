<script lang="ts">
	import { Pencil } from '@lucide/svelte';

	import type { K8sOpenAPISchema } from '$lib/components/custom/schema-form';
	import * as Sheet from '$lib/components/ui/sheet';

	import EditCronJobForm from './edit-form.svelte';

	let {
		name,
		schema,
		object,
		onOpenChangeComplete,
		onsuccess
	}: {
		name: string;
		schema: K8sOpenAPISchema;
		object: Record<string, unknown>;
		onOpenChangeComplete: () => void;
		onsuccess?: () => void;
	} = $props();

	let open = $state(false);

	function handleClose() {
		open = false;
		onsuccess?.();
	}
	console.log('object', object);
</script>

<Sheet.Root bind:open {onOpenChangeComplete}>
	<Sheet.Trigger class="flex items-center gap-2">
		<Pencil size={16} />
		Edit
	</Sheet.Trigger>
	<Sheet.Content
		class="fixed top-1/2 left-1/2 h-[90vh] w-[90vw] max-w-4xl min-w-[800px] -translate-x-1/2 -translate-y-1/2 rounded-lg border bg-background shadow-lg"
	>
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full flex-col">
				<div class="flex-1 overflow-y-auto p-6">
					{#if name && schema && object}
						<EditCronJobForm {name} {schema} {object} onsuccess={handleClose} />
					{:else}
						<div class="flex h-full items-center justify-center">
							<p class="text-muted-foreground">No cronjob selected.</p>
						</div>
					{/if}
				</div>
			</div>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>

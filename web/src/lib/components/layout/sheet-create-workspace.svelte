<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet';

	import CreateWorkspaceForm from './create-workspace-form.svelte';

	let { open = $bindable(false) }: { open: boolean } = $props();

	function handleClose() {
		open = false;
	}
</script>

<Sheet.Root bind:open onOpenChange={handleClose}>
	<Sheet.Content class="inset-y-auto bottom-0 h-9/10 rounded-tl-lg sm:max-w-4/5">
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full">
				<div class="flex-1 p-6 overflow-y-auto">
					<CreateWorkspaceForm onsuccess={handleClose} />
				</div>
                <!-- Optional: Add an image or sidebar here if desired, similar to Create Scope -->
                <!-- For now, we'll keep it simple as there are no "Plans" with images yet -->
			</div>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>

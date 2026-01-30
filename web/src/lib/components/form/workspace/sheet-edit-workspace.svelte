<script lang="ts">
	import BasicTierImage from '$lib/assets/basic-tier.jpg';
	import EditWorkspaceForm from '$lib/components/form/workspace/edit-workspace-form.svelte';
	import * as Sheet from '$lib/components/ui/sheet';

	let {
		open = $bindable(false),
		name
	}: {
		open: boolean;
		name: string;
	} = $props();

	function handleClose() {
		open = false;
	}
</script>

<Sheet.Root bind:open onOpenChange={handleClose}>
	<Sheet.Content class="inset-y-auto bottom-0 h-9/10 rounded-tl-lg sm:max-w-4/5">
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full">
				<div class="flex-1 overflow-y-auto p-6">
					{#if name}
						<EditWorkspaceForm {name} onsuccess={handleClose} />
					{:else}
						<div class="flex h-full items-center justify-center">
							<p class="text-muted-foreground">No workspace selected.</p>
						</div>
					{/if}
				</div>
				<!-- Workspace Image -->
				<div class="relative hidden w-2/5 lg:block">
					<img
						src={BasicTierImage}
						alt="Workspace"
						class="absolute inset-0 size-full object-cover"
					/>
				</div>
			</div>
		</Sheet.Header>
	</Sheet.Content>
</Sheet.Root>

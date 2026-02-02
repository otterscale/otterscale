<script lang="ts">
	import { type TenantOtterscaleIoV1Alpha1Workspace } from '@otterscale/types';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import BasicTierImage from '$lib/assets/basic-tier.jpg';
	import * as Sheet from '$lib/components/ui/sheet';

	import CreateWorkspaceForm from './create-form.svelte';

	let {
		open = $bindable(false),
		cluster,
		onsuccess
	}: {
		open: boolean;
		cluster: string;
		onsuccess?: (workspace?: TenantOtterscaleIoV1Alpha1Workspace) => void;
	} = $props();

	function handleClose(workspace?: TenantOtterscaleIoV1Alpha1Workspace) {
		open = false;
		if (workspace?.metadata?.name) {
			onsuccess?.(workspace);
			goto(
				resolve(
					`/(auth)/${cluster}/Workspace/workspaces?group=tenant.otterscale.io&version=v1alpha1&name=${workspace.metadata.name}`
				)
			);
		}
	}
</script>

<Sheet.Root bind:open onOpenChange={handleClose}>
	<Sheet.Content class="inset-y-auto bottom-0 h-9/10 rounded-tl-lg sm:max-w-4/5">
		<Sheet.Header class="h-full p-0">
			<div class="flex h-full">
				<div class="flex-1 overflow-y-auto p-6">
					<CreateWorkspaceForm onsuccess={handleClose} />
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

<script lang="ts">
	import SheetEditWorkspace from '$lib/components/form/workspace/sheet-edit-workspace.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	let open = $state(false);
	let workspaceName = $state('default'); // Default value for testing
</script>

<div class="container mx-auto py-10">
	<div class="mb-12 space-y-4">
		<h1 class="mb-4 text-2xl font-bold">Workspace Edit Dialog Test</h1>

		<div class="grid w-full max-w-sm items-center gap-1.5">
			<Label for="workspace-name">Workspace Name</Label>
			<Input
				type="text"
				id="workspace-name"
				placeholder="Enter workspace name"
				bind:value={workspaceName}
			/>
		</div>

		<Button onclick={() => (open = true)} disabled={!workspaceName}>Edit Workspace</Button>
		<SheetEditWorkspace bind:open name={workspaceName} />
	</div>
</div>

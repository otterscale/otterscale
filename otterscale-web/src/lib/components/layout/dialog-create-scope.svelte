<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

	let { open = $bindable(false) }: { open: boolean } = $props();

	let name = $state('');

	function handleSubmit() {
		if (name.trim()) {
			console.log('Creating scope:', name);
			open = false;
			name = '';
		}
	}

	function handleClose() {
		open = false;
		name = '';
	}
</script>

<Dialog.Root bind:open onOpenChange={handleClose}>
	<Dialog.Content class="sm:max-w-[475px]">
		<Dialog.Header>
			<Dialog.Title>Create Scope</Dialog.Title>
			<Dialog.Description>
				Organize users into groups with specific permissions. Resources are isolated between scopes.
			</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit}>
			<div class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="name" class="text-right">Name</Label>
					<Input
						id="name"
						bind:value={name}
						placeholder="Enter scope name"
						class="col-span-3"
						required
					/>
				</div>
			</div>

			<Dialog.Footer>
				<Button type="button" variant="outline" onclick={handleClose}>Cancel</Button>
				<Button type="submit" disabled={!name.trim()}>Create</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

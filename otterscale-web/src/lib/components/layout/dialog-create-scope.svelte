<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';

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
			<Dialog.Title>{m.create_scope()}</Dialog.Title>
			<Dialog.Description>{m.create_scope_description()}</Dialog.Description>
		</Dialog.Header>

		<form onsubmit={handleSubmit}>
			<div class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="name" class="text-right">{m.scope_name()}</Label>
					<Input
						id="name"
						bind:value={name}
						placeholder={m.scope_name_description()}
						class="col-span-3"
						required
					/>
				</div>
			</div>

			<Dialog.Footer>
				<Button type="button" variant="outline" onclick={handleClose}>{m.cancel()}</Button>
				<Button type="submit" disabled={!name.trim()}>{m.create()}</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

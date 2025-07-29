<script lang="ts">
	import { getContext } from 'svelte';
	import { createClient, type Transport } from '@connectrpc/connect';
	import { ScopeService, type CreateScopeRequest } from '$lib/api/scope/v1/scope_pb';
	import { Button } from '$lib/components/ui/button';
	import * as Dialog from '$lib/components/ui/dialog';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { m } from '$lib/paraglide/messages';
	import { toast } from 'svelte-sonner';
	import { triggerUpdateScopes } from '$lib/stores';

	let { open = $bindable(false) }: { open: boolean } = $props();

	const transport: Transport = getContext('transport');
	const scopeClient = createClient(ScopeService, transport);

	const DEFAULT_REQUEST = { name: '' } as CreateScopeRequest;

	let createScopeRequest = $state(DEFAULT_REQUEST);

	function handleSubmit() {
		if (createScopeRequest.name.trim()) {
			scopeClient
				.createScope(createScopeRequest)
				.then((r) => {
					toast.success(m.create_scope_success({ name: r.name }));
					triggerUpdateScopes.set(true);
				})
				.catch((e) => {
					toast.error(m.create_scope_error({ name: createScopeRequest.name, error: e.toString() }));
				});

			open = false;
			createScopeRequest = DEFAULT_REQUEST;
		}
	}

	function handleClose() {
		open = false;
		createScopeRequest = DEFAULT_REQUEST;
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
						bind:value={createScopeRequest.name}
						placeholder={m.scope_name_description()}
						class="col-span-3"
						required
					/>
				</div>
			</div>

			<Dialog.Footer>
				<Button type="button" variant="outline" onclick={handleClose}>{m.cancel()}</Button>
				<Button type="submit" disabled={!createScopeRequest.name.trim()}>{m.create()}</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<script lang="ts" module>
	import type { SMBShare_SecurityConfig_User } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages.js';
</script>

<script lang="ts">
	let {
		user = $bindable<SMBShare_SecurityConfig_User>(),
		invalid = $bindable<boolean>()
	}: { user?: SMBShare_SecurityConfig_User; invalid?: boolean } = $props();

	const defaults = {} as SMBShare_SecurityConfig_User;

	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	let open = $state(false);
	function close() {
		open = false;
	}

	$effect(() => {
		invalid = !user;
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="primary" class="w-full">{m.create()}/{m.edit()}</Modal.Trigger>
	<Modal.Content>
		<Form.Label>{m.name()}</Form.Label>
		<SingleInput.General type="text" bind:value={request.username} required={invalid} />

		<Form.Label>{m.password()}</Form.Label>
		<SingleInput.General type="password" bind:value={request.password} required={invalid} />

		<Button
			onclick={() => {
				user = request;
				reset();
				close();
			}}
		>
			{m.confirm()}
		</Button>
	</Modal.Content>
</Modal.Root>

<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { SMBShare_User } from '$lib/api/storage/v1/storage_pb';
	import * as Form from '$lib/components/custom/form';
	import { Single as SingleInput } from '$lib/components/custom/input';
	import Button from '$lib/components/ui/button/button.svelte';
	import { m } from '$lib/paraglide/messages.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		user = $bindable<SMBShare_User>(),
		invalid = $bindable<boolean>()
	}: { user?: SMBShare_User; invalid?: boolean } = $props();

	const defaults = {} as SMBShare_User;
	let request = $state(defaults);
	function reset() {
		request = defaults;
	}

	$effect(() => {
		invalid = !user;
	});
</script>

{#if user}
	<div class="rounded-lg border p-2">
		<div class="flex items-center gap-2 rounded-lg p-2">
			<div class={cn('flex size-8 items-center justify-center rounded-full border-2')}>
				<Icon icon="ph:user" class="size-5" />
			</div>

			<div class="flex flex-col gap-1">
				<p class="text-xs text-muted-foreground">{m.user()}</p>
				<p class="text-sm">{user.username}</p>
			</div>
		</div>
	</div>
{:else}
	<div
		class="rounded-lg border border-red-300 bg-destructive/10 p-4 text-center text-xs text-destructive"
	>
		There is no user. Please create a user for this share.
	</div>
{/if}

<Form.Label>{m.name()}</Form.Label>
<SingleInput.General type="text" bind:value={request.username} required={invalid} />

<Form.Label>{m.password()}</Form.Label>
<SingleInput.General type="password" bind:value={request.password} required={invalid} />

<Button
	onclick={() => {
		user = request;
		reset();
	}}
>
	<Icon icon="ph:plus" />
</Button>

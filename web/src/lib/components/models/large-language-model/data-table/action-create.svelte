<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { getContext, onMount } from 'svelte';

	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages.js';
</script>

<script lang="ts">
	let open = $state(false);
	function close() {
		open = false;
	}

	const reloadManager: ReloadManager = getContext('reloadManager');
	onMount(() => {
		reloadManager.force();
	});
</script>

<Modal.Root bind:open>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create()}</Modal.Header>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					close();
				}}
			>
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

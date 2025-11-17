<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages.js';
</script>

<script lang="ts">
	let { scope, reloadManager }: { scope: string; reloadManager: ReloadManager } = $props();

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Modal.Root bind:open>
	<Modal.Trigger class="default">
		<Icon icon="ph:plus" />
		{m.create()}
		{scope}
	</Modal.Trigger>
	<Modal.Content>
		<Modal.Header>{m.create()}</Modal.Header>
		<Modal.Footer>
			<Modal.Cancel>
				{m.cancel()}
			</Modal.Cancel>
			<Modal.Action
				onclick={() => {
					reloadManager.force();
					close();
				}}
			>
				{m.confirm()}
			</Modal.Action>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

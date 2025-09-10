<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Code from '$lib/components/custom/code';
	import { SingleStep as Modal } from '$lib/components/custom/modal';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();

	let open = $state(false);

	const code = virtualMachine.startupScript;
</script>

<Modal.Root bind:open>
	<Modal.Trigger variant="creative">
		<Icon icon="ph:file" />
		{m.startup_script()}
	</Modal.Trigger>
	<Modal.Content class="min-w-[50vw]">
		<Modal.Header>
			{m.startup_script()}
		</Modal.Header>
		<Code.Root lang="svelte" class="w-full" variant="secondary" hideLines {code}>
			<Code.CopyButton />
		</Code.Root>
		<Modal.Footer>
			<Modal.Cancel>{m.cancel()}</Modal.Cancel>
		</Modal.Footer>
	</Modal.Content>
</Modal.Root>

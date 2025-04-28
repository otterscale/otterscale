<script lang="ts">
	// External dependencies
	import type { Node } from '@xyflow/svelte';

	// Internal UI components
	import * as Drawer from '$lib/components/ui/drawer';

	// Internal utilities and types
	import { type Machine, type Model } from '$gen/api/stack/v1/stack_pb';
	import OrchestrationInformationMAAS from './orchestration-information-maas.svelte';
	import OrchestrationInformationJUJU from './orchestration-information-juju.svelte';

	let {
		open = $bindable(),
		node
	}: {
		open: boolean;
		node: Node;
	} = $props();
</script>

<Drawer.Root direction="right" bind:open>
	<Drawer.Content
		class="inset-x-auto inset-y-0 right-0 max-h-[100vh] w-4/5 overflow-x-hidden overflow-y-scroll rounded-tr-none"
	>
		{#if node.type == 'MAAS'}
			<OrchestrationInformationMAAS machine={node.data as Machine} />
		{:else if node.type == 'JUJU'}
			<OrchestrationInformationJUJU model={node.data as Model} />
		{/if}
	</Drawer.Content>
</Drawer.Root>

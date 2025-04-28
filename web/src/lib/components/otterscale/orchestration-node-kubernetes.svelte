<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';
	import Icon from '@iconify/svelte';
	import { nodeIcon } from '$lib/node';
	import * as HoverCard from '$lib/components/ui/hover-card/index.js';

	type $$Props = NodeProps;

	export let data: $$Props['data'];

	const nodeType = 'kubernetes';
</script>

<div class="animate-pulse flex-col space-y-2">
	<HoverCard.Root>
		<HoverCard.Trigger>
			<div
				class="bg-black-50 group flex h-[100px] w-[100px] items-center justify-center rounded-full shadow-md hover:bg-accent hover:text-accent-foreground"
			>
				<Icon icon={nodeIcon(nodeType)} class="h-10 w-10 group-hover:scale-110" />
			</div>
		</HoverCard.Trigger>
		<HoverCard.Content>
			{data.uuid}
		</HoverCard.Content>
	</HoverCard.Root>
	<Handle type="target" position={Position.Left} class="invisible" />
	<Handle type="source" position={Position.Right} class="invisible" />
</div>

<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Handle, type NodeProps } from '@xyflow/svelte';

	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { type PodInfo } from '$lib/api/essential/v1/essential_pb';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { dynamicPaths } from '$lib/path';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { data, selected, targetPosition, sourcePosition }: Omit<NodeProps, 'data'> & { data: PodInfo } = $props();
</script>

<div
	class={cn(
		'bg-card relative flex h-[150px] w-[300px] rounded-lg border p-2 hover:shadow',
		selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
	)}
>
	<div
		class={cn(
			'bg-card hover:bg-muted absolute top-1 right-1 translate-x-1/2 -translate-y-1/2 rounded-full border p-2 shadow hover:cursor-default',
			selected ? 'bg-primary-foreground ring-primary ring-1' : 'bg-card ring-0',
		)}
	>
		<Icon
			icon="simple-icons:maas"
			class="size-5"
			onclick={(e) => {
				e.stopPropagation();
				goto(`${dynamicPaths.machinesMetal(page.params.scope).url}`);
			}}
		/>
	</div>
	<div class="flex items-center justify-center p-4">
		<div class="flex items-center gap-2">
			<div class="bg-muted-foreground/50 size-fit rounded-full p-2">
				<Icon icon="ph:computer-tower" class="size-5" />
			</div>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<p class="max-w-[200px] truncate text-base text-nowrap whitespace-nowrap">{data.machineName}</p>
				</Tooltip.Trigger>
				<Tooltip.Content>
					{data.machineName}
				</Tooltip.Content>
			</Tooltip.Root>
		</div>
	</div>
	{#if targetPosition}
		<Handle type="target" position={targetPosition} class="invisible" />
	{/if}
	{#if sourcePosition}
		<Handle type="source" position={sourcePosition} class="invisible" />
	{/if}
</div>

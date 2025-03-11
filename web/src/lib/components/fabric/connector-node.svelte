<script lang="ts">
	import { Handle, Position, type NodeProps } from '@xyflow/svelte';
	import { Badge } from '$lib/components/ui/badge';
	import * as Avatar from '$lib/components/ui/avatar';
	import Icon from '@iconify/svelte';
	import { Separator } from '$lib/components/ui/separator';
	import { avatarFallback, avatarURL, type pbWorkload } from '$lib/pb';
	import { formatTimeAgo } from '$lib/formatter';
	import { connectorIcon, connectorLabel, connectorLabelIcon } from '$lib/connector';

	type $$Props = NodeProps;

	export let data: $$Props['data'] & {
		workload?: pbWorkload;
	};
</script>

<div class="flex-col space-y-2">
	<!-- TODO TAG -->
	<!-- <div class="flex items-end justify-end space-x-2 px-2">
		 {#each data.tags as string[] as tag}
			<Badge variant="default" class="hover:scale-105">
				<Icon icon="material-symbols:folder-outline" class="h-4 w-4 sm:flex" />
				<span class="pl-2 font-thin">{tag}</span>
			</Badge>
		{/each} 
	</div> -->
	<div
		class="group rounded-md border-2 px-4 py-2 shadow-md hover:bg-accent hover:text-accent-foreground"
	>
		<div class="flex flex-col items-start space-y-2">
			<div class="flex w-full items-center space-x-3">
				<Icon
					icon={connectorIcon(data.type as string)}
					class="hidden h-10 w-10 group-hover:scale-110 sm:flex"
				/>
				<div class="flex flex-col items-start">
					<div class="flex items-center text-sm text-muted-foreground">
						{data.kind}
					</div>
					<div class="text-md font-bold text-foreground">{data.name}</div>
				</div>
			</div>
			<Separator />
			<div class="flex items-center space-x-2">
				<Badge variant="secondary" class="group-hover:border-inherit">
					<Avatar.Root class="size-4">
						<Avatar.Image src={avatarURL(data.workload?.avatar as string)} />
						<Avatar.Fallback>{avatarFallback(data.workload?.user as string)}</Avatar.Fallback>
					</Avatar.Root>
					<span class="pl-2 font-bold text-muted-foreground">{data.workload?.user}</span>
				</Badge>
				<Badge variant="secondary" class="group-hover:border-inherit">
					<Icon icon="mingcute:time-line" class="hidden size-4 sm:flex" />
					<span class="pl-2 text-muted-foreground">{formatTimeAgo(data.updated as Date)}</span>
				</Badge>
				<Badge variant="secondary" class="group-hover:border-inherit">
					<Icon icon={connectorLabelIcon(data.type as string)} class="hidden size-4 sm:flex" />
					<span class="pl-2 text-muted-foreground">{connectorLabel(data.type as string)}</span>
				</Badge>
			</div>
		</div>
		{#if data.kind == 'destination'}
			<Handle type="target" position={Position.Left} class="h-3 w-3 rounded-full bg-primary" />
		{/if}
		{#if data.kind == 'source'}
			<Handle type="source" position={Position.Right} class="h-3 w-3 rounded-full bg-primary" />
		{/if}
	</div>
</div>

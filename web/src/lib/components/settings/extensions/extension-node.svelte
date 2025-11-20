<script lang="ts" module>
	import Icon from '@iconify/svelte';

	import type { Extension } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as Card from '$lib/components/ui/card';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let { extension, alignment }: { extension: Extension; alignment: 'left' | 'right' } = $props();
</script>

<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
	<div
		class={cn(
			alignment == 'right'
				? 'relative flex flex-row-reverse items-center gap-8'
				: 'relative flex items-center gap-8'
		)}
	>
		<div
			class="absolute left-1/2 z-10 flex size-13 -translate-x-1/2 transform items-center justify-center rounded-full border bg-muted/70 font-bold ring-1 ring-muted/30 ring-offset-1"
		>
			<Avatar.Root class="size-9">
				<Avatar.Image src={extension.icon} />
				<Avatar.Fallback>
					<Icon icon="ph:puzzle-piece" class="size-6" />
				</Avatar.Fallback>
			</Avatar.Root>
		</div>
		<div class={alignment == 'right' ? 'w-1/2 pr-16' : 'w-1/2 pl-16'}>
			<Card.Root class={cn(extension.current ? 'bg-secondary/50' : 'bg-destructive/10', 'p-0')}>
				<Card.Content class="space-y-3 p-5">
					<div class="flex justify-between gap-2">
						<div class="flex flex-row-reverse items-start justify-end gap-2">
							<Icon
								icon={extension.current ? 'ph:check-circle' : 'ph:minus-circle'}
								class={cn(extension.current ? 'text-green-500' : 'text-red-500', 'size-6')}
							/>
							<div class="space-y-1">
								<h3 class="text-sm font-bold">{extension.name}</h3>
								<p class="text-xs text-muted-foreground">
									{extension.current ? extension.current.version : extension.latest?.version}
								</p>
							</div>
						</div>
					</div>

					<p class="text-sm font-light text-muted-foreground">
						{extension.description}
					</p>
				</Card.Content>
			</Card.Root>
		</div>
		<div class="w-1/2"></div>
	</div>
</div>

<script lang="ts" module>
	import { createClient, type Transport } from '@connectrpc/connect';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import { writable, type Writable } from 'svelte/store';

	import type { Plugin } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import { OrchestratorService } from '$lib/api/orchestrator/v1/orchestrator_pb';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as Card from '$lib/components/ui/card';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		scope,
		facility,
		plugin,
		alignment,
	}: { scope: string; facility: string; plugin: Plugin; alignment: 'left' | 'right' } = $props();

	const transport: Transport = getContext('transport');
	const orchestratorClient = createClient(OrchestratorService, transport);

	const modelPlugins: Writable<Plugin[]> = writable([]);
	const generalPlugins: Writable<Plugin[]> = writable([]);

	orchestratorClient
		.listModelPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			modelPlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch model plugins:', error);
		});
	orchestratorClient
		.listGeneralPlugins({ scope: scope, facility: facility })
		.then((respoonse) => {
			generalPlugins.set(respoonse.plugins);
		})
		.catch((error) => {
			console.error('Failed to fetch general plugins:', error);
		});
</script>

<div class="relative mx-auto hidden w-full space-y-12 pt-10 md:block">
	<div
		class={cn(
			alignment == 'right'
				? 'relative flex flex-row-reverse items-center gap-8'
				: 'relative flex items-center gap-8',
		)}
	>
		<div
			class="bg-muted/70 ring-muted/30 absolute left-1/2 z-10 flex size-13 -translate-x-1/2 transform items-center justify-center rounded-full border font-bold ring-1 ring-offset-1"
		>
			<Avatar.Root class="size-9">
				<Avatar.Image src={plugin.latest?.icon} />
				<Avatar.Fallback>
					<Icon icon="ph:puzzle-piece" class="size-6" />
				</Avatar.Fallback>
			</Avatar.Root>
		</div>
		<div class={alignment == 'right' ? 'w-1/2 pr-16' : 'w-1/2 pl-16'}>
			<Card.Root class={cn(plugin.current ? 'bg-secondary/50' : 'bg-destructive/10', 'p-0')}>
				<Card.Content class="space-y-3 p-5">
					<div class="flex justify-between gap-2">
						<div class="flex flex-row-reverse items-start justify-end gap-2">
							<Icon
								icon={plugin.current ? 'ph:check-circle' : 'ph:minus-circle'}
								class={cn(plugin.current ? 'text-green-500' : 'text-red-500', 'size-6')}
							/>
							<div class="space-y-1">
								<h3 class="text-base font-bold">{plugin.latest?.name}</h3>
								<p class="text-muted-foreground text-xs">
									{plugin.latest?.version}
								</p>
							</div>
						</div>
						{#if plugin.latest?.ref}
							<Tooltip.Provider>
								<Tooltip.Root>
									<Tooltip.Trigger>
										<Icon icon="ph:archive-bold" class="size-4" />
									</Tooltip.Trigger>
									<Tooltip.Content>
										{plugin.latest?.ref}
									</Tooltip.Content>
								</Tooltip.Root>
							</Tooltip.Provider>
						{/if}
					</div>

					<p class="text-muted-foreground text-sm font-light">
						{plugin.latest?.description}
					</p>
				</Card.Content>
			</Card.Root>
		</div>
		<div class="w-1/2"></div>
	</div>
</div>

<script lang="ts" module>
	import Label from '$lib/components/ui/label/label.svelte';
	import Switch from '$lib/components/ui/switch/switch.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import type { ReloadManager } from './utils.svelte';
</script>

<script lang="ts">
	import * as Select from '$lib/components/ui/select';

	let { reloadManager }: { reloadManager: ReloadManager } = $props();

	let interval = $state('Normal');
</script>

<div class="flex items-center justify-start gap-2">
	<div class="bg-muted relative flex h-9 w-9 items-center justify-center rounded-sm p-2">
		{#if reloadManager.isReloading}
			<div class="absolute" transition:fade={{ duration: 100 }}>
				<Icon
					icon="ph:clock-countdown"
					class={cn(
						'size-5 animate-spin',
						'transition-all duration-300',
						reloadManager.interval && reloadManager.interval >= 5 * 60
							? 'duration-1000'
							: reloadManager.interval && reloadManager.interval >= 1 * 60
								? 'duration-700'
								: 'duration-500'
					)}
				/>
			</div>
		{:else}
			<div class="absolute" transition:fade={{ duration: 500 }}>
				<Icon icon="ph:clock-countdown-fill" class={cn('text-muted-foreground size-5')} />
			</div>
		{/if}
	</div>
	<Label class="bg-muted flex h-9 items-center justify-center rounded-sm p-2">Switch</Label>
	<Switch
		bind:checked={reloadManager.state}
		onCheckedChange={() => {
			if (reloadManager.state) {
				reloadManager.restart();
			} else {
				reloadManager.stop();
			}
		}}
	/>
	<Label class="bg-muted flex h-9 items-center justify-center rounded-sm p-2">Interval</Label>
	<Select.Root type="single" bind:value={interval}>
		<Select.Trigger class="w-fit">{interval ?? 'Select'}</Select.Trigger>
		<Select.Content>
			<Select.Item
				value="Slow"
				class="text-xs"
				onclick={() => {
					reloadManager.interval = 5 * 60;
					reloadManager.restart();
				}}
			>
				Slow
			</Select.Item>
			<Select.Item
				value="Normal"
				class="text-xs"
				onclick={() => {
					reloadManager.interval = 1 * 60;
					reloadManager.restart();
				}}
			>
				Normal
			</Select.Item>
			<Select.Item
				value="Fast"
				class="text-xs"
				onclick={() => {
					reloadManager.interval = 5;
					reloadManager.restart();
				}}
			>
				Fast
			</Select.Item>
		</Select.Content>
	</Select.Root>
</div>

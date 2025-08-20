<script lang="ts">
	import * as ToggleGroup from '$lib/components/ui/toggle-group/index.js';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import { Speed, TimerManager } from './utils.svelte';

	let { timerManager }: { timerManager: TimerManager } = $props();
	let interval = $state(timerManager.interval?.toString());

	timerManager.start();
</script>

<div class="flex items-center justify-start gap-2">
	<div class="bg-muted relative flex h-9 w-9 items-center justify-center rounded-sm p-2">
		{#if timerManager.isProcessing}
			<div class="absolute" transition:fade={{ duration: 100 }}>
				<Icon
					icon="ph:clock-countdown"
					class={cn('size-5 animate-spin transition duration-2000')}
				/>
			</div>
		{:else}
			<div class="absolute" transition:fade={{ duration: 500 }}>
				<Icon icon="ph:clock-countdown-fill" class="text-muted-foreground size-5" />
			</div>
		{/if}
	</div>

	<ToggleGroup.Root variant="outline" type="single" bind:value={interval}>
		<ToggleGroup.Item
			value={Speed.SLOW.toString()}
			onclick={() => {
				timerManager.interval =
					timerManager.interval && timerManager.interval === Speed.SLOW ? undefined : Speed.SLOW;
				timerManager.restart();
			}}
		>
			<Icon icon="ph:cell-signal-none" class="size-4" />
		</ToggleGroup.Item>
		<ToggleGroup.Item
			value={Speed.NORMAL.toString()}
			onclick={() => {
				timerManager.interval =
					timerManager.interval && timerManager.interval === Speed.NORMAL
						? undefined
						: Speed.NORMAL;
				timerManager.restart();
			}}
		>
			<Icon icon="ph:cell-signal-medium" class="size-4" />
		</ToggleGroup.Item>
		<ToggleGroup.Item
			value={Speed.FAST.toString()}
			onclick={() => {
				timerManager.interval =
					timerManager.interval && timerManager.interval === Speed.FAST ? undefined : Speed.FAST;
				timerManager.restart();
			}}
		>
			<Icon icon="ph:cell-signal-full" class="size-4" />
		</ToggleGroup.Item>
	</ToggleGroup.Root>
</div>

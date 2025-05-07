<script lang="ts">
	import { onMount } from 'svelte';
	import Icon from '@iconify/svelte';
	import * as Alert from '$lib/components/ui/alert/index.js';
	let { data, duration }: { data: any[]; duration: number } = $props();
	let interval: ReturnType<typeof setInterval>;
	let currentMessageIndex = $state(0);
	function startInterval() {
		interval = setInterval(() => {
			currentMessageIndex = (currentMessageIndex + 1) % data.length;
		}, duration);
	}
	onMount(() => {
		startInterval();
		return () => clearInterval(interval);
	});
</script>

<div>
	{#each data as datum, index (index)}
		{#if currentMessageIndex % data.length === index}
			<Alert.Root
				class="flex w-full items-center justify-between gap-4 hover:cursor-grabbing"
				onmouseenter={() => clearInterval(interval)}
				onmouseleave={() => startInterval()}
			>
				<span class="flex items-center gap-2">
					<Icon icon="ph:info" class="size-7" />
					<span>
						<Alert.Title class="text-sm">{datum.message}</Alert.Title>
						<Alert.Description class="text-xs">{datum.details}</Alert.Description>
					</span>
				</span>
				<span class="flex flex-col items-center justify-between gap-2">
					<Icon
						icon="ph:caret-up"
						class="hover:cursor-pointer"
						onclick={() => {
							currentMessageIndex = (currentMessageIndex - 1 + data.length) % data.length;
						}}
					/>
					<Icon
						icon="ph:caret-down"
						class="hover:cursor-pointer"
						onclick={() => {
							currentMessageIndex = (currentMessageIndex + 1) % data.length;
						}}
					/>
				</span>
			</Alert.Root>
		{/if}
	{/each}
</div>

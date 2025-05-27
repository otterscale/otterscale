<script lang="ts">
	import Icon from '@iconify/svelte';
	import { fade } from 'svelte/transition';
	import { fromDate, getLocalTimeZone, toTime } from '@internationalized/date';

	import TimePicker from './time-picker.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';

	let { value = $bindable() }: { value: Date } = $props();

	let time = $state(toTime(fromDate(value, getLocalTimeZone())));
	const isUpdated = $derived(
		toTime(fromDate(value, getLocalTimeZone())).toString() !== time.toString()
	);
	function update() {
		let temporaryTime = new Date(value);
		temporaryTime.setHours(time.hour);
		temporaryTime.setMinutes(time.minute);
		temporaryTime.setSeconds(time.second);
		value = temporaryTime;
	}

	let hourReference = $state<HTMLInputElement | null>(null);
	let minuteReference = $state<HTMLInputElement | null>(null);
	let secondReference = $state<HTMLInputElement | null>(null);
</script>

<span class="flex items-center gap-1">
	{#if isUpdated}
		<div class="relative h-6 w-6">
			<div class="absolute inset-0" in:fade|local>
				<Tooltip.Provider>
					<Tooltip.Root delayDuration={0}>
						<Tooltip.Trigger>
							<button onclick={update}>
								<Icon icon="ph:clock-clockwise" class="size-6 animate-pulse" />
							</button>
						</Tooltip.Trigger>
						<Tooltip.Content>
							<p>Click to update timestamp</p>
						</Tooltip.Content>
					</Tooltip.Root>
				</Tooltip.Provider>
			</div>
		</div>
	{:else}
		<div class="relative h-6 w-6">
			<div class="absolute inset-0" in:fade|local>
				<Icon icon="ph:clock" class="size-6" />
			</div>
		</div>
	{/if}

	<TimePicker bind:ref={hourReference} picker="hours" bind:time />
	<p>:</p>
	<TimePicker bind:ref={minuteReference} picker="minutes" bind:time />
	<p>:</p>
	<TimePicker bind:ref={secondReference} picker="seconds" bind:time />
</span>

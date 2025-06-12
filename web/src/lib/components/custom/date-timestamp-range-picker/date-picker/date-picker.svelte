<script lang="ts">
	import Icon from '@iconify/svelte';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import Button from '$lib/components/ui/button/button.svelte';

	import { Calendar } from '$lib/components/ui/calendar/index.js';
	import { fromDate, getLocalTimeZone, toCalendarDate } from '@internationalized/date';

	let { value: value = $bindable() }: { value: Date } = $props();

	let date = $state(toCalendarDate(fromDate(value, getLocalTimeZone())));
</script>

<span class="flex items-center gap-1">
	<Icon icon="ph:calendar-blank" class="size-6" />
	<Popover.Root>
		<Popover.Trigger>
			<Button variant="outline" class="bg-transparent font-mono">{date}</Button>
		</Popover.Trigger>
		<Popover.Content class="w-fit p-0">
			<Calendar
				bind:value={date}
				type="single"
				preventDeselect
				onValueChange={() => {
					if (toCalendarDate(fromDate(value, getLocalTimeZone())) !== date) {
						value = date.toDate(getLocalTimeZone());
					}
				}}
			/>
		</Popover.Content>
	</Popover.Root>
</span>

<script lang="ts">
	import { RangeCalendar } from '$lib/components/ui/range-calendar/index.js';
	import { fromDate, getLocalTimeZone, toCalendarDate } from '@internationalized/date';

	let {
		start = $bindable(),
		end = $bindable()
	}: {
		start: Date;
		end: Date;
		class?: string;
	} = $props();

	let value = $state({
		start: toCalendarDate(fromDate(start, getLocalTimeZone())),
		end: toCalendarDate(fromDate(end, getLocalTimeZone()))
	});
</script>

<RangeCalendar
	bind:value
	class="p-0"
	onValueChange={() => {
		start = value.start.toDate(getLocalTimeZone());
		end = value.end.toDate(getLocalTimeZone());
	}}
/>

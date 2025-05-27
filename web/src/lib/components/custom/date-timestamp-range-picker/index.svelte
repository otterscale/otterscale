<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog/index';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '@iconify/svelte';
	import { getLocalTimeZone, today } from '@internationalized/date';
	import { DialogStateController } from '../utils.svelte';
	import { RangeDatePicker } from './range-date-picker';
	import { RangeDateTimePicker } from './range-datetime-picker';
	import { DatePicker } from './date-picker';
	import { TimestampPicker } from './timestamp-picker';
	import { type TimeRange } from './index';

	let { value = $bindable() }: { value: TimeRange } = $props();

	let controller = $state(new DialogStateController(false));

	const duration = 7;
	const defaultEnd = today(getLocalTimeZone());
	const defaultStart = defaultEnd.subtract({ days: duration });

	let start = $state(defaultStart.toDate(getLocalTimeZone()));
	let end = $state(defaultEnd.toDate(getLocalTimeZone()));
	function setValue() {
		value.start = start;
		value.end = end;
	}
</script>

<p class="flex h-8 items-center rounded-lg bg-muted p-4">Time Interval</p>
<AlertDialog.Root bind:open={controller.state}>
	<AlertDialog.Trigger>
		<Button variant="outline" class="hover:cursor-pointer">
			<Icon icon="ph:calendar" />
		</Button>
	</AlertDialog.Trigger>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<div class="flex flex-col items-center justify-between gap-2 rounded-lg bg-muted p-4">
				<div class="flex w-full items-center justify-between gap-2">
					<Badge variant="outline">Start</Badge>
					<div class="**:font-mono flex items-center gap-2 rounded-md px-2">
						{#key start}
							<DatePicker bind:value={start} />
						{/key}
						{#key start}
							<TimestampPicker bind:value={start} />
						{/key}
					</div>
				</div>
				<div class="flex w-full items-center justify-between gap-2">
					<Badge>End</Badge>
					<div class="**:font-mono flex items-center gap-2 rounded-md px-2">
						{#key end}
							<DatePicker bind:value={end} />
						{/key}
						{#key end}
							<TimestampPicker bind:value={end} />
						{/key}
					</div>
				</div>
			</div>
		</AlertDialog.Header>
		<div class="flex justify-between gap-4">
			<div class="w-fit">
				{#key start}
					{#key end}
						<RangeDatePicker bind:start bind:end />
					{/key}
				{/key}
			</div>
			<div class="border-l text-muted-foreground"></div>
			<div class="w-fit">
				<RangeDateTimePicker bind:start bind:end />
			</div>
		</div>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
			<AlertDialog.Action
				onclick={() => {
					controller.close();
					setValue();
				}}
			>
				Confirm
			</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

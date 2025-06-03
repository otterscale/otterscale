<script lang="ts">
	import * as AlertDialog from '$lib/components/custom/alert-dialog/index';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import Icon from '@iconify/svelte';
	import { DialogStateController } from '../utils.svelte';
	import { RangeDatePicker } from './range-date-picker';
	import { RangeDateTimePicker } from './range-datetime-picker';
	import { DatePicker } from './date-picker';
	import { TimestampPicker } from './timestamp-picker';
	import { type TimeRange } from './index';

	let controller = $state(new DialogStateController(false));

	let { value = $bindable() }: { value: TimeRange } = $props();

	const DEFAULT_START = value.start;
	const DEFAULT_END = value.end;

	let start = $state(DEFAULT_START);
	let end = $state(DEFAULT_END);

	function setValue() {
		value = {
			start: start,
			end: end
		} as TimeRange;
	}
	function resetValue() {
		start = DEFAULT_START;
		end = DEFAULT_END;
	}
</script>

<span class="flex items-center gap-2">
	<p class="bg-muted flex h-8 items-center rounded-lg p-4">Range</p>
	<AlertDialog.Root bind:open={controller.state}>
		<AlertDialog.Trigger>
			<Button variant="outline" class="hover:cursor-pointer">
				<Icon icon="ph:calendar" />
			</Button>
		</AlertDialog.Trigger>
		<AlertDialog.Content>
			<AlertDialog.Header>
				<div class="bg-muted flex flex-col items-center justify-between gap-2 rounded-lg p-4">
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
				<div class="text-muted-foreground border-l"></div>
				<div class="w-fit">
					<RangeDateTimePicker bind:start bind:end />
				</div>
			</div>
			<AlertDialog.Footer>
				<AlertDialog.Cancel onclick={resetValue}>Cancel</AlertDialog.Cancel>
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
</span>

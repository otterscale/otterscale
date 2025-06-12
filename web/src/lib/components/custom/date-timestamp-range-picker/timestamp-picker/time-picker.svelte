<script lang="ts" module>
	import type { Time as TimeType } from '@internationalized/date';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLInputAttributes } from 'svelte/elements';

	export type TimePickerInputProps = WithElementRef<HTMLInputAttributes> & {
		type?: string;
		value?: string;
		name?: string;
		files?: FileList | undefined;
		picker: TimePickerType;
		time: TimeType | undefined;
		setTime?: (time: TimeType) => void;
	};
</script>

<script lang="ts">
	import { Input } from '$lib/components/ui/input';
	import { cn } from '$lib/utils';
	import { Time } from '@internationalized/date';
	import {
		type TimePickerType,
		getArrowByType,
		getDateByType,
		setDateByType
	} from './utils.svelte';

	let {
		class: className,
		type = 'number',
		value,
		files = $bindable(),

		id,
		name,
		time = $bindable(new Time(0, 0)),
		setTime,
		picker,

		onkeydown,
		onchange,

		ref = $bindable(null),

		...restProps
	}: TimePickerInputProps = $props();

	let flag = $state<boolean>(false);

	let calculatedValue = $derived.by(() => getDateByType(time, picker));

	$effect(() => {
		if (flag) {
			const timer = setTimeout(() => {
				flag = false;
			}, 2000);
			return () => clearTimeout(timer);
		}
	});

	function calculateNewValue(key: string) {
		const cursor = calculatedValue.padStart(2, '0');
		const raw = flag ? cursor[0] + key : key + cursor[1];

		let number = Number.parseInt(raw, 10);

		number =
			picker === 'hours' ? Math.min(Math.max(number, 0), 23) : Math.min(Math.max(number, 0), 59);

		return String(number).padStart(2, '0');
	}

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key === 'Tab') return;

		const isNavKey = ['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight'].includes(e.key);
		const isDigit = e.key >= '0' && e.key <= '9';

		if (isNavKey || isDigit) {
			e.preventDefault();
		}

		if (['ArrowUp', 'ArrowDown'].includes(e.key)) {
			const step = e.key === 'ArrowUp' ? 1 : -1;
			const newValue = getArrowByType(calculatedValue, step, picker);

			if (flag) flag = false;

			const temporaryTime = time.copy();
			time = setDateByType(temporaryTime, newValue, picker);
			setTime?.(time);
		}

		if (isDigit) {
			const newValue = calculateNewValue(e.key);

			flag = !flag;

			const tempTime = time.copy();
			time = setDateByType(tempTime, newValue, picker);
			setTime?.(time);
		}
	}
</script>

<Input
	bind:ref
	id={id || picker}
	name={name || picker}
	class={cn(
		'h-9 w-9 p-2 text-center font-mono text-xs tabular-nums caret-transparent outline-0 [&::-webkit-inner-spin-button]:appearance-none',
		className
	)}
	value={value || calculatedValue}
	onchange={(e) => {
		e.preventDefault();
		onchange?.(e);
	}}
	{type}
	inputmode="decimal"
	onkeydown={(e) => {
		handleKeyDown(e);
		onkeydown?.(e);
	}}
	{...restProps}
/>

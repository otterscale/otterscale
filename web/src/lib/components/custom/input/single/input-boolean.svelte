<script lang="ts" module>
	import { FormValidator } from '$lib/components/custom/form';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Select from '$lib/components/ui/select';
	import * as Switch from '$lib/components/ui/switch';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { Switch as SwitchPrimitive, type WithoutChildrenOrChild } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { BooleanOption } from './types';
	import { BORDER_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';

	const options: BooleanOption[] = [
		{ value: null, label: 'Null', icon: 'ph:empty' },
		{
			value: false,
			label: 'False',
			icon: 'ph:x'
		},
		{
			value: true,
			label: 'True',
			icon: 'ph:circle'
		}
	];
</script>

<script lang="ts">
	let {
		id,
		ref = $bindable(null),
		class: className,
		required,
		value: checked = $bindable(undefined),
		descriptor,
		format = 'switch',
		...restProps
	}: WithoutChildrenOrChild<SwitchPrimitive.RootProps> & {
		descriptor?: (v: any) => string;
		format?: 'switch' | 'checkbox';
	} = $props();

	const isNotFilled = $derived(required && [null, undefined].includes(checked));

	let proxyChecked = $state(false);

	const formValidator: FormValidator = getContext('FormValidator');
	$effect(() => {
		formValidator.set(id, isNotFilled);
	});
</script>

<div class="flex items-center gap-2">
	<div
		aria-invalid={isNotFilled}
		class={cn(
			BORDER_INPUT_CLASSNAME,
			'flex items-center gap-2 ring-1',
			isNotFilled ? 'ring-destructive' : '',
			format === 'checkbox' && 'border-none shadow-none ring-0',
			className
		)}
	>
		{#if format === 'switch'}
			<span class="pl-3">
				<Icon icon={typeToIcon['boolean']} />
			</span>
		{/if}

		{#if required}
			{@const isValid = [true, false].includes(checked)}
			{@const isNull = [null, undefined].includes(checked)}

			{#if isValid}
				{#if descriptor}
					<p class="text-muted-foreground text-xs">{descriptor(checked)}</p>
				{:else if checked === true}
					<Badge variant="default">True</Badge>
				{:else if checked === false}
					<Badge variant="outline">False</Badge>
				{/if}
			{:else if isNull}
				<Badge variant="destructive">Required</Badge>
			{:else}
				<Badge variant="destructive">Invalid</Badge>
			{/if}
		{:else}
			{@const isValid = [null, undefined, true, false].includes(checked)}

			{#if isValid}
				{#if descriptor}
					<p class="text-muted-foreground text-xs">{descriptor(checked)}</p>
				{:else if checked === true}
					<Badge variant="default">True</Badge>
				{:else if checked === false}
					<Badge variant="outline">False</Badge>
				{:else if checked === null || checked === undefined}
					<Badge variant="secondary">Null</Badge>
				{/if}
			{:else}
				<Badge variant="destructive">Invalid</Badge>
			{/if}
		{/if}
	</div>

	{#if checked === undefined}
		<Switch.Root
			bind:ref
			bind:checked={proxyChecked}
			data-slot="input-boolean"
			{...restProps}
			onCheckedChange={() => {
				checked = proxyChecked;
			}}
		/>
	{:else}
		<Switch.Root bind:ref bind:checked data-slot="input-boolean" {...restProps} />
	{/if}
</div>

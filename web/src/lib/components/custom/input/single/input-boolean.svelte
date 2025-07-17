<script lang="ts" module>
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as Select from '$lib/components/ui/select';
	import * as Switch from '$lib/components/ui/switch';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { Switch as SwitchPrimitive, type WithoutChildrenOrChild } from 'bits-ui';
	import { z } from 'zod';
	import InputValidation from './input-validation.svelte';
	import type { BooleanOption } from './types';
</script>

<script lang="ts">
	import {
		BORDER_INPUT_CLASSNAME,
		InputValidator,
		RING_INVALID_INPUT_CLASSNAME,
		RING_VALID_INPUT_CLASSNAME,
		typeToIcon
	} from './utils.svelte';

	let {
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

	const validator = new InputValidator(z.boolean());
	const validation = $derived(validator.validate(checked));

	let proxyChecked = $state(false);
</script>

<div class="flex items-center gap-2">
	{@render Viewer()}

	{#if required}
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
	{:else}
		{@const options: BooleanOption[] = [
				{ value: null, label: 'Null', icon: 'ph:empty' },
				{
					value: false,
					label: 'False',
					icon: 'ph:x'
				},
				{
					value: true,
					label: 'True',
					icon: 'ph:check'
				}
			]}
		<Select.Root type="single" bind:value={checked}>
			<Select.Trigger class="h-9 w-fit">Select</Select.Trigger>
			<Select.Content>
				<Select.Group>
					{#each options as option}
						<Select.Item value={option.value}>
							<Icon
								icon={option.icon ? option.icon : 'ph:empty'}
								class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
							/>
							{option.label}
						</Select.Item>
					{/each}
				</Select.Group>
			</Select.Content>
		</Select.Root>
	{/if}
</div>
{@render Errors()}

{#snippet Viewer()}
	{@const isNotFilled = required && [null, undefined].includes(checked)}

	<div
		class={cn(
			BORDER_INPUT_CLASSNAME,
			isNotFilled ? RING_INVALID_INPUT_CLASSNAME : RING_VALID_INPUT_CLASSNAME,
			'w-full justify-between',
			format === 'checkbox' && 'bg-muted border-none shadow-none ring-0',
			className
		)}
	>
		<span class="flex h-9 items-center gap-2">
			<span class="pl-3">
				<Icon icon={typeToIcon['boolean']} />
			</span>

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
			{:else if !required}
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
		</span>
	</div>
{/snippet}

{#snippet Errors()}
	{@const isInvalid = ![null, undefined].includes(checked) && !validation.isValid}

	{#if isInvalid}
		<InputValidation {isInvalid} errors={validation.errors} />
	{/if}
{/snippet}

<script lang="ts">
	import * as Switch from '$lib/components/ui/switch';
	import { BORDER_INPUT_CLASSNAME, typeToIcon } from './utils.svelte';
	import Icon from '@iconify/svelte';
	import InputRequired from './input-required.svelte';
	import InputValidation from './input-validation.svelte';
	import { z } from 'zod';
	import { InputValidator } from './utils.svelte';
	import { Switch as SwitchPrimitive, type WithoutChildrenOrChild } from 'bits-ui';
	import { cn } from '$lib/utils.js';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { Single as SingleSelect } from '$lib/components/custom/select';

	let {
		ref = $bindable(null),
		class: className,
		required,
		value: checked = $bindable(undefined),
		descriptor,
		...restProps
	}: WithoutChildrenOrChild<SwitchPrimitive.RootProps> & {
		descriptor?: (v: any) => string;
	} = $props();

	let proxyChecked = $state(false);

	const validator = new InputValidator(z.boolean());
	const validation = $derived(validator.validate(checked));

	const isInvalid = $derived(!(checked == null || checked == undefined) && !validation.isValid);
	const isNotFilled = $derived(required && (checked === null || checked == undefined));
</script>

<div class="flex items-center gap-2">
	<div
		class={cn(
			BORDER_INPUT_CLASSNAME,
			isNotFilled ? 'ring-destructive ring-1' : '',
			'w-full justify-between',
			className
		)}
	>
		<span class="flex h-9 items-center gap-2">
			<span class="pl-3">
				<Icon icon={typeToIcon['boolean']} />
			</span>

			{#if required}
				{@const isValid = checked === true || checked === false}
				{@const isNull = checked === null || checked === undefined}
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
				{@const isValid =
					checked === true || checked === false || checked === null || checked === undefined}
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
		{@const options: SingleSelect.OptionType[] = [
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
		<SingleSelect.Root {options} bind:value={checked}>
			<SingleSelect.Trigger class="h-9">Select</SingleSelect.Trigger>
			<SingleSelect.Content>
				<SingleSelect.Options>
					<SingleSelect.List>
						<SingleSelect.Group>
							{#each options as option}
								<SingleSelect.Item {option}>
									<Icon
										icon={option.icon ? option.icon : 'ph:empty'}
										class={cn('size-5', option.icon ? 'visibale' : 'invisible')}
									/>
									{option.label}
									<SingleSelect.Check {option} />
								</SingleSelect.Item>
							{/each}
						</SingleSelect.Group>
					</SingleSelect.List>
				</SingleSelect.Options>
			</SingleSelect.Content>
		</SingleSelect.Root>
	{/if}
</div>
{#if isInvalid}
	<InputValidation {isInvalid} errors={validation.errors} />
{/if}

<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext, hasContext } from 'svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';

	import { Input } from '$lib/components/ui/input';
	import { m } from '$lib/paraglide/messages';
	import { cn } from '$lib/utils.js';

	import type { InputType } from './types';
	import { typeToIcon } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		value = $bindable(),
		type = 'text',
		class: className,
		required,
		oninput,
		transformer = (value) => value,
		validator = () => true,
		invalid = $bindable(),
		validateRule,
		maxLength,
		...restProps
	}: WithElementRef<
		Omit<HTMLInputAttributes, 'type' | 'files'> & { type?: InputType | undefined }
	> & {
		transformer?: (value: any) => void;
		validator?: (value?: any) => boolean;
		invalid?: boolean | null | undefined;
		validateRule?:
			| 'rfc1123'
			| 'alphanum-dash'
			| 'lower-alphanum-dash'
			| 'alphanum-dash-start-alpha'
			| 'lower-alphanum-dash-start-alpha'
			| 'lower-alphanum-dash-dot';
		maxLength?: number;
	} = $props();

	const isNull = $derived(required && (value === null || value === undefined || value === ''));
	const isValidated = $derived(value ? validator(value) : true);

	const rfc1123Regex = /^[a-z0-9]([-a-z0-9.]*[a-z0-9])?$/;
	const alphanumDashRegex = /^[a-zA-Z0-9-]+$/;
	const lowerAlphanumDashRegex = /^[a-z0-9-]+$/;

	const alphanumDashStartAlphaRegex = /^[a-zA-Z][a-zA-Z0-9-]*$/;
	const lowerAlphanumDashStartAlphaRegex = /^[a-z][a-z0-9-]*$/;
	const lowerAlphanumDashDotRegex = /^[.a-z0-9][-a-z0-9.]*$/;

	const isRuleValid = $derived(
		!value ||
			!validateRule ||
			(validateRule === 'rfc1123' && rfc1123Regex.test(value)) ||
			(validateRule === 'alphanum-dash' && alphanumDashRegex.test(value)) ||
			(validateRule === 'lower-alphanum-dash' && lowerAlphanumDashRegex.test(value)) ||
			(validateRule === 'alphanum-dash-start-alpha' && alphanumDashStartAlphaRegex.test(value)) ||
			(validateRule === 'lower-alphanum-dash-start-alpha' &&
				lowerAlphanumDashStartAlphaRegex.test(value)) ||
			(validateRule === 'lower-alphanum-dash-dot' && lowerAlphanumDashDotRegex.test(value))
	);

	const isMaxLengthValid = $derived(!maxLength || !value || value.length <= maxLength);

	const isInvalid = $derived(isNull || !isValidated || !isRuleValid || !isMaxLengthValid);
	$effect(() => {
		invalid = isInvalid;
	});
</script>

<div>
	<div class="relative">
		{#if type}
			<span class="absolute top-1/2 left-3 -translate-y-1/2 items-center">
				<Icon icon={hasContext('icon') ? getContext('icon') : typeToIcon[type]} />
			</span>
		{/if}

		<Input
			bind:ref
			data-slot="input-general"
			class={cn(
				'read-only:focus-none pl-9 ring-1 read-only:cursor-default',
				isInvalid
					? 'text-destructive/60 placeholder:text-xs placeholder:text-destructive/60 focus:text-foreground focus:placeholder:invisible'
					: '',
				isInvalid ? 'ring-destructive focus:ring-primary' : '',
				className
			)}
			{type}
			bind:value
			placeholder={isInvalid ? 'Required' : ''}
			oninput={(e) => {
				value = transformer(value);
				oninput?.(e);
			}}
			{...restProps}
		/>
	</div>
	{#if validateRule === 'rfc1123' && !isRuleValid && value}
		<p class="mt-1 text-[12px] font-medium text-destructive">
			{m.validate_rfc1123()}
		</p>
	{/if}
	{#if validateRule === 'alphanum-dash' && !isRuleValid && value}
		<p class="mt-1 text-[12px] font-medium text-destructive">
			{m.validate_alphanum_dash()}
		</p>
	{/if}
	{#if validateRule === 'lower-alphanum-dash' && !isRuleValid && value}
		<p class="mt-1 text-[12px] font-medium text-destructive">
			{m.validate_lower_alphanum_dash()}
		</p>
	{/if}
	{#if validateRule === 'alphanum-dash-start-alpha' && !isRuleValid && value}
		<p class="mt-1 text-[12px] font-medium text-destructive">
			{m.validate_alphanum_dash_start_alpha()}
		</p>
	{/if}
	{#if validateRule === 'lower-alphanum-dash-start-alpha' && !isRuleValid && value}
		<p class="mt-1 text-[12px] font-medium text-destructive">
			{m.validate_lower_alphanum_dash_start_alpha()}
		</p>
	{/if}
	{#if validateRule === 'lower-alphanum-dash-dot' && !isRuleValid && value}
		<p class="mt-1 text-[12px] font-medium text-destructive">
			{m.validate_lower_alphanum_dash_dot()}
		</p>
	{/if}
	{#if !isMaxLengthValid && value}
		<p class="mt-1 text-[12px] font-medium text-destructive">
			{m.validate_max_length({ max: maxLength })}
		</p>
	{/if}
</div>

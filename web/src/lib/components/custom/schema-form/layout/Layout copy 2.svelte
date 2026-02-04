<script lang="ts">
	import { getFormContext, layoutAttributes, uiOptionProps, type ComponentProps } from '@sjsf/form';
	import { cn } from '$lib/utils';
	import { tv } from 'tailwind-variants';

	const props: ComponentProps['layout'] = $props();
	const { type, config, children, errors } = props;

	const ctx = getFormContext();
	const attributes = $derived(layoutAttributes(ctx, config, 'layout', 'layouts', type, {}));

	const isMeta = $derived(
		type === 'field-meta' || type === 'array-field-meta' || type === 'object-field-meta'
	);

	const isItem = $derived(type === 'array-item');
	const isGrowable = $derived(
		type === 'array-item-content' ||
			type === 'object-property-key-input' ||
			type === 'object-property-content'
	);
	const isObjectProperty = $derived(type === 'object-property');
	const isMultiFieldControls = $derived(type === 'multi-field-controls');

	// --- Inline Variants from shadcn4-theme ---

	const buttonGroupVariants = tv({
		base: "flex w-fit items-stretch has-[>[data-slot=button-group]]:gap-2 [&>*]:focus-visible:relative [&>*]:focus-visible:z-10 has-[select[aria-hidden=true]:last-child]:[&>[data-slot=select-trigger]:last-of-type]:rounded-r-md [&>[data-slot=select-trigger]:not([class*='w-'])]:w-fit [&>input]:flex-1",
		variants: {
			orientation: {
				horizontal:
					"[&>*:not(:first-child)]:rounded-l-none [&>*:not(:first-child)]:border-l-0 [&>*:not(:last-child)]:rounded-r-none",
				vertical:
					"flex-col [&>*:not(:first-child)]:rounded-t-none [&>*:not(:first-child)]:border-t-0 [&>*:not(:last-child)]:rounded-b-none",
			},
		},
		defaultVariants: {
			orientation: "horizontal",
		},
	});

	const fieldVariants = tv({
		base: "group/field data-[invalid=true]:text-destructive flex w-full gap-3",
		variants: {
			orientation: {
				vertical: "flex-col [&>*]:w-full [&>.sr-only]:w-auto",
				horizontal: [
					"flex-row items-center",
					"[&>[data-slot=field-label]]:flex-auto",
					"has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px has-[>[data-slot=field-content]]:items-start",
				],
				responsive: [
					"@md/field-group:flex-row @md/field-group:items-center @md/field-group:[&>*]:w-auto flex-col [&>*]:w-full [&>.sr-only]:w-auto",
					"@md/field-group:[&>[data-slot=field-label]]:flex-auto",
					"@md/field-group:has-[>[data-slot=field-content]]:items-start @md/field-group:has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px",
				],
			},
		},
		defaultVariants: {
			orientation: "vertical",
		},
	});

	// Helper to extract specific props for variants
	const getButtonGroupProps = (attrs: any) => {
		const { orientation, class: className, ...rest } = uiOptionProps('shadcn4ButtonGroup')(attrs, config, ctx);
		return { orientation, className, rest };
	};

	const getFieldProps = (attrs: any) => {
		const { orientation, class: className, ...rest } = uiOptionProps('shadcn4Field')(attrs, config, ctx);
		return { orientation, className, rest };
	};

    const getFieldSetProps = (attrs: any) => {
         const { class: className, ...rest } = uiOptionProps('shadcn4FieldSet')(attrs, config, ctx);
         return { className, rest };
    };

</script>

{#if (type === 'field-content' || isMeta) && Object.keys(attributes).length < 2}
	{@render children()}
{:else if type === 'array-item-controls'}
	{@const { orientation, className, rest } = getButtonGroupProps(attributes)}
	<div
		role="group"
		data-slot="button-group"
		data-orientation={orientation || 'horizontal'}
		class={cn(buttonGroupVariants({ orientation: orientation || 'horizontal' }), className)}
		{...rest}
	>
		{@render children()}
	</div>
{:else if type === 'array-field' || type === 'object-field'}
    {@const { className, rest } = getFieldSetProps(attributes)}
	<fieldset
		data-slot="field-set"
		class={cn(
			"flex flex-col gap-6",
			"has-[>[data-slot=checkbox-group]]:gap-3 has-[>[data-slot=radio-group]]:gap-3",
			className
		)}
		{...rest}
	>
		{@render children()}
	</fieldset>
{:else if type === 'field'}
	{@const { orientation, className, rest } = getFieldProps(attributes)}
	<div
		role="group"
		data-slot="field"
		data-orientation={orientation || 'vertical'}
		class={cn(fieldVariants({ orientation: orientation || 'vertical' }), className)}
		data-invalid={errors.length > 0}
		{...rest}
	>
		{@render children()}
	</div>
{:else if type === 'field-title-row'}
	<div class="flex w-full items-center justify-between" {...attributes}>
		{@render children()}
	</div>
{:else if type === 'array-field-title-row' || type === 'object-field-title-row'}
	<legend
		data-slot="field-legend"
		data-variant="legend"
		class="mb-3 font-medium data-[variant=legend]:text-base data-[variant=label]:text-sm flex w-full items-center justify-between"
		{...attributes}
	>
		{@render children()}
	</legend>
{:else if type === 'array-items' || type === 'object-properties' || type === 'multi-field' || type === 'multi-field-content'}
	<!-- FieldGroup -->
	<div
		data-slot="field-group"
		class={cn(
			"group/field-group @container/field-group flex w-full flex-col gap-7 data-[slot=checkbox-group]:gap-3 [&>[data-slot=field-group]]:gap-4",
			attributes.class
		)}
		{...attributes}
	>
		{@render children()}
	</div>
{:else}
	<div
		class={cn({
			grow: isGrowable,
			'flex items-center gap-2': isMultiFieldControls,
			'flex items-start gap-1.5': isItem,
			'grid grid-cols-1 grid-rows-[1fr] items-start gap-x-1.5 [&:has(>:nth-child(2))]:grid-cols-[1fr_1fr_auto]':
				isObjectProperty
		})}
		{...attributes}
	>
		{@render children()}
	</div>
{/if}

<script lang="ts">
	import {
		Text,
		getComponent,
		getFieldAction,
		getFormContext,
		retrieveTranslate,
		retrieveUiOption,
		type ComponentProps
	} from '@sjsf/form';

	import { createObjectContext, setObjectContext } from '@sjsf/form/fields/object/context.svelte';

	const field = 'objectField';

	const ctx = getFormContext();

	let { config, value = $bindable(), uiOption, translate }: ComponentProps[typeof field] = $props();

	const objCtx = createObjectContext({
		ctx,
		config: () => config,
		value: () => value,
		setValue: (v) => (value = v),
		translate
	});
	setObjectContext(objCtx);

	const ObjectProperty = $derived(getComponent(ctx, 'objectPropertyField', config));
	const Template = $derived(getComponent(ctx, 'objectTemplate', config));
	const Button = $derived(getComponent(ctx, 'button', config));

	const action = $derived(getFieldAction(ctx, config, field));
</script>

{#snippet addButton()}
	<Button
		type="object-property-add"
		{config}
		errors={objCtx.errors()}
		disabled={false}
		onclick={objCtx.addProperty}
	>
		<Text {config} id="add-object-property" {translate} />
	</Button>
{/snippet}
{#snippet renderAction()}
	{@render action?.(
		ctx,
		config,
		{
			get current() {
				return value;
			},
			set current(v) {
				value = v;
			}
		},
		objCtx.errors()
	)}
{/snippet}
<Template
	type="template"
	{value}
	{config}
	{uiOption}
	errors={objCtx.errors()}
	addButton={objCtx.canExpand() ? addButton : undefined}
	action={action && renderAction}
>
	{#each objCtx.propertiesOrder() as property (property)}
		{@const isAdditional = objCtx.isAdditionalProperty(property)}
		{@const cfg = objCtx.propertyConfig(config, property, isAdditional)}
		<ObjectProperty
			type="field"
			{property}
			{isAdditional}
			bind:value={
				() => value?.[property],
				(v) => {
					const c = value;
					if (!c) {
						value = { [property]: v };
					} else {
						c[property] = v;
					}
				}
			}
			config={cfg}
			uiOption={(opt) => retrieveUiOption(ctx, cfg, opt)}
			translate={retrieveTranslate(ctx, cfg)}
		/>
	{/each}
</Template>

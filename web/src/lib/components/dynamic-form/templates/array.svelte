<script lang="ts">
	import { getComponent, getFormContext, type ComponentProps } from '@sjsf/form';

	import { getTemplateProps } from '@sjsf/form/templates/get-template-props';

	const ctx = getFormContext();

	const templateType = 'arrayTemplate';

	const {
		children,
		addButton,
		action,
		uiOption,
		config,
		errors
	}: ComponentProps[typeof templateType] = $props();

	const Layout = $derived(getComponent(ctx, 'layout', config));
	const Title = $derived(getComponent(ctx, 'title', config));
	const Description = $derived(getComponent(ctx, 'description', config));
	const ErrorsList = $derived(getComponent(ctx, 'errorsList', config));

	const { title, description, showMeta } = $derived(getTemplateProps(uiOption, config));
</script>

<div class="*:data-[layout=array-field]:gap-3">
	<Layout type="array-field" {config} {errors}>
		{#if showMeta && (title || description)}
			<Layout type="array-field-meta" {config} {errors}>
				{#if title}
					<Layout type="array-field-title-row" {config} {errors}>
						<Title {templateType} {title} {config} {errors} />
						{@render action?.()}
					</Layout>
				{/if}
				{#if description}
					<Description {templateType} {description} {config} {errors} />
				{/if}
			</Layout>
		{/if}
		<Layout type="array-items" {config} {errors}>
			{@render children()}
		</Layout>
		{@render addButton?.()}
		{#if errors.length > 0}
			<ErrorsList {errors} {config} />
		{/if}
	</Layout>
</div>

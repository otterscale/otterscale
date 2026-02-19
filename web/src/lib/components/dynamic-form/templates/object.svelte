<script lang="ts">
	import { getComponent, getFormContext, type ComponentProps } from '@sjsf/form';

	import { getTemplateProps } from '@sjsf/form/templates/get-template-props';

	const templateType = 'objectTemplate';

	const {
		config,
		children,
		addButton,
		action,
		errors,
		uiOption
	}: ComponentProps[typeof templateType] = $props();

	const ctx = getFormContext();

	const Layout = $derived(getComponent(ctx, 'layout', config));
	const Title = $derived(getComponent(ctx, 'title', config));
	const Description = $derived(getComponent(ctx, 'description', config));
	const ErrorsList = $derived(getComponent(ctx, 'errorsList', config));

	const { title, description, showMeta } = $derived(getTemplateProps(uiOption, config));
</script>

<div class="*:data-[layout=object-field]:gap-3">
	<Layout type="object-field" {config} {errors}>
		{#if showMeta && (title || description)}
			<Layout type="object-field-meta" {config} {errors}>
				{#if title}
					<div class="*:data-[layout=object-field-title-row]:mb-0">
						<Layout type="object-field-title-row" {config} {errors}>
							<Title {templateType} {title} {config} {errors} />
							{@render action?.()}
						</Layout>
					</div>
				{/if}
				{#if description}
					<Description {templateType} {description} {config} {errors} />
				{/if}
			</Layout>
		{/if}
		<Layout type="object-properties" {config} {errors}>
			{@render children()}
		</Layout>
		{@render addButton?.()}
		{#if errors.length > 0}
			<ErrorsList {errors} {config} />
		{/if}
	</Layout>
</div>

<script lang="ts">
	import { getComponent, getFormContext, type ComponentProps } from '@sjsf/form';

	const { children, optionSelector, config, errors, action }: ComponentProps['multiFieldTemplate'] =
		$props();

	const ctx = getFormContext();

	const Layout = $derived(getComponent(ctx, 'layout', config));
</script>

<div class="*:data-[layout=multi-field]:gap-3">
	<Layout type="multi-field" {config} {errors}>
		<div class="flex items-start gap-3">
			<Layout type="multi-field-content" {config} {errors}>
				{@render children()}
			</Layout>
			<Layout type="multi-field-controls" {config} {errors}>
				{@render optionSelector()}
				{@render action?.()}
			</Layout>
		</div>
	</Layout>
</div>

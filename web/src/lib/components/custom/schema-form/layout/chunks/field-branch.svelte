<script lang="ts">
	import { type ComponentProps, getFormContext, uiOptionProps } from '@sjsf/form';
	import { getThemeContext } from '@sjsf/shadcn4-theme';

	type Props = ComponentProps['layout'] & { attributes: any };
	const { children, config, attributes, errors }: Props = $props();

	const ctx = getFormContext();
	const themeCtx = getThemeContext();
	const { Field } = $derived(themeCtx.components);
</script>

<Field data-invalid={errors.length > 0} {...uiOptionProps('shadcn4Field')(attributes, config, ctx)}>
	{@render children()}
</Field>

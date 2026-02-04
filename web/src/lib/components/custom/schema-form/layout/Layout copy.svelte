<script lang="ts">
	import { getFormContext, layoutAttributes, type ComponentProps } from '@sjsf/form';
	import ArrayItemControls from './chunks/ArrayItemControls.svelte';
	import DefaultBranch from './chunks/DefaultBranch.svelte';
	import FieldBranch from './chunks/FieldBranch.svelte';
	import FieldGroupBranch from './chunks/FieldGroupBranch.svelte';
	import FieldLegendBranch from './chunks/FieldLegendBranch.svelte';
	import FieldSetBranch from './chunks/FieldSetBranch.svelte';
	import FieldTitleRow from './chunks/FieldTitleRow.svelte';
	import SimpleContent from './chunks/SimpleContent.svelte';

	const props: ComponentProps['layout'] = $props();
	const { type, config } = props;

	const ctx = getFormContext();
	const attributes = $derived(layoutAttributes(ctx, config, 'layout', 'layouts', type, {}));

	const isMeta = $derived(
		type === 'field-meta' || type === 'array-field-meta' || type === 'object-field-meta'
	);

	const BRANCHES = {
		simple: SimpleContent,
		'array-item-controls': ArrayItemControls,
		fieldset: FieldSetBranch,
		field: FieldBranch,
		'field-title-row': FieldTitleRow,
		'field-legend': FieldLegendBranch,
		'field-group': FieldGroupBranch,
		default: DefaultBranch
	};

	function getBranchType(
		t: string,
		meta: boolean,
		attrs: any
	): keyof typeof BRANCHES {
		if ((t === 'field-content' || meta) && Object.keys(attrs).length < 2) {
			return 'simple';
		}

		const mapping: Record<string, keyof typeof BRANCHES> = {
			'array-item-controls': 'array-item-controls',
			'array-field': 'fieldset',
			'object-field': 'fieldset',
			field: 'field',
			'field-title-row': 'field-title-row',
			'array-field-title-row': 'field-legend',
			'object-field-title-row': 'field-legend',
			'array-items': 'field-group',
			'object-properties': 'field-group',
			'multi-field': 'field-group',
			'multi-field-content': 'field-group'
		};

		return mapping[t] || 'default';
	}

	const branchKey = $derived(getBranchType(type, isMeta, attributes));
	const Branch = $derived(BRANCHES[branchKey]);

	const branchProps = $derived(
		branchKey === 'simple' ? props : { ...props, attributes }
	);
	// const branchProps = $derived({ ...props, attributes });
</script>

<Branch {...branchProps} />

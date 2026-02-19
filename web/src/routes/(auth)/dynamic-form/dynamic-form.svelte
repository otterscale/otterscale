<script lang="ts">
	import {
		Content,
		createForm,
		Form,
		type Schema,
		setFormContext,
		type UiSchemaRoot
	} from '@sjsf/form';
	import { createFocusOnFirstError } from '@sjsf/form/focus-on-first-error';
	import { overrideByRecord } from '@sjsf/form/lib/resolver';
	import type { Ref } from '@sjsf/form/lib/svelte.svelte';

	import { arrayItemField } from '$lib/components/dynamic-form/fields';
	import * as defaults from '$lib/components/dynamic-form/defaults';

	import MultiStepField, { setStepperContext } from './multiple-step-field.svelte';

	let step = $state.raw(0);
	const stepperCtx: Ref<number> = {
		get current() {
			return step;
		},
		set current(v) {
			step = v;
		}
	};
	setStepperContext(stepperCtx);

	const schema: Schema = {
		type: 'array',
		items: [
			{
				title: 'Basic',
				type: 'object',
				properties: {
					listOfStrings: {
						type: 'array',
						title: 'A list of strings',
						items: {
							type: 'string',
							default: 'bazinga'
						}
					},
					multipleChoicesList: {
						type: 'array',
						title: 'A multiple choices list',
						items: {
							type: 'string',
							enum: ['foo', 'bar', 'fuzz', 'qux']
						},
						uniqueItems: true
					},
					fixedItemsList: {
						type: 'array',
						title: 'A list of fixed items',
						items: [
							{
								title: 'A string value',
								type: 'string',
								default: 'lorem ipsum'
							},
							{
								title: 'a boolean value',
								type: 'boolean'
							}
						],
						additionalItems: {
							title: 'Additional item',
							type: 'number'
						}
					}
				}
			},
			{
				title: 'Advanced',
				description: 'advanced form',
				type: 'object',
				properties: {
					defaultsAndMinItems: {
						type: 'array',
						title: 'List and item level defaults',
						minItems: 5,
						default: ['carp', 'trout', 'bream'],
						items: {
							type: 'string',
							default: 'unidentified'
						}
					},
					nestedList: {
						type: 'array',
						title: 'Nested list',
						items: {
							type: 'array',
							title: 'Inner list',
							items: {
								type: 'string',
								default: 'lorem ipsum'
							}
						}
					},
					unorderable: {
						title: 'Unorderable items',
						type: 'array',
						items: {
							type: 'string',
							default: 'lorem ipsum'
						}
					}
				}
			},
			{
				title: 'Expert',
				description: 'form for export',
				type: 'object',
				properties: {
					copyable: {
						title: 'Copyable items',
						type: 'array',
						items: {
							type: 'string',
							default: 'lorem ipsum'
						}
					},
					unremovable: {
						title: 'Unremovable items',
						type: 'array',
						items: {
							type: 'string',
							default: 'lorem ipsum'
						}
					},
					noToolbar: {
						title: 'No add, remove and order buttons',
						type: 'array',
						items: {
							type: 'string',
							default: 'lorem ipsum'
						}
					},
					fixedNoToolbar: {
						title: 'Fixed array without buttons',
						type: 'array',
						items: [
							{
								title: 'A number',
								type: 'number',
								default: 42
							},
							{
								title: 'A boolean',
								type: 'boolean',
								default: false
							}
						],
						additionalItems: {
							title: 'A string',
							type: 'string',
							default: 'lorem ipsum'
						}
					}
				}
			}
		]
	};

	const uiSchema = {
		'ui:components': {
			tupleField: MultiStepField
		},
		listOfStrings: {
			'ui:options': {
				translations: {
					'add-array-item': 'Add string'
				}
			},
			items: {
				'ui:options': {
					stringEmptyValue: ''
				}
			}
		},
		multipleChoicesList: {
			'ui:options': {}
		},
		unorderable: {
			'ui:options': {
				orderable: false
			}
		},
		copyable: {
			'ui:options': {
				copyable: true
			}
		},
		unremovable: {
			'ui:options': {
				removable: false
			}
		},
		noToolbar: {
			'ui:options': {
				addable: false,
				orderable: false,
				removable: false
			}
		},
		fixedNoToolbar: {
			'ui:options': {
				addable: false,
				orderable: false,
				removable: false
			}
		}
	} satisfies UiSchemaRoot;

	const theme = overrideByRecord(defaults.theme, {
		arrayItemField: arrayItemField
	});

	const form = createForm({
		...defaults,
		theme,
		schema,
		uiSchema,
		onSubmit: (data) => {
			console.log(data);
			form.reset();
			stepperCtx.current = 0;
		},
		onSubmitError(result, e, form) {
			if (result.errors.length === 0) {
				return;
			}
			step = result.errors[0].path[0] as number;
			createFocusOnFirstError()(result, e, form);
		}
	});
	setFormContext(form);
</script>

<Form attributes={{ novalidate: true }}>
	<Content />
</Form>

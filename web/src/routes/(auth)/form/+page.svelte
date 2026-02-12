<script lang="ts">
	import {
		Content,
		createForm,
		getValueSnapshot,
		type Schema,
		setFormContext,
		type UiSchema
	} from '@sjsf/form';
	import { overrideByRecord } from '@sjsf/form/lib/resolver';

	import * as defaults from '$lib/form-defaults';

	import ArrayField from './fields/array.svelte';

	const schema: Schema = {
		definitions: {
			Thing: {
				type: 'object',
				properties: {
					name: {
						type: 'string',
						default: 'Default name'
					}
				}
			}
		},
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
			},
			minItemsList: {
				type: 'array',
				title: 'A list with a minimal number of items',
				minItems: 3,
				items: {
					$ref: '#/definitions/Thing'
				}
			},
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
			},
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
	};

	const uiSchema: UiSchema = {
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
			'ui:options': {
				layouts: {
					'field-content': {
						style: 'display: flex; flex-direction: column; gap: 0.2rem'
					}
				}
			}
		},
		unorderable: {
			'ui:options': {
				orderable: false
			}
		},
		copyable: {
			'ui:options': {
				copyable: true,
				layouts: {
					'array-item': {
						style: ''
					}
				}
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
	};

	const theme = overrideByRecord(defaults.theme, {
		arrayField: ArrayField
	});

	const form = createForm({
		...defaults,
		theme,
		schema,
		uiSchema,
		onSubmit: console.log
	});

	setFormContext(form);
</script>

<Content />

<pre><code>{JSON.stringify(getValueSnapshot(form), null, 2)}</code></pre>

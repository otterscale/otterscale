<script lang="ts">
	import { BasicForm, createForm, type Schema } from '@sjsf/form';

	import * as defaults from '$lib/form-defaults';

	const schema = {
		type: 'object',
		properties: {
			checkbox: {
				type: 'object',
				properties: {
					default: {
						type: 'boolean'
					},
					error: {
						type: 'boolean'
					}
				},
				required: ['default']
			},
			checkboxes: {
				type: 'object',
				properties: {
					default: {
						type: 'array',
						items: {
							type: 'string',
							enum: ['foo', 'bar', 'fuzz', 'qux']
						},
						uniqueItems: true
					},
					error: {
						type: 'array',
						items: {
							type: 'string',
							enum: ['foo', 'bar', 'fuzz', 'qux']
						},
						uniqueItems: true
					}
				},
				required: ['default']
			},
			file: {
				type: 'object',
				properties: {
					default: {
						type: 'string',
						format: 'data-url'
					},
					error: {
						type: 'string',
						format: 'data-url'
					}
				},
				required: ['default']
			},
			multiFile: {
				type: 'object',
				properties: {
					default: {
						type: 'array',
						items: {
							type: 'string',
							format: 'data-url'
						}
					},
					error: {
						type: 'array',
						items: {
							type: 'string',
							format: 'data-url'
						}
					}
				},
				required: ['default']
			},
			number: {
				type: 'object',
				properties: {
					default: {
						type: 'number'
					},
					error: {
						type: 'number'
					}
				},
				required: ['default']
			},
			select: {
				type: 'object',
				properties: {
					default: {
						type: 'string',
						enum: ['foo', 'bar', 'fuzz', 'qux']
					},
					error: {
						type: 'string',
						enum: ['foo', 'bar', 'fuzz', 'qux']
					}
				},
				required: ['default']
			},
			text: {
				type: 'object',
				properties: {
					default: {
						type: 'string'
					},
					error: {
						type: 'string'
					}
				},
				required: ['default']
			},
			combobox: {
				type: 'object',
				properties: {
					default: {
						type: 'string',
						enum: ['foo', 'bar', 'fuzz', 'qux']
					},
					error: {
						type: 'string',
						enum: ['foo', 'bar', 'fuzz', 'qux']
					}
				},
				required: ['default']
			},
			datePicker: {
				type: 'object',
				properties: {
					default: {
						type: 'string'
					},
					error: {
						type: 'string'
					}
				},
				required: ['default']
			},
			dateRangePicker: {
				type: 'object',
				properties: {
					default: {
						type: 'object',
						properties: {
							start: {
								type: 'string'
							},
							end: {
								type: 'string'
							}
						}
					},
					error: {
						type: 'object',
						properties: {
							start: {
								type: 'string'
							},
							end: {
								type: 'string'
							}
						}
					}
				},
				required: ['default']
			},
			multiSelect: {
				type: 'object',
				properties: {
					default: {
						type: 'array',
						items: {
							type: 'string',
							enum: ['foo', 'bar', 'fuzz', 'qux']
						},
						uniqueItems: true
					},
					error: {
						type: 'array',
						items: {
							type: 'string',
							enum: ['foo', 'bar', 'fuzz', 'qux']
						},
						uniqueItems: true
					}
				},
				required: ['default']
			},
			radioButtons: {
				type: 'object',
				properties: {
					default: {
						type: 'string',
						enum: ['foo', 'bar', 'fuzz', 'qux']
					},
					error: {
						type: 'string',
						enum: ['foo', 'bar', 'fuzz', 'qux']
					}
				},
				required: ['default']
			},
			radio: {
				type: 'object',
				properties: {
					default: {
						type: 'string',
						enum: ['foo', 'bar', 'fuzz', 'qux']
					},
					error: {
						type: 'string',
						enum: ['foo', 'bar', 'fuzz', 'qux']
					}
				},
				required: ['default']
			},
			range: {
				type: 'object',
				properties: {
					default: {
						type: 'number'
					},
					error: {
						type: 'number'
					}
				},
				required: ['default']
			},
			rangeSlider: {
				type: 'object',
				properties: {
					default: {
						type: 'object',
						properties: {
							start: {
								type: 'number'
							},
							end: {
								type: 'number'
							}
						}
					},
					error: {
						type: 'object',
						properties: {
							start: {
								type: 'number'
							},
							end: {
								type: 'number'
							}
						}
					}
				},
				required: ['default']
			},
			switch: {
				type: 'object',
				properties: {
					default: {
						type: 'boolean'
					},
					error: {
						type: 'boolean'
					}
				},
				required: ['default']
			},
			textarea: {
				type: 'object',
				properties: {
					default: {
						type: 'string'
					},
					error: {
						type: 'string'
					}
				},
				required: ['default']
			}
		}
	} as const satisfies Schema;

	const form = createForm({
		...defaults,
		schema,
		onSubmit: console.log
	});
</script>

<div class="p-8"><BasicForm {form} /></div>

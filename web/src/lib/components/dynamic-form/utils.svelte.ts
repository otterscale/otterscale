import traverse from 'json-schema-traverse';
import  { JSONSchemaFaker } from 'json-schema-faker'
import { openapiSchemaToJsonSchema as openAPISchemaToJSONSchema } from '@openapi-contrib/openapi-schema-to-json-schema';

type Version = 'draft-04' | 'draft-07';
const keywords: Record<Version, Set<string>> = {
    'draft-04': new Set([
        '$ref',
        'title',
        'description',
        'default',
        'multipleOf',
        'maximum',
        'exclusiveMaximum',
        'minimum',
        'exclusiveMinimum',
        'maxLength',
        'minLength',
        'pattern',
        'additionalItems',
        'items',
        'maxItems',
        'minItems',
        'uniqueItems',
        'maxProperties',
        'minProperties',
        'required',
        'additionalProperties',
        'properties',
        'patternProperties',
        'dependencies',
        'enum',
        'type',
        'allOf',
        'anyOf',
        'oneOf',
        'not',
        'definitions'
    ]),
    'draft-07': new Set([
        '$id',
        '$schema',
        '$ref',
        '$comment',
        'title',
        'description',
        'default',
        'readOnly',
        'writeOnly',
        'examples',
        'multipleOf',
        'maximum',
        'exclusiveMaximum',
        'minimum',
        'exclusiveMinimum',
        'maxLength',
        'minLength',
        'pattern',
        'additionalItems',
        'items',
        'maxItems',
        'minItems',
        'uniqueItems',
        'contains',
        'maxProperties',
        'minProperties',
        'required',
        'additionalProperties',
        'definitions',
        'properties',
        'patternProperties',
        'dependencies',
        'propertyNames',
        'const',
        'enum',
        'type',
        'format',
        'allOf',
        'anyOf',
        'oneOf',
        'not',
        'if',
        'then',
        'else'
    ])
};

function toVersionedJSONSchema(schema: any, version: Version = 'draft-07'): any {
    const result = structuredClone(schema)
    traverse(result, {
        allKeys: true,
        cb: (currentSchema) => {
            Object.keys(currentSchema).forEach((key) => {
                if (!keywords[version].has(key)) {
                    delete currentSchema[key];
                }
            });
        }
    });
    return result;
}

function filterRequiredSchema(schema: any): any {
    const result = structuredClone(schema)
    traverse(result, {
        cb: (currentSchema) => {
            if (currentSchema.properties) {
                const requiredFields = currentSchema.required || [];
                Object.keys(currentSchema.properties).forEach(key => {
                    if (!requiredFields.includes(key)) delete currentSchema.properties[key];
                });
            }
        }
    })
    return result
}

function getMockData(schema: any) {
    return JSONSchemaFaker.generate(schema);
}

export { openAPISchemaToJSONSchema, toVersionedJSONSchema, filterRequiredSchema, getMockData };
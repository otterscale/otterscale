import type { Column } from '@tanstack/table-core';

type FieldType = {
    description: string;
    type: string;
    format?: string;
};

type FieldsType = {
    [key: string]: FieldType
};

type ValuesType = {
    [key: string]: JsonValue
}

export { FieldType, FieldsType, ValuesType }
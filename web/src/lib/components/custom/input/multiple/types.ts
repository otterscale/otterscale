type InputType =
    | 'color'
    | 'datetime-local'
    | 'date'
    | 'time'
    | 'url'
    | 'email'
    | 'tel'
    | 'text'
    | 'number'
    | 'search';

type valueSetterType = (values: any[]) => void

export type {
    InputType,
    valueSetterType
}
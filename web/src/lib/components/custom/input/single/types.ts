import type { InputType as MultipleInputType } from '../multiple/types'

type InputType =
    MultipleInputType
    | 'boolean'
    | 'password';

export type {
    InputType,
}
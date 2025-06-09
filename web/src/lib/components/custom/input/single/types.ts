import type { InputType as MultipleInputType } from '../multiple/types'

type InputType = Exclude<MultipleInputType, 'color'>

export type {
    InputType,
}
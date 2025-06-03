type OptionType = {
    value: string;
    label: string;
    icon?: string;
};

type valuesSetterType = (options: OptionType[]) => void

export type {
    OptionType,
    valuesSetterType
}
type OptionType = {
    value: string;
    label: string;
    icon?: string;
};
type valuesSetterType = (options: any[]) => void
type valuesGetterType = () => any[]

export type {
    OptionType, valuesGetterType, valuesSetterType
};

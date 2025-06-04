type OptionType = {
    value: any;
    label: string;
    icon?: string;
};

type valueSetterType = (option: OptionType) => void;

export type {
    OptionType,
    valueSetterType,
}
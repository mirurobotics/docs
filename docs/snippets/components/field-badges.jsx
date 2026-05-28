export const NullableBadge = ({ size = "sm" }) => {
    return (
        <Tooltip tip="Property does not need to be set">
            <Badge icon="circle" color="gray" size={size}>nullable</Badge>
        </Tooltip>
    );
};

export const EditableBadge = ({ size = "sm" }) => {
    return (
        <Tooltip tip="Property can be directly modified">
            <Badge icon="pencil" color="blue" size={size}>editable</Badge>
        </Tooltip>
    );
};

export const MutableBadge = ({ size = "sm" }) => {
    return (
        <Tooltip tip="Property is automatically updated by the system; cannot be modified directly">
            <Badge icon="feather" color="orange" size={size}>mutable</Badge>
        </Tooltip>
    );
};

export const ImmutableBadge = ({ size = "sm" }) => {
    return (
        <Tooltip tip="Property cannot be modified">
            <Badge icon="lock" color="gray" size={size}>immutable</Badge>
        </Tooltip>
    );
};

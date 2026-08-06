export const RequiredBadge = ({ size = "sm" }) => {
    return (
        <Tooltip tip="Required when creating a release">
            <Badge icon="asterisk" color="red" size={size}>Required</Badge>
        </Tooltip>
    );
};

export const OptionalBadge = ({ size = "sm" }) => {
    return (
        <Tooltip tip="Optional when creating a release">
            <Badge icon="circle" color="gray" size={size}>optional</Badge>
        </Tooltip>
    );
};

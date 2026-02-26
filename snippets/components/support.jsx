export const SupportedBadge = ({ size = "sm" }) => (
    <Tooltip tip="Actively maintained with bug fixes and security patches">
        <Badge icon="circle-check" color="green" size={size}>Supported</Badge>
    </Tooltip>
);

export const DeprecatedBadge = ({ size = "sm" }) => (
    <Tooltip tip="Functional but no longer receiving updates; migrate to a supported version">
        <Badge icon="clock" color="orange" size={size}>Deprecated</Badge>
    </Tooltip>
);

export const EndOfLifeBadge = ({ size = "sm" }) => (
    <Tooltip tip="No longer supported; may lose server compatibility">
        <Badge icon="ban" color="red" size={size}>End of life</Badge>
    </Tooltip>
);
/*
 * ImmutableBadge lives in /snippets/badges/immutable.mdx. MDX snippets render on the
 * server, but Mintlify only expands them on their own line, so badges used inline
 * (in headings, tables or next to another badge) stay here.
 */

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

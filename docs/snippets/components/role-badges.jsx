export const OwnerBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="You must be the owner to execute this action."
            cta="Workspace roles"
            href="/admin/users/roles"
        >
            <Badge icon="crown" color="yellow" size={size}>owner</Badge>
        </Tooltip>
    );
};

export const AdminBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="You must be the owner or an admin to execute this action."
            cta="Workspace roles"
            href="/admin/users/roles"
        >
            <Badge icon="shield-check" color="blue" size={size}>admin</Badge>
        </Tooltip>
    );
};
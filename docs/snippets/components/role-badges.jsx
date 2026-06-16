export const OwnerBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="You must be the owner to execute this action."
            cta="Workspace roles"
            href="/admin/users/access-control"
        >
            <Badge icon="crown" color="yellow" size={size}>owner</Badge>
        </Tooltip>
    );
};

export const MemberBadge = ({ size = "md" }) => {
    return (
        <Tooltip
            tip="Every workspace member can execute this action."
            cta="Workspace roles"
            href="/admin/users/access-control"
        >
            <Badge icon="user" color="gray" size={size}>member</Badge>
        </Tooltip>
    );
};

export const AdminBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="Workspace owners and administrators can execute this action."
            cta="Workspace roles"
            href="/admin/users/access-control"
        >
            <Badge icon="shield-check" color="blue" size={size}>admin</Badge>
        </Tooltip>
    );
};

export const GroupManagerBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="Members who are managers of a group that contains this resource can execute this action."
            cta="Workspace roles"
            href="/admin/users/access-control"
        >
            <Badge icon="group" color="blue" size={size}>group manager</Badge>
        </Tooltip>
    );
};

export const PublisherBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="Members with the publisher role can execute this action."
            cta="Workspace roles"
            href="/admin/users/access-control"
        >
            <Badge icon="git-merge" color="green" size={size}>publisher</Badge>
        </Tooltip>
    );
};

export const ProvisionerBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="Members with the provisioner role can execute this action."
            cta="Workspace roles"
            href="/admin/users/access-control"
        >
            <Badge icon="bot" color="purple" size={size}>provisioner</Badge>
        </Tooltip>
    );
};

export const OperatorBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="Members with the operator role can execute this action."
            cta="Workspace roles"
            href="/admin/users/access-control"
        >
            <Badge icon="wrench" color="orange" size={size}>operator</Badge>
        </Tooltip>
    );
};

export const ViewerBadge = ({ size = "md" }) => {
    return (
        <Tooltip 
            tip="Members with the viewer role can execute this action."
            cta="Workspace roles"
            href="/admin/users/access-control"
        >
            <Badge icon="eye" color="gray" size={size}>viewer</Badge>
        </Tooltip>
    );
};


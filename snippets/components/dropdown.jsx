export const DropdownItem = ({ children }) => {
    return (
        <div className="my-1 flex items-start gap-2">
            <span className="text-gray-400">•</span>
            <span className="flex-1">{children}</span>
        </div>
    );
};

export const DropdownGroup = ({ children }) => {
    const childrenArray = Array.isArray(children) ? children : [children];

    return (
        <div className="py-3">
            {childrenArray.map((child, index) => (
                <div key={index}>
                    {child}
                    {index < childrenArray.length - 1 && (
                        <div className="border-t border-white opacity-10 my-1" />
                    )}
                </div>
            ))}
        </div>
    );
};

export const Dropdown = ({
    title,
    defaultOpen = false,
    children
}) => {
    const [isOpen, setIsOpen] = React.useState(defaultOpen);

    function cn(...inputs) {
        return inputs.filter(Boolean).join(' ');
    }

    return (
        <div>
            <button
                className="w-full flex justify-between items-center py-3 bg-transparent 
                border-none cursor-pointer text-left"
                onClick={() => setIsOpen(!isOpen)}
                aria-expanded={isOpen}
            >
                <span className="font-bold text-gray-100 text-sm">
                    {title}
                </span>
                <svg
                    className={cn(
                        "transition-transform duration-200 opacity-50 flex-shrink-0",
                        isOpen && 'rotate-90'
                    )}
                    width="16"
                    height="16"
                    viewBox="0 0 16 16"
                    fill="none"
                    stroke="currentColor"
                    strokeWidth="2"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                >
                    <polyline points="6 4 10 8 6 12"></polyline>
                </svg>
            </button>

            {isOpen && <div className="pb-3">{children}</div>}
        </div>
    );
};

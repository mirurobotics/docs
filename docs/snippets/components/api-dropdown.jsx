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
                <span className="font-bold text-gray-100 text-md">
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
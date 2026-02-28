
import React from 'react';

export const CardNewTab = ({ children, ...props }) => {
  const ref = React.useRef(null);

  React.useEffect(() => {
    if (!ref.current) return;
    const anchor = ref.current.querySelector('a[href]');
    if (anchor) {
      anchor.setAttribute('target', '_blank');
      anchor.setAttribute('rel', 'noopener noreferrer');
    }
  }, []);

  return (
    <div ref={ref} style={{ display: 'contents' }}>
      <Card {...props}>
        {children}
      </Card>
    </div>
  );
};
